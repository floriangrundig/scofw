package main

import (
	"log"
	"os"

	"github.com/floriangrundig/scofw/config"
	"github.com/floriangrundig/scofw/fw"
	"github.com/floriangrundig/scofw/reporter"
	"github.com/libgit2/git2go"
)

var ()

type ScoFw struct {
	config *config.Config
}

// to statically link lib: www.petethompson.net/blog/golang/2015/10/04/getting-going-with-git2go/

func main() {

	var config = config.New()

	log.Println("Checking path:", config.BaseDir)

	// the default fw-engine uses git
	// there might be some other engines which don't need git
	// in the latter case we should make the engine configurable via cli params
	repo, err := git.OpenRepository(config.BaseDir)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("...OK (git repository)")
	gitReporter.SetRepo(repo)

	if _, err := os.Stat(".sco"); os.IsNotExist(err) {
		err = os.Mkdir(".sco", 0777)
		if err != nil {
			log.Fatal(err)
		}
		log.Println("Created .sco folder which is used by sco filewatcher...")
		log.Println("Please add the .sco folder to your git ignore")
		// TODO create folders "a" and "b" with and subfolders corresponding
		// to the current watched dir and its subdirectories
	}

	var fw = fw.New(config)

	fw.Start()
}
