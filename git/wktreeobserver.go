package wktreeobserver

import (
	"fmt"
	"log"

	"github.com/floriangrundig/scofw/config"
	"github.com/floriangrundig/scofw/util"
	"github.com/libgit2/git2go"
	"github.com/satori/go.uuid"
)

type WorkTreeObserver struct {
	config *config.Config
	util   *util.Util
}

func New(config *config.Config, util *util.Util) *WorkTreeObserver {
	return &WorkTreeObserver{
		config: config,
		util:   util,
	}
}

func (observer *WorkTreeObserver) Start() {
	repo, err := git.OpenRepository(observer.config.BaseDir)
	if err != nil {
		log.Fatal(err)
	}

	ref, err := repo.Head() // TODO is this really what we've checked out?
	if err != nil {
		log.Fatal(err)
	}

	parent := fmt.Sprint(ref.Target())
	log.Println("Current work tree:", parent)

	if !observer.hasMappingToCurrentGitCommit(parent) {
		newSession := observer.createNewMappingToCurrentGitCommit(parent)
		log.Println("Creating new session for current work tree:", newSession)
		observer.config.Persist()
	} else {
		log.Println("Continue with session:", observer.getCurrentSession(parent))
	}

	// update config if parent is not known -> create new uuid subdir which is the new working dir
}

func (observer *WorkTreeObserver) hasMappingToCurrentGitCommit(parent string) bool {
	_, exists := observer.config.ScoRuntimeData.GitCommits[parent]
	return exists
}

func (observer *WorkTreeObserver) getCurrentSession(parent string) string {
	data, _ := observer.config.ScoRuntimeData.GitCommits[parent]
	return string(data)
}

func (observer *WorkTreeObserver) createNewMappingToCurrentGitCommit(parent string) string {
	u1 := uuid.NewV4()

	observer.config.ScoRuntimeData.GitCommits[parent] = u1.String()
	return u1.String()
}
