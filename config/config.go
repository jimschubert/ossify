package config

import (
	"encoding/json"
	"github.com/pkg/errors"
	"github.com/shibukawa/configdir"
	"path"
	"path/filepath"
)

type Config struct {
	LicensePath string `json:"licensePath"`
	ConventionPath string `json:"conventionPath"`
}

var DefaultConfig = Config {
	LicensePath: "",
	ConventionPath: "",
}

var configName = "settings.json"
//noinspection GoNameStartsWithPackageName
var ConfigManager *Manager

type LoadConfig func() (*Config, error)
type SaveConfig func(config *Config) error

type Manager struct {
	Load LoadConfig
	Save SaveConfig
}

func applicationConfigDir() configdir.ConfigDir {
	return configdir.New("", "ossify")
}

func loadConfig() (*Config, error) {
	var config Config
	configDirs := applicationConfigDir()
	// optional: local path has the highest priority
	configDirs.LocalPath, _ = filepath.Abs(".")
	existing := configDirs.QueryFolderContainsFile(configName)

	if existing != nil {
		data, _ := existing.ReadFile(configName)
		err := json.Unmarshal(data, &config)
		if err != nil {
			return nil, err
		}
	} else {
		config = DefaultConfig
		config.LicensePath = path.Join(configDirs.QueryCacheFolder().Path, "licenses")
		config.ConventionPath = path.Join(configDirs.QueryCacheFolder().Path, "conventions")
		err := saveConfig(&config)
		return nil, err
	}

	return &config, nil
}

func saveConfig(config *Config) error {
	configDirs := applicationConfigDir()
	configDirs.LocalPath, _ = filepath.Abs(".")
	existing := configDirs.QueryFolderContainsFile(configName)
	data, err := json.MarshalIndent(&config, "", "   ")
	if err != nil {
		return err
	}

	if existing != nil {
		if err = existing.WriteFile(configName, data); err != nil {
			return err
		}
	} else {
		folders := configDirs.QueryFolders(configdir.Global)
		if len(folders) > 0 {
			folder := folders[0]
			if err = folder.WriteFile(configName, data); err != nil {
				return err
			}
		} else {
			return errors.New("no configuration folders available, cannot proceed")
		}
	}
	return nil
}

func init(){
	ConfigManager = &Manager{
		loadConfig,
		saveConfig,
	}
	if _, err := ConfigManager.Load(); err != nil {
		panic(err)
	}
}
