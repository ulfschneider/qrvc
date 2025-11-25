PLATFORMS := windows linux darwin
ARCHS := amd64 arm64
BINARY := qrvc
DIST := dist

build-all:
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

	# Load variables from .env to see if this is running on my local dev machine
  . ./.env; \
  if [ -n "$$AT_HOME" ]; then \
    echo "I´m at home, therefore copying $(DIST)/darwin/arm64/qrvc to ~/go/bin/" \
    cp "$(DIST)/darwin/arm64/qrvc" ~/go/bin/ \
  fi \

	echo "Ready 👋"
