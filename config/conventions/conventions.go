package conventions

import (
	"encoding/json"
	"github.com/jimschubert/ossify/config"
	"github.com/jimschubert/ossify/model"
	"github.com/pkg/errors"
	"io/ioutil"
	"path"
)

var ConventionPath = config.DefaultConfig.ConventionPath
func Load() (*[]model.Convention, error) {
	if ConventionPath == "" {
		return nil, errors.New("Invalid convention path.")
	}

	var conventions = make([]model.Convention, 2)
	copy(conventions, DefaultConventions)

	var files, err = ioutil.ReadDir(ConventionPath)
	if err != nil {
		return nil, err
	}

	for _, file := range files {
		if !file.IsDir() {
			var convention model.Convention

			var bytes []byte
			bytes, err = ioutil.ReadFile(path.Join(ConventionPath, file.Name()))
			if err != nil {
				continue
			}

			err = json.Unmarshal(bytes, &convention)
			conventions = append(conventions, convention)
		}
	}

	// returns all available conventions and last known error
	return &conventions, err
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
		{model.Required, model.Directory, "README.md"},
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
		{model.Required, model.Directory, "README.md"},
	},
}

//{
//	"name": "Standard",
//  "rules" : [
//		{ "level": "optional", "type": "directory", "value": "src" },
//		{ "level": "required", "type": "file", "value": "CONTRIBUTING.md" }
//    ]
//}
