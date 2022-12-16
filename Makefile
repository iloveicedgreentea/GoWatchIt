.PHONY: build test run

build: 
	go build -o ./build/server

test:
	@go vet
	@unset LOG_LEVEL && cd handlers && go test 
	@unset LOG_LEVEL && cd homeassistant && go test 
	@unset LOG_LEVEL && cd plex && go test 
	@unset LOG_LEVEL && cd mqtt && go test 
	@unset LOG_LEVEL && cd ezbeq && go test 

run: build
	@./build/server