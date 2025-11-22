#! /bin/bash
echo "Building qrvc"
# check for vulnerabilities
echo "Checking for vulnerabilities"
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
    GOOS=$os GOARCH=$arch go build -o "$dist/$os/$arch/${app}${suffix}" .
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

echo "Ready 👋"
