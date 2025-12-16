package bomembedded_test

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	bomembedded "github.com/ulfschneider/qrvc/internal/adapters/bom/embedded"
)

func TestBom(t *testing.T) {
	envVersion := os.Getenv("VERSION")

	if envVersion != "" {
		bomProvider := bomembedded.BomProvider{}
		b, err := bomProvider.Bom()
		assert.NoError(t, err)
		assert.NotEmpty(t, b)
	}
}

func TestBOMToJSON(t *testing.T) {
	envVersion := os.Getenv("VERSION")
	if envVersion != "" {
		bomProvider := bomembedded.BomProvider{}
		j, err := bomProvider.MarshalToJSON()
		assert.NoError(t, err)
		assert.NotEmpty(t, j)
	}
}
