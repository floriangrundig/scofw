package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"github.com/floriangrundig/scofw/config"
	"github.com/floriangrundig/scofw/util"
)

type GitRuntimeData struct {
	GitCommits map[string]string
}

type Config struct {
	scoConfig          *config.Config
	util               *util.Util
	gitRuntimeDataFile string
	GitRuntimeData
}

func New(config *config.Config, util *util.Util) *Config {

	gitConfig := Config{
		scoConfig:          config,
		util:               util,
		gitRuntimeDataFile: filepath.Join(config.ScoDir, "commits_sessions.json"),
	}

	if !gitConfig.gitRuntimeDataFileExists() {
		gitConfig.initializeGitRuntimeDataFile()
	}

	gitRuntimeData := gitConfig.loadGitRuntimeDataFileContent()

	gitConfig.setGitRuntimeData(gitRuntimeData)

	return &gitConfig
}

func (config *Config) Persist() {
	b, err := json.MarshalIndent(&config.GitRuntimeData, "", "  ")
	if err != nil {
		fmt.Println("error:", err)
	}

	err = ioutil.WriteFile(config.gitRuntimeDataFile, b, config.scoConfig.ScoDirPermissions)
	if err != nil {
		fmt.Println("error:", err)
	}

}

func (config *Config) setGitRuntimeData(rt GitRuntimeData) {
	config.GitRuntimeData = rt
}

func (config *Config) gitRuntimeDataFileExists() bool {
	if _, err := os.Stat(config.gitRuntimeDataFile); os.IsNotExist(err) {
		return false
	}
	return true
}

func (config *Config) loadGitRuntimeDataFileContent() GitRuntimeData {
	file, err := os.Open(config.gitRuntimeDataFile)
	if err != nil {
		log.Fatal("error:", err)
	}

	decoder := json.NewDecoder(file)
	gitRuntimeData := GitRuntimeData{}

	err = decoder.Decode(&gitRuntimeData)
	if err != nil {
		log.Fatal("error:", err)
	}
	return gitRuntimeData
}

func (config *Config) initializeGitRuntimeDataFile() {
	log.Printf("Creating %s\n", config.gitRuntimeDataFile)
	file := config.gitRuntimeDataFile

	gitRuntimeData := GitRuntimeData{
		GitCommits: make(map[string]string),
	}

	b, err := json.MarshalIndent(gitRuntimeData, "", "  ")
	if err != nil {
		log.Fatal("error:", err)
	}

	err = ioutil.WriteFile(file, b, config.scoConfig.ScoDirPermissions)
	if err != nil {
		log.Fatal("error:", err)
	}

}
