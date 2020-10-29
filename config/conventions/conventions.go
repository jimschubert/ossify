package conventions

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"os"
	"path"

	"github.com/jimschubert/ossify/config"
	"github.com/jimschubert/ossify/model"
)

func Load() (*[]model.Convention, error) {
	c, e := config.ConfigManager.Load()
	if e != nil {
		return nil, e
	}
	conventionPath := c.ConventionPath
	if conventionPath == "" {
		return nil, errors.New("invalid convention path")
	}

	var conventions = make([]model.Convention, 2)
	copy(conventions, DefaultConventions)

	if _, err := os.Stat(conventionPath); os.IsExist(err) {
		var files, err = ioutil.ReadDir(conventionPath)
		if err != nil {
			return nil, err
		}

		for _, file := range files {
			if !file.IsDir() {
				var convention model.Convention

				var bytes []byte
				bytes, err = ioutil.ReadFile(path.Join(conventionPath, file.Name()))
				if err != nil {
					continue
				}

				err = json.Unmarshal(bytes, &convention)
				conventions = append(conventions, convention)
			}
		}
	}

	// returns all available conventions and last known error
	return &conventions, nil
}

var DefaultConventions = []model.Convention{
	StandardDistributionConvention,
	GoConvention,
}

var StandardDistributionConvention = model.Convention{
	Name: "Standard Distribution",
	Rules: []model.Rule{
		{model.Required, model.Directory, "dist"},
		{model.Required, model.Directory, "docs"},
		{model.Optional, model.Directory, "lib"},
		{model.Required, model.Directory, "src"},
		{model.Required, model.Directory, "test"},
		{model.Optional, model.Directory, "tools"},
		{model.Required, model.File, "LICENSE"},
		{model.Required, model.File, "README.md"},
	},
}

var GoConvention = model.Convention{
	Name: "Go",
	Rules: []model.Rule{
		{model.Optional, model.Directory, "configs"},
		{model.Optional, model.Directory, "init"},
		{model.Optional, model.Directory, "scripts"},
		{model.Required, model.Directory, "docs"},
		{model.Optional, model.Directory, "tools"},
		{model.Optional, model.Directory, "deployments"},
		{model.Optional, model.Directory, "test"},
		{model.Optional, model.Directory, "build"},
		{model.Optional, model.Directory, "vendor"},
		{model.Prohibited, model.Directory, "src"},
		{model.Required, model.File, "LICENSE"},
		{model.Required, model.File, "README.md"},
	},
}

//{
//	"name": "Standard",
//  "rules" : [
//		{ "level": "optional", "type": "directory", "value": "src" },
//		{ "level": "required", "type": "file", "value": "CONTRIBUTING.md" }
//    ]
//}
