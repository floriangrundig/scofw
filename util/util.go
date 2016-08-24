package util

import (
	"io"
	"io/ioutil"
	log_ "log"
	"os"
	"path/filepath"

	"github.com/floriangrundig/scofw/config"
)

var (
	log *log_.Logger
)

type Util struct {
	config *config.Config
}

func CreateInternalScoDirectory(path string, permissions os.FileMode) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		errr := os.MkdirAll(path, permissions)
		if errr != nil {
			log.Fatal(errr)
		}
		log_.Printf("Created: %s", path)
	}
	// create log directory
	logDir := filepath.Join(path, "logs")
	if _, err := os.Stat(logDir); os.IsNotExist(err) {
		errr := os.MkdirAll(logDir, permissions)
		if errr != nil {
			log.Fatal(errr)
		}
		log_.Printf("Created: %s", logDir)
	}
}

func New(config *config.Config) *Util {
	log = config.Logger
	return &Util{
		config: config,
	}
}

// TODO rename - since it creates a sco internal folder
func (util *Util) CreateScoFolder(pathElements ...string) {
	file := filepath.Join(util.config.ScoDir, filepath.Join(pathElements...))
	if _, err := os.Stat(file); os.IsNotExist(err) {
		err = os.MkdirAll(file, util.config.ScoDirPermissions)
		log.Println("Creating directory", file)
		if err != nil {
			log.Fatal(err)
		}
	}
}

func (util *Util) ReadScoFile(pathElements ...string) (*[]byte, error) {
	file := filepath.Join(util.config.ScoDir, filepath.Join(pathElements...))

	content, err := ioutil.ReadFile(file)
	return &content, err
}

// TODO rename - since it is removing only sco internal files
func (util *Util) RemoveFile(pathElements ...string) {
	path := filepath.Join(util.config.ScoDir, filepath.Join(pathElements...))
	if _, err := os.Stat(path); os.IsNotExist(err) {
	} else {
		err = os.RemoveAll(path)
		if err != nil {
			log.Println("ERROR", err)
		}

	}

}

// TODO rename - since it is only writing sco internal files
func (util *Util) WriteFile(content *[]byte, pathElements ...string) {
	path := filepath.Join(util.config.ScoDir, filepath.Join(pathElements...))
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

func (util *Util) CopyFile(src string, destpathElements ...string) error {
	dst := filepath.Join(util.config.ScoDir, filepath.Join(destpathElements...))
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
