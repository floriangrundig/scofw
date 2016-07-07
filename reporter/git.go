package gitReporter

import (
	"log"
	"strings"

	"github.com/floriangrundig/scofw/config"
	"github.com/floriangrundig/scofw/fw"
	"github.com/libgit2/git2go"
)

type GitReporter struct {
	config           *config.Config
	repo             *git.Repository
	fileEventChannel chan *fw.FileEvent
}

func New(config *config.Config, fileEventChannel chan *fw.FileEvent) *GitReporter {

	// the default fw-engine uses git
	// there might be some other engines which don't need git
	// in the latter case we should make the engine configurable via cli params
	repo, err := git.OpenRepository(config.BaseDir)
	if err != nil {
		log.Fatal(err)
	}

	return &GitReporter{
		config:           config,
		repo:             repo,
		fileEventChannel: fileEventChannel,
	}
}

func (gr *GitReporter) Start() {
	log.Println("Starting git reporter")

	go func() {

		for {
			// wait for incoming events
			event, ok := <-gr.fileEventChannel

			if !ok {
				// channel is closed
				log.Println("Shutting down git reporter")
			}

			log.Printf("Received events %s %d\n", event.Name, event.Op)

			if event.Op != fw.Chmod {

				go func() {

					ref, err := gr.repo.Head()
					if err != nil {
						log.Fatal(err)
					}

					oidHead := ref.Target()
					// log.Println("HEAD:", oidHead)

					commit, err := gr.repo.LookupCommit(oidHead)
					if err != nil {
						log.Fatal(err)
					}
					commitTree, err := commit.Tree()
					if err != nil {
						log.Fatal(err)
					}

					options, err := git.DefaultDiffOptions()
					if err != nil {
						log.Fatal(err)
					}

					// Specifying full patch indices.
					options.IdAbbrev = 40
					options.Flags |= git.DiffIncludeUntracked

					gitDiff, err := gr.repo.DiffTreeToWorkdir(commitTree, &options)
					if err != nil {
						log.Fatal(err)
					}

					numDeltas, err := gitDiff.NumDeltas()
					if err != nil {
						log.Fatal(err)
					}
					for d := 0; d < numDeltas; d++ {

						delta, err := gitDiff.GetDelta(d)
						if err != nil {
							log.Fatal(err)
						}

						if strings.HasSuffix(event.Name, delta.NewFile.Path) {
							if delta.Status == git.DeltaUnmodified {
								log.Println("Delta: unmodified")
							}
							if delta.Status == git.DeltaUntracked {
								log.Println("Delta: untracked")
							}
							if delta.Status == git.DeltaAdded {
								log.Println("Delta: added")
							}
							if delta.Status == git.DeltaDeleted {
								log.Println("Delta: deleted")
							}
							if delta.Status == git.DeltaRenamed {
								log.Println("Delta: renamed")
							}
							if delta.Status == git.DeltaModified {
								log.Println("Delta: modified")
							}
							if delta.Status == git.DeltaCopied {
								log.Println("Delta: copied")
							}
							if delta.Status == git.DeltaTypeChange {
								log.Println("Delta: type changed")
							}

							patch, err := gitDiff.Patch(d)
							if err != nil {
								log.Fatal(err)
							}
							patchString, err := patch.String()
							if err != nil {
								log.Fatal(err)
							}
							log.Printf("\n%s", patchString)
							patch.Free()
						}

					}

					log.Printf("reporting modification [%s] of file: %s", event.Op, event.Name)

				}()
			}
		}
	}()
}
