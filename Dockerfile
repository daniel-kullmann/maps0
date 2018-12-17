# NAME danielkullmann/maps0-standalone
FROM alpine
MAINTAINER daniel.kullmann@gmx.de

ENV DEBIAN_FRONTEND noninteractive
ENV LANG en_GB.UTF-8
ENV LC_ALL en_GB.UTF-8

RUN mkdir -p /var/www/app/{tiles,gpx}

WORKDIR /var/www/app
ADD go-maps-backend/simple-offline-map /var/www/app/simple-offline-map
ADD config.ini /var/www/app/simple-offline-map
ADD frontend /var/www/app/html


EXPOSE 9191

CMD [ "/var/www/app/simple-offline-map -f /var/www/app/config.ini" ]

