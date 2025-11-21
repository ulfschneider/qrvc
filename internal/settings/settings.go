package settings

import (
	"github.com/spf13/pflag"

	"github.com/mazznoer/csscolorparser"
)

const ApplicationName = "QRVC"

type Settings struct {
	InputFilePath        *string
	VcardVersion         *string
	QrCodeOutputFilePath *string
	VcardOutputFilePath  *string
	BackgroundColor      *csscolorparser.Color
	ForegroundColor      *csscolorparser.Color
}

func PrepareSettings() (*Settings, error) {

	settings := Settings{}
	settings.InputFilePath = pflag.String("i", "", "The path and name of the vcard input file. When this argument is not provided, the programm will ask for vcard details interactively.")

	settings.VcardVersion = pflag.String("v", "3.0", "The vcard version to create. This will only be used when no vcard input file is given.")

	settings.QrCodeOutputFilePath = pflag.String("q", "vcard.png", "The path and name of the QR code output file to write.")

	settings.VcardOutputFilePath = pflag.String("o", "vcard.vcf", "When no vcard input file is provided, the interactiveley generated vcard will be stored under the name given here.")

	foregroundColor := pflag.String("f", "black", "Define the foreground color of the QR code. This can be a hex RGB color value or a CSS color name.")

	backgroundColor := pflag.String("b", "transparent", "Define the background color of the QR code. This can be a hex RGB color value or a CSS color name.")

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
