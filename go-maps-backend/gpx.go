package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
)

const GpxBase = "../gpx_store_files/"

type Gpx = struct {
	Name         string
	Date         string
	Description  string
	Track_points [][]float64
}

func Map(input [][]float64, f func([]float64) string) []string {
	result := make([]string, len(input))
	for index, element := range input {
		result[index] = f(element)
	}
	return result
}

func CreateTrackPointXml(trackPoint []float64) string {
	return fmt.Sprintf(
		`      <trkpt lat="%f" lon="%f">
      </trkpt>`,
		trackPoint[0], trackPoint[1])
}

func CreateGpxContent(gpx *Gpx) string {
	trackPointsXml := strings.Join(Map(gpx.Track_points, CreateTrackPointXml), "\n")
	return fmt.Sprintf(`<?xml version="1.0" encoding="UTF-8"?>
<gpx creator="maps0" version="1.1" xmlns="http://www.topografix.com/GPX/1/1" xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance" xsi:schemaLocation="http://www.topografix.com/GPX/1/1 http://www.topografix.com/GPX/1/1/gpx.xsd http://www.garmin.com/xmlschemas/TrackPointExtension/v1 http://www.garmin.com/xmlschemas/TrackPointExtensionv1.xsd http://www.garmin.com/xmlschemas/GpxExtensions/v3 http://www.garmin.com/xmlschemas/GpxExtensionsv3.xsd" xmlns:gpxtpx="http://www.garmin.com/xmlschemas/TrackPointExtension/v1" xmlns:gpxx="http://www.garmin.com/xmlschemas/GpxExtensions/v3">
  <metadata>
    <name>%s</name>
    <desc>%s</desc>
    <time>%s</time>
  </metadata>
  <trk>
    <name>%s</name>
    <desc>%s</desc>
    <trkseg>
%s
    </trkseg>
  </trk>
</gpx>`, gpx.Name, gpx.Description, gpx.Date, gpx.Name, gpx.Description, trackPointsXml)
}

func SaveGpx(w http.ResponseWriter, r *http.Request) {
	if len(r.Header["X-Csrftoken"]) == 0 || r.Header["X-Csrftoken"][0] != token {
		w.WriteHeader(403)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.Header().Add("Access-Control-Allow-Origin", "*")
	buffer := make([]uint8, 1024*1024)
	count, err := r.Body.Read(buffer)
	if err != nil && err != io.EOF {
		log.Fatal(err)
	}
	var gpx Gpx
	err = json.Unmarshal(buffer[0:count], &gpx)
	gpx.Name = strings.Trim(gpx.Name, " \t\n")
	content := CreateGpxContent(&gpx)
	fileName := gpx.Date + "-" + gpx.Name + ".gpx"
	fh, err := os.Create(GpxBase + fileName)
	if err != nil {
		log.Fatal(err)
	}
	fh.Write([]byte(content))
	fh.Close()
	w.Write([]byte("{}"))
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
