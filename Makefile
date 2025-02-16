.DEFAULT_GOAL := help

.PHONY: run
run:  ## Runs the project using 'go run'.
	@go run ./mhelp/mhelp.go

.PHONY: test
test:  ## Tests the project using 'go test'.
	@go test

.PHONY: build
build: build_windows build_darwin build_linux  ## Builds the project for all platforms using 'go build'.

.PHONY: clean
clean:  ## Cleans the builds folder using 'rm'.
	@rm ./builds/*

.PHONY: build_windows
build_windows:  ## Builds the project for windows platforms using 'go build'.
	@GOOS=windows GOARCH=amd64 go build -o ./builds/help_windows_amd64 ./mhelp/mhelp.go

.PHONY: build_darwin
build_darwin:  ## Builds the project for darwin platforms using 'go build'.
	@GOOS=darwin  GOARCH=arm64 go build -o ./builds/help_darwin_arm64  ./mhelp/mhelp.go

.PHONY: build_linux
build_linux:  ## Builds the project for linux platforms using 'go build'.
	@GOOS=linux   GOARCH=amd64 go build -o ./builds/help_linux_amd64   ./mhelp/mhelp.go

.PHONY: help
help:  ## Show help and exit.
	@./builds/help_darwin_arm64 -filepath $(MAKEFILE_LIST)
