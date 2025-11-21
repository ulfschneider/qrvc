package main

import (
	"fmt"
	"os"
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

	fmt.Println("\nWriting the result")

	if file, err := os.Create(*settings.VcardOutputFilePath); err != nil {
		return err
	} else {
		defer file.Close()
		if _, err := file.WriteString(vcardContent); err != nil {
			return err
		} else {
			fmt.Println("Vcard has been written to", *settings.VcardOutputFilePath)
		}
	}

	if err := qrcode.WriteColorFile(vcardContent, qrcode.Medium, 256, *settings.BackgroundColor, *settings.ForegroundColor, *settings.QrCodeOutputFilePath); err != nil {
		return err
	} else {
		fmt.Println("QR code has been written to", *settings.QrCodeOutputFilePath)
	}

	return nil
}

func finalize(err error) {
	if err == promptui.ErrInterrupt {
		fmt.Println("You stopped via CTRL-C")
	} else if err != nil {
		fmt.Println(err)
	}
	fmt.Println("👋")
}

func main() {
	var err error
	var s *settings.Settings

	defer func() {
		finalize(err)
	}()

	if s, err = settings.PrepareSettings(); err != nil {
		return
	}

	fmt.Println("\nPreparing a QR code from a vcard")
	fmt.Println("You get a list of options by starting the program in the form: qrvc -h")

	err = run(s)
}

// TODO format cli help
// TODO test
