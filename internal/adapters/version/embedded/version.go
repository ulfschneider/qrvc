package versionembedded

import (
	"embed"
	"runtime/debug"
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

func (vp *VersionProvider) Version() string {
	f, err := generated.Open(versionPath)
	if err != nil {
		return ""
	}
	defer f.Close()

	version, err := generated.ReadFile(versionPath)
	if err != nil {
		return ""
	}
	return string(version)
}

func (vp *VersionProvider) Commit() string {
	if info, ok := debug.ReadBuildInfo(); ok {
		for _, setting := range info.Settings {
			if setting.Key == "vcs.revision" {
				return setting.Value
			}
		}
	}
	return ""
}

func (vp *VersionProvider) Time() string {
	if info, ok := debug.ReadBuildInfo(); ok {
		for _, setting := range info.Settings {
			if setting.Key == "vcs.time" {
				return setting.Value
			}
		}
	}
	return ""
}
