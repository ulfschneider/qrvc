package configcli

import (
	"errors"
	"image/color"
	"os"
	"path/filepath"
	"strings"

	"github.com/mazznoer/csscolorparser"
	"github.com/skip2/go-qrcode"
	"github.com/spf13/pflag"

	notifiercli "github.com/ulfschneider/qrvc/internal/adapters/notifier/cli"
	"github.com/ulfschneider/qrvc/internal/application/config"
	"github.com/ulfschneider/qrvc/internal/application/services"
)

type SettingsProvider struct {
	flagSet        *pflag.FlagSet
	versionService services.VersionService
	userNotifier   notifiercli.UserNotifier
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
	Bom        bool
	AppVersion bool
}

func NewSettingsProvider(versionService services.VersionService) SettingsProvider {
	flagSet := pflag.NewFlagSet(os.Args[0], pflag.ExitOnError)
	cliNotifier := notifiercli.NewUserNotifier()
	return SettingsProvider{flagSet: flagSet, versionService: versionService, userNotifier: cliNotifier}
}

func (sp *SettingsProvider) Load() (CLIFileSettings, error) {

	silent := sp.flagSet.BoolP("silent", "s", false, "The silent mode will not interactively ask for input and instead requires a vCard input file.")

	readVCardPath := sp.flagSet.StringP("input", "i", "", "The path and name of the vCard input file. When you provide a file name without extension, .vcf will automatically added as an extension.")

	writePath := sp.flagSet.StringP("output", "o", "", "The path and name for the output. Please do not add any file extension, as those will be added automatically.\nWill receive the extension .png for the QR code and .vcf for the vCard. The input file basename will be used by default.")

	vCardVersion := sp.flagSet.StringP("cardversion", "c", "3.0", "The vCard version to create.")

	foregroundColor := sp.flagSet.StringP("foreground", "f", "black", "The foreground color of the QR code. This can be a hex RGB color value (like \"#000\") or a CSS color name (like black).")

	backgroundColor := sp.flagSet.StringP("background", "b", "white", "The background color of the QR code. This can be a hex RGB color value (like \"#fff\") or a CSS color name (like white, or transparent).")

	border := sp.flagSet.BoolP("border", "r", false, "Whether the QR code has a border or not.")

	size := sp.flagSet.IntP("size", "z", 400, "The size of the resulting QR code in width and height of pixels.")

	bom := sp.flagSet.BoolP("bom", "m", false, "List the Software Bill of Materials of this tool in CycloneDX format.")

	appVersion := sp.flagSet.BoolP("version", "v", false, "Show the qrvc version.")

	sp.formatFlagUsage()          //adjust help format before parsing
	sp.flagSet.Parse(os.Args[1:]) //process flags

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
	settings.CLI.AppVersion = *appVersion

	return settings, nil
}

func (sp *SettingsProvider) formatFlagUsage() {
	sp.flagSet.SortFlags = false
	sp.flagSet.Usage = func() {
		version := sp.versionService.Version()
		time := sp.versionService.Time()

		if version != "" {
			sp.userNotifier.NotifyfLoud("qrvc %s %s\n", version, time)
		} else {
			sp.userNotifier.NotifyLoud("qrvc")
		}
		sp.userNotifier.NotifyLoud("qrvc is a tool to prepare a QR code from a vCard")
		sp.userNotifier.Section()
		sp.userNotifier.NotifyLoud("Usage: qrvc [flags]")
		sp.userNotifier.Section()
		sp.userNotifier.NotifyLoud("Flags:")
		sp.flagSet.VisitAll(func(f *pflag.Flag) {
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
