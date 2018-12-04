#!/bin/sh

uwsgi \
  --chdir=$(dirname $0)/maps_backend \
  --module=maps_backend.wsgi:application \
  --env DJANGO_SETTINGS_MODULE=maps_backend.settings \
  --master \
  --pidfile=$(dirname $0)/uwsgi.pid \
  --socket=127.0.0.1:9000 \
  --processes=3 \
  --max-requests=5000 \
  --vacuum \
  --daemonize=$(dirname $0)/uwsgi.log

