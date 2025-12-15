package embeddedversion

import (
	"embed"
)

type EmbeddedVersionProvider struct {
}

func NewVersionProvider() EmbeddedVersionProvider {
	return EmbeddedVersionProvider{}
}

//go:embed generated/*
var generated embed.FS

// version
const versionPath = "generated/version.txt"

func (vp *EmbeddedVersionProvider) Version() (string, error) {
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
