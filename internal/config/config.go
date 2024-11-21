package config

import (
	"encoding/json"
	"fmt"
	"os"
)

const configFileName = ".gatorconfig.json"

type Config struct {
	DbUrl           string `json:"db_url"`
	CurrentUserName string `json:"current_user_name"`
}

func Read() (Config, error) {
	path, err := getConfigFilePath()
	if err != nil {
		return Config{}, err
	}

	file, err := os.Open(path)
	if err != nil {
		return Config{}, err
	}
	defer file.Close()

	config := Config{}
	err = json.NewDecoder(file).Decode(&config)
	if err != nil {
		return Config{}, err
	}

	return config, nil
}

func (c *Config) SetUser(name string) error {
	c.CurrentUserName = name
	err := write(*c)
	if err != nil {
		return err
	}
	return nil
}

func getConfigFilePath() (string, error) {
	path, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%s/%s", path, configFileName), nil
}

func write(cfg Config) error {
	path, err := getConfigFilePath()
	if err != nil {
		return err
	}

	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	err = json.NewEncoder(file).Encode(cfg)
	if err != nil {
		return err
	}

	return nil
}
