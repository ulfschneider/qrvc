# go path
GOPATH := $(shell go env GOPATH)
GOBINPATH := $(GOPATH)/bin

# Tools used during build
TOOLS :=	github.com/google/go-licenses golang.org/x/vuln/cmd/govulncheck github.com/CycloneDX/cyclonedx-gomod/cmd/cyclonedx-gomod github.com/fzipp/gocyclo/cmd/gocyclo github.com/gordonklaus/ineffassign github.com/client9/misspell/cmd/misspell github.com/goreleaser/goreleaser/v2

# License handling
ALLOWED_LICENSES := MIT,BSD-2-Clause,BSD-3-Clause,Apache-2.0
IGNORE_LICENSES := qrvc,golang.org

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
MAIN_BRANCH := main

## help: show a list of available make commands
.PHONY: help
help:
	@echo "Usage:"
	@sed -n 's/^##//p' $(MAKEFILE_LIST) | column -t -s ':' | sed -e 's/^/ /'

## build: build the application for all targets. This requires a successful run of make release upfront.
.PHONY: build
build:
	@echo "Building qrvc"

	goreleaser release --snapshot --clean

	## verify-main: verify
.PHONY: verify-main
verify-main:
	@if [ "$$(git rev-parse --abbrev-ref HEAD)" != "$(MAIN_BRANCH)" ]; then \
		echo "Error: you are not on branch $(MAIN_BRANCH)!"; \
		exit 1; \
	fi
	@echo "On branch $(MAIN_BRANCH) ✅"

## release: tag the current state as a release in Git
.PHONY: release
release:
	@if [ "$$(git rev-parse --abbrev-ref HEAD)" != "$(MAIN_BRANCH)" ]; then \
		echo "Error: you are not on branch $(MAIN_BRANCH)!"; \
		exit 1; \
	fi

	@if [ -z "$(VERSION)" ]; then \
		echo "ERROR: You must pass VERSION=x.y.z to make a release."; exit 1; \
	fi

	@if [ -n "$$(git status --porcelain)" ]; then \
		echo "ERROR: Working tree is not clean. Commit or stash changes first."; \
		git status --porcelain; \
		exit 1; \
	fi

	@echo
	@ $(MAKE) version

	@echo
	@ $(MAKE) bom

	@echo
	@ $(MAKE) check

	@echo
	@ $(MAKE) test

	@echo
	@echo "Adding generated content to Git"
	git add $(BOM_GENERATED_FOLDER)
	git add $(VERSION_GENERATED_FOLDER)
	#git commit -m "Add BOM and version for release $(NORMALIZED_VERSION)"

	@echo
	@echo "Creating or updating tag $(NORMALIZED_VERSION)"
	#git tag $(NORMALIZED_VERSION)

	@echo
	@echo "Pushing release tag"
	#git push origin

	@echo
	@echo "👋 Release $(NORMALIZED_VERSION) complete."


## test: run all the automated tests
.PHONY: test
test:
	@echo "Automated tests"
	@go test ./...


## version: prepare version for embedding into the build
.PHONY: version
ifneq ($(strip $(VERSION)),)
version:
	@echo "Preparing version $(NORMALIZED_VERSION)"
	@printf "%s" "$(NORMALIZED_VERSION)" > $(VERSION_FILE)
else
version:
	@echo "No version information"
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
