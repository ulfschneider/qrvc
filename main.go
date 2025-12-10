// Package main provides a command line program that allows to process vCard data from a file or from a command line form and store that data as a QR Code.
package main

import (
	"errors"
	"fmt"

	"github.com/spf13/afero"
	"github.com/ulfschneider/qrvc/internal/appmeta"
	"github.com/ulfschneider/qrvc/internal/cli"
	"github.com/ulfschneider/qrvc/internal/qrcard"
	"github.com/ulfschneider/qrvc/internal/settings"

	"github.com/charmbracelet/huh"
)

func runQRCard(settings *settings.Settings) error {
	filesystem := afero.NewOsFs()

	vcardContent, err := qrcard.PrepareVCard(settings.InputFilePath, *settings.VCardVersion, *settings.Silent, filesystem)

	if err != nil {
		return err
	}

	err = qrcard.WriteResults(vcardContent, settings, filesystem)
	return err
}

func runBOM() error {
	bom, err := appmeta.LoadEmbeddedBOM()
	if err != nil {
		return err
	}

	json, err := appmeta.MarshalBOMToJSON(bom)
	if err != nil {
		return err
	} else {
		fmt.Println(string(json))
	}
	return nil
}

func finalize(args *settings.Settings, err error) {
	if errors.Is(err, huh.ErrUserAborted) {
		// User pressed Ctrl-C
		fmt.Println("You stopped with CTRL-C")
		return
	} else if err != nil {
		fmt.Println(errors.Unwrap(err))
	}
	if args != nil && !*args.Bom {
		fmt.Println("👋")
	}

}

func main() {

	var err error
	var args *settings.Settings

	defer func() {
		finalize(args, err)
	}()

	if args, err = settings.PrepareSettings(); err != nil {
		return
	}

	if !*args.Bom {
		fmt.Println("You are running qrvc, a tool to prepare a QR code from a vCard.")
		fmt.Println("Get a list of options by starting the program in the form:", cli.SprintValue("qrvc -h"))
		fmt.Println("Stop the program by pressing", cli.SprintValue("CTRL-C"))
		fmt.Println()

		err = runQRCard(args)
	} else {
		err = runBOM()
	}
}
