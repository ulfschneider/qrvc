package vcard

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"slices"
	"strings"

	"github.com/ulfschneider/qrvc/internal/cli"
	"github.com/ulfschneider/qrvc/internal/settings"

	"github.com/charmbracelet/huh"
	"github.com/emersion/go-vcard"
	"github.com/pkg/errors"
)

type VCardFormData struct {
	Name         vcard.Name
	Gender       vcard.Sex
	Title        string
	Organization string
	Department   string
	Address      vcard.Address
	Email        string
	Url          string
	CellPhone    string
	WorkPhone    string
	HomePhone    string
	Ready        bool
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
		Name:         *card.Name(),
		Gender:       sex,
		Title:        card.Value(vcard.FieldTitle),
		Organization: organization,
		Department:   department,
		Address:      *card.Address(),
		Email:        card.Value(vcard.FieldEmail),
		Url:          card.Value(vcard.FieldURL),
		CellPhone:    typedVcardFieldValue(card, vcard.FieldTelephone, vcard.TypeCell),
		WorkPhone:    typedVcardFieldValue(card, vcard.FieldTelephone, vcard.TypeWork),
		HomePhone:    typedVcardFieldValue(card, vcard.FieldTelephone, vcard.TypeHome),
		Ready:        true,
	}

	return &data
}

func transferFormData(card *vcard.Card, formData *VCardFormData) {
	card.SetName(&formData.Name)
	card.SetGender(vcard.Sex(formData.Gender), "")
	card.SetValue(vcard.FieldTitle, formData.Title)
	card.SetValue(vcard.FieldOrganization, formData.Organization+";"+formData.Department)
	card.SetAddress(&formData.Address)
	card.SetValue(vcard.FieldEmail, formData.Email)
	card.SetValue(vcard.FieldURL, formData.Url)
	setTypedVcardFieldValue(card, vcard.FieldTelephone, vcard.TypeCell, formData.CellPhone)
	setTypedVcardFieldValue(card, vcard.FieldTelephone, vcard.TypeWork, formData.WorkPhone)
	setTypedVcardFieldValue(card, vcard.FieldTelephone, vcard.TypeHome, formData.HomePhone)
}

func prepareForm(formData *VCardFormData) *huh.Form {

	vCardForm := huh.NewForm(
		huh.NewGroup(
			huh.NewInput().Title("Given (first) name").Value(&formData.Name.GivenName),
			huh.NewInput().Title("Additional (middle) name").Value(&formData.Name.AdditionalName),
			huh.NewInput().Title("Family name").Value(&formData.Name.FamilyName),
			huh.NewInput().Title("Honorific prefix (e.g. Capt.)").Value(&formData.Name.HonorificPrefix),
			huh.NewInput().Title("Honorific suffix (e.g. Sr.)").Value(&formData.Name.HonorificSuffix),
		),
		huh.NewGroup(
			huh.NewSelect[vcard.Sex]().Title("Gender").Options(
				huh.NewOption("Male", vcard.SexMale).Selected(vcard.SexMale == formData.Gender),
				huh.NewOption("Female", vcard.SexFemale).Selected(vcard.SexFemale == formData.Gender),
				huh.NewOption("Other", vcard.SexOther).Selected(vcard.SexOther == formData.Gender),
				huh.NewOption("Unspecified", vcard.SexUnspecified).Selected(formData.Gender != vcard.SexMale && formData.Gender != vcard.SexFemale && formData.Gender != vcard.SexUnspecified),
			).Value(&formData.Gender),
		),
		huh.NewGroup(
			huh.NewInput().Title("Job title").Value(&formData.Title),
			huh.NewInput().Title("Organization or company").Value(&formData.Organization),
			huh.NewInput().Title("Department").Value(&formData.Department),
		),

		huh.NewGroup(
			huh.NewInput().Title("Mail").Value(&formData.Email),
			huh.NewInput().Title("Web address").Value(&formData.Url),
			huh.NewInput().Title("Cell phone").Value(&formData.CellPhone),
			huh.NewInput().Title("Work phone").Value(&formData.WorkPhone),
			huh.NewInput().Title("Private phone").Value(&formData.HomePhone),
		),

		huh.NewGroup(
			huh.NewInput().Title("Post office box").Value(&formData.Address.PostOfficeBox),
			huh.NewInput().Title("Street address").Value(&formData.Address.StreetAddress),
			huh.NewInput().Title("Extended street address (e.g. building, floor)").Value(&formData.Address.ExtendedAddress),
			huh.NewInput().Title("City").Value(&formData.Address.Locality),
			huh.NewInput().Title("Postal code").Value(&formData.Address.PostalCode),
			huh.NewInput().Title("Country").Value(&formData.Address.Country),
		),

		huh.NewGroup(
			huh.NewConfirm().
				Title("Are you ready?").
				Affirmative("Yes, print the result!").
				Negative("No, I´m not ready.").
				Value(&formData.Ready),
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

func ensureNilSafety(card *vcard.Card) {
	if card.Name() == nil {
		name := vcard.Name{}
		card.SetName(&name)
	}
	if card.Address() == nil {
		address := vcard.Address{}
		card.SetAddress(&address)
	}
}

func openVcard(settings *settings.Settings) (*os.File, error) {
	file, err := os.Open(*settings.InputFilePath)
	if err != nil {
		if filepath.Ext(*settings.InputFilePath) == "" {
			//try .vcf
			alternateFilePath := *settings.InputFilePath + ".vcf"
			file, err = os.Open(alternateFilePath)
			if err == nil {
				//when the .vcf was possible to read, this will be the new InputFilePath
				*settings.InputFilePath = alternateFilePath
			}
		}
	}

	if err != nil {
		return nil, errors.Wrap(err, "Error when trying to open file "+cli.SprintValue(*settings.InputFilePath))
	}
	return file, err
}

func cardInstance(settings *settings.Settings) (*vcard.Card, error) {

	if *settings.InputFilePath == "" {
		//no path to a vcard file, create a new card
		card := make(vcard.Card)
		card.SetValue(vcard.FieldVersion, *settings.VCardVersion)
		ensureNilSafety(&card)
		return &card, nil
	} else {
		// we have a path to a vcard file, try to read it
		file, err := openVcard(settings)
		if err != nil {
			return nil, err
		}

		defer file.Close()

		fmt.Println("Reading vCard file", cli.SprintValue(*settings.InputFilePath))
		cli.Println()

		if card, err := decodeVcard(file); err != nil {
			return nil, err
		} else {
			ensureNilSafety(&card)
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

	transferFormData(card, formData)

	return encodeVcard(card)

}
