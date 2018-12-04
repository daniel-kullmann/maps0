# NAME danielkullmann/maps0-sandalone
FROM python:alpine
MAINTAINER daniel.kullmann@gmx.de

ENV DEBIAN_FRONTEND noninteractive
ENV LANG en_GB.UTF-8
ENV LC_ALL en_GB.UTF-8

RUN apk add --no-cache ca-certificates curl vim wget nginx

RUN pip3 install pip --upgrade && \
    pip3 install wheel && \
    pip3 install -U setuptools

ADD requirements.txt /var/www/requirements.txt
RUN pip3 install -r /var/www/requirements.txt
RUN pip3 install uwsgi

RUN echo 'alias python=python3' >> ~/.bashrc && echo 'alias pip=pip3' >> ~/.bashrc

WORKDIR /var/www/app
ADD . /var/www/app

RUN ./manage.py migrate
RUN mkdir -p /srv/logs
RUN sed -i -e 's%http://localhost:8000/api%/api%g' /var/www/app/frontend/code.js

EXPOSE 80

CMD [ "/var/www/app/start-application.sh" ]

