package config

import (
	"log"
	"os"

	"gopkg.in/alecthomas/kingpin.v2"
)

type Config struct {
	BaseDir           string
	ScoDir            string
	ScoConfigFile     string
	ScoDirPermissions os.FileMode
}

var (
	repoPath = kingpin.Flag("path", "Path to reposiory to watch.").Default(".").String()
)

func New() *Config {
	kingpin.Parse()

	config := Config{
		BaseDir:           *repoPath,
		ScoDir:            ".sco",
		ScoConfigFile:     ".sco.json",
		ScoDirPermissions: 0700,
	}

	// create sco internal directory
	if _, err := os.Stat(config.ScoDir); os.IsNotExist(err) {
		errr := os.MkdirAll(config.ScoDir, config.ScoDirPermissions)
		if errr != nil {
			log.Fatal(errr)
		}
	}

	return &config
}
