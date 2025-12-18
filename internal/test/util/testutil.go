package testutil

import (
	"image"
	"image/draw"
	"strings"

	"github.com/emersion/go-vcard"

	qrcodec "github.com/ulfschneider/qrvc/internal/adapters/codec/qr"
	vcardcodec "github.com/ulfschneider/qrvc/internal/adapters/codec/vcard"
	configcli "github.com/ulfschneider/qrvc/internal/adapters/config/cli"
	"github.com/ulfschneider/qrvc/internal/application/config"
	"github.com/ulfschneider/qrvc/internal/application/services"
	qrcard "github.com/ulfschneider/qrvc/internal/domain"
)

func NormalizeNewLines(s string) string {
	return strings.ReplaceAll(s, "\r\n", "\n")
}

var ExpectedVCF = NormalizeNewLines(`BEGIN:VCARD
VERSION:3.0
ADR:Post office box;Extended street address;Street address;City;;Postal code;Country
EMAIL:Email address
GENDER:N
N:Family name;Given name;Additional name;Honorific prefix;Honorific suffix
ORG:Organization or company;Department
TEL;TYPE=cell:Cell phone
TEL;TYPE=work:Work phone
TEL;TYPE=home:Home phone
TITLE:Job title
URL:Web address
END:VCARD
`)

func CreateCard() vcard.Card {
	card := vcard.Card{}
	card.SetValue(vcard.FieldVersion, "3.0")
	card.SetAddress(&vcard.Address{
		PostOfficeBox:   "Post office box",
		ExtendedAddress: "Extended street address",
		StreetAddress:   "Street address",
		Locality:        "City",
		PostalCode:      "Postal code",
		Country:         "Country"})
	card.SetValue(vcard.FieldEmail, "Email address")
	card.SetValue(vcard.FieldURL, "Web address")
	card.SetValue(vcard.FieldTitle, "Job title")
	card.SetValue(vcard.FieldOrganization, "Organization or company;Department")
	card.SetGender(vcard.SexNone, "")
	card.SetName(&vcard.Name{
		GivenName:       "Given name",
		FamilyName:      "Family name",
		AdditionalName:  "Additional name",
		HonorificPrefix: "Honorific prefix",
		HonorificSuffix: "Honorific suffix"})
	qrcard.SetTypedVcardFieldValue(card, vcard.FieldTelephone, vcard.TypeCell, "Cell phone")
	qrcard.SetTypedVcardFieldValue(card, vcard.FieldTelephone, vcard.TypeWork, "Work phone")
	qrcard.SetTypedVcardFieldValue(card, vcard.FieldTelephone, vcard.TypeHome, "Home phone")

	return card
}

func CreateQRCode(card vcard.Card, settings config.QRCodeSettings) image.Image {
	codec := qrcodec.NewCodec()
	img, _ := codec.Encode(card, settings)
	return img
}

func EncodeCard(card vcard.Card) []byte {
	codec := vcardcodec.NewCodec()
	content, _ := codec.Encode(card)
	return content
}

type testVersionProvider struct{}

func CreateVersionProvider() *testVersionProvider {
	return &testVersionProvider{}
}

func (vp *testVersionProvider) Version() string {
	return "TEST VERSION"
}

func (vp *testVersionProvider) Commit() string {
	return "Commit"
}

func (vp *testVersionProvider) Time() string {
	return "Time"
}

func LoadTestSettings() configcli.CLIFileSettings {
	var versionService = services.NewVersionService(CreateVersionProvider())
	var settingsProvider = configcli.NewSettingsProvider(versionService)

	settings, _ := settingsProvider.Load()
	return settings
}

func ToRGBA(img image.Image) *image.RGBA {
	b := img.Bounds()
	rgba := image.NewRGBA(b)
	draw.Draw(rgba, b, img, b.Min, draw.Src)
	return rgba
}
