TOOLS :=	github.com/google/go-licenses golang.org/x/vuln/cmd/govulncheck github.com/CycloneDX/cyclonedx-gomod/cmd/cyclonedx-gomod

PLATFORMS := windows linux darwin
ARCHS := amd64 arm64

BINARY_NAME := qrvc
DIST_FOLDER := dist
GENERATED_FOLDER := internal/version/generated
LICENSES_FOLDER:= $(GENERATED_FOLDER)/licenses
SBOM_FILE  := $(GENERATED_FOLDER)/sbom.json
VERSION_FILE := $(GENERATED_FOLDER)/version.txt

RELEASE_BRANCH := release-tmp-$(VERSION)

## help: show a list of available make commands
.PHONY: help
help:
	@echo "Usage:"
	@sed -n 's/^##//p' $(MAKEFILE_LIST) | column -t -s ':' | sed -e 's/^/ /'

## build: build the application for all targets. To build for a release, do not forget the set the version.
.PHONY: build
build:
	@echo "Building qrvc"

	@ $(MAKE) update-tools

	@ $(MAKE) update

	@rm -rf $(DIST)

	@ for platform in $(PLATFORMS); do \
	    for arch in $(ARCHS); do \
			if [ "$$platform" = "windows" ]; then \
           target=$(DIST_FOLDER)/$$platform/$$arch/$(BINARY_NAME).exe; \
         else \
           target=$(DIST_FOLDER)/$$platform/$$arch/$(BINARY_NAME); \
         fi; \
			mkdir -p $(DIST_FOLDER)/$$platform/$$arch; \
			echo; \
			echo "Building $$target"; \
			GOOS=$$platform GOARCH=$$arch go build -o $$target . ; \
		 done; \
	done

	@#if the environment variable AT_HOME is defined in the .env file and it is not empty, execute the code
	@ . ./.env; \
	if [ -n "$$AT_HOME" ]; then \
	   echo;\
      echo "I´m at home, therefore copying $(DIST_FOLDER)/darwin/arm64/qrvc to ~/go/bin/"; \
      cp "$(DIST_FOLDER)/darwin/arm64/qrvc" ~/go/bin/; \
   fi

	@ echo "👋 Binaries are built"

## release: tag the current state as a release in Git
.PHONY: release
release:
	@if [ -n "$$(git status --porcelain)" ]; then \
		echo "ERROR: Working tree is not clean. Commit or stash changes first."; \
		git status --porcelain; \
		exit 1; \
	fi

	@if [ -z "$(VERSION)" ]; then \
		echo "ERROR: You must pass VERSION=v<MAJOR>.<MINOR>.<FIX> to make a release. Do not forget the v prefix for your release!"; exit 1; \
	fi

	@echo "Creating temporary release branch $(RELEASE_BRANCH)"
	git checkout -b $(RELEASE_BRANCH)

	@printf "%s" "$(VERSION)" > $(VERSION_FILE)

	@ $(MAKE) sbom

	@echo "Adding generated content to release branch"
	git add -f $(GENERATED_FOLDER)
	git commit -m "Add SBOM and version for release $(VERSION)"

	@echo "Creating or updating tag $(VERSION)"
	git tag -f $(VERSION)

	@echo "Pushing release tag"
	git push -f origin $(VERSION)

	@echo "Cleaning up temporary branch"
	git checkout -
	git branch -D $(RELEASE_BRANCH)

	@echo "👋 Release $(VERSION) complete."


## sbom: check and prepare licenses and sbom for embedding them into the build
.PHONY: sbom
sbom:
	@echo "Preparing licenses"
	rm -rf $(LICENSES_FOLDER);
	go-licenses check ./... --allowed_licenses=MIT,BSD-2-Clause,BSD-3-Clause,Apache-2.0 --ignore qrvc,golang.org
	go-licenses save ./... --save_path=$(LICENSES_FOLDER) --ignore qrvc,golang.org
	@echo "Preparing SBOM"
	@cyclonedx-gomod app -json=true -licenses=true -output=$(SBOM_FILE)


## update: update all dependencies perform a check and prepare the sbom
.PHONY: update
update:
	@echo "Updating dependencies"
	go get -u ./...
	@ $(MAKE) check
	@ $(MAKE) sbom

## update-tools: update the tools that are required for building
.PHONY: update-tools
update-tools:
	@echo "Updating build tools"
	@for t in $(TOOLS); do \
		echo "Updating $$t..."; \
		go install "$$t@latest"; \
	done

## check: tidy up the go.mod file and check for vulnerabilities
.PHONY: check
check:
	@echo "Tidying up the mod file and doing a vulnerability check"
	go mod tidy
	go mod verify
	govulncheck ./...

## check-verbose: this is like check but with verbose logging
check-verbose:
	@echo "Tidying up the mod file and doing a vulnerability check with verbose logging"
	go mod tidy
	go mod verify
	govulncheck -show verbose ./...
