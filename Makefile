PLATFORMS := windows linux darwin
ARCHS := amd64 arm64
BINARY := qrvc
DIST := dist

.PHONY: build-all

build-all:
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


	@ . ./.env; \
	if [ -n "$$AT_HOME" ]; then \
    echo "I´m at home, therefore copying $(DIST)/darwin/arm64/qrvc to ~/go/bin/"; \
    cp "$(DIST)/darwin/arm64/qrvc" ~/go/bin/; \
  fi

	@ echo "Ready 👋"
