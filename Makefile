# go path
GOPATH := $(shell go env GOPATH)
GOBINPATH := $(GOPATH)/bin

# Tools used during build
TOOLS :=	github.com/google/go-licenses golang.org/x/vuln/cmd/govulncheck github.com/CycloneDX/cyclonedx-gomod/cmd/cyclonedx-gomod github.com/fzipp/gocyclo/cmd/gocyclo github.com/gordonklaus/ineffassign github.com/client9/misspell/cmd/misspell

# License handling
ALLOWED_LICENSES := MIT,BSD-2-Clause,BSD-3-Clause,Apache-2.0
IGNORE_LICENSES := qrvc,golang.org

# build targets
PLATFORMS := windows linux darwin
ARCHS := amd64 arm64

# Names for building
BINARY_NAME := qrvc
DIST_FOLDER := dist
BOM_GENERATED_FOLDER := internal/adapters/bom/embedded/generated
LICENSES_FOLDER:= $(BOM_GENERATED_FOLDER)/licenses
BOM_FILE  := $(BOM_GENERATED_FOLDER)/bom.json
VERSION_GENERATED_FOLDER := internal/adapters/version/embedded/generated
VERSION_FILE := $(VERSION_GENERATED_FOLDER)/version.txt
GIT_HASH_FILE := $(VERSION_GENERATED_FOLDER)/commit.txt
GIT_TIMESTAMP_FILE := $(VERSION_GENERATED_FOLDER)/time.txt

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
	@ $(MAKE) update

	@echo
	@ $(MAKE) version

	@echo
	@ $(MAKE) bom

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
      echo "I´m at home, therefore copying $(DIST_FOLDER)/darwin/arm64/qrvc to $(GOBINPATH)"; \
      cp "$(DIST_FOLDER)/darwin/arm64/qrvc" $(GOBINPATH); \
   fi

	@ echo "👋 Binaries are built"

## release: tag the current state as a release in Git
.PHONY: release
release:
	@if [ -z "$(VERSION)" ]; then \
		echo "ERROR: You must pass VERSION=x.y.z to make a release."; exit 1; \
	fi

	@ $(MAKE) check

	@if [ -n "$$(git status --porcelain)" ]; then \
		echo "ERROR: Working tree is not clean. Commit or stash changes first."; \
		git status --porcelain; \
		exit 1; \
	fi

	@echo
	@echo "Creating temporary release branch $(RELEASE_BRANCH)"
	git checkout -b $(RELEASE_BRANCH)

	@echo
	@ $(MAKE) version

	@echo
	@ $(MAKE) bom

	@echo
	@ $(MAKE) test

	@echo
	@echo "Adding generated content to release branch"
	git add --force $(BOM_GENERATED_FOLDER)
	git add --force $(VERSION_GENERATED_FOLDER)
	git commit -m "Add BOM and version for release $(NORMALIZED_VERSION)"

	@echo
	@echo "Creating or updating tag $(NORMALIZED_VERSION)"
	git tag --force $(NORMALIZED_VERSION)

	@echo
	@echo "Pushing release tag"
	git push --force origin $(NORMALIZED_VERSION)

	@echo
	@echo "Cleaning up temporary branch"
	git checkout -
	git branch --delete --force $(RELEASE_BRANCH)

	@echo
	@echo "👋 Release $(NORMALIZED_VERSION) complete."


## test: run all the automated tests
.PHONY: test
test:
	@echo "Automated tests"
	@go test ./...

## gitinfo: prepare git last commit hash and commit time for embedding into the build
.PHONY: gitinfo
gitinfo:
	@if git rev-parse --is-inside-work-tree >/dev/null 2>&1; then \
		printf "%s" "$$(git rev-parse --short HEAD)" > $(GIT_HASH_FILE); \
		printf "%s" "$$(git show -s --format=%cI HEAD)" > $(GIT_TIMESTAMP_FILE); \
	else \
		: > $(GIT_HASH_FILE); \
		: > $(GIT_TIMESTAMP_FILE); \
	fi

## version: prepare version for embedding into the build
.PHONY: version
ifneq ($(strip $(VERSION)),)
version:
	@echo "Preparing version $(NORMALIZED_VERSION)"
	@printf "%s" "$(NORMALIZED_VERSION)" > $(VERSION_FILE)
	@ $(MAKE) gitinfo
else
version:
	@echo "No version information"
	@ $(MAKE) gitinfo
endif

## bom: check and prepare licenses and bom for embedding them into the build
.PHONY: bom
bom:
	@echo "Preparing licenses"
	rm -rf $(LICENSES_FOLDER);
	go-licenses check ./... --allowed_licenses=$(ALLOWED_LICENSES) --ignore=$(IGNORE_LICENSES)
	go-licenses save ./... --save_path=$(LICENSES_FOLDER) --ignore=$(IGNORE_LICENSES)
	@echo
	@echo "Preparing BOM"
	@cyclonedx-gomod app -json=true -licenses=true -output=$(BOM_FILE)


## update: update dependencies and then do a check
.PHONY: update
update:
	@echo "Updating dependencies"
	go get -u ./...
	@ $(MAKE) check

## update-tools: update the tools that are required for building
.PHONY: update-tools
update-tools:
	@echo "Updating build tools"
	@for t in $(TOOLS); do \
		echo "Updating $$t..."; \
		go install "$$t@latest"; \
	done

## check: tidy up the go.mod file and check dependencies and code
.PHONY: check
check:
	@ $(MAKE) update-tools
	@echo
	@echo "Tidying up the mod file, format code, check code quality, and check dependencies for vulnerabilities"
	go mod tidy
	go mod verify
	gofmt -s -w  .
	go vet ./...
	gocyclo -over 15 .
	ineffassign ./...
	misspell -error ./...
	govulncheck ./...
