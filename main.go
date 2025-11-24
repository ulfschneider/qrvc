package main

import (
	"qrvc/internal/cli"
	"qrvc/internal/out"
	"qrvc/internal/settings"
	"qrvc/internal/vcard"

	"github.com/manifoldco/promptui"
)

func run(settings *settings.Settings) error {
	vcardContent, err := vcard.PrepareVcard(settings)

	if err != nil {
		return err
	}

	return out.PrintResults(vcardContent, settings)
}

func finalize(err error) {
	if err == promptui.ErrInterrupt {
		cli.Println("You stopped via CTRL-C")
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
