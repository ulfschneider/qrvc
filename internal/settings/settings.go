package settings

import (
	"errors"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/ulfschneider/qrvc/internal/cli"
	"github.com/ulfschneider/qrvc/internal/version"

	"github.com/spf13/pflag"

	"github.com/mazznoer/csscolorparser"
)

type Settings struct {
	InputFilePath        *string
	VCardVersion         *string
	QRCodeOutputFilePath *string
	VCardOutputFilePath  *string
	Border               *bool
	Size                 *int
	Silent               *bool
	BackgroundColor      *csscolorparser.Color
	ForegroundColor      *csscolorparser.Color
}

func PrepareSettings() (*Settings, error) {

	settings := Settings{}

	settings.InputFilePath = pflag.StringP("input", "i", "", "The path and name of the vCard input file. When you provide a file name without extension, .vcf will automatically added as an extension.")

	settings.VCardVersion = pflag.StringP("version", "v", "3.0", "The vCard version to create.")

	outputFilePath := pflag.StringP("output", "o", "", "The path and name for the output. Please do not add any file extension, as those will be added automatically.\nWill receive the extension .png for the QR code and .vcf for the vCard. The input file basename will be used by default.")

	foregroundColor := pflag.StringP("foreground", "f", "black", "The foreground color of the QR code. This can be a hex RGB color value or a CSS color name.")

	backgroundColor := pflag.StringP("background", "b", "transparent", "The background color of the QR code. This can be a hex RGB color value or a CSS color name.")

	settings.Border = pflag.BoolP("border", "r", false, "Whether the QR code has a border or not.")

	settings.Size = pflag.IntP("size", "z", 400, "The size of the resulting QR code in width and height of pixels.")

	settings.Silent = pflag.BoolP("silent", "s", false, "The silent mode will not interactively ask for input and instead requires an input file.")

	pflag.CommandLine.SortFlags = false

	pflag.Usage = func() {
		fmt.Printf("qrvc %s\n", version.Version)
		fmt.Println("qrvc is a tool to prepare a QR code from a vCard")
		fmt.Println("\nUsage: qrvc [flags]")

		fmt.Println("\nFlags:")
		pflag.CommandLine.VisitAll(func(f *pflag.Flag) {
			if f.Shorthand != "" {
				// prints: -h, --help
				fmt.Printf("\n-%s, --%s (%s)", cli.SprintValue(f.Shorthand), cli.SprintValue(f.Name), f.Value.Type())
			} else {
				// prints: --name (no short version)
				fmt.Printf("\n--%s (%s)", cli.SprintValue(f.Name), f.Value.Type())
			}

			if f.DefValue != "" {
				fmt.Printf("\n%s (Default: %s)\n", f.Usage, cli.SprintValue(f.DefValue))
			} else {
				fmt.Printf("\n%s\n", f.Usage)
			}
		})

	}

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
