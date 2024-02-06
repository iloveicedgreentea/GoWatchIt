.PHONY: build test run docker-run docker-build
SHELL := /bin/bash
build: 
	cd cmd && go build -o ../build/server
test:
	./test.sh
docker-build:
	docker buildx build --load --tag gowatchit-local .
docker-push:
	docker buildx build --push --platform linux/amd64 --tag ghcr.io/iloveicedgreentea/gowatchit:test . 
docker-run:
	LOG_FILE=false LOG_LEVEL=debug docker-compose up
run: build
	LOG_FILE=false LOG_LEVEL=debug ./build/server