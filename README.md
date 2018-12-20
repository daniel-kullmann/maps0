# Simple offline maps application

This is a small maps application with support for an offline cache for the map tiles.

## Features

* Shows a map via the Leaflet library, using tiles from OpenStreetMap
* Show GPX tracks
* Create new tracks
* Can optionally use a local tile cache so that offline usage is possible

Planned features (TODO list):
* Use React or another framework to better structure the code
* Management of GPX files
  * Don't show all gpx files; allow selection of gpx files, and then switching them on and off
  * Automatic and manual tagging of GPX files
* Adding your own tracks
  * TODO Determine and save metadata (date, description, ...)
* Enhance showing of GPX:
  * Customizing markers and colors of shown GPX tracks
* Add customizing of tiles cache (maximum age of tiles, maximum size of cache)
* Add refreshing of cache
  * For that, we need to store cache metadata, like the etags and max-age
    headers that come from the OSM tile server.


## Configuration

TODO

## Installation

Setup (you only have to do this once):
* Install golang
* Make sure that the `GOBIN` environment variable is set to a directory that is available on `PATH`, e.g. `export GOBIN=$HOME/bin`
* Install go-bindata: run `go get -u github.com/jteeuwen/go-bindata/...`
* Make sure that the `go-bindata` command is available by running it

To build the executable, you need to execute these commands:
```
cd go-maps-backend
go generate
go build -o simple-offline-map
```
Instead of `go generate`, you can also run `go-bindata -prefix ../frontend/ ../frontend/...`.

You can also just do a `make` to build the binary.

The resulting binary is located in `./go-maps-backend/simple-offline-map`.


## Build a docker image

I also added the necessary files to build a docker image. Since this docker
image is based on Alpine, and Alpine uses a different loader .so than regular
linux systems, I use a separate docker image just to build the binary. This
results in a really small resulting docker image (~19MB), while the docker image
for building is ~600MB.

Simply call `make docker-image` for creating the docker-image, or look into the
`Makefile` under the `docker-image` heading to see how the image is built.

