package qrcodec_test

import (
	"image"
	"image/draw"
	"testing"

	"github.com/mazznoer/csscolorparser"
	"github.com/skip2/go-qrcode"
	"github.com/stretchr/testify/assert"

	qrcodec "github.com/ulfschneider/qrvc/internal/adapters/codec/qr"
	vcardcodec "github.com/ulfschneider/qrvc/internal/adapters/codec/vcard"
	testutil "github.com/ulfschneider/qrvc/internal/adapters/test/util"

	"github.com/ulfschneider/qrvc/internal/application/config"
)

func TestQRCodec(t *testing.T) {
	card := testutil.CreateCard()
	cardCodec := vcardcodec.NewCodec()
	vcf, _ := cardCodec.Encode(card)

	backgroundColor, _ := csscolorparser.Parse("transparent")
	foregroundColor, _ := csscolorparser.Parse("orange")
	testSettings := config.QRCodeSettings{Border: false, Size: 300, RecoveryLevel: qrcode.Low, BackgroundColor: foregroundColor, ForegroundColor: backgroundColor}

	expectedImg, err := makeQRCode(vcf, testSettings)
	assert.NoError(t, err)

	qrCodec := qrcodec.NewCodec()
	resultImg, _ := qrCodec.Encode(card, testSettings)

	assert.Equal(t, toRGBA(expectedImg).Pix, toRGBA(resultImg).Pix)
}

func makeQRCode(vcf []byte, settings config.QRCodeSettings) (image.Image, error) {
	q, err := qrcode.New(string(vcf), settings.RecoveryLevel)
	if err != nil {
		return nil, err
	}

	q.DisableBorder = !settings.Border
	q.ForegroundColor = settings.ForegroundColor
	q.BackgroundColor = settings.BackgroundColor

	img := q.Image(settings.Size)

	return img, nil
}

func toRGBA(img image.Image) *image.RGBA {
	b := img.Bounds()
	rgba := image.NewRGBA(b)
	draw.Draw(rgba, b, img, b.Min, draw.Src)
	return rgba
}
