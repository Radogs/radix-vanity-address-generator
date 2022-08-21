.DEFAULT_GOAL := help
.PHONY: start build build-docker prepare fresh

%::
	make
	@echo "> type one of the targets above"
	@echo

## update: Retrieves missing dependencies and cleans up the mod file
tidy:
	@echo "> Updating mod-file..."
	go mod tidy
	@echo "> Done"

## test: Runs unit tests
test:
	@echo "> Running Unit tests ..."
	go test ./... -count=1 -failfast
	@echo "> Done"

## format: Runs gofmt on all code
format:
	@echo "> Running gofmt ..."
	gofmt -s -w .
	@echo "> Done"

## build: Run Dockerfile locally
build:
	@echo "> Running gofmt ..."
	docker build -t radogs/radix-vanity-address-generator:latest .
	@echo "> Done"

## run: Run the vanity address generator in Docker
run:
	@echo "> Running gofmt ..."
	docker run -i radogs/radix-vanity-address-generator
	@echo "> Done"

## release: Push the vanity address generator to Dockerhub
release:
	@echo "> Running gofmt ..."
	./release.sh
	@echo "> Done"

makefile: help
help: Makefile
	@echo "> Choose a make command from the following:"
	@echo
	@sed -n 's/^##//p' $< | column -t -s ':' |  sed -e 's/^/ /'
	@echo
