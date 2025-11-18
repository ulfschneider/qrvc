package vcard

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"qrvc/internal/settings"
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
	field      string
	fieldType  string
	fieldLabel label
}

var inputProperties = []inputProperty{
	{field: vcard.FieldName},
	{field: vcard.FieldGender, fieldLabel: Gender},
	{field: vcard.FieldTitle, fieldLabel: Title},
	{field: vcard.FieldOrganization, fieldLabel: Organization},
	{field: vcard.FieldAddress},
	{field: vcard.FieldEmail, fieldLabel: EMail},
	{field: vcard.FieldURL, fieldLabel: URL},
	{field: vcard.FieldTelephone, fieldType: vcard.TypeCell, fieldLabel: CellPhone},
	{field: vcard.FieldTelephone, fieldType: vcard.TypeWork, fieldLabel: WorkPhone},
	{field: vcard.FieldTelephone, fieldType: vcard.TypeHome, fieldLabel: PrivatePhone},
}

var labels = map[label]string{
	GivenName:       "Given name (e.g. Harry)",
	FamilyName:      "Family name (e.g. Potter)",
	AdditionalName:  "Additional name (e.g. Fitzgerald)",
	HonorificPrefix: "Honorific prefix (e.g. Capt.)",
	HonorificSuffix: "Honorific suffix (e.g. Sr.)",
	Gender:          "Gender",
	Title:           "Job title",
	Organization:    "Organization or company",
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

func mustScanString(label string) string {
	templates := &promptui.PromptTemplates{
		Prompt:  "{{ . }}: ",
		Valid:   "{{ . }}: ",
		Invalid: "{{ . }}: ",
		Success: "{{ . }}: ",
	}
	prompt := promptui.Prompt{
		Label:     formatLabel(label),
		Default:   "",
		Templates: templates,
	}

	value, err := prompt.Run()
	if err != nil {
		panic(err)
	}

	return value
}

func mustScanGender() string {

	fmt.Println()

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

	genders := []string{"Male", "Female", "Other", "Unknown"}

	prompt := promptui.Select{
		Label:     "Select the gender or leave it unset",
		Items:     genders,
		CursorPos: 4,
		Templates: templates,
	}

	_, result, err := prompt.Run()
	if err != nil {
		panic(err)
	}
	fmt.Println("")

	switch result {
	case "Male":
		return "M"
	case "Female":
		return "F"
	case "Other":
		return "O"
	default:
		return ""
	}
}

func mustScanName(label string) *vcard.Name {
	name := vcard.Name{}

	if label != "" {
		fmt.Println(label)
	} else {
		fmt.Println("")
	}

	name.GivenName = mustScanString(labels[GivenName])
	name.FamilyName = mustScanString(labels[FamilyName])
	name.AdditionalName = mustScanString(labels[AdditionalName])
	name.HonorificPrefix = mustScanString(labels[HonorificPrefix])
	name.HonorificSuffix = mustScanString(labels[HonorificSuffix])

	return &name
}

func mustScanAddress(label string) *vcard.Address {
	address := vcard.Address{}

	if label != "" {
		fmt.Println(label)
	} else {
		fmt.Println("")
	}

	address.PostOfficeBox = mustScanString(labels[PostOfficeBox])
	address.StreetAddress = mustScanString(labels[StreetAddress])
	address.ExtendedAddress = mustScanString(labels[ExtendedAddress])
	address.Locality = mustScanString(labels[Locality])
	address.PostalCode = mustScanString(labels[PostalCode])
	address.Country = mustScanString(labels[Country])

	return &address
}

func mustScanProperty(card *vcard.Card, prop *inputProperty) *vcard.Card {

	switch prop.field {
	case vcard.FieldName:
		name := mustScanName("")
		card.AddName(name)
	case vcard.FieldAddress:
		address := mustScanAddress("")
		card.AddAddress(address)
	case vcard.FieldGender:
		gender := mustScanGender()
		card.SetValue(prop.field, gender)
	default:
		s := mustScanString(labels[prop.fieldLabel])
		if prop.fieldType == "" {
			card.SetValue(prop.field, s)
		} else {
			card.Add(prop.field, &vcard.Field{
				Value: s,
				Params: map[string][]string{
					"TYPE": {strings.ToUpper(prop.fieldType)},
				},
			})
		}
	}

	name := card.Name()
	if name.GivenName != "" || name.FamilyName != "" {
		card.SetValue("FN", strings.Join(strings.Fields(name.HonorificPrefix+" "+name.GivenName+" "+name.AdditionalName+" "+name.FamilyName+" "+name.HonorificSuffix), " "))
	} else {
		card.SetValue("FN", card.Value("ORG"))
	}

	return card
}

func mustEncode(card *vcard.Card) string {
	var buf bytes.Buffer
	enc := vcard.NewEncoder(&buf)
	if err := enc.Encode(*card); err != nil {
		panic(err)
	}

	return buf.String()
}

func MustPrepareVcard(args *settings.Settings) string {

	if *args.InputFilePath == "" {
		// use the interactive mode to ask for vcard properties
		// and produce vcard content out of it
		fmt.Println("\nReading input to create the vcard")
		fmt.Println("Press ENTER to proceed or CTRL-C to cancel")
		card := make(vcard.Card)
		card.SetValue(vcard.FieldVersion, *args.VcardVersion)
		for _, prop := range inputProperties {
			mustScanProperty(&card, &prop)
		}
		return mustEncode(&card)
	} else {
		// use the input file as vcard content
		fmt.Println("\nReading vcard file")
		file, err := os.Open(*args.InputFilePath)
		if err != nil {
			panic(err)
		}
		defer file.Close()

		b, err := io.ReadAll(file)
		if err != nil {
			panic(err)
		}
		return string(b)
	}
}
