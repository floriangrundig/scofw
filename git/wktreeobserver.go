package wktreeobserver

import (
	"fmt"
	"log"

	"github.com/floriangrundig/scofw/config"
	gitconfig "github.com/floriangrundig/scofw/git/config"
	"github.com/libgit2/git2go"
	"github.com/satori/go.uuid"
)

type GitRuntimeData struct {
	GitCommits map[string]string
}

type WorkTreeObserver struct {
	scoConfig          *config.Config
	config             *gitconfig.Config
	gitRuntimeDataFile string
	GitRuntimeData
}

func New(config *config.Config, gitConfig *gitconfig.Config) *WorkTreeObserver {

	return &WorkTreeObserver{
		scoConfig:          config,
		config:             gitConfig,
		gitRuntimeDataFile: "commits_sessions.json",
	}
}

func (observer *WorkTreeObserver) Start() {

	repo, err := git.OpenRepository(observer.scoConfig.BaseDir)
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
		observer.config.SetCurrentScoSession(newSession)
		observer.config.Persist()
	} else {
		session := observer.getCurrentSession(parent)
		observer.config.SetCurrentScoSession(session)
		log.Println("Continue with session:", session)
	}

	// update config if parent is not known -> create new uuid subdir which is the new working dir
}

func (observer *WorkTreeObserver) hasMappingToCurrentGitCommit(parent string) bool {
	_, exists := observer.config.GitRuntimeData.GitCommits[parent]
	return exists
}

func (observer *WorkTreeObserver) getCurrentSession(parent string) string {
	data, _ := observer.config.GitRuntimeData.GitCommits[parent]
	return string(data)
}

func (observer *WorkTreeObserver) createNewMappingToCurrentGitCommit(parent string) string {
	u1 := uuid.NewV4()

	observer.config.GitRuntimeData.GitCommits[parent] = u1.String()
	return u1.String()
}
