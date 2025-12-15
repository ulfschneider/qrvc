# qrvc

![Version](https://img.shields.io/github/v/tag/ulfschneider/qrvc?sort=semver&label=version)
![Go Reference](https://pkg.go.dev/badge/github.com/ulfschneider/qrvc.svg)
![Go Report Card](https://goreportcard.com/badge/github.com/ulfschneider/qrvc)
![License](https://img.shields.io/github/license/ulfschneider/qrvc)

qrvc is command line tool to prepare a QR code from a vCard.

## Install

```sh
go install github.com/ulfschneider/qrvc@latest
```

## Usage

After installation, start the tool with the `-h` flag to get information about how to use it:

```sh
qrvc -h
```

## Build

The project contains a Makefile which allows you to build the tool yourself:

```sh
make build
```

To get a list of all available make targets, call:

```sh
make
```

## License

MIT
