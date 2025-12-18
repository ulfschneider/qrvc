package main_test

import (
	"image"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	qrcodec "github.com/ulfschneider/qrvc/internal/adapters/codec/qr"
	testutil "github.com/ulfschneider/qrvc/internal/test/util"
)

func TestSmoke(t *testing.T) {

	testFolder := t.TempDir()

	cmd := exec.Command("go", "build", "-o", filepath.Join(testFolder, "qrvc"), ".")
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("Build failed: %v\n%s", err, out)
	}

	card := testutil.CreateCard()
	content := testutil.EncodeCard(card)
	os.WriteFile(filepath.Join(testFolder, "vcard.vcf"), content, fs.ModePerm)

	//create qr code with default settings
	cmd = exec.Command(filepath.Join(testFolder, "qrvc"), "-s", "-i", filepath.Join(testFolder, "vcard.vcf"), "-o", filepath.Join(testFolder, "result"))
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("Smoketest failed: %v\n%s", err, out)
	}

	expectedVcf, err := os.ReadFile(filepath.Join(testFolder, "vcard.vcf"))
	assert.NoError(t, err)
	assert.NotEmpty(t, string(expectedVcf))

	actualVcf, err := os.ReadFile(filepath.Join(testFolder, "result.vcf"))
	assert.NoError(t, err)
	assert.NotEmpty(t, string(actualVcf))

	assert.Equal(t, testutil.NormalizeNewLines(string(expectedVcf)), testutil.NormalizeNewLines(string(actualVcf)))

	file, err := os.Open(filepath.Join(testFolder, "result.png"))
	assert.NoError(t, err)
	defer file.Close()

	//load default test settings
	settings := testutil.LoadTestSettings()

	codec := qrcodec.NewCodec()
	expectedQR, err := codec.Encode(card, settings.App.QRSettings)
	assert.NoError(t, err)
	assert.NotEmpty(t, expectedQR)
	actualQR, format, err := image.Decode(file)
	assert.NoError(t, err)
	assert.Equal(t, "png", format)
	assert.Equal(t, testutil.ToRGBA(expectedQR), testutil.ToRGBA(actualQR))

}
