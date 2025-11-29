PLATFORMS := windows linux darwin
ARCHS := amd64 arm64
BINARY := qrvc
DIST := dist
LICENSES := licenses

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

	go mod verify

	@echo

	@rm -rf $(LICENSES)

	go-licenses save ./... --save_path=licenses --ignore qrvc,golang.org

	@echo

	govulncheck ./...

	@echo

	cyclonedx-gomod mod -json=true -licenses=true -output=internal/sbom/sbom.json

	@rm -rf $(DIST)

	@echo

	@ for platform in $(PLATFORMS); do \
	    for arch in $(ARCHS); do \
			if [ "$$platform" = "windows" ]; then \
           target=$(DIST)/$$platform/$$arch/$(BINARY).exe; \
         else \
           target=$(DIST)/$$platform/$$arch/$(BINARY); \
         fi; \
	      echo "Building $$target"; \
			mkdir -p $(DIST)/$$platfom/$$arch; \
	      GOOS=$$platform GOARCH=$$arch go build -o $$target .; \
		  done; \
		done


	@#if the environment variable AT_HOME is defined in the .env file and it is not empty, execute the code
	@ . ./.env; \
	if [ -n "$$AT_HOME" ]; then \
	 echo;\
    echo "I´m at home, therefore copying $(DIST)/darwin/arm64/qrvc to ~/go/bin/"; \
    cp "$(DIST)/darwin/arm64/qrvc" ~/go/bin/; \
  fi

	@ echo "Ready 👋"

## version: show the current application version, change in .version file
.PHONY: version
version:
	@ . ./.version; \
	echo $$VERSION

## update: update all dependencies
.PHONY: update
update:
	go get -u ./...
	go mod tidy
	go mod verify
