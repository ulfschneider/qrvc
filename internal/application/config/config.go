package config

import (
	"github.com/mazznoer/csscolorparser"
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
	BackgroundColor csscolorparser.Color
	ForegroundColor csscolorparser.Color
}
