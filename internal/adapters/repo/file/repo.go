package repofile

import (
	"image/png"
	"path/filepath"

	"github.com/emersion/go-vcard"
	"github.com/pkg/errors"
	"github.com/spf13/afero"

	configcli "github.com/ulfschneider/qrvc/internal/adapters/config/cli"
	notifiercli "github.com/ulfschneider/qrvc/internal/adapters/notifier/cli"
	"github.com/ulfschneider/qrvc/internal/application/config"
	"github.com/ulfschneider/qrvc/internal/application/ports"
)

func NewRepo(
	fileSystem afero.Fs,
	cardCodec ports.VCardCodec,
	qrCodec ports.QRCodec,
	fileSettings configcli.FileSettings,
	appSettings config.Settings,
) Repository {

	return Repository{
		fileSystem:   fileSystem,
		cardCodec:    cardCodec,
		qrCodec:      qrCodec,
		userNotifier: notifiercli.NewUserNotifier(),
		fileSettings: fileSettings,
		appSettings:  appSettings,
	}
}

type Repository struct {
	fileSystem   afero.Fs
	cardCodec    ports.VCardCodec
	qrCodec      ports.QRCodec
	userNotifier notifiercli.UserNotifier
	fileSettings configcli.FileSettings
	appSettings  config.Settings
}

func (fr *Repository) ReadOrCreateVCard() (vcard.Card, error) {

	if fr.fileSettings.ReadVCardPath == "" && fr.userNotifier.Silent() == false {
		//no path to a vcard file, create a new card
		card := make(vcard.Card)
		card.SetValue(vcard.FieldVersion, fr.appSettings.VCardVersion)
		ensureNilSafety(card)
		return card, nil
	} else if fr.fileSettings.ReadVCardPath == "" && fr.userNotifier.Silent() == true {
		//no path to a vcard file, but tool runs in silent mode
		return nil, errors.New("Missing input file path")
	} else {

		fr.fitReadVCardPath()

		fr.userNotifier.Section()
		fr.userNotifier.Notifyf("Reading vCard file %s", fr.fileSettings.ReadVCardPath)
		fr.userNotifier.Section()

		data, err := afero.ReadFile(fr.fileSystem, fr.fileSettings.ReadVCardPath)
		if err != nil {
			return nil, err
		}

		if card, err := fr.cardCodec.Decode(data); err != nil {
			return nil, err
		} else {
			ensureNilSafety(card)
			return card, nil
		}
	}
}

func (fr *Repository) fitReadVCardPath() {
	if filepath.Ext(fr.fileSettings.ReadVCardPath) == "" {
		//try .vcf
		alternateFilePath := fr.fileSettings.ReadVCardPath + ".vcf"
		_, err := fr.fileSystem.Stat(alternateFilePath)
		if err == nil {
			//when the .vcf ending does not produce an error, this will be the inputFilePath to use
			fr.fileSettings.ReadVCardPath = alternateFilePath
		}
	}
}

func (fr *Repository) WriteVCard(card vcard.Card) error {
	file, err := fr.fileSystem.Create(fr.fileSettings.WriteVCardPath)
	if err != nil {
		return err
	}
	defer file.Close()

	vCardContent, err := fr.cardCodec.Encode(card)
	if err != nil {
		return err
	}

	if _, err := file.Write(vCardContent); err != nil {
		return err
	} else {
		fr.userNotifier.Notifyf("The vCard has been written to %s", fr.fileSettings.WriteVCardPath)
	}

	return nil
}

func (fr *Repository) WriteQRCode(card vcard.Card) error {
	img, err := fr.qrCodec.Encode(card, fr.appSettings.QRSettings)
	if err != nil {
		return err
	}

	file, err := fr.fileSystem.Create(fr.fileSettings.WriteQRCodePath)
	if err != nil {
		return err
	}
	defer file.Close()

	if err := png.Encode(file, img); err != nil {
		return err
	} else {
		fr.userNotifier.Notifyf("The QR code has been written to %s", fr.fileSettings.WriteQRCodePath)
	}

	return nil

}

func ensureNilSafety(card vcard.Card) {
	if card.Name() == nil {
		name := vcard.Name{}
		card.SetName(&name)
	}
	if card.Address() == nil {
		address := vcard.Address{}
		card.SetAddress(&address)
	}
}
