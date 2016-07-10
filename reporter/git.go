	"io/ioutil"
	"path/filepath"

	gitconfig "github.com/floriangrundig/scofw/git/config"
	"github.com/floriangrundig/scofw/util"
	gitConfig        *gitconfig.Config
	util             *util.Util
func New(config *config.Config, gitConfig *gitconfig.Config, util *util.Util, fileEventChannel chan *fw.FileEvent) *GitReporter {
		gitConfig:        gitConfig,
		util:             util,
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

		// we wait infinitely - TODO how to shut down gracefully?
			// flg: I don't know why but when editing with atom editor lot's of chmod-events are triggered - we're not interested in those

				log.Printf("Received event %s %d\n", event.Name, event.Op)
				lastChange, alreadyTracked := gr.gitConfig.LastChanges[event.Name]

				if !alreadyTracked {
					go gr.handleFirstChange(event)
				} else {
					go gr.handleChange(event, lastChange)
				}


func verifyNoError(err error) {
	if err != nil {
		log.Fatal(err)
	}

}

func (gr *GitReporter) handleFirstChange(event *fw.FileEvent) {
	log.Println("This is the first change detected for", event.Name)

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

	gitDiff, err := gr.repo.DiffTreeToWorkdir(commitTree, &options)
	verifyNoError(err)

	numDeltas, err := gitDiff.NumDeltas()
	verifyNoError(err)

	for d := 0; d < numDeltas; d++ {

		delta, err := gitDiff.GetDelta(d)
		verifyNoError(err)

		// we only need to look at our file
		if strings.HasSuffix(event.Name, delta.NewFile.Path) {

			gr.logStatusDiffDelta(&delta)

			baseFolder := filepath.Join("diffs", gr.gitConfig.CurrentScoSession, filepath.Dir(event.Name))
			baseFile := filepath.Join(baseFolder, filepath.Base(event.Name))

			gr.util.CreateScoFolder(baseFolder)

			var contentA []byte
			var contentB []byte
			emptyContent := []byte("")

			if delta.Status != git.DeltaUntracked {
				log.Println("This is the first change detected for", event.Name)
				blob := gr.getOriginalBlob(commitTree, event)
				contentA = blob.Contents()

			} else {
				// we create an empty file in diffs/.../a since this file event belongs to a new file
				contentA = emptyContent
			}

			if delta.Status != git.DeltaDeleted {
				contentB, err = ioutil.ReadFile(event.Name)
				verifyNoError(err)
			} else {
				contentB = emptyContent // TODO is this really identical to delete?
			}

			// TODO if event.Name referes to a new file -> the patch contains "new file mode 100644" -> we should change the file mode to the original settings
			patch, err := gr.repo.PatchFromBuffers(event.Name, event.Name, contentA, contentB, &options)
			defer patch.Free()
			verifyNoError(err)
			patchString, err := patch.String()
			_, err = patch.String()
			verifyNoError(err)

			// TOOD use channel to publish change...
			log.Printf("\n%s", patchString)

			// we store contentB as a snapshot of that file -> all further diffs will be made between workspace file and snapshot
			gr.util.WriteFile(&contentB, baseFile)
			gr.gitConfig.LastChanges[event.Name] = uint32(event.Op)
			gr.gitConfig.Persist()
			return
		}
	}

	log.Printf("ERROR: No matching git change to file: %s", event.Op, event.Name)

}

func (gr *GitReporter) getOriginalBlob(commitTree *git.Tree, event *fw.FileEvent) *git.Blob {
	treeEntry, err := commitTree.EntryByPath(event.Name)
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

	// Specifying full patch indices.
	options.IdAbbrev = 40
	options.Flags |= git.DiffIncludeUntracked

	baseFolder := filepath.Join("diffs", gr.gitConfig.CurrentScoSession, filepath.Dir(event.Name))
	baseFile := filepath.Join(baseFolder, filepath.Base(event.Name))

	gr.util.CreateScoFolder(baseFile)

	var contentA *[]byte
	var contentB []byte
	emptyContent := []byte("")

	if event.Op == fw.Create || event.Op == fw.Write {
		contentA, err = gr.util.ReadScoFile(baseFile)
		verifyNoError(err)
	} else if event.Op == fw.Remove {
		// we create an empty file in diffs/.../a since this file event belongs to a new file
		contentA = &emptyContent
	} else {
		contentA = &emptyContent // TODO: how to handle renamed files? Maybe we should treat them as removed?
	}

	if event.Op != fw.Remove {
		contentB, err = ioutil.ReadFile(event.Name)
		verifyNoError(err)
	} else {
		contentB = emptyContent // TODO is this really identical to delete?
	}

	// TODO if event.Name referes to a new file -> the patch contains "new file mode 100644" -> we should change the file mode to the original settings
	patch, err := gr.repo.PatchFromBuffers(event.Name, event.Name, *contentA, contentB, &options)
	defer patch.Free()
	verifyNoError(err)
	patchString, err := patch.String()
	_, err = patch.String()
	verifyNoError(err)

	// TOOD use channel to publish change...
	log.Printf("\n%s", patchString)

	// we store contentB as a snapshot of that file -> all further diffs will be made between workspace file and snapshot
	gr.util.WriteFile(&contentB, baseFile)

	return

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