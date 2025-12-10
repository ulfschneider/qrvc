package qrcard

import (
	"bytes"
	"fmt"
	"image/png"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"slices"
	"strings"

	"github.com/skip2/go-qrcode"
	"github.com/ulfschneider/qrvc/internal/cli"
	"github.com/ulfschneider/qrvc/internal/settings"

	"github.com/spf13/afero"

	"github.com/charmbracelet/huh"
	"github.com/emersion/go-vcard"
	"github.com/pkg/errors"
)

type qrcardFormData struct {
	name         vcard.Name
	gender       vcard.Sex
	title        string
	organization string
	department   string
	address      vcard.Address
	email        string
	url          string
	cellPhone    string
	workPhone    string
	homePhone    string
	ready        bool
}

func maybeGet(l []string, i int) string {
	if i < len(l) {
		return l[i]
	}
	return ""
}

func transferVCardIntoFormData(card *vcard.Card) *qrcardFormData {

	sex, _ := card.Gender()

	orgSplit := strings.SplitN(card.Value(vcard.FieldOrganization), ";", 2)
	organization := orgSplit[0]
	department := maybeGet(orgSplit, 1)

	data := qrcardFormData{
		name:         *card.Name(),
		gender:       sex,
		title:        card.Value(vcard.FieldTitle),
		organization: organization,
		department:   department,
		address:      *card.Address(),
		email:        card.Value(vcard.FieldEmail),
		url:          card.Value(vcard.FieldURL),
		cellPhone:    typedVcardFieldValue(card, vcard.FieldTelephone, vcard.TypeCell),
		workPhone:    typedVcardFieldValue(card, vcard.FieldTelephone, vcard.TypeWork),
		homePhone:    typedVcardFieldValue(card, vcard.FieldTelephone, vcard.TypeHome),
		ready:        true,
	}

	return &data
}

func transferFormDataIntoVCard(card *vcard.Card, formData *qrcardFormData) {
	card.SetName(&formData.name)
	card.SetGender(vcard.Sex(formData.gender), "")
	card.SetValue(vcard.FieldTitle, formData.title)
	card.SetValue(vcard.FieldOrganization, formData.organization+";"+formData.department)
	card.SetAddress(&formData.address)
	card.SetValue(vcard.FieldEmail, formData.email)
	card.SetValue(vcard.FieldURL, formData.url)
	setTypedVcardFieldValue(card, vcard.FieldTelephone, vcard.TypeCell, formData.cellPhone)
	setTypedVcardFieldValue(card, vcard.FieldTelephone, vcard.TypeWork, formData.workPhone)
	setTypedVcardFieldValue(card, vcard.FieldTelephone, vcard.TypeHome, formData.homePhone)
}

func prepareForm(formData *qrcardFormData) *huh.Form {

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

func openVcard(inputFilePath *string, filesystem afero.Fs) (fs.File, error) {
	file, err := filesystem.Open(*inputFilePath)
	if err != nil {
		if filepath.Ext(*inputFilePath) == "" {
			//try .vcf
			alternateFilePath := *inputFilePath + ".vcf"
			file, err = filesystem.Open(alternateFilePath)
			if err == nil {
				//when the .vcf was possible to read, this will be the new inputFilePath
				*inputFilePath = alternateFilePath
			}
		}
	}

	if err != nil {
		return nil, errors.Wrap(err, "Error when trying to open file "+cli.SprintValue(*inputFilePath))
	}
	return file, err
}

func makeVCardInstance(inputFilePath *string, vcardVersion string, filesystem afero.Fs) (*vcard.Card, error) {

	if *inputFilePath == "" {
		//no path to a vcard file, create a new card
		card := make(vcard.Card)
		card.SetValue(vcard.FieldVersion, vcardVersion)
		ensureNilSafety(&card)
		return &card, nil
	} else {
		// we have a path to a vcard file, try to read it
		file, err := openVcard(inputFilePath, filesystem)
		if err != nil {
			return nil, err
		}

		defer file.Close()

		fmt.Println("Reading vCard file", cli.SprintValue(*inputFilePath))
		cli.Println()

		if card, err := decodeVcard(file); err != nil {
			return nil, err
		} else {
			ensureNilSafety(&card)
			return &card, nil
		}
	}
}

// PrepareVCard reads a vCard file when `inputFilePath` is not empty and uses the data as a starting point for further user input,
// When `inputFilePath` is empty, an empty the user can provide vCard data to create a new vCard.
// The user can edit vCard content in a command line form and the resulting vCard content is returned from the function in a `string`.
// `vcardVersion` defines what vCard version is assigned to the resulting vCard content.
// When `silent` is true, the `inputFilePath` must not be empty and the file described by the path must exist.
// In that case the user will not be asked for further input, instead the read file content is returned as a string.
// `filesystem` is used to access the vCard content file.
func PrepareVCard(inputFilePath *string, vcardVersion string, silent bool, filesystem afero.Fs) (string, error) {
	card, err := makeVCardInstance(inputFilePath, vcardVersion, filesystem)
	if err != nil {
		return "", err
	}

	if silent {
		return encodeVcard(card)
	}

	formData := transferVCardIntoFormData(card)

	for {
		form := prepareForm(formData)
		if err := form.Run(); err != nil {
			return "", err
		}
		if formData.ready {
			break
		}
	}

	transferFormDataIntoVCard(card, formData)

	return encodeVcard(card)

}

// WriteResults stores the `vcardContent` in a vCard file and creates a QR Code that is as well stored.
// The `filesystem` is used to create the files.
// `settings` provide the namess for the files and style information for the QR Code.
func WriteResults(vcardContent string, settings *settings.Settings, filesystem afero.Fs) error {

	if file, err := filesystem.Create(*settings.VCardOutputFilePath); err != nil {
		return err
	} else {
		defer file.Close()
		if _, err := file.WriteString(vcardContent); err != nil {
			return err
		} else {
			fmt.Println("The vCard has been written to", cli.SprintValue(*settings.VCardOutputFilePath))
		}
	}

	q, err := qrcode.New(vcardContent, qrcode.Low)
	if err != nil {
		return err
	}

	q.DisableBorder = !*settings.Border
	q.ForegroundColor = *settings.ForegroundColor
	q.BackgroundColor = *settings.BackgroundColor
	q.Level = qrcode.Low

	img := q.Image(*settings.Size)

	if file, err := os.Create(*settings.QRCodeOutputFilePath); err != nil {
		return err
	} else {
		defer file.Close()
		if err := png.Encode(file, img); err != nil {
			return err
		} else {
			fmt.Println("The QR code has been written to", cli.SprintValue(*settings.QRCodeOutputFilePath))
		}
	}
	return nil
}
