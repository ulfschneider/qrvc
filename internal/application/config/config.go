package config

import (
	"image/color"

	"github.com/skip2/go-qrcode"
)

type Settings struct {
	Silent       bool
	VCardVersion string
	QRSettings   QRCodeSettings
}

type QRCodeSettings struct {
	Border          bool
	Size            int
	RecoveryLevel   qrcode.RecoveryLevel
	BackgroundColor color.Color
	ForegroundColor color.Color
}
