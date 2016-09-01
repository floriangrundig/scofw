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
	"github.com/mitchellh/go-homedir"
)

type apps struct {
	projectConfigs []struct {
		name   string
		config config.Config
	}
}

// to statically link lib: www.petethompson.net/blog/golang/2015/10/04/getting-going-with-git2go/

func main() {

	// parse CLI parameters
	config.ParseCLI()

	home, err := homedir.Dir()

	if err != nil {
		log.Fatal("Unable to resolve home directory: ", err)
	}

	// create sco internal directory

	/*
		TODO:
		1. 		read sco configuration from a configurable path (defaults to ~/.sco)
		1.1 	from that file get list of projects
		1.2 	for each listed project start all components

		2.1		filebeat.yml file location is part of 1.1
		2.2 	location of internal sco folder for that project is part of 1.1 (defaults to project dir)

	*/

	globalConfig := config.ParseGlobalConfig(home)

	done := make(chan bool)

	for _, projectDefinition := range globalConfig.Projects {
		// set default scoDir if not set in global config
		if projectDefinition.ScoDir == "" {
			projectDefinition.ScoDir = filepath.Join(home, ".sco", filepath.Base(projectDefinition.Path))
		}

		util.CreateInternalScoDirectory(projectDefinition.ScoDir, 0700)

		log.Printf("Starting SCO for project \"%s\":\n", projectDefinition.Name)
		log.Printf("SCO internal directory: %s", projectDefinition.ScoDir)

		/* For each project
		* 1. Create sco internal directory
		* 2. Create log filepath
		* 3. Start Sco -> store config ()
		 */

		// TODO util.CreateInternalScoDirectory is not able to handle paths like "~/foo"

		file, err := os.OpenFile(filepath.Join(projectDefinition.ScoDir, "logs", "scofw.log"), os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0770)
		if err != nil {
			log.Fatalln("Failed to open log file scofw.log:", err)
		}
		defer file.Close()

		launchScoForProject(&projectDefinition, file, home)
		log.Println("Done...")
	}

	<-done
}

func launchScoForProject(projectDefinition *config.ProjectConfig, logfile io.Writer, home string) *config.Config {

	logger := createLogger(logfile, config.GetVerboseLoggingFlag())

	// stores global configuration
	config := config.New(projectDefinition, logger)

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
	gitReporter := gitReporter.New(config, gitConfig, util, wktreeObserver, fileEventChannel, fileChangedMessageChannel, home)

	// listen
	publisher := publisher.New(config, gitConfig, util, fileChangedMessageChannel)

	wktreeObserver.UpdateCurrentScoSession()

	publisher.Start()

	gitReporter.Start()

	fw.Start()
	return config
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
