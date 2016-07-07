package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"gopkg.in/alecthomas/kingpin.v2"
)

type ScoRuntimeData struct {
	GitCommits map[string]string
}

type Config struct {
	BaseDir           string
	ScoDir            string
	ScoConfigFile     string
	ScoDirPermissions os.FileMode
	ScoRuntimeData
}

var (
	repoPath = kingpin.Flag("path", "Path to reposiory to watch.").Default(".").String()
)

func New() *Config {
	kingpin.Parse()

	config := Config{
		BaseDir:           *repoPath,
		ScoDir:            ".sco",
		ScoConfigFile:     "config.json",
		ScoDirPermissions: 0700,
	}

	// create sco internal directory
	if !scoConfigFileExists(config.ScoDir) {
		errr := os.MkdirAll(config.ScoDir, config.ScoDirPermissions)
		if errr != nil {
			log.Fatal(errr)
		}
	}

	if !scoConfigFileExists(config.ScoDir, config.ScoConfigFile) {
		initializeScoConfigfile(config)
	}

	scoRuntimeData := loadScoConfigFileContent(config)

	config.setScoRuntimeData(scoRuntimeData)
	return &config
}

func (config *Config) Persist() {
	file := filepath.Join(config.ScoDir, config.ScoConfigFile)

	b, err := json.MarshalIndent(&config.ScoRuntimeData, "", "  ")
	if err != nil {
		fmt.Println("error:", err)
	}

	err = ioutil.WriteFile(file, b, config.ScoDirPermissions)
	if err != nil {
		fmt.Println("error:", err)
	}

}

func (config *Config) setScoRuntimeData(rt ScoRuntimeData) {
	config.ScoRuntimeData = rt
}

func scoConfigFileExists(pathElements ...string) bool {
	file := filepath.Join(pathElements...)
	if _, err := os.Stat(file); os.IsNotExist(err) {
		return false
	}
	return true
}

func loadScoConfigFileContent(config Config) ScoRuntimeData {
	filePath := filepath.Join(config.ScoDir, config.ScoConfigFile)
	file, err := os.Open(filePath)
	if err != nil {
		log.Println("error:", err)
	}

	decoder := json.NewDecoder(file)
	scoRuntimeData := ScoRuntimeData{}
	err = decoder.Decode(&scoRuntimeData)
	if err != nil {
		log.Println("error:", err)
	}
	return scoRuntimeData
}

func initializeScoConfigfile(config Config) {
	log.Printf("Creating %s\n", config.ScoConfigFile)
	file := filepath.Join(config.ScoDir, config.ScoConfigFile)
	scoRuntimeData := ScoRuntimeData{
		GitCommits: make(map[string]string),
	}

	b, err := json.MarshalIndent(scoRuntimeData, "", "  ")
	if err != nil {
		log.Println("error:", err)
	}

	err = ioutil.WriteFile(file, b, config.ScoDirPermissions)
	if err != nil {
		log.Println("error:", err)
	}

}
