package main

import (
	"errors"
	"fmt"
	"qrvc/internal/out"
	"qrvc/internal/settings"
	"qrvc/internal/vcard"

	"github.com/charmbracelet/huh"
)

func run(settings *settings.Settings) error {
	vcardContent, err := vcard.PrepareVcard(settings)

	if err != nil {
		return err
	}

	return out.StoreResults(vcardContent, settings)
}

func finalize(err error) {
	if errors.Is(err, huh.ErrUserAborted) {
		// User pressed Ctrl-C
		fmt.Println("You stopped with CTRL-C")
		return
	} else if err != nil {
		fmt.Println(err)
	}

	fmt.Println("👋")
}

func main() {

	var err error
	var args *settings.Settings

	defer func() {
		finalize(err)
	}()

	if args, err = settings.PrepareSettings(); err != nil {
		return
	}

	fmt.Println("You are running qrvc, a tool to prepare a QR code from a vCard.")
	fmt.Println("Get a list of options by starting the program in the form: qrvc -h")
	fmt.Println("Stop the program by pressing CTRL-C.")
	fmt.Println()

	err = run(args)
}

// TODO sbom for each build target
// TODO test
// TODO documentation
