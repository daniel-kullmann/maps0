package main

import (
	"github.com/gorilla/mux"
	"log"
	"net/http"
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
	log.Fatal(http.ListenAndServe(":9191", router))
}
