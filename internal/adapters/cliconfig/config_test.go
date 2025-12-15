package cliconfig_test

import (
	"testing"

	"github.com/mazznoer/csscolorparser"
	"github.com/skip2/go-qrcode"
	"github.com/stretchr/testify/assert"
	"github.com/ulfschneider/qrvc/internal/adapters/cliconfig"
	"github.com/ulfschneider/qrvc/internal/application/services"
)

type versionProvider struct {
}

func (vp versionProvider) Version() (string, error) {
	return "TEST VERSION", nil
}

func TestDefaultSettings(t *testing.T) {

	versionService := services.NewVersionService(versionProvider{})
	settingsProvider := cliconfig.NewCLIFileSettingsProvider(versionService)

	settings, err := settingsProvider.Load()
	assert.NoError(t, err)
	assert.NotEmpty(t, settings)

	//test adapter settings

	assert.Equal(t, "", settings.Files.ReadVCardPath)

	assert.Equal(t, "vcard.vcf", settings.Files.WriteVCardPath)

	assert.Equal(t, "vcard.png", settings.Files.WriteQRCodePath)

	assert.False(t, settings.CLI.Bom)

	//test application settings

	assert.Equal(t, "3.0", settings.App.VCardVersion)

	assert.False(t, settings.App.Silent)

	bgColor, err := csscolorparser.Parse("transparent")
	assert.NoError(t, err)
	assert.Equal(t, bgColor, settings.App.QRSettings.BackgroundColor)

	fgColor, err := csscolorparser.Parse("black")
	assert.NoError(t, err)
	assert.Equal(t, fgColor, settings.App.QRSettings.ForegroundColor)

	assert.False(t, settings.App.QRSettings.Border)

	assert.Equal(t, 400, settings.App.QRSettings.Size)

	assert.Equal(t, qrcode.Low, settings.App.QRSettings.RecoveryLevel)

}
