# NAME daniel-kullmann/simple-offline-map
FROM alpine
MAINTAINER daniel.kullmann@gmx.de

ENV DEBIAN_FRONTEND noninteractive
ENV LANG en_GB.UTF-8
ENV LC_ALL en_GB.UTF-8

RUN apk add ca-certificates
RUN mkdir -p /var/www/app/tile /var/www/app/gpx

WORKDIR /var/www/app
ADD simple-offline-map /var/www/app/simple-offline-map
ADD deploy/config.ini /var/www/app/config.ini
#ADD frontend /var/www/app/html

EXPOSE 9191

CMD [ "/var/www/app/simple-offline-map", "-config", "/var/www/app/config.ini" ]

