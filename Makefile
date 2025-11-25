PLATFORMS := windows linux darwin
ARCHS := amd64 arm64
BINARY := qrvc
DIST := dist


## help: show a list of available make commands
.PHONY: help
help:
	@echo "Usage:"
	@sed -n 's/^##//p' $(MAKEFILE_LIST) | column -t -s ':' | sed -e 's/^/ /'

## build: build the application for all targets
.PHONY: build
build:
	@ . ./.version; \
	echo "Building qrvc $$VERSION"

	go mod tidy

	govulncheck ./...

	@rm -rf $(DIST)

	@ for platform in $(PLATFORMS); do \
	    for arch in $(ARCHS); do \
	      echo "Building for $$platform $$arch"; \
		     if [ "$$platform" = "windows" ]; then \
	        GOOS=$$platform GOARCH=$$arch go build -o $(DIST)/$$platform/$$arch/$(BINARY).exe .; \
				 else \
	        GOOS=$$platform GOARCH=$$arch go build -o $(DIST)/$$platform/$$arch/$(BINARY) .; \
				 fi; \
		  done; \
		done


	@#if the environment variable AT_HOME is defined in the .env file and it is not empty, execute the code
	@ . ./.env; \
	if [ -n "$$AT_HOME" ]; then \
    echo "I´m at home, therefore copying $(DIST)/darwin/arm64/qrvc to ~/go/bin/"; \
    cp "$(DIST)/darwin/arm64/qrvc" ~/go/bin/; \
  fi

	@ echo "Ready 👋"

## version: show the current application version
.PHONY: version
version:
	@ . ./.version; \
	echo $$VERSION
