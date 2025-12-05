package out

import (
	"fmt"
	"image/png"
	"os"

	"github.com/ulfschneider/qrvc/internal/cli"
	"github.com/ulfschneider/qrvc/internal/settings"

	"github.com/skip2/go-qrcode"
)

func StoreResults(vcardContent string, settings *settings.Settings) error {

	if file, err := os.Create(*settings.VCardOutputFilePath); err != nil {
		return err
	} else {
		defer file.Close()
		if _, err := file.WriteString(vcardContent); err != nil {
			return err
		} else {
			fmt.Println("The vCard has been written to", cli.SprintValue(*settings.VCardOutputFilePath))
		}
	}

	q, err := qrcode.New(vcardContent, qrcode.Low)
	if err != nil {
		return err
	}

	q.DisableBorder = !*settings.Border
	q.ForegroundColor = *settings.ForegroundColor
	q.BackgroundColor = *settings.BackgroundColor
	q.Level = qrcode.Low

	img := q.Image(*settings.Size)

	if file, err := os.Create(*settings.QRCodeOutputFilePath); err != nil {
		return err
	} else {
		defer file.Close()
		if err := png.Encode(file, img); err != nil {
			return err
		} else {
			fmt.Println("The QR code has been written to", cli.SprintValue(*settings.QRCodeOutputFilePath))
		}
	}
	return nil
}
