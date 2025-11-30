package licenses

import (
	"embed"
)

//go:embed generated/*
var LicensesFS embed.FS
