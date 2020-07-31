# NAME daniel-kullmann/simple-offline-map

########################################################################
# Build Stage
FROM golang:alpine as builder
MAINTAINER daniel.kullmann@gmx.de

ENV DEBIAN_FRONTEND noninteractive
ENV LANG en_GB.UTF-8
ENV LC_ALL en_GB.UTF-8
ENV GOBIN /usr/bin

RUN apk add git build-base
RUN go get -u github.com/jteeuwen/go-bindata/... && \
    go get -u github.com/gorilla/mux && \
    go get -u github.com/mattn/go-sqlite3

RUN mkdir /var/app/
ADD frontend /var/app/frontend/
ADD go-maps-backend /var/app/go-maps-backend/
WORKDIR /var/app/

RUN cd /var/app/go-maps-backend/ && go generate && go build -o simple-offline-map *.go


########################################################################
# App Image Stage
FROM alpine as app-image
MAINTAINER daniel.kullmann@gmx.de

ENV DEBIAN_FRONTEND noninteractive
ENV LANG en_GB.UTF-8
ENV LC_ALL en_GB.UTF-8

RUN apk add ca-certificates
RUN mkdir -p /var/www/app/tile /var/www/app/gpx

WORKDIR /var/www/app
COPY --from=builder /var/app/go-maps-backend/simple-offline-map /var/www/app/simple-offline-map
ADD deploy/config.ini /var/www/app/config.ini

EXPOSE 9191

CMD [ "/var/www/app/simple-offline-map", "-config", "/var/www/app/config.ini" ]

