package filerepo_test

import (
	"image"
	"image/draw"
	"testing"

	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"github.com/ulfschneider/qrvc/internal/adapters/cliconfig"
	"github.com/ulfschneider/qrvc/internal/adapters/filerepo"
	"github.com/ulfschneider/qrvc/internal/adapters/qrcodec"
	"github.com/ulfschneider/qrvc/internal/adapters/testutil"
	"github.com/ulfschneider/qrvc/internal/adapters/vcardcodec"

	"github.com/ulfschneider/qrvc/internal/application/services"
)

func createTestSettings() cliconfig.CLIFileSettings {

	versionService := services.NewVersionService(testutil.CreateVersionProvider())
	settingsProvider := cliconfig.NewSettingsProvider(versionService)
	settings, _ := settingsProvider.Load()

	return settings
}

func createTestRepo(fs afero.Fs, settings cliconfig.CLIFileSettings) filerepo.Repository {
	cardCodec := vcardcodec.NewCodec()
	qrCodec := qrcodec.NewCodec()
	repo := filerepo.NewRepo(fs, &cardCodec, &qrCodec, settings.Files, settings.App)
	return repo
}

func TestMakeVCardInstanceFromNonExistingFile(t *testing.T) {
	filesystem := afero.NewMemMapFs()
	settings := createTestSettings()
	settings.Files.ReadVCardPath = "vcard"
	repo := createTestRepo(filesystem, settings)

	//vcard file does not exist, there is nothing to read, must return an error
	vcard, err := repo.ReadOrCreateVCard()
	assert.Error(t, err)
	assert.Nil(t, vcard)

	//file does still not exist
	_, err = filesystem.Stat(settings.Files.ReadVCardPath)
	assert.Error(t, err)
}

func TestMakeVCardInstanceFromExistingFile(t *testing.T) {
	filesystem := afero.NewMemMapFs()
	settings := createTestSettings()
	filePath := "vcard.vcf"
	settings.Files.ReadVCardPath = filePath
	repo := createTestRepo(filesystem, settings)

	//create file
	f, err := filesystem.Create(filePath)
	assert.NoError(t, err)
	assert.NotEmpty(t, f)

	expectedCard := testutil.CreateCard()
	expectedContent := testutil.EncodeCard(expectedCard)

	f.Write(expectedContent)

	//vcard.vcf file does exist
	actualCard, err := repo.ReadOrCreateVCard()
	assert.NoError(t, err)
	assert.NotEmpty(t, actualCard)
	assert.Equal(t, expectedCard, actualCard)

	//vcard (.vcf is added automatically) can as well be used to access the file
	settings.Files.ReadVCardPath = "vcard"
	actualCard, err = repo.ReadOrCreateVCard()
	assert.NoError(t, err)
	assert.NotEmpty(t, actualCard)
	assert.Equal(t, expectedCard, actualCard)

	//file does still exist
	info, err := filesystem.Stat(filePath)
	assert.NoError(t, err)
	assert.NotEmpty(t, info)
}

func TestWriteVCard(t *testing.T) {
	filesystem := afero.NewMemMapFs()
	settings := createTestSettings()
	settings.Files.ReadVCardPath = "vcard.vcf"
	settings.Files.WriteVCardPath = "vcard.vcf"
	settings.Files.WriteQRCodePath = "vcard.png"
	repo := createTestRepo(filesystem, settings)

	expectedCard := testutil.CreateCard()
	err := repo.WriteVCard(expectedCard)
	assert.NoError(t, err)

	actualCard, err := repo.ReadOrCreateVCard()
	assert.NoError(t, err)
	assert.NotEmpty(t, actualCard)
	assert.Equal(t, expectedCard, actualCard)

	err = repo.WriteQRCode(expectedCard)
	assert.NoError(t, err)
	expectedCode := testutil.CreateQRCode(expectedCard, settings.App.QRSettings)
	file, err := filesystem.Open("vcard.png")
	actualCode, format, err := image.Decode(file)
	assert.NoError(t, err)
	assert.Equal(t, "png", format)
	assert.Equal(t, toRGBA(expectedCode), toRGBA(actualCode))
}

func toRGBA(img image.Image) *image.RGBA {
	b := img.Bounds()
	rgba := image.NewRGBA(b)
	draw.Draw(rgba, b, img, b.Min, draw.Src)
	return rgba
}
