
build: 
	cd cmd/gowatchit && go build -o ../../build/gowatchit
test:
	./test.sh
docker-build:
	docker buildx build --platform linux/arm64 --load --tag gowatchit-local . -f ./Dockerfile
docker-push:
	docker buildx build --push --platform linux/amd64 --tag ghcr.io/iloveicedgreentea/gowatchit:test . 
docker-run:
	LOG_FILE=true LOG_LEVEL=debug docker-compose -f docker-compose.yml up
run: 
	LOG_ENV=local LOG_FILE=true LOG_LEVEL=debug go run ./cmd/gowatchit/
run-ui:
	cd web && bun run dev