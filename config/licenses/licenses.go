package licenses

import (
	"encoding/json"
	"fmt"
	"github.com/jimschubert/ossify/model"
	"io/ioutil"
	"os"
	"path"
)

func Load() (*model.Licenses, error) {
	// should this be configurable?
	location := "data/licenses/licenses.json"
	var licenses *model.Licenses
	bytes, err := ioutil.ReadFile(location)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(bytes, &licenses)
	return licenses, err
}

func PrintLicenseText(id string, customTemplateLocation string) error {
	// should this be configurable?
	location := path.Join("data/licenses/texts/plain/", id)
	customLocation := path.Join(customTemplateLocation, id)

	var useCustom bool
	// user-defined license templates take precedence over built-ins.
	if _, customErr := os.Stat(customLocation); os.IsNotExist(customErr) {
		if _, err := os.Stat(location); os.IsNotExist(err) {
			return err
		} else {
			useCustom = false
		}
	} else {
		useCustom = true
	}

	var b []byte
	if useCustom {
		customContent, err := ioutil.ReadFile(customLocation)
		if err != nil {
			return err
		}
		b = customContent
	} else {
		embeddedContent, err := ioutil.ReadFile(location)
		if err != nil {
			return err
		}
		b = embeddedContent
	}

	str := string(b)
	fmt.Println(str)

	return nil
}
