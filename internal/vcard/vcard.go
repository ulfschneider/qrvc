package vcard

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"qrvc/internal/settings"
	"slices"
	"strings"
	"text/template"

	"github.com/emersion/go-vcard"
	"github.com/manifoldco/promptui"
)

type label int

const (
	Undefined label = iota
	GivenName
	FamilyName
	AdditionalName
	HonorificPrefix
	HonorificSuffix
	Gender
	Title
	Organization
	Department
	PostOfficeBox
	StreetAddress
	ExtendedAddress
	Locality
	PostalCode
	Country
	EMail
	URL
	CellPhone
	WorkPhone
	PrivatePhone
)

type inputProperty struct {
	fieldName  string
	fieldType  string
	fieldLabel label
}

var inputProperties = []inputProperty{
	{fieldName: vcard.FieldName},
	{fieldName: vcard.FieldGender},
	{fieldName: vcard.FieldOrganization},
	{fieldName: vcard.FieldAddress},
	{fieldName: vcard.FieldEmail, fieldLabel: EMail},
	{fieldName: vcard.FieldURL, fieldLabel: URL},
	{fieldName: vcard.FieldTelephone, fieldType: vcard.TypeCell, fieldLabel: CellPhone},
	{fieldName: vcard.FieldTelephone, fieldType: vcard.TypeWork, fieldLabel: WorkPhone},
	{fieldName: vcard.FieldTelephone, fieldType: vcard.TypeHome, fieldLabel: PrivatePhone},
}

var labels = map[label]string{
	GivenName:       "Given name (e.g. Harry)",
	FamilyName:      "Family name (e.g. Potter)",
	AdditionalName:  "Additional name (e.g. James)",
	HonorificPrefix: "Honorific prefix (e.g. Capt.)",
	HonorificSuffix: "Honorific suffix (e.g. Sr.)",
	Gender:          "Gender",
	Title:           "Job title",
	Organization:    "Organization or company",
	Department:      "Department",
	PostOfficeBox:   "Post office box",
	StreetAddress:   "Street address",
	ExtendedAddress: "Extended address (e.g. building, floor)",
	Locality:        "City",
	PostalCode:      "Postal code",
	Country:         "Country",
	EMail:           "E-mail address",
	URL:             "Web address (URL)",
	CellPhone:       "Cell phone",
	WorkPhone:       "Work phone",
	PrivatePhone:    "Private phone",
}

var labelWidth = maxLabelWidth()

func maxLabelWidth() int {
	var width int
	for _, label := range labels {
		if len(label) > width {
			width = len(label)
		}
	}

	return width
}

func padRight(label string, width int, pad rune) string {
	padding := make([]rune, width-len(label))
	for i := range padding {
		padding[i] = pad
	}
	return label + string(padding) + " "
}

func formatLabel(label string) string {
	return fmt.Sprintf("%s", padRight(label, labelWidth, '.'))
}

func scanString(label, defaultValue string) (string, error) {

	templates := &promptui.PromptTemplates{
		Prompt:  "{{ . }}: ",
		Valid:   "{{ . }}: ",
		Invalid: "{{ . }}: ",
		Success: "{{ . }}: ",
	}
	prompt := promptui.Prompt{
		Label:     formatLabel(label),
		Default:   defaultValue,
		Templates: templates,
	}

	value, err := prompt.Run()
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(value), nil
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
		if f.Params.HasType(strings.ToUpper(wantType)) {
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
			fmt.Println(fieldName, wantType, "has type")
			f.Value = value
			return
		}
	}

	// no field of that type was found, add one
	card.Add(fieldName, &vcard.Field{
		Value: value,
		Params: map[string][]string{
			"TYPE": {strings.ToUpper(wantType)},
		},
	})

}

func scanGender(card *vcard.Card) error {
	funcMap := template.FuncMap{
		"formatSelected": func(option string) string {
			if option == "Unknown" {
				return ""
			}
			return option
		},
		"faint": func(option string) string {
			return "\033[2m" + option + "\033[0m"
		},
	}

	templates := &promptui.SelectTemplates{
		Label:    "{{ . }}:",
		Active:   "> {{ . }}",
		Inactive: "  {{ . }}",
		Selected: formatLabel(labels[Gender]) + ": {{ formatSelected . }}",
		FuncMap:  funcMap,
	}

	if templates.FuncMap == nil {
		templates.FuncMap = make(template.FuncMap)
	}

	genders := []string{"Male", "Female", "Other", "Unspecified"}
	defaultValue, _ := card.Gender()

	var cursorPos int
	switch defaultValue {
	case vcard.SexMale:
		cursorPos = 0
	case vcard.SexFemale:
		cursorPos = 1
	case vcard.SexOther:
		cursorPos = 2
	default:
		cursorPos = 3
	}

	prompt := promptui.Select{
		Label:     "Select the gender or leave it unset",
		Items:     genders,
		CursorPos: cursorPos,
		Templates: templates,
	}

	_, result, err := prompt.Run()
	if err != nil {
		return err
	}

	switch result {
	case "Male":
		card.SetGender(vcard.SexMale, "")

	case "Female":
		card.SetGender(vcard.SexFemale, "")
	case "Other":
		card.SetGender(vcard.SexOther, "")
	default:
		card.SetGender(vcard.SexUnspecified, "")
	}
	return nil
}

func scanName(card *vcard.Card) error {
	name := card.Name()
	if name == nil {
		name = &vcard.Name{}
	}

	if givenName, err := scanString(labels[GivenName], name.GivenName); err != nil {
		return err
	} else {
		name.GivenName = givenName
	}

	if familyName, err := scanString(labels[FamilyName], name.FamilyName); err != nil {
		return err
	} else {
		name.FamilyName = familyName
	}

	if additionalName, err := scanString(labels[AdditionalName], name.AdditionalName); err != nil {
		return err
	} else {
		name.AdditionalName = additionalName
	}

	if honorifixPrefix, err := scanString(labels[HonorificPrefix], name.HonorificPrefix); err != nil {
		return err
	} else {
		name.HonorificPrefix = honorifixPrefix
	}

	if honorificSuffix, err := scanString(labels[HonorificSuffix], name.HonorificSuffix); err != nil {
		return err
	} else {
		name.HonorificSuffix = honorificSuffix
	}

	if len(card.Names()) == 0 {
		card.AddName(name)
	} else {
		card.Names()[0] = name
	}

	return nil
}

func scanOrg(card *vcard.Card) error {

	title, err := scanString(labels[Title], card.Value(vcard.FieldTitle))
	if err != nil {
		return err
	}

	card.SetValue(vcard.FieldTitle, title)

	orgSplit := strings.Split(card.Value(vcard.FieldOrganization), ";")
	orgDefault := ""
	departmentDefault := ""
	if len(orgSplit) > 0 {
		orgDefault = orgSplit[0]
	}
	if len(orgSplit) > 1 {
		departmentDefault = orgSplit[1]
	}

	org, err := scanString(labels[Organization], orgDefault)
	if err != nil {
		return err
	}

	department, err := scanString(labels[Department], departmentDefault)
	if err != nil {
		return err
	}

	card.SetValue(vcard.FieldOrganization, org+";"+department)

	return nil
}

func scanAddress(card *vcard.Card) error {
	address := card.Address()
	if address == nil {
		address = &vcard.Address{}
	}

	if postOfficeBox, err := scanString(labels[PostOfficeBox], address.PostOfficeBox); err != nil {
		return err
	} else {
		address.PostOfficeBox = postOfficeBox
	}

	if streetAddress, err := scanString(labels[StreetAddress], address.StreetAddress); err != nil {
		return err
	} else {
		address.StreetAddress = streetAddress
	}

	if extendedAddress, err := scanString(labels[ExtendedAddress], address.ExtendedAddress); err != nil {
		return err
	} else {
		address.ExtendedAddress = extendedAddress
	}

	if locality, err := scanString(labels[Locality], address.Locality); err != nil {
		return err
	} else {
		address.Locality = locality
	}

	if postalCode, err := scanString(labels[PostalCode], address.PostalCode); err != nil {
		return err
	} else {
		address.PostalCode = postalCode
	}

	if country, err := scanString(labels[Country], address.Country); err != nil {
		return err
	} else {
		address.Country = country
	}

	if len(card.Addresses()) == 0 {
		card.AddAddress(address)
	} else {
		card.Addresses()[0] = address
	}

	return nil
}

func scanVcardProperty(card *vcard.Card, prop *inputProperty) error {

	switch prop.fieldName {
	case vcard.FieldName:
		fmt.Println()
		if err := scanName(card); err != nil {
			return err
		}
	case vcard.FieldAddress:
		fmt.Println()
		if err := scanAddress(card); err != nil {
			return err
		}
		fmt.Println()
	case vcard.FieldGender:
		fmt.Println()
		if err := scanGender(card); err != nil {
			return err
		}
	case vcard.FieldOrganization:
		fmt.Println()
		if err := scanOrg(card); err != nil {
			return err
		}
	default:
		if s, err := scanString(labels[prop.fieldLabel], typedVcardFieldValue(card, prop.fieldName, prop.fieldType)); err != nil {
			return err
		} else {
			setTypedVcardFieldValue(card, prop.fieldName, prop.fieldType, s)
		}
	}

	name := card.Name()
	if name.GivenName != "" || name.FamilyName != "" {
		card.SetValue("FN", strings.Join(strings.Fields(name.HonorificPrefix+" "+name.GivenName+" "+name.AdditionalName+" "+name.FamilyName+" "+name.HonorificSuffix), " "))
	} else {
		card.SetValue("FN", card.Value("ORG"))
	}

	return nil
}

func encodeVcard(card *vcard.Card) (string, error) {
	var buf bytes.Buffer
	enc := vcard.NewEncoder(&buf)
	if err := enc.Encode(*card); err != nil {
		return "", err
	}

	return buf.String(), nil
}

func askRepeatReadingInput() (bool, error) {
	fmt.Println()

	label := "Want to change something? Repeat?"

	templates := &promptui.SelectTemplates{
		Label:    "{{ . }}:",
		Active:   "> {{ . }}",
		Inactive: "  {{ . }}",
		Selected: label + " {{ . }}",
	}

	prompt := promptui.Select{
		Label:     label,
		Items:     []string{"Yes", "No"},
		Templates: templates,
		CursorPos: 1,
	}

	_, result, err := prompt.Run()
	if err != nil {
		return false, err
	}

	if result == "Yes" {
		return true, nil
	} else {
		return false, nil
	}
}

func PrepareVcard(args *settings.Settings) (string, error) {

	if *args.InputFilePath == "" {
		// use the interactive mode to ask for vcard properties
		// and produce vcard content out of it

		card := make(vcard.Card)

		for {
			fmt.Println("\nReading input to create the vcard")
			fmt.Println("Press ENTER to proceed or CTRL-C to cancel")

			card.SetValue(vcard.FieldVersion, *args.VcardVersion)
			for _, prop := range inputProperties {
				if err := scanVcardProperty(&card, &prop); err != nil {
					return "", err
				}

			}
			if readInput, err := askRepeatReadingInput(); err != nil {
				return "", err
			} else if !readInput {
				break
			}

		}
		return encodeVcard(&card)
	} else {
		// use the input file as vcard content
		fmt.Println("\nReading vcard file")
		file, err := os.Open(*args.InputFilePath)
		if err != nil {
			return "", err
		}
		defer file.Close()

		b, err := io.ReadAll(file)
		if err != nil {
			return "", err
		}
		return string(b), nil
	}
}
