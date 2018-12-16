package main

import (
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/tile/{s}/{z}/{x}/{y}.png", GetTile).Methods("GET")
	log.Fatal(http.ListenAndServe(":9191", router))
}



}
