package licenses

import (
	"embed"
	_ "embed"
	"encoding/json"
	"fmt"
	"os"
	"path"

	"github.com/jimschubert/ossify/internal/model"
)

//go:embed data
var licenseContent embed.FS

func Load() (*model.Licenses, error) {
	var licenses *model.Licenses
	bytes, err := licenseContent.ReadFile("data/licenses.json")
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(bytes, &licenses)
	return licenses, err
}

func PrintLicenseText(id string, customTemplateLocation string) error {
	location := path.Join("data/texts/plain/", id)
	customLocation := path.Join(customTemplateLocation, id)

	var b []byte
	// user-defined license templates take precedence over built-ins.
	if _, customErr := os.Stat(customLocation); os.IsNotExist(customErr) {
		embeddedContent, err := licenseContent.ReadFile(location)
		if err != nil {
			return err
		}
		b = embeddedContent
	} else {
		customContent, err := os.ReadFile(customLocation)
		if err != nil {
			return err
		}
		b = customContent
	}

	str := string(b)
	fmt.Println(str)

	return nil
}
