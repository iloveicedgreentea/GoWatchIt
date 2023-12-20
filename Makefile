.PHONY: build test run docker-run docker-build

build: 
	cd cmd && go build -o ../build/server

test:
	# @go vet
	@unset LOG_LEVEL && cd internal/config && go test -v
	@unset LOG_LEVEL && cd internal/handlers && go test -v
	@unset LOG_LEVEL && cd internal/homeassistant && go test -v
	@unset LOG_LEVEL && cd internal/denon && go test -v
	@unset LOG_LEVEL && cd internal/plex && go test -v
	@unset LOG_LEVEL && cd internal/mqtt && go test -v
	@unset LOG_LEVEL && cd internal/ezbeq && go test -v
docker-build:
	docker buildx build --load --tag plex-webhook-automation-local .
docker-push:
	docker buildx build --push --platform linux/amd64 --tag ghcr.io/iloveicedgreentea/plex-webhook-automation:test . 
docker-run: docker-build
	docker run -p 9999:9999 -e LOG_LEVEL=debug plex-webhook-automation-local
run: build
	LOG_LEVEL=debug ./build/server