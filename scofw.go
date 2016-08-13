package main

import (
	"io"
	"log"
	"os"
	"path/filepath"

	"github.com/floriangrundig/scofw/config"
	"github.com/floriangrundig/scofw/fw"
	"github.com/floriangrundig/scofw/git"
	gitconfig "github.com/floriangrundig/scofw/git/config"
	"github.com/floriangrundig/scofw/publisher"
	"github.com/floriangrundig/scofw/reporter"
	"github.com/floriangrundig/scofw/util"
)

// to statically link lib: www.petethompson.net/blog/golang/2015/10/04/getting-going-with-git2go/

func main() {
	scoDir := ".sco"

	file, err := os.OpenFile(filepath.Join(scoDir, "logs", "scofw.log"), os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0770)
	if err != nil {
		log.Fatalln("Failed to open log file scofw.log:", err)
	}

	// TODO hmmm - can I not close that defer into a separate function when I want to have an open filehandle for my app live time?
	defer file.Close()
	logger := createLogger(file, config.GetVerboseLoggingFlag())

	// stores global configuration
	config := config.New(scoDir, logger)

	// utility module for creating sco related files/directories
	util := util.New(config)

	// stores e.g. mapping between current git commit (HEAD) and sco-session
	gitConfig := gitconfig.New(config)

	// observes current working tree
	wktreeObserver := wktreeobserver.New(config, gitConfig)

	// Channel from filewatcher to reporter
	fileEventChannel := make(chan *fw.FileEvent)

	// Channel from reporter to publisher
	fileChangedMessageChannel := make(chan *publisher.Message)

	// file watcher -> reports file event changes into fileEventChannel
	fw := fw.New(config, fileEventChannel)

	// listen on fileEventChannel -> determines the diff and updates the current sco-wktree patch
	gitReporter := gitReporter.New(config, gitConfig, util, wktreeObserver, fileEventChannel, fileChangedMessageChannel)

	// listen
	publisher := publisher.New(config, gitConfig, util, fileChangedMessageChannel)

	publisher.Start()

	wktreeObserver.Start()

	gitReporter.Start()

	fw.Start()

	// TODO do not rely on fw.Start() to block use channel ...
}

func createLogger(file io.Writer, verbose bool) *log.Logger {

	var writer io.Writer

	if verbose {
		writer = io.MultiWriter(file, os.Stdout)
	} else {
		writer = io.Writer(file)
	}

	return log.New(writer, "", log.Ldate|log.Ltime)
}
