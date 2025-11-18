package main

import (
	"fmt"
	"os"
	"qrvc/internal/settings"
	"qrvc/internal/vcard"

	"github.com/skip2/go-qrcode"
)

func main() {

	settings := settings.MustPrepareSettings()

	fmt.Println("\nPreparing a QR code from a vcard")
	fmt.Println("You get a list of options by starting the program with the -h argument")

	vcardContent := vcard.MustPrepareVcard(settings)

	fmt.Println("\nWriting the result")

	if file, err := os.Create(*settings.VcardOutputFilePath); err != nil {
		panic(err)
	} else {
		defer file.Close()
		if _, err := file.WriteString(vcardContent); err != nil {
			panic(err)
		} else {
			fmt.Println("Vcard has been written to", *settings.VcardOutputFilePath)
		}
	}

	if err := qrcode.WriteColorFile(vcardContent, qrcode.Medium, 256, *settings.BackgroundColor, *settings.ForegroundColor, *settings.QrCodeOutputFilePath); err != nil {
		panic(err)
	} else {
		fmt.Println("QR code has been written to", *settings.QrCodeOutputFilePath)
	}

	fmt.Println("👋")
}

// TODO reuse already given input
// TODO format cli help
// TODO configure logging and improve panic logging
// TODO audit dependencies
