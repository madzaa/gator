package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

type Config struct {
	DbUrl           string `json:"db_url"`
	CurrentUserName string `json:"current_user_name"`
}

const configFileName = ".gatorconfig.json"

func Read() (Config, error) {
	configFilePath, err := getConfigFile()
	if err != nil {
		return Config{}, err
	}
	data, err := os.ReadFile(configFilePath)
	if err != nil {
		return Config{}, err
	}

	config := Config{}
	err = json.Unmarshal(data, &config)
	if err != nil {
		return Config{}, err
	}
	return config, nil
}

func (c *Config) SetUser(username string) error {
	c.CurrentUserName = username
	return write(c)
}

func (c *Config) String() string {
	return fmt.Sprintf("Config{\n DbUrl: %s\n CurrentUserName: %s\n}", c.DbUrl, c.CurrentUserName)
}

func write(cfg *Config) error {
	jsonData, err := json.Marshal(cfg)
	if err != nil {
		return err
	}
	configFilePath, err := getConfigFile()
	if err != nil {
		return err
	}
	return os.WriteFile(configFilePath, jsonData, 0644)
}

func getConfigFile() (string, error) {
	homedir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(homedir, configFileName), nil
}
