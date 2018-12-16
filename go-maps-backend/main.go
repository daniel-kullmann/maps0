package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"os"
	"strings"
)

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/tile/{s}/{x}/{y}/{z}.png", GetTile).Methods("GET")
	log.Fatal(http.ListenAndServe(":9191", router))
}

const TileBase = "../tile_cache/"

func GetTile(w http.ResponseWriter, r *http.Request) {
	var parts []string
	parts = strings.SplitN(r.URL.RequestURI(), "/", -1)[2:]
	dirName := TileBase + strings.Join(parts[0:3], "/")
	fileName := TileBase + strings.Join(parts, "/")
	if _, err := os.Stat(fileName); os.IsNotExist(err) {

		log.Print(fileName + ": fetch from OSM")
		url := fmt.Sprintf("https://%s.tile.openstreetmap.org/%s/%s/%s", parts[0], parts[1], parts[2], parts[3])
		ua := r.Header.Get("User-Agent")
		client := http.Client{}
		request, err := http.NewRequest("GET", url, nil)
		if err != nil {
			w.Header().Add("Content-Type", "text/plain")
			w.WriteHeader(404)
			//w.Write(err.String())
		}
		request.Header.Set("User-Agent", ua)
		response, err := client.Do(request)
		if err != nil {
			w.Header().Add("Content-Type", "text/plain")
			w.WriteHeader(404)
			response.Body.Close()
			log.Print("ERROR: Can't make request to " + url)
		} else if response.StatusCode != 200 {
			w.Header().Add("Content-Type", "text/plain")
			w.WriteHeader(404)
			w.Write([]byte(response.Status))
			response.Body.Close()
			log.Print("ERROR: Request resulted in " + response.Status + ", " + url)
		} else {
			w.Header().Add("Content-Type", "image/png")
			w.Header().Add("Access-Control-Allow-Origin", "*")
			os.MkdirAll(dirName, 0777)
			fh, err_fh := os.Create(fileName)
			if err_fh != nil {
				log.Print("Can't create output file " + fileName)
			}
			size := 0
			buffer := make([]uint8, 100*1024)
			for {
				count, err := response.Body.Read(buffer)
				if err != nil || count == 0 {
					break
				}
				size += count
				if err_fh == nil {
					fh.Write(buffer[0:count])
				}
			}
			fh.Close()
			response.Body.Close()
			log.Printf("Finished retrieving %s: %d", fileName, size)
		}
	} else {
		w.Header().Add("Content-Type", "image/png")
		w.Header().Add("Access-Control-Allow-Origin", "*")
		log.Print(fileName + ": locally cached")
		fh, err := os.Open(fileName)
		if err != nil {
			w.Header().Add("Content-Type", "text/plain")
			w.WriteHeader(404)
			log.Print("ERROR: Can't read file " + fileName)
		}
		buffer := make([]uint8, 100*1024)
		for count, err := fh.Read(buffer); err == nil && count != 0; count, err = fh.Read(buffer) {
			w.Write(buffer[0:count])
		}
	}

}
