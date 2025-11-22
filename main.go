package main

import (
	"os"
	"qrvc/internal/cli"
	"qrvc/internal/settings"
	"qrvc/internal/vcard"

	"github.com/manifoldco/promptui"
	"github.com/skip2/go-qrcode"
)

func run(settings *settings.Settings) error {

	vcardContent, err := vcard.PrepareVcard(settings)

	if err != nil {
		return err
	}

	if file, err := os.Create(*settings.VCardOutputFilePath); err != nil {
		return err
	} else {
		defer file.Close()
		if _, err := file.WriteString(vcardContent); err != nil {
			return err
		} else {
			cli.Println("The vCard has been written to", cli.SprintValue(*settings.VCardOutputFilePath))
		}
	}

	if err := qrcode.WriteColorFile(vcardContent, qrcode.Medium, 256, *settings.BackgroundColor, *settings.ForegroundColor, *settings.QRCodeOutputFilePath); err != nil {
		return err
	} else {
		cli.Println("The QR code has been written to", cli.SprintValue(*settings.QRCodeOutputFilePath))
	}

	return nil
}

func finalize(err error, args *settings.Settings) {
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
		finalize(err, args)
	}()

	if args, err = settings.PrepareSettings(); err != nil {
		return
	}

	cli.Println("You are running qrvc, a commandline tool to prepare a QR code from a vCard.")
	cli.Println("Get a list of options by starting the program in the form: qrvc -h")

	err = run(args)
}

// TODO format cli help
// TODO use a logger that reflects silent
// TODO test
