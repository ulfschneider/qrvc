PLATFORMS := windows linux darwin
ARCHS := amd64 arm64
BINARY := qrvc
DIST := dist
LICENSES := sbom/generated/licenses
SBOM  := sbom/generated/sbom.json

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

	@ $(MAKE) update

	@rm -rf $(DIST)

	@ for platform in $(PLATFORMS); do \
	    for arch in $(ARCHS); do \
			if [ "$$platform" = "windows" ]; then \
           target=$(DIST)/$$platform/$$arch/$(BINARY).exe; \
         else \
           target=$(DIST)/$$platform/$$arch/$(BINARY); \
         fi; \
			mkdir -p $(DIST)/$$platform/$$arch; \
			echo; \
			echo "Creating SBOM for $$target"; \
			GOOS=$$platform GOARCH=$$arch cyclonedx-gomod app -json=true -licenses=true -output=$(SBOM); \
			echo "Building $$target"; \
			GOOS=$$platform GOARCH=$$arch go build -o $$target . ; \
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

## update: update all dependencies and perform a check
.PHONY: update
update:
	go get -u ./...
	@ $(MAKE) check

## check: tidy up the go.mod file and do a vulnerability check
.PHONY: check
check:
	go mod tidy
	go mod verify
	go-licenses check ./... --allowed_licenses=MIT,BSD-2-Clause,BSD-3-Clause,Apache-2.0 --ignore qrvc,golang.org
	rm -rf $(LICENSES);
	go-licenses save ./... --save_path=$(LICENSES) --ignore qrvc,golang.org
	govulncheck ./...
