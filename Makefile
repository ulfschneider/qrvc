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

# Strip leading v, then prepend exactly one v
NORMALIZED_VERSION := v$(patsubst v%,%,$(VERSION))

MAIN_BRANCH := main
BREW_REPO := ulfschneider/homebrew-tap

## help: show a list of available make commands
.PHONY: help
help:
	@echo "Usage:"
	@sed -n 's/^##//p' $(MAKEFILE_LIST) | column -t -s ':' | sed -e 's/^/ /'

.PHONY: build
build:
	@echo "Building qrvc $(NORMALIZED_VERSION)"

	@if [ ! -f ./.env ]; then \
			echo "ERROR: .env not found"; exit 1; \
	fi
	@set -a; \
	. ./.env; \
	set +a; \

	@echo
	HOMEBREW_REPO=$(BREW_REPO)	goreleaser release --clean


.PHONY: verify-main
verify-main:
	@if [ "$$(git rev-parse --abbrev-ref HEAD)" != "$(MAIN_BRANCH)" ]; then \
		echo "Error: you are not on branch $(MAIN_BRANCH)!"; \
		exit 1; \
	fi

## release: tag the current state as a release in Git and distribute the binaries
.PHONY: release
release:
	$(MAKE) verify-main

	@if [ -z "$(VERSION)" ]; then \
		echo "ERROR: You must pass VERSION=x.y.z to make a release."; exit 1; \
	fi

	@echo
	@ $(MAKE) update

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
	@ $(MAKE) test

	@echo
	@echo "Adding generated content to Git"
	git add .
	git commit -m "Add version $(NORMALIZED_VERSION) and BOM for release"

	@echo
	@echo "Creating tag $(NORMALIZED_VERSION)"
	git tag $(NORMALIZED_VERSION)

	@echo
	@echo "Pushing release"
	git push origin

	@echo
	$(MAKE) build


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
