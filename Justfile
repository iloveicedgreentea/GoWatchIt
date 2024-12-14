
build: 
	cd cmd/gowatchit && go build -o ../../build/gowatchit
test:
	./test.sh
docker-build:
	docker buildx build --platform linux/amd64 --load --tag gowatchit-local . -f ./Dockerfile.dev
docker-push:
	docker buildx build --push --platform linux/amd64 --tag ghcr.io/iloveicedgreentea/gowatchit:test . 
docker-run:
	LOG_FILE=false LOG_LEVEL=debug docker-compose -f docker-compose-test.yml up
run: 
	LOG_FILE=false LOG_LEVEL=debug go run ./cmd/gowatchit/