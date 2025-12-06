PLATFORMS := windows linux darwin
ARCHS := amd64 arm64
BINARY := qrvc
DIST := dist
GENERATED := internal/version/generated/
LICENSES := $(GENERATED)licenses
SBOM  := $(GENERATED)sbom.json
VER := $(GENERATED)version.txt
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

## sbom: check and prepare licenses and sbom for embedding them into the build
.PHONY: sbom
sbom:
	@echo "Preparing licenses"
	rm -rf $(LICENSES);
	go-licenses check ./... --allowed_licenses=MIT,BSD-2-Clause,BSD-3-Clause,Apache-2.0 --ignore qrvc,golang.org
	go-licenses save ./... --save_path=$(LICENSES) --ignore qrvc,golang.org
	@echo "Preparing SBOM"
	@cyclonedx-gomod app -json=true -licenses=true -output=$(SBOM)


## update: update all dependencies perform a check and prepare the sbom
.PHONY: update
update:
	go get -u ./...
	@ $(MAKE) check
	@ $(MAKE) sbom

## check: tidy up the go.mod file and check for vulnerabilities
.PHONY: check
check:
	go mod tidy
	go mod verify
	govulncheck ./...

## release: tag the current state as a release in Git
.PHONY: release
release:
	@if [ -n "$$(git status --porcelain)" ]; then \
		echo "ERROR: Working tree is not clean. Commit or stash changes first."; \
		git status --porcelain; \
		exit 1; \
	fi

	@if [ -z "$(VERSION)" ]; then \
		echo "ERROR: You must pass VERSION=vx.y.z to make a release"; exit 1; \
	fi

	@echo "Creating temporary release branch $(RELEASE_BRANCH)"
	git checkout -b $(RELEASE_BRANCH)

	@printf "%s" "$(VERSION)" > $(VER)

	@echo "Adding generated content to release branch"
	git add -f $(GENERATED)
	git commit -m "Add SBOM and version for release $(VERSION)"

	@echo "Creating or updating tag $(VERSION)"
	git tag -f $(VERSION)

	@echo "Pushing release tag"
	git push -f origin $(VERSION)

	@echo "Cleaning up temporary branch"
	git checkout -
	git branch -D $(RELEASE_BRANCH)

	@echo "Release $(VERSION) complete."
