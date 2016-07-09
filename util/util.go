package util

import (
	"io"
	"io/ioutil"
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
		log.Println("Creating directory", file)
		if err != nil {
			log.Fatal(err)
		}
	}
}

func (util *Util) WriteFile(content *[]byte, folders ...string) {
	path := filepath.Join(util.config.ScoDir, filepath.Join(folders...))
	if _, err := os.Stat(path); os.IsNotExist(err) {
		err = ioutil.WriteFile(path, *content, util.config.ScoDirPermissions)
		log.Println("Creating file", path)
		if err != nil {
			log.Fatal(err)
		}
	} else {
		err = ioutil.WriteFile(path, *content, util.config.ScoDirPermissions)
		log.Println("Overwrite file", path)
		if err != nil {
			log.Fatal(err)
		}
	}
}

func (util *Util) CopyFile(src string, destfolders ...string) error {
	dst := filepath.Join(util.config.ScoDir, filepath.Join(destfolders...))
	s, err := os.Open(src)
	if err != nil {
		return err
	}
	// no need to check errors on read only file, we already got everything
	// we need from the filesystem, so nothing can go wrong now.
	defer s.Close()
	d, err := os.Create(dst)
	if err != nil {
		return err
	}
	if _, err := io.Copy(d, s); err != nil {
		d.Close()
		return err
	}
	return d.Close()
}
