package settings

import (
	"errors"
	"path/filepath"
	"qrvc/internal/cli"
	"strings"

	"github.com/spf13/pflag"

	"github.com/mazznoer/csscolorparser"
)

const ApplicationName = "QRVC"

type Settings struct {
	InputFilePath        *string
	VCardVersion         *string
	QRCodeOutputFilePath *string
	VCardOutputFilePath  *string
	Silent               *bool
	BackgroundColor      *csscolorparser.Color
	ForegroundColor      *csscolorparser.Color
}

func PrepareSettings() (*Settings, error) {

	settings := Settings{}
	settings.InputFilePath = pflag.StringP("input", "i", "", "The path and name of the vCard input file. ")

	settings.VCardVersion = pflag.StringP("version", "v", "3.0", "The vCard version to create.")

	outputFilePath := pflag.StringP("output", "o", "", "The path and name for the output. Will receive the extension .png for the QR code and .vcf for the vCard. Will use the input file basename by default.")

	foregroundColor := pflag.StringP("foreground", "f", "black", "The foreground color of the QR code. This can be a hex RGB color value or a CSS color name.")

	backgroundColor := pflag.StringP("background", "b", "transparent", "The background color of the QR code. This can be a hex RGB color value or a CSS color name.")

	settings.Silent = pflag.BoolP("silent", "s", false, "Silent mode, will not interactively ask for input.")

	pflag.Parse()

	//prepare output file names
	if *settings.InputFilePath != "" && *outputFilePath == "" {
		base := filepath.Base(*settings.InputFilePath)                 // "file.txt"
		*outputFilePath = strings.TrimSuffix(base, filepath.Ext(base)) // "file"
	}
	if *outputFilePath == "" {
		*outputFilePath = "vcard"
	}
	settings.QRCodeOutputFilePath = new(string)
	*settings.QRCodeOutputFilePath = *outputFilePath + ".png"
	settings.VCardOutputFilePath = new(string)
	*settings.VCardOutputFilePath = *outputFilePath + ".vcf"

	//verify silent mode
	if *settings.Silent && *settings.InputFilePath == "" {
		return nil, errors.New("You must provide an input file when running in silent mode")
	}
	cli.Silent = *settings.Silent

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
