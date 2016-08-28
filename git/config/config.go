package config

import (
	"encoding/json"
	"io/ioutil"
	log_ "log"
	"os"
	"path/filepath"
	"time"

	"github.com/floriangrundig/scofw/config"
)

var (
	log *log_.Logger // our logger
)

type FileModificationInfo struct {
	Op   uint32
	Date time.Time
}

type SessionData struct {
	FirstModificationDate time.Time
	Modifications         map[string][]FileModificationInfo
}

type GitRuntimeData struct {
	CurrentScoSession string
	GitCommits        map[string]string
	Sessions          map[string]SessionData
}

type Config struct {
	scoConfig          *config.Config
	gitRuntimeDataFile string
	GitRuntimeData
	CurrentScoSession string
}

func New(config *config.Config) *Config {
	log = config.Logger
	gitConfig := Config{
		scoConfig:          config,
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
		log.Println("error:", err)
	}

	err = ioutil.WriteFile(config.gitRuntimeDataFile, b, config.scoConfig.ScoDirPermissions)
	if err != nil {
		log.Println("error:", err)
	}

}

func (config *Config) setGitRuntimeData(rt GitRuntimeData) {
	config.GitRuntimeData = rt

}

func (config *Config) SetCurrentScoSession(session string) {
	config.CurrentScoSession = session                // in memory
	config.GitRuntimeData.CurrentScoSession = session // in Json
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
		Sessions:   make(map[string]SessionData),
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
