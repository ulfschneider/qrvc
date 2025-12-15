// Package main provides a command line program that allows to process vCard data from a file or from a command line form and store that data as a QR Code.
package main

import (
	"errors"

	"github.com/spf13/afero"

	"github.com/ulfschneider/qrvc/internal/adapters/cliconfig"
	"github.com/ulfschneider/qrvc/internal/adapters/clieditor"
	"github.com/ulfschneider/qrvc/internal/adapters/clinotifier"

	"github.com/ulfschneider/qrvc/internal/adapters/embeddedbom"
	"github.com/ulfschneider/qrvc/internal/adapters/embeddedversion"
	"github.com/ulfschneider/qrvc/internal/adapters/filerepo"
	"github.com/ulfschneider/qrvc/internal/adapters/qrencoder"
	"github.com/ulfschneider/qrvc/internal/application/services"

	"github.com/charmbracelet/huh"
)

func runQRCard(settings cliconfig.CLIFileSettings) error {
	qrEncoder := qrencoder.NewQRCardEncoder()
	repo := filerepo.NewFileRepo(
		afero.NewOsFs(),
		&qrEncoder,
		settings.Files,
		settings.App)

	editor := clieditor.NewCLIVCardEditor()

	cardService := services.NewQRCardService(settings.App, &repo, &editor)

	err := cardService.TransformCard()

	return err
}

func runBOM() error {
	bomProvider := embeddedbom.NewBomProvider()
	bomService := services.NewBomService(&bomProvider)

	err := bomService.WriteBomJSON()

	return err
}

func finalize(settings cliconfig.CLIFileSettings, err error) {
	userNotifier := clinotifier.NewCLINotifier()
	if errors.Is(err, huh.ErrUserAborted) {
		// User pressed Ctrl-C
		userNotifier.Notify("You stopped with CTRL-C")
	} else if err != nil {
		//any other error
		userNotifier.Notify(err)
	}
	if settings.CLI.Bom == false {
		//say good bye
		userNotifier.Section()
		userNotifier.Notify("👋")
	}

}

func loadConfig() (cliconfig.CLIFileSettings, error) {

	versionProvider := embeddedversion.NewVersionProvider()
	versionService := services.NewVersionService(&versionProvider)
	settingsProvider := cliconfig.NewCLIFileSettingsProvider(versionService)

	settings, err := settingsProvider.Load()
	if err != nil {
		return cliconfig.CLIFileSettings{}, err
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

	if !settings.CLI.Bom {
		userNotifier := clinotifier.NewCLINotifier()
		userNotifier.Notify("You are running qrvc, a tool to prepare a QR code from a vCard.")
		userNotifier.Notifyf("Get a list of options by starting the program in the form: %s", "qrvc -h")
		userNotifier.Notifyf("Stop the program by pressing %s", "CTRL-C")
		userNotifier.Section()
		err = runQRCard(settings)
	} else {
		err = runBOM()
	}
}
