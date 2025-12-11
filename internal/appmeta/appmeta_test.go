package appmeta_test

import (
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/ulfschneider/qrvc/internal/appmeta"
)

func ensureVPrefix(s string) string {
	if strings.HasPrefix(s, "v") {
		return s
	}
	return "v" + s
}

func TestVersion(t *testing.T) {
	envVersion := os.Getenv("VERSION")
	v, err := appmeta.LoadEmbeddedVersion()
	if envVersion != "" {
		assert.NoError(t, err)
		assert.Equal(t, ensureVPrefix(envVersion), v)
	}
}

func TestBom(t *testing.T) {
	envVersion := os.Getenv("VERSION")
	if envVersion != "" {
		b, err := appmeta.LoadEmbeddedBOM()
		assert.NoError(t, err)
		assert.NotEmpty(t, b)
	}
}

func TestBOMToJSON(t *testing.T) {
	envVersion := os.Getenv("VERSION")
	if envVersion != "" {
		b, _ := appmeta.LoadEmbeddedBOM()
		j, err := appmeta.MarshalBOMToJSON(b)
		assert.NoError(t, err)
		assert.NotEmpty(t, j)
	}
}
