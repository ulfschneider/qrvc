package vcardcodec_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/ulfschneider/qrvc/internal/adapters/testutil"
	"github.com/ulfschneider/qrvc/internal/adapters/vcardcodec"
)

func TestVCardCodec(t *testing.T) {
	card := testutil.CreateCard()
	codec := vcardcodec.NewCodec()
	vcf, _ := codec.Encode(card)
	assert.Equal(t, testutil.ExpectedVCF, testutil.NormalizeNewLines(string(vcf)))
}
