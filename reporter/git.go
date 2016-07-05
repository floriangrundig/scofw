package gitReporter

import (
	"log"
	"strings"

	"github.com/libgit2/git2go"
)

var (
	repo *git.Repository
)

func SetRepo(_repo *git.Repository) {
	repo = _repo
}

func FileModified(name string, op string) {

	if op != "Chmod " {

		go func() {

			ref, err := repo.Head()
			if err != nil {
				log.Fatal(err)
			}

			oidHead := ref.Target()
			// log.Println("HEAD:", oidHead)

			commit, err := repo.LookupCommit(oidHead)
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

			gitDiff, err := repo.DiffTreeToWorkdir(commitTree, &options)
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

				if strings.HasSuffix(name, delta.NewFile.Path) {
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

			log.Printf("reporting modification [%s] of file: %s", op, name)

		}()
	}
}
