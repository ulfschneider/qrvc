package embeddedversion

import (
	"embed"
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

func (vp *VersionProvider) Version() (string, error) {
	f, err := generated.Open(versionPath)
	if err != nil {
		return "", err
	}
	defer f.Close()

	version, err := generated.ReadFile(versionPath)
	if err != nil {
		return "", err
	}
	return string(version), nil
}
