# qrvc

![Version](https://img.shields.io/github/v/tag/ulfschneider/qrvc?sort=semver&label=version)
[![Go Reference](https://pkg.go.dev/badge/github.com/ulfschneider/qrvc.svg)](https://pkg.go.dev/github.com/ulfschneider/qrvc)
[![Go Report Card](https://goreportcard.com/badge/github.com/ulfschneider/qrvc)](https://goreportcard.com/report/github.com/ulfschneider/qrvc)
![License](https://img.shields.io/github/license/ulfschneider/qrvc)

qrvc is command line tool to prepare a QR code from a vCard.

## Installation

### Install with Homebrew on Mac and Linux

```sh
brew tap ulfschneider/tap
brew install --cask qrvc
```

### Install with Go on any machine that has Go on board

```sh
go install github.com/ulfschneider/qrvc@latest
```

### Manual install

You can also download the appropriate binary directly from GitHub Releases:

Visit [github.com/ulfschneider/qrvc/releases](https://github.com/ulfschneider/qrvc/releases)

1.	Download the archive matching your OS and architecture
2.	Extract it
3.	Move the binary to a directory included in your PATH (for example /usr/local/bin)

### Verify the installation

```sh
qrvc --version
```

This command should print out the qrvc version you are using.

## Usage

After installation, start the tool by typing:

```sh
qrvc
```

When using the `-h` flag, you get information about how to use it:

```sh
qrvc -h
```

## Issues

Please file issues at [github.com/ulfschneider/qrvc/issues](https://github.com/ulfschneider/qrvc/issues).

## License

MIT
