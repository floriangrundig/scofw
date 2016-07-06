package main

import (
	"github.com/floriangrundig/scofw/config"
	"github.com/floriangrundig/scofw/fw"
	"github.com/floriangrundig/scofw/git"
	"github.com/floriangrundig/scofw/reporter"
	"github.com/floriangrundig/scofw/util"
)

// to statically link lib: www.petethompson.net/blog/golang/2015/10/04/getting-going-with-git2go/

func main() {

	// stores global configuration
	config := config.New()

	// utility module for creating sco related files/directories
	util := util.New(config)

	// observes current working tree
	wktreeObserver := wktreeobserver.New(config, util)

	// Channel from filewatcher to reporter
	fileEventChannel := make(chan *fw.FileEvent)

	// file watcher -> reports file event changes into fileEventChannel
	fw := fw.New(config, fileEventChannel)

	// listen on fileEventChannel -> determines the diff and updates the current sco-wktree patch
	gitReporter := gitReporter.New(config, fileEventChannel)

	go wktreeObserver.Start()
	go gitReporter.Start()
	fw.Start()
}
