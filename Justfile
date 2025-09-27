help:
	
build: 
	cd cmd/gowatchit && go build -o ../../build/gowatchit
test:
	./test_web.sh
docker-build:
	docker buildx build --platform linux/arm64 --load --tag gowatchit-local . -f services/gowatchit/Dockerfile
docker-push:
	docker buildx build --push --platform linux/amd64 --tag ghcr.io/iloveicedgreentea/gowatchit:test . 
docker-run: docker-build
	LOG_FILE=true LOG_LEVEL=debug docker-compose -f services/gowatchit/docker-compose.yml up
run: 
	LOG_ENV=local LOG_FILE=false LOG_LEVEL=debug go run ./services/gowatchit
run-ui:
	cd web && bun run dev
live-test:
	./test_web.sh