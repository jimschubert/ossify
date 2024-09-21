package config

import (
	"encoding/json"
	"log"
	"os"
	"path/filepath"
)

type Config struct {
	LicensePath    string `json:"licensePath"`
	ConventionPath string `json:"conventionPath"`
}

var defaultConfig Config

var Version = "0.1"
var Commit = "n/a"
var Date = "n/a"

var configName = ".config/ossify/settings.json"

// noinspection GoNameStartsWithPackageName
var ConfigManager *Manager

type LoadConfig func() (*Config, error)
type SaveConfig func(config *Config) error

type Manager struct {
	Load LoadConfig
	Save SaveConfig
}

func fullConfigPath(c string) (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	fullPath := filepath.Join(home, c)
	_, err = os.Stat(fullPath)
	if os.IsNotExist(err) {
		err = os.MkdirAll(filepath.Dir(fullPath), 0755)
		if err != nil {
			return "", err
		}
	}

	return fullPath, nil
}

func loadConfig() (*Config, error) {
	var c Config
	fullPath, err := fullConfigPath(configName)
	if err != nil {
		return nil, err
	}
	content, err := os.ReadFile(fullPath)
	if os.IsNotExist(err) {
		c = defaultConfig
		if err := saveConfig(&c); err != nil {
			return nil, err
		} else {
			return &c, nil
		}
	}

	if err = json.Unmarshal(content, &c); err != nil {
		return nil, err
	}

	return &c, nil
}

func saveConfig(config *Config) error {
	fullPath, err := fullConfigPath(configName)
	if err != nil {
		return err
	}

	content, err := json.Marshal(config)
	if err != nil {
		return err
	}

	return os.WriteFile(fullPath, content, 0600)
}

func init() {
	// TODO: make this configurable and cross-platform
	licensePath, err := fullConfigPath(".config/ossify/licenses")
	if err != nil {
		log.Fatal("Failed to create configuration path(s).")
	}
	conventionsPath, err := fullConfigPath(".config/ossify/conventions")
	if err != nil {
		log.Fatal("Failed to create configuration path(s).")
	}
	defaultConfig = Config{
		LicensePath:    licensePath,
		ConventionPath: conventionsPath,
	}
	ConfigManager = &Manager{
		Load: loadConfig,
		Save: saveConfig,
	}
	if _, err := ConfigManager.Load(); err != nil {
		panic(err)
	}
}
