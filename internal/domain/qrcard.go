package qrcard

import (
	"slices"

	"github.com/emersion/go-vcard"
)

func TypedVcardFieldValue(card vcard.Card, fieldName, wantType string) string {
	if wantType == "" {
		return card.Value(fieldName)
	}

	typedFields := card[fieldName]
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

func SetTypedVcardFieldValue(card vcard.Card, fieldName, wantType, value string) {
	// we didn´t get a type
	if wantType == "" {
		card.SetValue(fieldName, value)
		return
	}

	// check if there is already a field of suitable type
	typedFields := card[fieldName]
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
