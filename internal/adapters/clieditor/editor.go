package clieditor

import (
	"strings"

	"github.com/charmbracelet/huh"
	"github.com/emersion/go-vcard"
	qrcard "github.com/ulfschneider/qrvc/internal/domain"
)

type CardEditor struct {
}

func NewCardEditor() CardEditor {
	return CardEditor{}
}

func (e *CardEditor) Edit(card vcard.Card) error {
	formData := transferVCardIntoFormData(card)

	for {
		form := prepareForm(formData)
		if err := form.Run(); err != nil {
			return err
		}
		if formData.ready {
			break
		}
	}

	transferFormDataIntoVCard(card, formData)
	return nil
}

type qrCardFormData struct {
	name         *vcard.Name
	gender       vcard.Sex
	title        string
	organization string
	department   string
	address      *vcard.Address
	email        string
	url          string
	cellPhone    string
	workPhone    string
	homePhone    string
	ready        bool
}

func maybeGet(s []string, i int) string {
	if i < len(s) {
		return s[i]
	}
	return ""
}

func transferVCardIntoFormData(card vcard.Card) qrCardFormData {

	sex, _ := card.Gender()

	orgSplit := strings.SplitN(card.Value(vcard.FieldOrganization), ";", 2)
	organization := maybeGet(orgSplit, 0)
	department := maybeGet(orgSplit, 1)

	data := qrCardFormData{
		name:         card.Name(),
		gender:       sex,
		title:        card.Value(vcard.FieldTitle),
		organization: organization,
		department:   department,
		address:      card.Address(),
		email:        card.Value(vcard.FieldEmail),
		url:          card.Value(vcard.FieldURL),
		cellPhone:    qrcard.TypedVcardFieldValue(card, vcard.FieldTelephone, vcard.TypeCell),
		workPhone:    qrcard.TypedVcardFieldValue(card, vcard.FieldTelephone, vcard.TypeWork),
		homePhone:    qrcard.TypedVcardFieldValue(card, vcard.FieldTelephone, vcard.TypeHome),
		ready:        true,
	}

	return data
}

func transferFormDataIntoVCard(card vcard.Card, formData qrCardFormData) {
	card.SetName(formData.name)
	card.SetGender(vcard.Sex(formData.gender), "")
	card.SetValue(vcard.FieldTitle, formData.title)
	card.SetValue(vcard.FieldOrganization, formData.organization+";"+formData.department)
	card.SetAddress(formData.address)
	card.SetValue(vcard.FieldEmail, formData.email)
	card.SetValue(vcard.FieldURL, formData.url)
	qrcard.SetTypedVcardFieldValue(card, vcard.FieldTelephone, vcard.TypeCell, formData.cellPhone)
	qrcard.SetTypedVcardFieldValue(card, vcard.FieldTelephone, vcard.TypeWork, formData.workPhone)
	qrcard.SetTypedVcardFieldValue(card, vcard.FieldTelephone, vcard.TypeHome, formData.homePhone)
}

func prepareForm(formData qrCardFormData) *huh.Form {

	vCardForm := huh.NewForm(
		huh.NewGroup(
			huh.NewInput().Title("Given (first) name").Value(&formData.name.GivenName),
			huh.NewInput().Title("Additional (middle) name").Value(&formData.name.AdditionalName),
			huh.NewInput().Title("Family name").Value(&formData.name.FamilyName),
			huh.NewInput().Title("Honorific prefix (e.g. Capt.)").Value(&formData.name.HonorificPrefix),
			huh.NewInput().Title("Honorific suffix (e.g. Sr.)").Value(&formData.name.HonorificSuffix),
		),
		huh.NewGroup(
			huh.NewSelect[vcard.Sex]().Title("Gender").Options(
				huh.NewOption("Male", vcard.SexMale).Selected(vcard.SexMale == formData.gender),
				huh.NewOption("Female", vcard.SexFemale).Selected(vcard.SexFemale == formData.gender),
				huh.NewOption("Other", vcard.SexOther).Selected(vcard.SexOther == formData.gender),
				huh.NewOption("Unspecified", vcard.SexUnspecified).Selected(formData.gender != vcard.SexMale && formData.gender != vcard.SexFemale && formData.gender != vcard.SexUnspecified),
			).Value(&formData.gender),
		),
		huh.NewGroup(
			huh.NewInput().Title("Job title").Value(&formData.title),
			huh.NewInput().Title("Organization or company").Value(&formData.organization),
			huh.NewInput().Title("Department").Value(&formData.department),
		),

		huh.NewGroup(
			huh.NewInput().Title("Mail").Value(&formData.email),
			huh.NewInput().Title("Web address").Value(&formData.url),
			huh.NewInput().Title("Cell phone").Value(&formData.cellPhone),
			huh.NewInput().Title("Work phone").Value(&formData.workPhone),
			huh.NewInput().Title("Private phone").Value(&formData.homePhone),
		),

		huh.NewGroup(
			huh.NewInput().Title("Post office box").Value(&formData.address.PostOfficeBox),
			huh.NewInput().Title("Street address").Value(&formData.address.StreetAddress),
			huh.NewInput().Title("Extended street address (e.g. building, floor)").Value(&formData.address.ExtendedAddress),
			huh.NewInput().Title("City").Value(&formData.address.Locality),
			huh.NewInput().Title("Postal code").Value(&formData.address.PostalCode),
			huh.NewInput().Title("Country").Value(&formData.address.Country),
		),

		huh.NewGroup(
			huh.NewConfirm().
				Title("Are you ready?").
				Affirmative("Yes, print the result!").
				Negative("No, I´m not ready.").
				Value(&formData.ready),
		),
	).WithTheme(huh.ThemeBase16())
	return vCardForm
}
