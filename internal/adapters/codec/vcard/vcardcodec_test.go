package vcardcodec_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	vcardcodec "github.com/ulfschneider/qrvc/internal/adapters/codec/vcard"
	testutil "github.com/ulfschneider/qrvc/internal/test/util"
)

func TestVCardCodec(t *testing.T) {
	card := testutil.CreateCard()
	codec := vcardcodec.NewCodec()
	vcf, _ := codec.Encode(card)
	assert.Equal(t, testutil.ExpectedVCF, testutil.NormalizeNewLines(string(vcf)))
}
