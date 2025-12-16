package qrcard_test

import (
	"testing"

	"github.com/emersion/go-vcard"
	"github.com/stretchr/testify/assert"
	qrcard "github.com/ulfschneider/qrvc/internal/domain"
)

func TestTypedValues(t *testing.T) {
	card := vcard.Card{}
	expectedCellPhone := "cell phone"
	expectedWorkPhone := "work phone"
	expectedHomePhone := "home phone"
	qrcard.SetTypedVcardFieldValue(card, vcard.FieldTelephone, vcard.TypeCell, expectedCellPhone)
	qrcard.SetTypedVcardFieldValue(card, vcard.FieldTelephone, vcard.TypeWork, expectedWorkPhone)
	qrcard.SetTypedVcardFieldValue(card, vcard.FieldTelephone, vcard.TypeHome, expectedHomePhone)

	assert.Equal(t, expectedCellPhone, qrcard.TypedVcardFieldValue(card, vcard.FieldTelephone, vcard.TypeCell))
	assert.Equal(t, expectedWorkPhone, qrcard.TypedVcardFieldValue(card, vcard.FieldTelephone, vcard.TypeWork))
	assert.Equal(t, expectedHomePhone, qrcard.TypedVcardFieldValue(card, vcard.FieldTelephone, vcard.TypeHome))
}
