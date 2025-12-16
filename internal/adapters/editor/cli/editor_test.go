package editorcli

import (
	"testing"

	"github.com/emersion/go-vcard"
	"github.com/stretchr/testify/assert"

	testutil "github.com/ulfschneider/qrvc/internal/adapters/test/util"
	qrcard "github.com/ulfschneider/qrvc/internal/domain"
)

func TestEditor(t *testing.T) {

	card := testutil.CreateCard()

	formData := transferVCardIntoFormData(card)
	assert.Equal(t, card.Name(), formData.name)
	assert.Equal(t, card.Address(), formData.address)
	assert.Equal(t, "Email address", formData.email)
	assert.Equal(t, "Web address", formData.url)
	assert.Equal(t, "Cell phone", formData.cellPhone)
	assert.Equal(t, "Work phone", formData.workPhone)
	assert.Equal(t, "Home phone", formData.homePhone)
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
	assert.Equal(t, "cell phone number", qrcard.TypedVcardFieldValue(card, vcard.FieldTelephone, vcard.TypeCell))
	assert.Equal(t, "work phone number", qrcard.TypedVcardFieldValue(card, vcard.FieldTelephone, vcard.TypeWork))
	assert.Equal(t, "home phone number", qrcard.TypedVcardFieldValue(card, vcard.FieldTelephone, vcard.TypeHome))
	assert.Equal(t, "extended address", card.Address().ExtendedAddress)
	assert.Equal(t, "street address", card.Address().StreetAddress)
	assert.Equal(t, "post office box", card.Address().PostOfficeBox)
	assert.Equal(t, "city", card.Address().Locality)
	assert.Equal(t, "postal code", card.Address().PostalCode)
	assert.Equal(t, "country", card.Address().Country)

}
