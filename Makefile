.PHONY: build docker-image run-docker-image

build:
	cd go-maps-backend && go generate && go build -o simple-offline-map *.go

docker-image:
	docker build --tag "daniel-kullmann/simple-offline-map" .

run-docker-image: docker-image
	docker run -d --rm --publish "9000:9191" "daniel-kullmann/simple-offline-map"

run-docker-compose: docker-image
	docker-compose up -d
