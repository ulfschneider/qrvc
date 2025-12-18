package ports

import (
	"image"

	"github.com/CycloneDX/cyclonedx-go"
	"github.com/emersion/go-vcard"
	"github.com/ulfschneider/qrvc/internal/application/config"
)

type VCardEditor interface {
	Edit(card vcard.Card) error
}

type Repository interface {
	ReadOrCreateVCard() (vcard.Card, error)
	WriteVCard(card vcard.Card) error
	WriteQRCode(card vcard.Card) error
}

type QRCodec interface {
	Encode(card vcard.Card, settings config.QRCodeSettings) (image.Image, error)
}

type VCardCodec interface {
	Encode(card vcard.Card) ([]byte, error)
	Decode(vcf []byte) (vcard.Card, error)
}

type VersionProvider interface {
	Version() string
	Commit() string
	Time() string
}

type UserNotifier interface {
	Notify(message ...string)
	NotifyLoud(message ...string)
	Notifyf(format string, values ...any)
	NotifyfLoud(format string, values ...any)
	Section()
	SectionLoud()
	SetSilent(isSilent bool)
	Silent() bool
}

type BomProvider interface {
	Bom() (*cyclonedx.BOM, error)
	MarshalToJSON() ([]byte, error)
	WriteBomJSON() error
}
