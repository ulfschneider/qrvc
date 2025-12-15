package qrcodec_test

import (
	"image"
	"image/draw"
	"testing"

	"github.com/emersion/go-vcard"
	"github.com/mazznoer/csscolorparser"
	"github.com/skip2/go-qrcode"
	"github.com/stretchr/testify/assert"

	"github.com/ulfschneider/qrvc/internal/adapters/qrcodec"
	"github.com/ulfschneider/qrvc/internal/adapters/vcardcodec"
	"github.com/ulfschneider/qrvc/internal/application/config"
)

func makeTestCard() vcard.Card {
	card := vcard.Card{}
	card.SetValue(vcard.FieldVersion, "3.0")
	card.SetAddress(&vcard.Address{
		PostOfficeBox:   "Post office box",
		ExtendedAddress: "Extended street address",
		StreetAddress:   "Street address",
		Locality:        "City",
		PostalCode:      "Postal code",
		Country:         "Country"})
	card.SetValue(vcard.FieldEmail, "EMAIL")
	card.SetValue(vcard.FieldURL, "Web address")
	card.SetValue(vcard.FieldTitle, "Job title")
	card.SetValue(vcard.FieldOrganization, "Organization")
	card.SetGender(vcard.SexNone, "")
	card.SetName(&vcard.Name{
		GivenName:       "Given name",
		FamilyName:      "Family name",
		AdditionalName:  "Additional name",
		HonorificPrefix: "Prefix",
		HonorificSuffix: "Suffix"})

	return card
}

func TestQRCodec(t *testing.T) {
	card := makeTestCard()
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
