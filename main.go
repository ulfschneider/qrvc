package main

import (
	"errors"
	"fmt"
	"qrvc/internal/cli"
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
		cli.Println(err)
	}

	cli.Println("👋")
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

	cli.Println("You are running qrvc, a tool to prepare a QR code from a vCard.")
	cli.Println("Get a list of options by starting the program in the form: qrvc -h")

	err = run(args)
}

// TODO test
// TODO documentation
