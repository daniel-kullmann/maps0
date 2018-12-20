# NAME daniel-kullmann/golang-builder
FROM golang:alpine
MAINTAINER daniel.kullmann@gmx.de

ENV DEBIAN_FRONTEND noninteractive
ENV LANG en_GB.UTF-8
ENV LC_ALL en_GB.UTF-8
ENV GOBIN /usr/bin

EXPOSE 9191

RUN apk add git build-base
RUN go get -u github.com/jteeuwen/go-bindata/... && \
    go get -u github.com/gorilla/mux && \
    go get -u github.com/mattn/go-sqlite3

RUN mkdir /var/app/
ADD frontend /var/app/frontend/
ADD go-maps-backend /var/app/go-maps-backend/
WORKDIR /var/app/

RUN cd /var/app/go-maps-backend/ && go generate && go build -o simple-offline-map *.go

CMD [ /bin/true ]

