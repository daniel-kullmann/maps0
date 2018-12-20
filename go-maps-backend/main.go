package main

import (
	"flag"
	"github.com/gorilla/mux"
	"log"
	"mime"
	"net/http"
	"os"
	"path"
	"strings"
)

//go:generate go-bindata -prefix ../frontend/ ../frontend/...

var (
	ConfigFileName = "${HOME}/.config/simple-offline-map/config.ini"
	CacheBaseDir   = "${HOME}/.local/share/simple-offline-map"
	TileBase       = CacheBaseDir + "/tiles/"
	FileBase       = "./frontend/"
	GpxBase        = CacheBaseDir + "/gpx/"
	DataBasePath   = CacheBaseDir + "/db.sqlite3"
)

func main() {
	router := mux.NewRouter()
	setup()
	router.HandleFunc("/tile/{s}/{z}/{x}/{y}.png", GetTile).Methods("GET")
	router.HandleFunc("/api/gpx/", LoadGpxList).Methods("GET")
	router.HandleFunc("/api/gpx/save/", SaveGpx).Methods("POST")
	router.HandleFunc("/api/gpx/get/{name}", GetGpx).Methods("GET")
	router.HandleFunc("/api/settings/token/", GetToken).Methods("GET")
	router.HandleFunc("/api/settings/", GetSettings).Methods("GET")
	router.HandleFunc("/api/settings/set_all_settings/", SetSettings).Methods("POST")
	router.HandleFunc("/{name:.*}", GetFile).Methods("GET")
	log.Fatal(http.ListenAndServe(":9191", router))
}

func FileExists(name string) bool {
	_, err := os.Stat(name)
	return !os.IsNotExist(err)
}

func ReadConfigFile() map[string]string {
	result := make(map[string]string)
	fileInfo, err := os.Stat(ConfigFileName)
	if err != nil {
		return result
	}
	fh, err := os.Open(ConfigFileName)
	if err != nil {
		return result
	}
	log.Print("Reading config from " + ConfigFileName)
	buffer := make([]uint8, fileInfo.Size())
	_, err = fh.Read(buffer)
	var content = string(buffer)
	lines := strings.Split(content, "\n")
	for index, line := range lines {
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			log.Printf("Unexpected line in config file %d: %s", index, line)
			continue
		}
		log.Print(parts)
		switch strings.Trim(parts[0], " \t\n") {
		case "FileBase":
			FileBase = os.ExpandEnv(strings.Trim(parts[1], " \t\n"))
		case "TileBase":
			TileBase = os.ExpandEnv(strings.Trim(parts[1], " \t\n"))
		case "DataBasePath":
			DataBasePath = os.ExpandEnv(strings.Trim(parts[1], " \t\n"))
		case "GpxBase":
			GpxBase = os.ExpandEnv(strings.Trim(parts[1], " \t\n"))
		default:
			log.Println("ERROR: Unknown setting + " + parts[0])
		}
	}
	return result
}

func FixFileNames() {
	ConfigFileName = os.ExpandEnv(ConfigFileName)
	CacheBaseDir = os.ExpandEnv(CacheBaseDir)
	TileBase = os.ExpandEnv(TileBase)
	FileBase = os.ExpandEnv(FileBase)
	GpxBase = os.ExpandEnv(GpxBase)
	DataBasePath = os.ExpandEnv(DataBasePath)
}

func setup() {
	// setup config variables to sensible settings
	// for every config variable (TileBase, GpxBase), check the following:
	// 1. Check whether they were given as command line flags; then use them like that
	// 2. Check whether a config file exists (${HOME}/.config/simple-map/config.ini
	// 3. use files in ${HOME}/.local/share/simple-map/
	// 4. just for FileBase: check ./frontend/
	var TileBaseFlag string
	var GpxBaseFlag string
	var DataBasePathFlag string
	var FileBaseFlag string
	flag.StringVar(&ConfigFileName, "config", "", "use this file for configuration settings")
	flag.StringVar(&TileBaseFlag, "tile", "", "use this directory as map tile cache")
	flag.StringVar(&GpxBaseFlag, "gpx", "", "use this directory as gpx file store")
	flag.StringVar(&DataBasePathFlag, "db", "", "use this file as database (sqlite3)")
	flag.StringVar(&FileBaseFlag, "files", "", "use this directory for frontend files (index.html, css files, js files, etc)")
	flag.Parse()
	ReadConfigFile()
	if TileBaseFlag != "" {
		TileBase = os.ExpandEnv(TileBaseFlag)
	}
	if GpxBaseFlag != "" {
		GpxBase = os.ExpandEnv(GpxBaseFlag)
	}
	if DataBasePathFlag != "" {
		DataBasePath = os.ExpandEnv(DataBasePathFlag)
	}
	if FileBaseFlag != "" {
		FileBase = os.ExpandEnv(FileBaseFlag)
	}
	FixFileNames()

	log.Print("INFO: FileBase set to " + FileBase)
	log.Print("INFO: TileBase set to " + TileBase)
	log.Print("INFO: GpxBase set to " + GpxBase)
	log.Print("INFO: DataBasePath set to " + DataBasePath)

	if !FileExists(CacheBaseDir) {
		os.MkdirAll(CacheBaseDir, 0755)
	}
	initDatabase()
}

func GetFile(w http.ResponseWriter, r *http.Request) {
	name := mux.Vars(r)["name"]
	log.Print("Serve " + name)
	if name == "" {
		name = "index.html"
	}
	asset, err := Asset(name)
	if err == nil {
		contentType := mime.TypeByExtension(path.Ext(name))
		w.Header().Add("Content-Type", contentType)
		w.Write(asset)
	} else {
		w.Header().Add("Content-Type", "text/plain")
		w.WriteHeader(404)
		w.Write([]byte("Not found: " + name))
	}
}
