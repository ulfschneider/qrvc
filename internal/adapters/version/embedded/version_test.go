package versionembedded_test

import (
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	versionembedded "github.com/ulfschneider/qrvc/internal/adapters/version/embedded"
)

func ensureVPrefix(s string) string {
	if strings.HasPrefix(s, "v") {
		return s
	}
	return "v" + s
}

func TestVersion(t *testing.T) {
	envVersion := os.Getenv("VERSION")

	if envVersion != "" {
		versionProvider := versionembedded.VersionProvider{}
		v, err := versionProvider.Version()
		assert.NoError(t, err)
		assert.Equal(t, ensureVPrefix(envVersion), v)
	}
}
