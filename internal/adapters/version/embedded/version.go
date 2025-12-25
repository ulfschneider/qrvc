package versionembedded

import (
	"embed"
	"strings"
)

type VersionProvider struct {
}

func NewVersionProvider() VersionProvider {
	return VersionProvider{}
}

//go:embed generated/*
var generated embed.FS

// version
const versionPath = "generated/version.txt"

func readEmbeddedData(filePath string) string {
	f, err := generated.Open(filePath)
	if err != nil {
		return ""
	}
	defer f.Close()

	data, err := generated.ReadFile(filePath)
	if err != nil {
		return ""
	}
	return strings.TrimSpace(string(data))
}

func (vp *VersionProvider) Version() string {
	return readEmbeddedData(versionPath)
}
