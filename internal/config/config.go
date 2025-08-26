package config

import (
	"encoding/json"
	"fmt"
	"log"
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
		log.Printf("Config.Read error: failed to get config file path: %v\n", err)
		return Config{}, err
	}
	data, err := os.ReadFile(configFilePath)
	if err != nil {
		log.Printf("Config.Read error: failed to read config file: %v\n", err)
		return Config{}, err
	}

	config := Config{}
	err = json.Unmarshal(data, &config)
	if err != nil {
		log.Printf("Config.Read error: failed to unmarshal config: %v\n", err)
		return Config{}, err
	}
	return config, nil
}

func (c *Config) SetUser(username string) error {
	c.CurrentUserName = username
	err := write(c)
	if err != nil {
		log.Printf("Config.SetUser error: failed to write config: %v\n", err)
	}
	return err
}

func (c *Config) String() string {
	return fmt.Sprintf("Config{\n DbUrl: %s\n CurrentUserName: %s\n}", c.DbUrl, c.CurrentUserName)
}

func write(cfg *Config) error {
	jsonData, err := json.Marshal(cfg)
	if err != nil {
		log.Printf("Config.write error: failed to marshal config: %v\n", err)
		return err
	}
	configFilePath, err := getConfigFile()
	if err != nil {
		log.Printf("Config.write error: failed to get config file path: %v\n", err)
		return err
	}
	err = os.WriteFile(configFilePath, jsonData, 0644)
	if err != nil {
		log.Printf("Config.write error: failed to write config file: %v\n", err)
		return err
	}
	return nil
}

func getConfigFile() (string, error) {
	homedir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(homedir, configFileName), nil
}
