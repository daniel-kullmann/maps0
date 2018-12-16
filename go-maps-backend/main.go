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
	router.HandleFunc("/api/settings/get_all_settings/", GetSettings).Methods("GET")
	router.HandleFunc("/api/settings/set_all_settings/", SetSettings).Methods("POST")
	log.Fatal(http.ListenAndServe(":9191", router))
}

func GetSettings(w http.ResponseWriter, r *http.Request) {
}

func GetToken(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	w.Write([]byte("{\"token\": \"qpoewpfqjwekjfpqoiwjreapfqurewpfvjqlkn4ef1432lkjf1043flkj\"}"))
}

func SetSettings(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	w.Write([]byte("{}"))
}
