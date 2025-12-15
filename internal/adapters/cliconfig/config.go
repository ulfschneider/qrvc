package cliconfig

import (
	"errors"
	"image/color"
	"path/filepath"
	"strings"

	"github.com/mazznoer/csscolorparser"
	"github.com/skip2/go-qrcode"
	"github.com/spf13/pflag"

	"github.com/ulfschneider/qrvc/internal/adapters/clinotifier"
	"github.com/ulfschneider/qrvc/internal/application/config"
	"github.com/ulfschneider/qrvc/internal/application/services"
)

type SettingsProvider struct {
	versionService services.VersionService
	userNotifier   clinotifier.UserNotifier
}

type CLIFileSettings struct {
	App   config.Settings
	Files FileSettings
	CLI   CLISettings
}

type FileSettings struct {
	ReadVCardPath   string
	WriteVCardPath  string
	WriteQRCodePath string
}

type CLISettings struct {
	Bom bool
}

func NewSettingsProvider(versionService services.VersionService) SettingsProvider {
	cliNotifier := clinotifier.NewUserNotifier()
	return SettingsProvider{versionService: versionService, userNotifier: cliNotifier}
}

func (sp *SettingsProvider) Load() (CLIFileSettings, error) {

	silent := pflag.BoolP("silent", "s", false, "The silent mode will not interactively ask for input and instead requires a vCard input file.")

	readVCardPath := pflag.StringP("input", "i", "", "The path and name of the vCard input file. When you provide a file name without extension, .vcf will automatically added as an extension.")

	writePath := pflag.StringP("output", "o", "", "The path and name for the output. Please do not add any file extension, as those will be added automatically.\nWill receive the extension .png for the QR code and .vcf for the vCard. The input file basename will be used by default.")

	vCardVersion := pflag.StringP("version", "v", "3.0", "The vCard version to create.")

	foregroundColor := pflag.StringP("foreground", "f", "black", "The foreground color of the QR code. This can be a hex RGB color value (like \"#000\") or a CSS color name (like black).")

	backgroundColor := pflag.StringP("background", "b", "transparent", "The background color of the QR code. This can be a hex RGB color value (like \"#fff\") or a CSS color name (like white).")

	border := pflag.BoolP("border", "r", false, "Whether the QR code has a border or not.")

	size := pflag.IntP("size", "z", 400, "The size of the resulting QR code in width and height of pixels.")

	bom := pflag.BoolP("bom", "m", false, "List the Software Bill of Materials of this tool in CycloneDX format.")

	sp.formatFlagUsage() //adjust help format before parsing
	pflag.Parse()        //process flags

	settings := CLIFileSettings{}
	settings.App = config.Settings{}
	settings.App.QRSettings = config.QRCodeSettings{}
	settings.Files = FileSettings{}
	settings.CLI = CLISettings{}

	settings.App.Silent = *silent
	sp.userNotifier.SetSilent(settings.App.Silent)

	settings.Files.ReadVCardPath = *readVCardPath
	if settings.App.Silent && settings.Files.ReadVCardPath == "" {
		return CLIFileSettings{}, errors.New("You must provide an input file when running in silent mode")
	}

	//adjust names according to readVCard
	if settings.Files.ReadVCardPath != "" && *writePath == "" {
		base := filepath.Base(settings.Files.ReadVCardPath)       // "file.txt"
		*writePath = strings.TrimSuffix(base, filepath.Ext(base)) // "file"
	}
	if *writePath == "" {
		*writePath = "vcard"
	}
	settings.Files.WriteQRCodePath = *writePath + ".png"
	settings.Files.WriteVCardPath = *writePath + ".vcf"

	settings.App.VCardVersion = *vCardVersion

	settings.App.QRSettings.Border = *border
	settings.App.QRSettings.Size = *size

	//bring the colors into the correct format
	if color, err := sp.parseColor(*foregroundColor); err != nil {
		return CLIFileSettings{}, err
	} else {
		settings.App.QRSettings.ForegroundColor = color
	}
	if color, err := sp.parseColor(*backgroundColor); err != nil {
		return CLIFileSettings{}, err
	} else {
		settings.App.QRSettings.BackgroundColor = color
	}

	settings.App.QRSettings.RecoveryLevel = qrcode.Low

	settings.CLI.Bom = *bom

	return settings, nil
}

func (sp *SettingsProvider) formatFlagUsage() {
	pflag.CommandLine.SortFlags = false
	pflag.Usage = func() {
		version, _ := sp.versionService.Version()

		if version != "" {
			sp.userNotifier.NotifyfLoud("qrvc %s\n", version)
		} else {
			sp.userNotifier.NotifyLoud("qrvc")
		}
		sp.userNotifier.NotifyLoud("qrvc is a tool to prepare a QR code from a vCard")
		sp.userNotifier.Section()
		sp.userNotifier.NotifyLoud("Usage: qrvc [flags]")
		sp.userNotifier.Section()
		sp.userNotifier.NotifyLoud("Flags:")
		pflag.CommandLine.VisitAll(func(f *pflag.Flag) {
			sp.userNotifier.Section()
			if f.Shorthand != "" {
				// prints: -h, --help (type)
				sp.userNotifier.NotifyfLoud("-%s, --%s ("+f.Value.Type()+")", f.Shorthand, f.Name)
			} else {
				// prints: --flag (type)
				sp.userNotifier.NotifyfLoud("--%s ("+f.Value.Type()+")", f.Name)
			}

			if f.DefValue != "" {
				sp.userNotifier.NotifyfLoud(f.Usage+" (Default: %s)", f.DefValue)
			} else {
				sp.userNotifier.NotifyLoud(f.Usage)
			}
		})
	}
}

func (sp *SettingsProvider) parseColor(color string) (color.Color, error) {
	if c, err := csscolorparser.Parse(color); err != nil {
		return nil, err
	} else {
		return c, nil
	}
}
