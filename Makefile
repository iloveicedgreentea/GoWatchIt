.PHONY: build test run

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

run: build
	@./build/server