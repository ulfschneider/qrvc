package filerepo

import (
	"bytes"
	"image/png"
	"io"
	"io/fs"
	"path/filepath"

	"github.com/emersion/go-vcard"
	"github.com/pkg/errors"
	"github.com/spf13/afero"
	"github.com/ulfschneider/qrvc/internal/adapters/cliconfig"
	"github.com/ulfschneider/qrvc/internal/adapters/clinotifier"

	"github.com/ulfschneider/qrvc/internal/application/config"
	"github.com/ulfschneider/qrvc/internal/application/ports"
)

func NewFileRepo(
	fileSystem afero.Fs,
	qrEncoder ports.QREncoder,
	fileSettings cliconfig.FileSettings,
	appSettings config.Settings,
) FileRepository {

	return FileRepository{
		fileSystem:   fileSystem,
		qrEncoder:    qrEncoder,
		userNotifier: clinotifier.NewCLINotifier(),
		fileSettings: fileSettings,
		appSettings:  appSettings,
	}
}

type FileRepository struct {
	fileSystem   afero.Fs
	qrEncoder    ports.QREncoder
	userNotifier clinotifier.CLINotifier
	fileSettings cliconfig.FileSettings
	appSettings  config.Settings
}

func (fr *FileRepository) ReadOrCreateVCard() (*vcard.Card, error) {

	if fr.fileSettings.ReadVCardPath == "" && fr.userNotifier.Silent() == false {
		//no path to a vcard file, create a new card
		card := make(vcard.Card)
		card.SetValue(vcard.FieldVersion, fr.appSettings.VCardVersion)
		ensureNilSafety(&card)
		return &card, nil
	} else if fr.fileSettings.ReadVCardPath == "" && fr.userNotifier.Silent() == true {
		//no path to a vcard file, but tool runs in silent mode
		return nil, errors.New("Missing input file path")
	} else {
		// we have a path to a vcard file, try to read it
		file, err := fr.openVcard()
		if err != nil {
			return nil, err
		}

		defer file.Close()

		fr.userNotifier.Section()
		fr.userNotifier.Notifyf("Reading vCard file %s", fr.fileSettings.ReadVCardPath)
		fr.userNotifier.Section()

		if card, err := fr.decodeVcard(file); err != nil {
			return nil, err
		} else {
			ensureNilSafety(&card)
			return &card, nil
		}
	}

}

func (fr *FileRepository) WriteVCard(card *vcard.Card) error {
	file, err := fr.fileSystem.Create(fr.fileSettings.WriteVCardPath)
	if err != nil {
		return err
	}
	defer file.Close()

	vCardContent, err := fr.encodeVcard(card)
	if err != nil {
		return err
	}

	if _, err := file.WriteString(vCardContent); err != nil {
		return err
	} else {
		fr.userNotifier.Notifyf("The vCard has been written to %s", fr.fileSettings.WriteVCardPath)
	}

	return nil
}

func (fr *FileRepository) WriteQRCode(card *vcard.Card) error {
	img, err := fr.qrEncoder.Encode(card, fr.appSettings.QRSettings)
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

func (fr *FileRepository) openVcard() (fs.File, error) {
	file, err := fr.fileSystem.Open(fr.fileSettings.ReadVCardPath)
	if err != nil {
		if filepath.Ext(fr.fileSettings.ReadVCardPath) == "" {
			//try .vcf
			alternateFilePath := fr.fileSettings.ReadVCardPath + ".vcf"
			file, err = fr.fileSystem.Open(alternateFilePath)
			if err == nil {
				//when the .vcf was possible to read, this will be the new inputFilePath
				fr.fileSettings.ReadVCardPath = alternateFilePath
			}
		}
	}

	if err != nil {
		return nil, err
	}
	return file, err
}

func (fr *FileRepository) decodeVcard(reader io.Reader) (vcard.Card, error) {
	dec := vcard.NewDecoder(reader)
	card, err := dec.Decode()
	if err != nil {
		return nil, err
	}
	return card, nil
}

func (fr *FileRepository) encodeVcard(card *vcard.Card) (string, error) {
	var buf bytes.Buffer
	enc := vcard.NewEncoder(&buf)
	if err := enc.Encode(*card); err != nil {
		return "", err
	}

	return buf.String(), nil
}
