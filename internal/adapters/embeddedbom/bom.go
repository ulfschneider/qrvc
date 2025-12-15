package embeddedbom

import (
	"bytes"
	"embed"
	"path"
	"regexp"

	"github.com/CycloneDX/cyclonedx-go"
	"github.com/package-url/packageurl-go"
	"github.com/pkg/errors"
	"github.com/ulfschneider/qrvc/internal/adapters/clinotifier"
)

type BomProvider struct {
	userNotifier clinotifier.UserNotifier
}

func NewBomProvider() BomProvider {
	return BomProvider{userNotifier: clinotifier.NewUserNotifier()}
}

//go:embed generated/*
var generated embed.FS

// sbom
const sbomPath = "generated/bom.json"
const licensesPath = "generated/licenses"

var licenseRegex = regexp.MustCompile(`(?i)^(license|licence|copying|notice|readme)(\.[^.]+)?$`)

func (bp *BomProvider) Bom() (*cyclonedx.BOM, error) {

	f, err := generated.Open(sbomPath)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	decoder := cyclonedx.NewBOMDecoder(f, cyclonedx.BOMFileFormatJSON)

	var bom cyclonedx.BOM

	if err := decoder.Decode(&bom); err != nil {
		return nil, err
	}

	embeddedLicenses := map[string]string{}
	if err := bp.loadEmbeddedLicenses("", &embeddedLicenses); err != nil {
		return nil, err
	}

	if err := bp.injectLicenseText(&bom, &embeddedLicenses); err != nil {
		return nil, err
	}

	return &bom, nil
}

func (bp *BomProvider) MarshalToJSON() ([]byte, error) {

	bom, err := bp.Bom()

	if err != nil {
		return nil, err
	}

	var buffer bytes.Buffer
	if bom == nil {
		return buffer.Bytes(), errors.New("Given BOM is nil")
	}

	enc := cyclonedx.NewBOMEncoder(&buffer, cyclonedx.BOMFileFormatJSON)
	enc.SetPretty(true) // pretty-print with indentation
	if err := enc.Encode(bom); err != nil {
		return buffer.Bytes(), err
	} else {
		return buffer.Bytes(), nil
	}
}

func (bp *BomProvider) WriteBomJSON() error {
	json, err := bp.MarshalToJSON()

	if err != nil {
		return err
	}

	bp.userNotifier.NotifyLoud(string(json))
	return nil
}

func (bp *BomProvider) loadEmbeddedLicenses(dir string, licenses *map[string]string) error {

	entries, err := generated.ReadDir(path.Join(licensesPath, dir))
	if err != nil {
		return err
	}

	for _, entry := range entries {
		name := entry.Name()
		if !entry.IsDir() {
			//the entry is a file
			if licenseRegex.MatchString(name) {
				filePath := path.Join(licensesPath, dir, name)
				data, err := generated.ReadFile(filePath)
				if err != nil {
					return err
				}
				(*licenses)[dir] = string(data)
				return nil
			}
		} else {
			//the entry is a folder
			moduleDir := path.Join(dir, name)
			if err := bp.loadEmbeddedLicenses(moduleDir, licenses); err != nil {
				return err
			}
		}

	}
	return nil
}

func (bp *BomProvider) injectLicenseText(bom *cyclonedx.BOM, licenseMap *map[string]string) error {
	if bom.Components == nil {
		return nil
	}

	for i := range *bom.Components {
		c := &(*bom.Components)[i]

		module := bp.extractModulePath(c)

		if module == "" {
			continue
		}

		if text, ok := (*licenseMap)[module]; ok {
			// Licenses field
			id, err := bp.extractLicenseId(c)
			if err != nil {
				return err
			}

			c.Licenses = &cyclonedx.Licenses{
				{
					License: &cyclonedx.License{
						ID: id,
						Text: &cyclonedx.AttachedText{
							ContentType: "text/plain",
							Encoding:    "plain",
							Content:     text,
						}},
				},
			}
		}
	}
	return nil
}

func (bp *BomProvider) extractLicenseId(c *cyclonedx.Component) (string, error) {

	if c.Evidence != nil && c.Evidence.Licenses != nil {
		evidence := *c.Evidence.Licenses
		if len(evidence) > 1 {
			return "", errors.New("The component " + c.BOMRef + " contains more than one license")
		}
		for _, l := range evidence {
			if l.License != nil && l.License.ID != "" {
				return l.License.ID, nil
			}
		}
	}

	return "", nil
}

func (bp *BomProvider) extractModulePath(c *cyclonedx.Component) string {
	if c.PackageURL != "" {

		purl, err := packageurl.FromString(c.PackageURL)

		if err == nil {
			return path.Join(purl.Namespace, purl.Name) // this is the module path
		}
	}
	return c.Name
}
