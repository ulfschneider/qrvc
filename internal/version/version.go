package version

import (
	"embed"
	"path"
	"regexp"
	"strings"

	"github.com/pkg/errors"

	cdx "github.com/CycloneDX/cyclonedx-go"
	purl "github.com/package-url/packageurl-go"
)

var Version, _ = loadEmbeddedVersion()
var bom, _ = loadEmbeddedSBOM()
var BOM, _ = sprintf(bom)

//go:embed generated/*
var generated embed.FS

//version

const versionPath = "generated/version.txt"

func loadEmbeddedVersion() (string, error) {
	f, err := generated.Open(versionPath)
	if err != nil {
		return "", errors.Wrap(err, "Could not open generated version file")
	}
	defer f.Close()

	version, err := generated.ReadFile(versionPath)
	if err != nil {
		return "", errors.Wrap(err, "Could not read generated version file")
	}
	return string(version), nil
}

//sbom

const sbomPath = "generated/sbom.json"
const licensesPath = "generated/licenses"

var licenseRegex = regexp.MustCompile(`(?i)^(license|licence|copying|notice|readme)(\.[^.]+)?$`)

func loadEmbeddedSBOM() (*cdx.BOM, error) {

	f, err := generated.Open(sbomPath)
	if err != nil {
		return nil, errors.Wrap(err, "Failure when trying to access the SBOM")
	}
	defer f.Close()

	decoder := cdx.NewBOMDecoder(f, cdx.BOMFileFormatJSON)

	var bom cdx.BOM

	if err := decoder.Decode(&bom); err != nil {
		return nil, errors.Wrap(err, "Failure when decoding the SBOM")
	}

	embeddedLicenses := map[string]string{}
	if err := loadEmbeddedLicenses("", &embeddedLicenses); err != nil {
		return nil, errors.Wrap(err, "Failure when loading embedded licenses")
	}

	if err := injectLicenseText(&bom, &embeddedLicenses); err != nil {
		return nil, err
	}

	return &bom, nil
}

func sprintf(bom *cdx.BOM) (string, error) {

	if bom == nil {
		return "", errors.New("Given bom is nil")
	}

	var sb strings.Builder
	enc := cdx.NewBOMEncoder(&sb, cdx.BOMFileFormatJSON)
	enc.SetPretty(true) // pretty-print with indentation
	if err := enc.Encode(bom); err != nil {
		return "", err
	} else {
		return sb.String(), nil
	}
}

func loadEmbeddedLicenses(dir string, licenses *map[string]string) error {

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
			if err := loadEmbeddedLicenses(moduleDir, licenses); err != nil {
				return err
			}
		}

	}
	return nil
}

func injectLicenseText(bom *cdx.BOM, licenseMap *map[string]string) error {
	if bom.Components == nil {
		return nil
	}

	for i := range *bom.Components {
		c := &(*bom.Components)[i]

		module := extractModulePath(c)

		if module == "" {
			continue
		}

		if text, ok := (*licenseMap)[module]; ok {
			// Licenses field
			id, err := extractLicenseId(c)
			if err != nil {
				return err
			}

			c.Licenses = &cdx.Licenses{
				{
					License: &cdx.License{
						ID: id,
						Text: &cdx.AttachedText{
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

func extractLicenseId(c *cdx.Component) (string, error) {

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

func extractModulePath(c *cdx.Component) string {
	if c.PackageURL != "" {

		purl, err := purl.FromString(c.PackageURL)

		if err == nil {
			return path.Join(purl.Namespace, purl.Name) // this is the module path
		}
	}
	return c.Name
}
