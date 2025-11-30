package sbom

import (
	"embed"
	_ "embed"
)

//go:embed generated/*
var SBOMFS embed.FS
