package conventions

import (
	"encoding/json"
	"errors"
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
		var files, err = os.ReadDir(conventionPath)
		if err != nil {
			return nil, err
		}

		for _, file := range files {
			if !file.IsDir() {
				var convention model.Convention

				var bytes []byte
				bytes, err = os.ReadFile(path.Join(conventionPath, file.Name()))
				if err != nil {
					// TODO: warn
					continue
				}

				err = json.Unmarshal(bytes, &convention)
				if err != nil {
					// TODO: warn
					continue
				}
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
		{Level: model.Required, Type: model.Directory, Value: "dist"},
		{Level: model.Required, Type: model.Directory, Value: "docs"},
		{Level: model.Optional, Type: model.Directory, Value: "lib"},
		{Level: model.Required, Type: model.Directory, Value: "src"},
		{Level: model.Required, Type: model.Directory, Value: "test"},
		{Level: model.Optional, Type: model.Directory, Value: "tools"},
		{Level: model.Required, Type: model.File, Value: "LICENSE"},
		{Level: model.Required, Type: model.File, Value: "README.md"},
	},
}

var GoConvention = model.Convention{
	Name: "Go",
	Rules: []model.Rule{
		{Level: model.Optional, Type: model.Directory, Value: "configs"},
		{Level: model.Optional, Type: model.Directory, Value: "init"},
		{Level: model.Optional, Type: model.Directory, Value: "scripts"},
		{Level: model.Required, Type: model.Directory, Value: "docs"},
		{Level: model.Optional, Type: model.Directory, Value: "tools"},
		{Level: model.Optional, Type: model.Directory, Value: "deployments"},
		{Level: model.Optional, Type: model.Directory, Value: "test"},
		{Level: model.Optional, Type: model.Directory, Value: "build"},
		{Level: model.Optional, Type: model.Directory, Value: "vendor"},
		{Level: model.Prohibited, Type: model.Directory, Value: "src"},
		{Level: model.Required, Type: model.File, Value: "LICENSE"},
		{Level: model.Required, Type: model.File, Value: "README.md"},
	},
}

// {
//	"name": "Standard",
//  "rules" : [
//		{ "level": "optional", "type": "directory", "value": "src" },
//		{ "level": "required", "type": "file", "value": "CONTRIBUTING.md" }
//    ]
// }
