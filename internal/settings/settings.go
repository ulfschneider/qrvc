package settings

import (
	"github.com/spf13/pflag"

	"github.com/mazznoer/csscolorparser"
)

const ApplicationName = "QRVC"

type Settings struct {
	InputFilePath        *string
	VCardVersion         *string
	QRCodeOutputFilePath *string
	VCardOutputFilePath  *string
	BackgroundColor      *csscolorparser.Color
	ForegroundColor      *csscolorparser.Color
}

func PrepareSettings() (*Settings, error) {

	settings := Settings{}
	settings.InputFilePath = pflag.String("i", "", "The path and name of the vCard input file. ")

	settings.VCardVersion = pflag.String("v", "3.0", "The vCard version to create.")

	settings.QRCodeOutputFilePath = pflag.String("q", "vcard.png", "The path and name of the generated QR code file.")

	settings.VCardOutputFilePath = pflag.String("o", "vcard.vcf", "The path and name of generated vCard file.")

	foregroundColor := pflag.String("f", "black", "The foreground color of the QR code. This can be a hex RGB color value or a CSS color name.")

	backgroundColor := pflag.String("b", "transparent", "The background color of the QR code. This can be a hex RGB color value or a CSS color name.")

	pflag.Parse()

	//bring the colors into the correct format

	if c, err := csscolorparser.Parse(*foregroundColor); err != nil {
		return nil, err
	} else {
		settings.ForegroundColor = &c
	}

	if c, err := csscolorparser.Parse(*backgroundColor); err != nil {
		return nil, err
	} else {
		settings.BackgroundColor = &c
	}

	return &settings, nil
}
