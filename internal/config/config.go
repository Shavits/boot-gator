package config

import (
	"encoding/json"
	"fmt"
	"os"
)

const configFileName = ".gatorconfig.json"

type Config struct {
	DbURL           string `json:"db_url"`
	CurrentUserName string `json:"current_user_name"`
}

func Read() (Config, error) {
	configFilePath, err := getConfigFilePath()
	if err != nil {
		return Config{}, err
	}
	configFile, err := os.Open(configFilePath)
	if err != nil {
		return Config{}, err
	}
	defer configFile.Close()

	dec := json.NewDecoder(configFile)
	var config Config
	if err := dec.Decode(&config); err != nil {
		return Config{}, err
	}
	return config, nil

}

func SetUser(userName string) error {
	config, err := Read()
	if err != nil {
		return err
	}
	config.CurrentUserName = userName
	return write(config)
}

func getConfigFilePath() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("failed to get config file path: %s", err)
	}

	return fmt.Sprintf("%s/%s", homeDir, configFileName), nil

}

func write(cfg Config) error {
	configFilePath, err := getConfigFilePath()
	if err != nil {
		return err
	}

	data, err := json.Marshal(cfg)
	if err != nil {
		return fmt.Errorf("error marshaling json: %s", err)
	}

	if err := os.WriteFile(configFilePath, data, 0600); err != nil {
		return fmt.Errorf("failed to write file: %s", err)
	}

	return nil

}
