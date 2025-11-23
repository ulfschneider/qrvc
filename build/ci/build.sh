#! /bin/bash
if [ -f ".version" ]; then
    source ".version"
fi


echo "Building qrvc $VERSION"

#tidy up the go.mod file
go mod tidy

# check for vulnerabilities
govulncheck ./...

#application name
app="qrvc"

#distribution folder
dist="./dist"
rm -rf "$dist"

build_for_target() {
    os=$1
    arch=$2
    suffix=$3

    echo "Building $os/$arch"
    GOOS=$os GOARCH=$arch go build -ldflags="-X qrvc/internal/version.Version=$VERSION" -o "$dist/$os/$arch/${app}${suffix}" .
}


# build Mac on Apple Silicon
build_for_target darwin arm64 ""

# build Mac on Intel/AMD
build_for_target darwin amd64 ""

# build Win on Intel/AMD
build_for_target windows amd64 .exe

# build Win on ARM
build_for_target windows arm64 .exe

# build Linux on Intel/AMD
build_for_target linux amd64 ""

# build Linux on ARM
build_for_target linux arm64 ""

# Load variables from .env to see if this is running on my local dev machine
if [ -f ".env" ]; then
    source ".env"
fi
if [ -n "$AT_HOME" ]; then
    echo "I´m at home, therefore copying $dist/darwin/arm64/qrvc to ~/go/bin/"
    cp "$dist/darwin/arm64/qrvc" ~/go/bin/
fi

echo "Ready 👋"
