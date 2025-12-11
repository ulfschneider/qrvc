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

	assert.Equal(t, "", *args.InputFilePath)

	assert.Equal(t, "3.0", *args.VCardVersion)

	assert.False(t, *args.Bom)

	assert.False(t, *args.Silent)

	bgColor, err := csscolorparser.Parse("transparent")
	assert.NoError(t, err)
	assert.Equal(t, bgColor, *args.OutputSettings.BackgroundColor)

	fgColor, err := csscolorparser.Parse("black")
	assert.NoError(t, err)
	assert.Equal(t, fgColor, *args.OutputSettings.ForegroundColor)

	assert.False(t, *args.OutputSettings.Border)

	assert.Equal(t, 400, *args.OutputSettings.Size)

	assert.Equal(t, "vcard.vcf", *args.OutputSettings.VCardFilePath)

	assert.Equal(t, "vcard.png", *args.OutputSettings.QRCodeFilePath)

}
