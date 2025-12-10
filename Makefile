# Tools used during build
TOOLS :=	github.com/google/go-licenses golang.org/x/vuln/cmd/govulncheck github.com/CycloneDX/cyclonedx-gomod/cmd/cyclonedx-gomod

# License handling
ALLOWED_LICENSES := MIT,BSD-2-Clause,BSD-3-Clause,Apache-2.0
IGNORE_LICENSES := qrvc,golang.org

# build targets
PLATFORMS := windows linux darwin
ARCHS := amd64 arm64

# Names for building
BINARY_NAME := qrvc
DIST_FOLDER := dist
GENERATED_FOLDER := internal/appmeta/generated
LICENSES_FOLDER:= $(GENERATED_FOLDER)/licenses
SBOM_FILE  := $(GENERATED_FOLDER)/sbom.json
VERSION_FILE := $(GENERATED_FOLDER)/version.txt

# Strip leading v, then prepend exactly one v
NORMALIZED_VERSION := v$(patsubst v%,%,$(VERSION))

# Name the temporyry release branch
RELEASE_BRANCH := release-tmp-$(NORMALIZED_VERSION)

## help: show a list of available make commands
.PHONY: help
help:
	@echo "Usage:"
	@sed -n 's/^##//p' $(MAKEFILE_LIST) | column -t -s ':' | sed -e 's/^/ /'

## build: build the application for all targets. To build for a release, do not forget the set the version.
.PHONY: build
build:
	@echo "Building qrvc"

	@echo
	@ $(MAKE) update-tools

	@echo
	@ $(MAKE) update

	@echo
	@ $(MAKE) check

	@echo
	@ $(MAKE) test

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
	@if [ -z "$(VERSION)" ]; then \
		echo "ERROR: You must pass VERSION=x.y.z to make a release."; exit 1; \
	fi

	@ $(MAKE) check

	@echo
	@ $(MAKE) test

	@if [ -n "$$(git status --porcelain)" ]; then \
		echo "ERROR: Working tree is not clean. Commit or stash changes first."; \
		git status --porcelain; \
		exit 1; \
	fi

	@echo
	@echo "Creating temporary release branch $(RELEASE_BRANCH)"
	git checkout -b $(RELEASE_BRANCH)

	@printf "%s" "$(NORMALIZED_VERSION)" > $(VERSION_FILE)

	@echo
	@ $(MAKE) sbom

	@echo
	@ $(MAKE) test_appmeta

	@echo
	@echo "Adding generated content to release branch"
	git add -f $(GENERATED_FOLDER)
	git commit -m "Add SBOM and version for release $(NORMALIZED_VERSION)"

	@echo
	@echo "Creating or updating tag $(NORMALIZED_VERSION)"
	git tag -f $(NORMALIZED_VERSION)

	@echo
	@echo "Pushing release tag"
	git push -f origin $(NORMALIZED_VERSION)

	@echo
	@echo "Cleaning up temporary branch"
	git checkout -
	git branch -D $(RELEASE_BRANCH)

	@echo
	@echo "👋 Release $(NORMALIZED_VERSION) complete."

## test: run all the automated tests, excluding appmeta tests
.PHONY: test
test:
	@echo "Automated tests"
	go test $(shell go list ./... | grep -v appmeta)

## test_appmeta: run the automated tests for appmeta
.PHONY: test_appmeta
test_appmeta:
	@echo "Test appmeta"
	go test ./internal/appmeta

## sbom: check and prepare licenses and sbom for embedding them into the build
.PHONY: sbom
sbom:
	@echo "Preparing licenses"
	rm -rf $(LICENSES_FOLDER);
	go-licenses check ./... --allowed_licenses=$(ALLOWED_LICENSES) --ignore=$(IGNORE_LICENSES)
	go-licenses save ./... --save_path=$(LICENSES_FOLDER) --ignore=$(IGNORE_LICENSES)
	@echo
	@echo "Preparing SBOM"
	@cyclonedx-gomod app -json=true -licenses=true -output=$(SBOM_FILE)


## update: update all dependencies
.PHONY: update
update:
	@echo "Updating dependencies"
	go get -u ./...

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
	go fmt ./...
	go vet ./...
	govulncheck ./...
