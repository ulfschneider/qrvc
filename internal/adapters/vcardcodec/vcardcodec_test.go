package vcardcodec_test

import (
	"strings"
	"testing"

	"github.com/emersion/go-vcard"
	"github.com/stretchr/testify/assert"
	"github.com/ulfschneider/qrvc/internal/adapters/vcardcodec"
)

func normalizeNewLines(s string) string {
	return strings.ReplaceAll(s, "\r\n", "\n")
}

var expectedVCF = normalizeNewLines(`BEGIN:VCARD
VERSION:3.0
ADR:Post office box;Extended street address;Street address;City;;Postal code;Country
EMAIL:Email address
GENDER:N
N:Family name;Given name;Additional name;Honorific prefix;Honorific suffix
ORG:Organization or company;Department
TITLE:Job title
URL:Web address
END:VCARD
`)

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

	return card
}

func TestVCardCodec(t *testing.T) {
	card := makeTestCard()
	codec := vcardcodec.NewCodec()
	vcf, _ := codec.Encode(card)
	assert.Equal(t, expectedVCF, normalizeNewLines(string(vcf)))
}
