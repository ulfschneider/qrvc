// Package main provides a command line program that allows to process vCard data from a file or from a command line form and store that data as a QR Code.
package main

import (
	"errors"

	"github.com/spf13/afero"

	bomembedded "github.com/ulfschneider/qrvc/internal/adapters/bom/embedded"
	qrcodec "github.com/ulfschneider/qrvc/internal/adapters/codec/qr"
	vcardcodec "github.com/ulfschneider/qrvc/internal/adapters/codec/vcard"
	configcli "github.com/ulfschneider/qrvc/internal/adapters/config/cli"
	editorcli "github.com/ulfschneider/qrvc/internal/adapters/editor/cli"
	notifiercli "github.com/ulfschneider/qrvc/internal/adapters/notifier/cli"
	repofile "github.com/ulfschneider/qrvc/internal/adapters/repo/file"
	versionembedded "github.com/ulfschneider/qrvc/internal/adapters/version/embedded"

	"github.com/ulfschneider/qrvc/internal/application/services"

	"github.com/charmbracelet/huh"
)

func runQRCard(settings configcli.CLIFileSettings) error {
	cardCodec := vcardcodec.NewCodec()
	qrCodec := qrcodec.NewCodec()
	repo := repofile.NewRepo(
		afero.NewOsFs(),
		&cardCodec,
		&qrCodec,
		settings.Files,
		settings.App)

	editor := editorcli.NewCardEditor()

	cardService := services.NewQRCardService(settings.App, &repo, &editor)

	err := cardService.TransformCard()

	return err
}

func runBOM() error {
	bomProvider := bomembedded.NewBomProvider()
	bomService := services.NewBomService(&bomProvider)

	err := bomService.WriteBomJSON()

	return err
}

func runVersion() {
	versionProvider := versionembedded.NewVersionProvider()
	version := versionProvider.Version()
	notifier := notifiercli.NewUserNotifier()

	if version != "" {
		notifier.NotifyfLoud("%s", version)
	} else {
		notifier.NotifyLoud("No version information available")
	}
}

func finalize(settings configcli.CLIFileSettings, err error) {
	userNotifier := notifiercli.NewUserNotifier()
	if errors.Is(err, huh.ErrUserAborted) {
		// User pressed Ctrl-C
		userNotifier.Notify("You stopped with CTRL-C")
	} else if err != nil {
		//any other error
		userNotifier.Notify(err)
	}
	if settings.CLI.Bom == false && settings.CLI.AppVersion == false && settings.App.Silent == false {
		//say good bye
		userNotifier.Section()
		userNotifier.Notify("👋")
	}

}

func loadConfig() (configcli.CLIFileSettings, error) {

	versionProvider := versionembedded.NewVersionProvider()
	versionService := services.NewVersionService(&versionProvider)
	settingsProvider := configcli.NewSettingsProvider(versionService)

	settings, err := settingsProvider.Load()
	if err != nil {
		return configcli.CLIFileSettings{}, err
	}
	return settings, nil
}

func main() {

	settings, err := loadConfig()
	if err != nil {
		return
	}
	defer func() {
		finalize(settings, err)
	}()

	if !settings.CLI.Bom && !settings.CLI.AppVersion {
		userNotifier := notifiercli.NewUserNotifier()
		userNotifier.Notify("You are running qrvc, a tool to prepare a QR code from a vCard.")
		userNotifier.Notifyf("Get a list of options by starting the program in the form: %s", "qrvc -h")
		userNotifier.Notifyf("Stop the program by pressing %s", "CTRL-C")
		userNotifier.Section()
		err = runQRCard(settings)
	} else if settings.CLI.Bom {
		err = runBOM()
	} else if settings.CLI.AppVersion {
		runVersion()
	}
}
