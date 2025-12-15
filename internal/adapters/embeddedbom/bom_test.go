package embeddedbom_test

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/ulfschneider/qrvc/internal/adapters/embeddedbom"
)

func TestBom(t *testing.T) {
	envVersion := os.Getenv("VERSION")

	if envVersion != "" {
		bomProvider := embeddedbom.BomProvider{}
		b, err := bomProvider.Bom()
		assert.NoError(t, err)
		assert.NotEmpty(t, b)
	}
}

func TestBOMToJSON(t *testing.T) {
	envVersion := os.Getenv("VERSION")
	if envVersion != "" {
		bomProvider := embeddedbom.BomProvider{}
		j, err := bomProvider.MarshalToJSON()
		assert.NoError(t, err)
		assert.NotEmpty(t, j)
	}
}
