package main

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"os"
)

const GpxBase = "../gpx_store_files/"

func SaveGpx(w http.ResponseWriter, r *http.Request) {
	log.Println(r)
}

func GetGpx(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	name := vars["name"]
	fh, err := os.Open(GpxBase + name)
	if err != nil {
		w.Header().Add("Content-Type", "text/plain")
		w.WriteHeader(404)
		w.Write([]byte("Could not open gpx file"))
	}
	w.Header().Add("Content-Type", "text/xml")
	w.Header().Add("Access-Control-Allow-Origin", "*")
	buffer := make([]uint8, 100*1024)
	for {
		count, err := fh.Read(buffer)
		if err != nil || count == 0 {
			break
		}
		w.Write(buffer[0:count])
	}
	fh.Close()
}

func LoadGpxList(w http.ResponseWriter, r *http.Request) {
	log.Println(r)
	fh, err := os.Open(GpxBase)
	if err != nil {
		w.Header().Add("Content-Type", "text/plain")
		w.WriteHeader(404)
		w.Write([]byte("Could not open gpx store directoy"))
	}
	names, err := fh.Readdirnames(-1)
	if err != nil {
		w.Header().Add("Content-Type", "text/plain")
		w.WriteHeader(404)
		w.Write([]byte("Could not read gpx store directoy"))
	}
	w.Header().Add("Content-Type", "application/json")
	w.Header().Add("Access-Control-Allow-Origin", "*")
	jsonResponse, err := json.Marshal(names)
	if err != nil {
		w.Header().Add("Content-Type", "text/plain")
		w.WriteHeader(404)
		w.Write([]byte("Could not marshal result"))
	}
	w.Write([]byte(jsonResponse))

}
