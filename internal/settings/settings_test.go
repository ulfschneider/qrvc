package settings_test

import (
	"testing"

	"github.com/mazznoer/csscolorparser"
	"github.com/stretchr/testify/assert"
	"github.com/ulfschneider/qrvc/internal/settings"
)

func TestDefaultSettings(t *testing.T) {
	args, err := settings.PrepareSettings()
	assert.NoError(t, err)
	assert.NotEmpty(t, args)

	bgColor, err := csscolorparser.Parse("transparent")
	assert.Equal(t, bgColor, *args.BackgroundColor)

	fgColor, err := csscolorparser.Parse("black")
	assert.Equal(t, fgColor, *args.ForegroundColor)

	assert.False(t, *args.Bom)

	assert.False(t, *args.Silent)

	assert.False(t, *args.Border)

	assert.Equal(t, 400, *args.Size)

	assert.Equal(t, "3.0", *args.VCardVersion)

	assert.Equal(t, "", *args.InputFilePath)

	assert.Equal(t, "vcard.vcf", *args.VCardOutputFilePath)

	assert.Equal(t, "vcard.png", *args.QRCodeOutputFilePath)

}
