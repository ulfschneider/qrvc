# qrvc

[![Go Reference](https://pkg.go.dev/badge/github.com/ulfschneider/qrvc.svg)](https://pkg.go.dev/github.com/ulfschneider/qrvc)
[![Go Report Card](https://goreportcard.com/badge/github.com/ulfschneider/qrvc)](https://goreportcard.com/report/github.com/ulfschneider/qrvc)
![License](https://img.shields.io/github/license/ulfschneider/qrvc)

qrvc is command line tool to prepare a QR code from a vCard.

## Install

```sh
go install github.com/ulfschneider/qrvc@latest
```

## Usage

After installation, start up the took with the `-h` flag to get information about how to use it.

```sh
qrvc -h
```

## Build

The project contains a Makefile which allows you to build the tool yourself. 

```sh
make build
```

To get a list of all available make targets, call

```sh
make
```

## License

MIT
