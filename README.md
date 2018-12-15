# Simple maps application

This is supposed to be a small maps application.

## Features

* Shows a map via the Leaflet library, using tiles from OpenStreetMap
* Possibility to show GPX tracks
* Uses a tile cache (via a django app) so that offline usage is possible

Planned features (TODO list):
* Use React or another framework to better structure the code
* Management of GPX files
  * Don't show all gpx files; allow selection of gpx files, and then switching them on and off
  * Automatic and manual tagging of GPX files
* Adding your own tracks
  * DONE Add a mode to the map (viewing, add new track)
  * DONE Capture clicks on map, and add appropriate markers and lines
  * DONE Save the resulting GPX in the end
  * TODO Determine and save metadata (date, description, ...)
* Enhance showing of GPX:
  * Customizing markers and colors of shown GPX tracks
  * DONE Show lines between track points
  * DONE Show markers only at beginning and end of track
* Stand-alone mode: create a docker image that contains the whole application,
  including the tiles cache
* Add customizing of tiles cache (maximum age of tiles, maximum size of cache)
* Add refreshing of cache
  * For that, we need to store cache metadata, like the etags and max-age
    headers that come from the OSM tile server.
* DONE Add nginx config to serve the frontend (static html) and backend (via uwsgi).
* DONE Add node backend of tile cache for better concurrency
* Add Go implementation of whole backend

## Installation

* Install django
* Run the following commands from the `maps_backend` folder:
```
./manage.py migrate
./manage.py runserver
```
* Visit the file `frontend/index.html` with your web browser

After the first run, it will be enough to restart the django application via
`./manage.py runserver` and o point the web browser at the same file again.
