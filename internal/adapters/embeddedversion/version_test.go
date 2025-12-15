package embeddedversion_test

import (
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/ulfschneider/qrvc/internal/adapters/embeddedversion"
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
		versionProvider := embeddedversion.VersionProvider{}
		v, err := versionProvider.Version()
		assert.NoError(t, err)
		assert.Equal(t, ensureVPrefix(envVersion), v)
	}
}
