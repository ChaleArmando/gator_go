package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

type Config struct {
	DbURL           string `json:"db_url"`
	CurrentUserName string `json:"current_user_name"`
}

const configFileName = "workspace/github.com/ChaleArmando/gator_go/.gatorconfig.json"

func Read() Config {
	var conf Config

	filePath, err := getConfigFilePath()
	if err != nil {
		return conf
	}

	file, err := os.Open(filePath)
	if err != nil {
		fmt.Println(err)
		return conf
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	decoder.Decode(&conf)
	return conf
}

func (conf *Config) SetUser(userName string) error {
	conf.CurrentUserName = userName
	return write(*conf)
}

func write(conf Config) error {
	filePath, err := getConfigFilePath()
	if err != nil {
		return err
	}

	jsonData, err := json.MarshalIndent(conf, "", "  ")
	if err != nil {
		return err
	}

	err = os.WriteFile(filePath, jsonData, 0664)
	if err != nil {
		return err
	}

	return nil
}

func getConfigFilePath() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	filePath := filepath.Join(homeDir, configFileName)
	return filePath, err
}
