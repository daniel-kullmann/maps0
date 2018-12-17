package main

import (
	"github.com/gorilla/mux"
	"log"
	"mime"
	"net/http"
	"os"
	"path"
)

func main() {
	router := mux.NewRouter()
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

const (
	FileBase = "../frontend/"
)

func GetFile(w http.ResponseWriter, r *http.Request) {
	name := mux.Vars(r)["name"]
	log.Print("Serve " + name)
	if name == "" {
		name = "index.html"
	}
	_, err := os.Stat(FileBase + name)
	if err == nil {
		contentType := mime.TypeByExtension(path.Ext(name))
		//var contentType = "text/plain"
		//switch path.Ext(name) {
		//case ".html":
		//	contentType = "text/html"
		//case ".js":
		//	contentType = "application/javascript"
		//case ".css":
		//	contentType = "application/css"
		//case ".png":
		//	contentType = "image/png"
		//default:
		//	log.Print("Unexpected " + name)
		//}
		w.Header().Add("Content-Type", contentType)
		log.Print("as " + contentType)
		fh, err := os.Open(FileBase + name)
		if err != nil {
			w.Header().Add("Content-Type", "text/plain")
			w.WriteHeader(500)
			log.Print("Could not open: " + name)
		}
		buffer := make([]uint8, 10*1024)
		for {
			count, err := fh.Read(buffer)
			if err != nil || count == 0 {
				break
			}
			w.Write(buffer[0:count])
		}
		fh.Close()
	} else {
		w.Header().Add("Content-Type", "text/plain")
		w.WriteHeader(404)
		log.Print("Not found: " + name)
	}
}
