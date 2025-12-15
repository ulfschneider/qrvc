package qrencoder

import (
	"bytes"
	"image"

	"github.com/emersion/go-vcard"
	"github.com/skip2/go-qrcode"
	"github.com/ulfschneider/qrvc/internal/application/config"
)

func NewQRCardEncoder() QRCardEncoder {
	return QRCardEncoder{}
}

type QRCardEncoder struct {
}

func (qe *QRCardEncoder) Encode(card *vcard.Card, settings config.QRCodeSettings) (image.Image, error) {
	vCardContent, err := qe.encodeToString(card)
	if err != nil {
		return nil, err
	}

	qr, err := qrcode.New(vCardContent, settings.RecoveryLevel)
	if err != nil {
		return nil, err
	}

	qr.DisableBorder = !settings.Border
	qr.ForegroundColor = settings.ForegroundColor
	qr.BackgroundColor = settings.BackgroundColor

	img := qr.Image(settings.Size)

	return img, nil
}

func (qe *QRCardEncoder) encodeToString(card *vcard.Card) (string, error) {
	var buf bytes.Buffer
	enc := vcard.NewEncoder(&buf)
	if err := enc.Encode(*card); err != nil {
		return "", err
	}

	return buf.String(), nil
}
