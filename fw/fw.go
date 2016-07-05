package fw

import (
	"log"

	"github.com/floriangrundig/scofw/config"
	"github.com/floriangrundig/scofw/reporter"
	"github.com/fsnotify/fsnotify"
)

type FileWatcher struct {
	config *config.Config
}

func New(config *config.Config) *FileWatcher {
	return &FileWatcher{
		config: config,
	}
}

func (fw *FileWatcher) Start() {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Println("Error while creating file watcher")
		log.Fatal(err)
	}

	log.Println("Watching directory (no subdirectories): " + fw.config.BaseDir)

	defer watcher.Close()

	done := make(chan bool)
	go func() {
		for {
			select {
			case event := <-watcher.Events:
				// log.Println("file watcher fired event:", event)
				op := ""
				if event.Op&fsnotify.Chmod == fsnotify.Chmod {
					op += "Chmod "
				}
				if event.Op&fsnotify.Create == fsnotify.Create {
					op += "Create "
				}
				if event.Op&fsnotify.Write == fsnotify.Write {
					op += "Write "
				}
				if event.Op&fsnotify.Remove == fsnotify.Remove {
					op += "Remove "
				}
				if event.Op&fsnotify.Rename == fsnotify.Rename {
					op += "Rename "
				}

				/*
				 * TODO
				 * We'll have a .sco subdirectory with current git hash as subfolder
				 * where we'll store the intermediate file versions.
				 *
				 * When a change of a file e.g. f1.txt is detected for the first time,
				 * we copy the file into the .sco/head-git-hash/.../f1.txt.tmp
				 * We use the git-diff to report the first change. If git-diff
				 * signal an untracked file we create an empty file with the real
				 * name f1.txt.base and use a normal diff with f1.txt.tmp.
				 * After performing the diff we move the f1.txt.tmp to f1.txt.base
				 * Maybe we store each patch (with corrected paths) as f1.txt.p1 ... f1.txt.p9999
				 *
				 * When further changes were detected we copy the file as f1.txt.tmp
				 * again and use a normal diff to detect a change with f1.txt.base -
				 * after that we move the f1.txt.tmp to f1.txt.base.
				 * Maybe we store each patch (with corrected paths) as f1.txt.p1 ... f1.txt.p9999
				 *
				 * However we have to make sure that we don't copy a file to as file.tmp
				 * if there's still such file - then we have to retry later ...
				 */
				gitReporter.FileModified(event.Name, op)
			case err := <-watcher.Errors:
				log.Println("error:", err)
			}
		}
	}()

	err = watcher.Add(fw.config.BaseDir)
	if err != nil {
		log.Fatal(err)
	}
	<-done
}
