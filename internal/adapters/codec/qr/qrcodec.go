package qrcodec

import (
	"image"

	"github.com/emersion/go-vcard"
	"github.com/skip2/go-qrcode"

	vcardcodec "github.com/ulfschneider/qrvc/internal/adapters/codec/vcard"
	"github.com/ulfschneider/qrvc/internal/application/config"
)

func NewCodec() Codec {
	return Codec{}
}

type Codec struct {
}

func (qe *Codec) Encode(card vcard.Card, settings config.QRCodeSettings) (image.Image, error) {
	cardCodec := vcardcodec.NewCodec()

	vCardContent, err := cardCodec.Encode(card)
	if err != nil {
		return nil, err
	}

	qr, err := qrcode.New(string(vCardContent), settings.RecoveryLevel)
	if err != nil {
		return nil, err
	}

	qr.DisableBorder = !settings.Border
	qr.ForegroundColor = settings.ForegroundColor
	qr.BackgroundColor = settings.BackgroundColor

	img := qr.Image(settings.Size)

	return img, nil
}
