package qrcard

import (
	"testing"

	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"

	"github.com/emersion/go-vcard"
)

const testVCardContent = `BEGIN:VCARD
VERSION:3.0
ADR:Post office box;Extended street address;Street address;City;;Postal code;Country
EMAIL:Mail
GENDER:O
N:Family name;Given name;Additional name;Honorific prefix;Honorific suffix
ORG:Organization or company;Department
TEL;TYPE=cell:Cell phone
TEL;TYPE=work:Work phone
TEL;TYPE=home:Private phone
TITLE:Job title
URL:Web address
END:VCARD
`

func TestMakeVCardInstanceFromNonExistingFile(t *testing.T) {
	filesystem := afero.NewMemMapFs()
	filePath := "vcard"

	//vcard file does not exist, there is nothing to read
	vcard, err := makeVCardInstance(&filePath, "3.0", filesystem)
	assert.Error(t, err)
	assert.Nil(t, vcard)

	//file does still not exist
	_, err = filesystem.Stat(filePath)
	assert.Error(t, err)
}

func TestMakeVCardInstanceFromExistingFile(t *testing.T) {
	filesystem := afero.NewMemMapFs()
	filePath := "vcard.vcf"
	simpleFilePath := "vcard"

	//create file
	f, err := filesystem.Create(filePath)
	assert.NoError(t, err)
	assert.NotEmpty(t, f)

	f.Write([]byte(testVCardContent))

	//vcard.vcf file does exist
	card, err := makeVCardInstance(&filePath, "3.0", filesystem)
	assert.NoError(t, err)
	assert.NotEmpty(t, card)

	//vcard (.vcf is added automatically) can as well be used to access the file
	card, err = makeVCardInstance(&simpleFilePath, "3.0", filesystem)
	assert.NoError(t, err)
	assert.NotEmpty(t, card)

	//file does still exist
	info, err := filesystem.Stat(filePath)
	assert.NoError(t, err)
	assert.NotEmpty(t, info)

	//verify vcard content
	assert.Equal(t, "3.0", card.Value(vcard.FieldVersion))

	assert.Equal(t, "Given name", card.Name().GivenName)
	assert.Equal(t, "Additional name", card.Name().AdditionalName)
	assert.Equal(t, "Family name", card.Name().FamilyName)
	assert.Equal(t, "Honorific prefix", card.Name().HonorificPrefix)
	assert.Equal(t, "Honorific suffix", card.Name().HonorificSuffix)

	sex, identity := card.Gender()
	assert.Equal(t, vcard.Sex("O"), sex)
	assert.Equal(t, "", identity)

	assert.Equal(t, "Job title", card.Value(vcard.FieldTitle))
	assert.Equal(t, "Organization or company;Department", card.Value(vcard.FieldOrganization))

	assert.Equal(t, "Mail", card.Value(vcard.FieldEmail))
	assert.Equal(t, "Web address", card.Value(vcard.FieldURL))
	assert.Equal(t, "Cell phone", typedVcardFieldValue(card, vcard.FieldTelephone, vcard.TypeCell))
	assert.Equal(t, "Work phone", typedVcardFieldValue(card, vcard.FieldTelephone, vcard.TypeWork))
	assert.Equal(t, "Private phone", typedVcardFieldValue(card, vcard.FieldTelephone, vcard.TypeHome))

	assert.Equal(t, "Post office box", card.Address().PostOfficeBox)
	assert.Equal(t, "Street address", card.Address().StreetAddress)
	assert.Equal(t, "Extended street address", card.Address().ExtendedAddress)
	assert.Equal(t, "City", card.Address().Locality)
	assert.Equal(t, "Postal code", card.Address().PostalCode)
	assert.Equal(t, "Country", card.Address().Country)
}

func TestFormHandling(t *testing.T) {
	filesystem := afero.NewMemMapFs()
	filePath := "vcard.vcf"

	//create file
	f, _ := filesystem.Create(filePath)

	f.Write([]byte(testVCardContent))

	//vcard.vcf file does exist
	card, _ := makeVCardInstance(&filePath, "3.0", filesystem)

	formData := transferVCardIntoFormData(card)
	assert.Equal(t, *card.Name(), formData.name)
	assert.Equal(t, *card.Address(), formData.address)
	assert.Equal(t, "Mail", formData.email)
	assert.Equal(t, "Web address", formData.url)
	assert.Equal(t, "Cell phone", formData.cellPhone)
	assert.Equal(t, "Work phone", formData.workPhone)
	assert.Equal(t, "Private phone", formData.homePhone)
	assert.Equal(t, "Organization or company", formData.organization)
	assert.Equal(t, "Department", formData.department)

	//modifiy form data
	formData.name.GivenName = "given"
	formData.name.AdditionalName = "additional"
	formData.name.FamilyName = "family"
	formData.name.HonorificPrefix = "prefix"
	formData.name.HonorificSuffix = "suffix"
	formData.gender = vcard.SexNone
	formData.title = "title"
	formData.organization = "organization"
	formData.department = "department"
	formData.email = "email address"
	formData.url = "url"
	formData.cellPhone = "cell phone number"
	formData.workPhone = "work phone number"
	formData.homePhone = "home phone number"
	formData.address.PostOfficeBox = "post office box"
	formData.address.StreetAddress = "street address"
	formData.address.ExtendedAddress = "extended address"
	formData.address.PostalCode = "postal code"
	formData.address.Locality = "city"
	formData.address.Country = "country"

	//bring the input form data back into the vcard
	transferFormDataIntoVCard(card, formData)

	assert.Equal(t, "given", card.Name().GivenName)
	assert.Equal(t, "additional", card.Name().AdditionalName)
	assert.Equal(t, "family", card.Name().FamilyName)
	assert.Equal(t, "prefix", card.Name().HonorificPrefix)
	assert.Equal(t, "suffix", card.Name().HonorificSuffix)
	sex, identity := card.Gender()
	assert.Equal(t, vcard.SexNone, sex)
	assert.Equal(t, "", identity)
	assert.Equal(t, "organization;department", card.Value(vcard.FieldOrganization))
	assert.Equal(t, "title", card.Value(vcard.FieldTitle))
	assert.Equal(t, "email address", card.Value(vcard.FieldEmail))
	assert.Equal(t, "url", card.Value(vcard.FieldURL))
	assert.Equal(t, "cell phone number", typedVcardFieldValue(card, vcard.FieldTelephone, vcard.TypeCell))
	assert.Equal(t, "work phone number", typedVcardFieldValue(card, vcard.FieldTelephone, vcard.TypeWork))
	assert.Equal(t, "home phone number", typedVcardFieldValue(card, vcard.FieldTelephone, vcard.TypeHome))
	assert.Equal(t, "extended address", card.Address().ExtendedAddress)
	assert.Equal(t, "street address", card.Address().StreetAddress)
	assert.Equal(t, "post office box", card.Address().PostOfficeBox)
	assert.Equal(t, "city", card.Address().Locality)
	assert.Equal(t, "postal code", card.Address().PostalCode)
	assert.Equal(t, "country", card.Address().Country)
}

func TestWriteResults(t *testing.T) {
	filesystem := afero.NewMemMapFs()
	filePath := "vcard.vcf"

	//create file
	f, _ := filesystem.Create(filePath)

	f.Write([]byte(testVCardContent))

	//vcard.vcf file does exist
	card, _ := makeVCardInstance(&filePath, "3.0", filesystem)

	vcardContent, err := encodeVcard(card)
	assert.NotEmpty(t, vcardContent)
	assert.NoError(t, err)

	// TODO test the write results
}
