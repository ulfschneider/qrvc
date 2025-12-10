package appmeta_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/ulfschneider/qrvc/internal/appmeta"
)

func TestVersion(t *testing.T) {
	v, err := appmeta.LoadEmbeddedVersion()
	assert.NoError(t, err)
	assert.NotEmpty(t, v)
}

func TestBom(t *testing.T) {
	b, err := appmeta.LoadEmbeddedBOM()
	assert.NoError(t, err)
	assert.NotEmpty(t, b)
}

func TestBOMToJSON(t *testing.T) {
	b, _ := appmeta.LoadEmbeddedBOM()
	j, err := appmeta.MarshalBOMToJSON(b)
	assert.NoError(t, err)
	assert.NotEmpty(t, j)
}
