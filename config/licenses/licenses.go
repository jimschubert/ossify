package licenses

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path"

	"github.com/gobuffalo/packr"

	"github.com/jimschubert/ossify/model"
)

func Load() (*model.Licenses, error) {
	// should this be configurable?
	box := packr.NewBox("../../data/licenses")
	var licenses *model.Licenses
	bytes, err := box.Find("licenses.json")
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(bytes, &licenses)
	return licenses, err
}

func PrintLicenseText(id string, customTemplateLocation string) error {
	// should this be configurable?
	box := packr.NewBox("../../data/licenses")
	location := path.Join("texts/plain/", id)
	customLocation := path.Join(customTemplateLocation, id)

	var b []byte
	// user-defined license templates take precedence over built-ins.
	if _, customErr := os.Stat(customLocation); os.IsNotExist(customErr) {
		embeddedContent, err := box.Find(location)
		if err != nil {
			return err
		}
		b = embeddedContent
	} else {
		customContent, err := ioutil.ReadFile(customLocation)
		if err != nil {
			return err
		}
		b = customContent
	}

	str := string(b)
	fmt.Println(str)

	return nil
}
