package sbom

import _ "embed"

//go:embed sbom.json
var sbom []byte

func SBOM() []byte {
	return sbom
}
