package vcard

import (
	"bytes"
	"io"
	"os"
	"qrvc/internal/cli"
	"qrvc/internal/settings"
	"slices"
	"strings"

	"github.com/charmbracelet/huh"
	"github.com/emersion/go-vcard"
)

type VCardFormData struct {
	GivenName       string
	AdditionalName  string
	FamilyName      string
	HonorificPrefix string
	HonorificSuffix string
	Gender          vcard.Sex
	Title           string
	Organization    string
	Department      string
	PostOfficeBox   string
	StreetAddress   string
	ExtendedAddress string
	Locality        string
	PostalCode      string
	Country         string
	Email           string
	Url             string
	CellPhone       string
	WorkPhone       string
	HomePhone       string
	Ready           bool
}

func maybeGet(l []string, i int) string {
	if i < len(l) {
		return l[i]
	}
	return ""
}

func prepareFormData(card *vcard.Card) *VCardFormData {

	sex, _ := card.Gender()

	orgSplit := strings.SplitN(card.Value(vcard.FieldOrganization), ";", 2)
	organization := orgSplit[0]
	department := maybeGet(orgSplit, 1)

	data := VCardFormData{
		GivenName:       card.Name().GivenName,
		AdditionalName:  card.Name().AdditionalName,
		FamilyName:      card.Name().FamilyName,
		HonorificPrefix: card.Name().HonorificPrefix,
		HonorificSuffix: card.Name().HonorificSuffix,
		Gender:          sex,
		Title:           card.Value(vcard.FieldTitle),
		Organization:    organization,
		Department:      department,
		PostOfficeBox:   card.Address().PostOfficeBox,
		StreetAddress:   card.Address().StreetAddress,
		ExtendedAddress: card.Address().ExtendedAddress,
		Locality:        card.Address().Locality,
		PostalCode:      card.Address().PostalCode,
		Country:         card.Address().Country,
		Email:           card.Value(vcard.FieldEmail),
		Url:             card.Value(vcard.FieldURL),
		CellPhone:       typedVcardFieldValue(card, vcard.FieldTelephone, vcard.TypeCell),
		WorkPhone:       typedVcardFieldValue(card, vcard.FieldTelephone, vcard.TypeWork),
		HomePhone:       typedVcardFieldValue(card, vcard.FieldTelephone, vcard.TypeHome),
		Ready:           true,
	}

	return &data
}

func prepareVCardData(card *vcard.Card, data *VCardFormData) {
	card.SetGender(vcard.Sex(data.Gender), "")
	card.SetValue(vcard.FieldOrganization, data.Organization+";"+data.Department)
	setTypedVcardFieldValue(card, vcard.FieldTelephone, vcard.TypeCell, data.CellPhone)
	setTypedVcardFieldValue(card, vcard.FieldTelephone, vcard.TypeWork, data.WorkPhone)
	setTypedVcardFieldValue(card, vcard.FieldTelephone, vcard.TypeHome, data.HomePhone)
}

func prepareForm(data *VCardFormData) *huh.Form {

	vCardForm := huh.NewForm(
		huh.NewGroup(
			huh.NewInput().Title("Given (first) name").Value(&data.GivenName),
			huh.NewInput().Title("Additional (middle) name").Value(&data.AdditionalName),
			huh.NewInput().Title("Family name").Value(&data.FamilyName),
			huh.NewInput().Title("Honorific prefix (e.g. Capt.)").Value(&data.HonorificPrefix),
			huh.NewInput().Title("Honorific suffix (e.g. Sr.)").Value(&data.HonorificSuffix),
		),
		huh.NewGroup(
			huh.NewSelect[vcard.Sex]().Title("Gender").Options(
				huh.NewOption("Male", vcard.SexMale).Selected(vcard.SexMale == data.Gender),
				huh.NewOption("Female", vcard.SexFemale).Selected(vcard.SexFemale == data.Gender),
				huh.NewOption("Other", vcard.SexOther).Selected(vcard.SexOther == data.Gender),
				huh.NewOption("Unspecified", vcard.SexUnspecified).Selected(data.Gender != vcard.SexMale && data.Gender != vcard.SexFemale && data.Gender != vcard.SexUnspecified),
			).Value(&data.Gender),
		),
		huh.NewGroup(
			huh.NewInput().Title("Job title").Value(&data.Title),
			huh.NewInput().Title("Organization or company").Value(&data.Organization),
			huh.NewInput().Title("Department").Value(&data.Department),
		),

		huh.NewGroup(
			huh.NewInput().Title("Mail").Value(&data.Email),
			huh.NewInput().Title("Web address").Value(&data.Url),
			huh.NewInput().Title("Cell phone").Value(&data.CellPhone),
			huh.NewInput().Title("Work phone").Value(&data.WorkPhone),
			huh.NewInput().Title("Private phone").Value(&data.HomePhone),
		),

		huh.NewGroup(
			huh.NewInput().Title("Post office box").Value(&data.PostOfficeBox),
			huh.NewInput().Title("Street address").Value(&data.StreetAddress),
			huh.NewInput().Title("Extended street address (e.g. building, floor)").Value(&data.ExtendedAddress),
			huh.NewInput().Title("City").Value(&data.Locality),
			huh.NewInput().Title("Postal code").Value(&data.PostalCode),
			huh.NewInput().Title("Country").Value(&data.Country),
		),

		huh.NewGroup(
			huh.NewConfirm().
				Title("Are you ready?").
				Affirmative("Yes, print the result!").
				Negative("No, I´m not ready.").
				Value(&data.Ready),
		),
	).WithTheme(huh.ThemeBase16())
	return vCardForm
}

func typedVcardFieldValue(card *vcard.Card, fieldName, wantType string) string {
	if wantType == "" {
		return card.Value(fieldName)
	}

	typedFields := (*card)[fieldName]
	if typedFields == nil {
		return ""
	}

	for _, f := range typedFields {
		if f.Params.HasType(wantType) {
			return f.Value
		}
	}

	return ""
}

func setTypedVcardFieldValue(card *vcard.Card, fieldName, wantType, value string) {
	// we didn´t get a type
	if wantType == "" {
		card.SetValue(fieldName, value)
		return
	}

	// check if there is already a field of suitable type
	typedFields := (*card)[fieldName]
	for _, f := range typedFields {
		if slices.Contains(f.Params.Types(), wantType) {
			f.Value = value
			return
		}
	}

	// no field of that type was found, add one
	card.Add(fieldName, &vcard.Field{
		Value: value,
		Params: map[string][]string{
			"TYPE": {wantType},
		},
	})
}

func encodeVcard(card *vcard.Card) (string, error) {
	var buf bytes.Buffer
	enc := vcard.NewEncoder(&buf)
	if err := enc.Encode(*card); err != nil {
		return "", err
	}

	return buf.String(), nil
}

func decodeVcard(reader io.Reader) (vcard.Card, error) {
	dec := vcard.NewDecoder(reader)
	card, err := dec.Decode()
	if err != nil {
		return nil, err
	}
	return card, nil
}

func ensureNullSafety(card *vcard.Card) {
	if card.Name() == nil {
		name := vcard.Name{}
		card.SetName(&name)
	}
	if card.Address() == nil {
		address := vcard.Address{}
		card.SetAddress(&address)
	}
}

func cardInstance(settings *settings.Settings) (*vcard.Card, error) {

	if *settings.InputFilePath == "" {
		card := make(vcard.Card)
		card.SetValue(vcard.FieldVersion, *settings.VCardVersion)
		ensureNullSafety(&card)
		return &card, nil
	} else {
		// use the input file as vcard content
		cli.Println("Reading vCard file", cli.SprintValue(*settings.InputFilePath))
		cli.Println()
		file, err := os.Open(*settings.InputFilePath)
		if err != nil {
			return nil, err
		}
		defer file.Close()

		if card, err := decodeVcard(file); err != nil {
			return nil, err
		} else {
			ensureNullSafety(&card)
			return &card, nil
		}
	}
}

func PrepareVcard(settings *settings.Settings) (string, error) {

	card, err := cardInstance(settings)
	if err != nil {
		return "", err
	}

	if *settings.Silent {
		return encodeVcard(card)
	}

	formData := prepareFormData(card)

	for {
		form := prepareForm(formData)
		if err := form.Run(); err != nil {
			return "", err
		}
		if formData.Ready {
			break
		}
	}
	prepareVCardData(card, formData)

	return encodeVcard(card)

}
