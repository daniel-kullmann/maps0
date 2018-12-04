#!/bin/bash

set -e -u

./uwsgi.sh

nginx -t -c $(pwd)/nginx.conf
nginx -c $(pwd)/nginx.conf &

