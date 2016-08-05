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

// TODO rename that struct since it's currently not observing anything - it's totally static
type WorkTreeObserver struct {
	scoConfig          *config.Config
	config             *gitconfig.Config
	gitRuntimeDataFile string
	repo               *git.Repository
	GitRuntimeData
}

func New(config *config.Config, gitConfig *gitconfig.Config) *WorkTreeObserver {

	repo, err := git.OpenRepository(config.BaseDir)
	if err != nil {
		log.Fatal(err)
	}

	return &WorkTreeObserver{
		scoConfig:          config,
		config:             gitConfig,
		gitRuntimeDataFile: "commits_sessions.json",
		repo:               repo,
	}
}

func (observer *WorkTreeObserver) Start() {
	observer.UpdateCurrentScoSession()
}

func (observer *WorkTreeObserver) UpdateCurrentScoSession() {

	// TODO it would be nice if we detect a new ref automatically
	// TODO handle bare repositories
	ref, err := observer.repo.Head() // TODO is this really what we've checked out?
	if err != nil {
		log.Fatal(err)
	}

	parent := fmt.Sprint(ref.Target())
	// log.Println("Current work tree:", parent)

	if !observer.hasMappingToCurrentGitCommit(parent) {
		newSession := observer.createNewMappingToCurrentGitCommit(parent)
		log.Println("Creating new session for current work tree:", newSession)
		observer.config.SetCurrentScoSession(newSession)
		observer.config.Persist()
	} else {
		session := observer.getCurrentSession(parent)
		observer.config.SetCurrentScoSession(session)
		// log.Println("Continue with session:", session)
	}
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
