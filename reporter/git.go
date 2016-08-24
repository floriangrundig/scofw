package gitReporter

import (
	"fmt"
	log_ "log"
	"os"
	"strings"

	"io/ioutil"

	"path/filepath"

	"github.com/floriangrundig/scofw/config"
	"github.com/floriangrundig/scofw/fw"
	wkTree "github.com/floriangrundig/scofw/git"
	gitconfig "github.com/floriangrundig/scofw/git/config"
	"github.com/floriangrundig/scofw/publisher"
	"github.com/floriangrundig/scofw/util"
	"github.com/libgit2/git2go"
)

var (
	log *log_.Logger
)

// TODO rename that struct
type GitReporter struct {
	home                      string
	config                    *config.Config
	gitConfig                 *gitconfig.Config
	repo                      *git.Repository
	util                      *util.Util
	observer                  *wkTree.WorkTreeObserver
	fileEventChannel          chan *fw.FileEvent
	fileChangedMessageChannel chan *publisher.Message
}

func New(config *config.Config, gitConfig *gitconfig.Config, util *util.Util, observer *wkTree.WorkTreeObserver, fileEventChannel chan *fw.FileEvent, fileChangedMessageChannel chan *publisher.Message, home string) *GitReporter {

	log = config.Logger
	// the default fw-engine uses git
	// there might be some other engines which don't need git
	// in the latter case we should make the engine configurable via cli params
	repo, err := git.OpenRepository(config.ProjectDir)
	if err != nil {
		log.Fatal(err)
	}

	return &GitReporter{
		home:                      home,
		config:                    config,
		gitConfig:                 gitConfig,
		repo:                      repo,
		util:                      util,
		observer:                  observer,
		fileEventChannel:          fileEventChannel,
		fileChangedMessageChannel: fileChangedMessageChannel,
	}
}

/*
 * We keep track of each file we already processed - we store the filepath and the last action.
 * If a file was renamed we have to store the old and new file name.
 *
 * We handle the following scenarios:
 * 1. The file is untracked by scofw (alreadyTracked==false).
 * 		In this case it's the first time we received a change event for that file.
 * 		We need to determine now what happend to that file:
 * 		a.) The file was changed (not removed or renamed):
 * 				i) If the file is already tracked by git we have two alternatives:
 * 					 1.)
 *						 -> copy the original file content to diffs/.../a
 *						 -> copy the current file content to diffs/.../b
 *  					 -> make a patch
 * 					 2.)
 * 						 -> copy the current file content to diffs/.../a
 *						 -> use diff git to get latest changes
 * 							* here we might encounter timing issues
 *				ii) If the file is not tracked by git:
 *					-> copy the new file to diffs/.../b
 *					-> create an empty file with same name under diffs/.../a
 *					-> make a patch
 *
 * 		b.) The file was removed
 *				i) If the file is already tracked by git we have to alternatives:
 *					1.)
 *						-> copy the original file content to diffs/.../a
 * 						-> make a patch
 * 					2.)
 *						-> use diff git to get latest changes
 * 							* here we might encounter timing issues
 *				ii) If the file was not tracked by git and scofw -> that shouldn't be possible
 *
 * 		c.) The file was renamed
 *				i) If the file is already tracked by git we have to alternatives:
 *					1.)
 *						-> copy the original file content to diffs/.../a
 *						-> copy new renamed file content to diffs/.../b
 * 						-> make a patch
 * 					2.)
 *						-> use diff git to get latest changes
 * 							* here we might encounter timing issues
 *				ii) If the file is not tracked by git -> do the same as in c.) i) 1.)
 *
 *
 * 2. The file is tracked by scofw (alreadyTracked==true).
 * 		In this case there should be already a previous version of the file in diffs/.../a and in diffs/.../b
 * 		execpt the last action was a delete - so we need to store the type of the last change per file.
 * 		We need to determine now what happend to that file:
 * 		a.) The file was changed (not removed or renamed):
 *	 		-> copy the current file content to diffs/.../b
 *  		-> make a patch
 *
 * 		b.) The file was removed
 * 			Since a file can not be removed twice there must current file under diffs/.../a and diffs/.../b
 *			-> remove file under diffs/.../b
 * 			-> make a patch
 *
 * 		c.) The file was renamed
 *			-> rename file under diffs/.../b (better copy new file to diffs/.../b and remove old file under diffs/.../b)
 * 			-> make a patch
 *
 *



 */

func (gr *GitReporter) Start() {
	log.Println("Starting git reporter")

	go func() {

		// we wait infinitely - TODO how to shut down gracefully?
		for {
			// wait for incoming events
			event, ok := <-gr.fileEventChannel

			if !ok {
				log.Println("Shutting down git reporter")
				break
			}

			// since we currently are not able to detect a commit we have to sync with current work tree manually
			gr.observer.UpdateCurrentScoSession()

			// flg: I don't know why but when editing with atom editor lot's of chmod-events are triggered - we're not interested in those
			if event.Op != fw.Chmod {

				log.Printf("Received event %s %d (%s)\n", event.Name, event.Op, FileEventToString(event))
				lastChanges, anyChangesForSession := gr.gitConfig.LastChanges[gr.gitConfig.CurrentScoSession]

				var alreadyTracked bool
				var lastChange uint32

				if !anyChangesForSession {
					alreadyTracked = false
				} else {

					lastChange, alreadyTracked = lastChanges[event.Name]
				}

				// TODO if we receive a delete on a folder we have to deal with it -> e.g. deleting all files we know of (even files we don't know -> they not tracked but they should be listed in the git tracking...)

				if !alreadyTracked {
					gr.handleFirstChange(event)
				} else {
					gr.handleChange(event, lastChange)
				}
			}

		}
	}()
}

func FileEventToString(event *fw.FileEvent) string {
	result := ""
	if event.Op&fw.Chmod == fw.Chmod {
		result = fmt.Sprint("| chmod |", result)
	}
	if event.Op&fw.Create == fw.Create {
		result = fmt.Sprint("| create |", result)
	}
	if event.Op&fw.Write == fw.Write {
		result = fmt.Sprint("| write |", result)
	}
	if event.Op&fw.Remove == fw.Remove {
		result = fmt.Sprint("| remove |", result)
	}
	if event.Op&fw.Rename == fw.Rename {
		result = fmt.Sprint("| rename |", result)
	}

	if result == "" {
		return "!!!unknown - this is not expected to happen!!!"
	}
	return result
}

func verifyNoError(err error) {
	if err != nil {
		log.Fatal(err)
	}

}

func joinFilePaths(p1, p2 string) string {

	if p1 == "." {
		return p2
	}
	if p2 == "." {
		return p1
	}
	return filepath.Join(p1, p2)
}

// TODO this function should go into its own package because its the only function which is really tied to git
func (gr *GitReporter) handleFirstChange(event *fw.FileEvent) {
	log.Println("This is the first change detected for", event.Name)

	baseFolder := filepath.Join("diffs", gr.gitConfig.CurrentScoSession, filepath.Dir(event.Name))
	baseFile := filepath.Join(baseFolder, filepath.Base(event.Name))
	gr.util.CreateScoFolder(baseFolder)

	var contentA []byte
	var contentB []byte
	emptyContent := []byte("")

	// TODO: I assume you can only work on Head or do we need a more sophisticated way of determing what I'm working on
	ref, err := gr.repo.Head()
	verifyNoError(err)

	oidHead := ref.Target()

	commit, err := gr.repo.LookupCommit(oidHead)
	verifyNoError(err)

	commitTree, err := commit.Tree()
	verifyNoError(err)

	options, err := git.DefaultDiffOptions()
	verifyNoError(err)

	// Specifying full patch indices.
	options.IdAbbrev = 40
	options.Flags |= git.DiffIncludeUntracked
	options.Flags |= git.DiffShowUntrackedContent

	gitDiff, err := gr.repo.DiffTreeToWorkdir(commitTree, &options)
	verifyNoError(err)

	numDeltas, err := gitDiff.NumDeltas()
	verifyNoError(err)

	var contentDeltaDetermined bool

	for d := 0; d < numDeltas; d++ {

		delta, err := gitDiff.GetDelta(d)
		verifyNoError(err)

		// we only need to look at our file
		if strings.HasSuffix(event.Name, delta.NewFile.Path) {

			gr.logStatusDiffDelta(&delta)

			if delta.Status != git.DeltaUntracked {
				blob := gr.getOriginalBlob(commitTree, event)
				contentA = blob.Contents()

			} else {
				log.Println("This file is not tracked by git", event.Name)
				// we create an empty file in diffs/.../a since this file event belongs to a new file
				contentA = emptyContent
			}

			if delta.Status != git.DeltaDeleted {
				contentB, err = ioutil.ReadFile(event.Name)
				verifyNoError(err)
			} else {
				contentB = emptyContent // TODO is this really identical to delete?
			}

			contentDeltaDetermined = true
			break

		}
	}

	if !contentDeltaDetermined {

		log.Printf("No matching git change for file: %s", event.Name)

		_, err := commitTree.EntryByPath(event.Name)
		if err == nil {
			blob := gr.getOriginalBlob(commitTree, event)
			log.Printf("Anyway %s is tracked by git and git doesn't detect a change - assuming nothing has changed", event.Name)
			contentA = blob.Contents()
		} else {
			log.Printf("Going to fallback - %s seems to be an untracked file", event.Name)
			contentA = emptyContent
		}

		if event.Op&fw.Remove != fw.Remove {
			log.Println("not a removal")
			if _, err := os.Stat(event.Name); os.IsNotExist(err) {
				contentB = emptyContent
			} else {
				contentB, err = ioutil.ReadFile(event.Name) // TODO can we be sure that this file is there (deleted?)?
				verifyNoError(err)
			}
		} else {
			contentB = emptyContent
		}
	}

	// TODO if event.Name referes to a new file -> the patch contains "new file mode 100644" -> we should change the file mode to the original settings
	oldPath := gr.toProjectRelativePath(event.Name)
	newPath := gr.toProjectRelativePath(event.Name)
	patch, err := gr.repo.PatchFromBuffers(oldPath, newPath, contentA, contentB, &options)
	defer patch.Free()
	verifyNoError(err)
	patchString, err := patch.String()
	_, err = patch.String()
	verifyNoError(err)

	// send to publisher
	gr.fileChangedMessageChannel <- &publisher.Message{
		FileEvent: event,
		Patch:     &patchString,
	}

	// we store contentB as a snapshot of that file -> all further diffs will be made between workspace file and snapshot

	if event.Op&fw.Remove != fw.Remove {
		gr.util.WriteFile(&contentB, baseFile)
	} else {
		gr.util.RemoveFile(baseFile)
	}

	gr.storeLastChange(event)
}

func (gr *GitReporter) toProjectRelativePath(path string) string {
	// TODO add teh
	relativePath, err := filepath.Rel(gr.config.ProjectDir, path)
	if err != nil {
		log.Println("Error while transforming project directory into relative directory:", err)
	}

	return relativePath
}

func (gr *GitReporter) storeLastChange(event *fw.FileEvent) {
	lastChanges, anyChangesForSession := gr.gitConfig.LastChanges[gr.gitConfig.CurrentScoSession]

	if !anyChangesForSession {
		gr.gitConfig.LastChanges[gr.gitConfig.CurrentScoSession] = make(map[string]uint32)
	}

	lastChanges, _ = gr.gitConfig.LastChanges[gr.gitConfig.CurrentScoSession]

	lastChanges[event.Name] = uint32(event.Op)
	gr.gitConfig.Persist()
}

func (gr *GitReporter) getOriginalBlob(commitTree *git.Tree, event *fw.FileEvent) *git.Blob {
	path := gr.toProjectRelativePath(event.Name) // event.Name is an absolute path

	// under linux files in the root directory are reported as ./filename which is an invalid git tree path -> we have to remove the "./"
	// if filepath.Dir(path) == "." {
	// path = filepath.Base(path)
	// }

	log.Println("Looking in commit tree for ", path)
	treeEntry, err := commitTree.EntryByPath(path)
	if err != nil {
		log.Fatal(err)
	}

	blob, err := gr.repo.LookupBlob(treeEntry.Id)
	if err != nil {
		log.Fatal(err)
	}

	return blob
}

func (gr *GitReporter) handleChange(event *fw.FileEvent, lastChange uint32) {
	log.Println("This is a change detected for", event.Name)

	options, err := git.DefaultDiffOptions()

	verifyNoError(err)

	// Specifying full patch indices. TODO what is needed here?
	options.IdAbbrev = 40
	options.Flags |= git.DiffIncludeUntracked

	baseFolder := filepath.Join("diffs", joinFilePaths(gr.gitConfig.CurrentScoSession, filepath.Dir(event.Name)))
	baseFile := filepath.Join(baseFolder, filepath.Base(event.Name))

	var contentA *[]byte
	var contentB []byte
	emptyContent := []byte("")

	if lastChange&uint32(fw.Remove) != uint32(fw.Remove) { // TODO: how to handle renamed files? (Beware -> IntelliJ stores the changes in a tmp file and renames that tmp file to the current file)
		log.Printf("Comparing current %s with last snapshot %s\n", event.Name, baseFile)
		contentA, err = gr.util.ReadScoFile(baseFile)
		verifyNoError(err)
	} else {
		// we create an empty file in diffs/.../a since this file event belongs to a new file
		contentA = &emptyContent // TODO this is not correct for IntelliJ -> when you revert your changes it's removed first and then created again... so we think it's a complete new file
	}

	if event.Op&fw.Remove != fw.Remove {

		if _, err = os.Stat(event.Name); os.IsNotExist(err) {
			contentB = emptyContent
		} else {
			contentB, err = ioutil.ReadFile(event.Name) // TODO can we be sure that this file is there (deleted?)?
			verifyNoError(err)
		}

	} else {
		contentB = emptyContent // TODO is this really identical to delete?
	}

	// TODO if event.Name referes to a new file -> the patch contains "new file mode 100644" -> we should change the file mode to the original settings
	oldPath := gr.toProjectRelativePath(event.Name)
	newPath := gr.toProjectRelativePath(event.Name)
	patch, err := gr.repo.PatchFromBuffers(oldPath, newPath, *contentA, contentB, &options)
	defer patch.Free()
	verifyNoError(err)
	patchString, err := patch.String()
	_, err = patch.String()
	verifyNoError(err)

	// publish event
	gr.fileChangedMessageChannel <- &publisher.Message{
		FileEvent: event,
		Patch:     &patchString,
	}

	// we store contentB as a snapshot of that file -> all further diffs will be made between workspace file and snapshot
	if event.Op&fw.Remove != fw.Remove {
		gr.util.WriteFile(&contentB, baseFile)
	} else {
		gr.util.RemoveFile(baseFile)
	}

	gr.storeLastChange(event)
}

func (gr *GitReporter) logStatusDiffDelta(delta *git.DiffDelta) {
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

}
