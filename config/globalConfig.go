package config

import (
	"encoding/json"
	standardLogger "log"
	"os"
)

type ProjectConfig struct {
	Name        string   `json:"name"`
	Path        string   `json:"path"`
	IgnoreFiles []string `json:"ignoreFiles"`
	ScoDir      string   `json:"scoDir"`
}

type GlobalConfig struct {
	Projects       []ProjectConfig `json:"projects"`
	IgnorePatterns []string        `json:"ignorePatterns"`
}

func ParseGlobalConfig(homedir string) *GlobalConfig {

	config := GlobalConfig{}
	file, err := os.Open(*ConfigFile)
	if err != nil {
		standardLogger.Fatal("Error while opening configuration file:", err)
	}

	decoder := json.NewDecoder(file)

	err = decoder.Decode(&config)
	if err != nil {
		standardLogger.Fatal("error:", err)
	}

	return &config
}
