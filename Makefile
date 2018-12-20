.PHONY: build docker-image run-docker-image

build:
	cd go-maps-backend && go generate && go build -o simple-offline-map *.go

docker-image:
	docker build -f Dockerfile.builder --tag "daniel-kullmann/golang-builder" .
	docker run --name "temp-golang-builder" "daniel-kullmann/golang-builder"
	docker cp "temp-golang-builder:/var/app/go-maps-backend/simple-offline-map" .
	docker rm "temp-golang-builder"
	docker build --tag "daniel-kullmann/simple-offline-map" .
	rm -f simple-offline-map

run-docker-image: docker-image
	docker run -d --rm --publish "9000:9191" "daniel-kullmann/simple-offline-map"


run-docker-compose: docker-image
	docker-compose up -d
