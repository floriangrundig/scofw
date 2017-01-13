package wktreeobserver

/**
* The WorkTreeObserver module provides functions to get the current git work tree hash.
* The git work tree hash is used to store our runtime data (file snapshots, modifications stats, ...)
 */
import (
	"fmt"
	log_ "log"

	"github.com/floriangrundig/scofw/config"
	gitconfig "github.com/floriangrundig/scofw/git/config"
	"github.com/libgit2/git2go"
	"github.com/satori/go.uuid"
)

var (
	log *log_.Logger
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
	log = config.Logger
	repo, err := git.OpenRepository(config.ProjectDir)
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

// UpdateCurrentScoSession detects the current git work tree.
// If that has changed it creates a new unique sco session
// and stores it in the config
func (observer *WorkTreeObserver) UpdateCurrentScoSession() {

	// TODO it would be nice if we detect a new ref automatically
	// TODO handle bare repositories (#10)
	ref, err := observer.repo.Head() // TODO is this really what we've checked out?
	if err != nil {
		log.Fatal(err)
	}

	target := fmt.Sprint(ref.Target())
	// log.Println("Current work tree:", parent)

	if !observer.hasMappingToCurrentGitCommit(target) {
		newSession := observer.createNewMappingToCurrentGitCommit(target)
		log.Println("Creating new session for current work tree:", newSession)
		observer.config.SetCurrentScoSession(newSession)
		observer.config.Persist()
	} else {
		session := observer.getCurrentSession(target)
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
