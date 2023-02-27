
build:
	@go build -o bin/GYM_API

run: build
	@./bin/GYM_API

test:
	@go test -v ./...