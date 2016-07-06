package util

import (
	"log"
	"os"
	"path/filepath"

	"github.com/floriangrundig/scofw/config"
)

type Util struct {
	config *config.Config
}

func New(config *config.Config) *Util {
	return &Util{
		config: config,
	}
}

// func (util *Util) Get

func (util *Util) CreateScoFolder(folders ...string) {
	file := filepath.Join(util.config.ScoDir, filepath.Join(folders...))
	if _, err := os.Stat(file); os.IsNotExist(err) {
		err = os.MkdirAll(file, util.config.ScoDirPermissions)
		if err != nil {
			log.Fatal(err)
		}
	}
}
