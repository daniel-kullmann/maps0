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

		url := fmt.Sprintf("https://%s.tile.openstreetmap.org/%s/%s/%s", parts[0], parts[1], parts[2], parts[3])
		client := http.DefaultClient
		response, err := client.Get(url)
		if err != nil {
			w.Header().Add("Content-Type", "text/plain")
			w.WriteHeader(404)
			response.Body.Close()
			log.Panic(err)
		} else {
			w.Header().Add("Content-Type", "image/png")
			w.Header().Add("Access-Control-Allow-Origin", "*")
			os.MkdirAll(dirName, 0666)
			fh, err_fh := os.Create(fileName)
			buffer := make([]uint8, 100*1024)
			for count, err := response.Body.Read(buffer); err == nil && count != 0; count, err = response.Body.Read(buffer) {
				w.Write(buffer[0:count])
				if err_fh == nil {
					fh.Write(buffer[0:count])
				}
			}
			fh.Close()
			response.Body.Close()
		}
	} else {
		w.Header().Add("Content-Type", "image/png")
		w.Header().Add("Access-Control-Allow-Origin", "*")
		log.Print(fileName + " does exist")
		fh, err := os.Open(fileName)
		if err != nil {
			w.Header().Add("Content-Type", "text/plain")
			w.WriteHeader(404)
			log.Panic(err)
		}
		buffer := make([]uint8, 100*1024)
		for count, err := fh.Read(buffer); err == nil && count != 0; count, err = fh.Read(buffer) {
			w.Write(buffer[0:count])
		}
	}

}
