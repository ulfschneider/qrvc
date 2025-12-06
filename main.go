package main

import (
	"errors"
	"fmt"

	"github.com/ulfschneider/qrvc/internal/cli"
	"github.com/ulfschneider/qrvc/internal/out"
	"github.com/ulfschneider/qrvc/internal/sbom"
	"github.com/ulfschneider/qrvc/internal/settings"
	"github.com/ulfschneider/qrvc/internal/vcard"

	"github.com/charmbracelet/huh"
)

func runVCard(settings *settings.Settings) error {
	vcardContent, err := vcard.PrepareVcard(settings)

	if err != nil {
		return err
	}

	return out.StoreResults(vcardContent, settings)
}

func runSbom() error {
	bom, err := sbom.LoadEmbeddedSBOM()
	if err != nil {
		return err
	}
	formattedBom, err := sbom.Sprintf(bom)
	if err != nil {
		return err
	}
	fmt.Println(formattedBom)
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

		err = runVCard(args)
	} else {
		err = runSbom()
	}
}

// TODO test
// TODO documentation
