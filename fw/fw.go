package fw

import (
	"log"
	"os"
	"path/filepath"

	"github.com/FlorianGrundig/scofw/config"
	"github.com/fsnotify/fsnotify"
)

// The types Event and Op and the Op-constants were copied from the fsnotify.go -> we wrap them here

// FileEvent represents a single file system notification.
type FileEvent struct {
	Name string // Relative path to the file or directory.
	Op   Op     // File operation that triggered the event.
}

// Op describes a set of file operations.
type Op uint32

// These are the generalized file operations that can trigger a notification.
const (
	Create Op = 1 << iota
	Write
	Remove
	Rename
	Chmod
)

type FileWatcher struct {
	config    *config.Config
	eventSink chan *FileEvent
}

func New(config *config.Config, eventSink chan *FileEvent) *FileWatcher {
	return &FileWatcher{
		config:    config,
		eventSink: eventSink,
	}
}

func (fw *FileWatcher) Start() {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Println("Error while creating file watcher")
		log.Fatal(err)
	}

	log.Println("Watching directory: " + fw.config.BaseDir)
	defer watcher.Close()

	walkFunc := func(path string, info os.FileInfo, err error) error {

		if err == nil {
			if info.IsDir() && !fw.config.GitIgnore.MatchesPath(path) {
				log.Println("Watching", path)
				err = watcher.Add(path)
				if err != nil {
					log.Fatal(err)
				}
			}
		} else {
			log.Println("Error while walking through directory tree in workspace", err)
		}

		return nil
	}

	done := make(chan bool)

	go func() {
		for {
			select {
			case event := <-watcher.Events:
				if !fw.config.GitIgnore.MatchesPath(event.Name) {
					shouldEmitEvent := true

					if event.Op&fsnotify.Create == fsnotify.Create {
						fileInfo, err := os.Stat(event.Name)

						if err != nil {
							log.Println("error", err)
						}

						go func() {
						if fileInfo.IsDir() {
							walkErr := filepath.Walk(event.Name, walkFunc)
							if walkErr != nil {
								log.Fatal(walkErr)
							}
						}

							// TODO if there're already some file in the new folder or its subfolder then we should emit an event
							shouldEmitEvent = false
						}
					} else if event.Op&fsnotify.Remove == fsnotify.Remove {
						if _, err := os.Stat(event.Name); os.IsNotExist(err) {
							// TODO store all watches and remove watch if file.Name matches...
							// watcher.Remove(event.Name)
						}

					}

					if shouldEmitEvent {
						fw.eventSink <- convertFsNotifyEvent(event)
					}
				}
			case err := <-watcher.Errors:
				log.Println("error:", err)
			}
		}
	}()

	walkErr := filepath.Walk(fw.config.BaseDir, walkFunc)
	if walkErr != nil {
		log.Fatal(walkErr)
	}

	<-done
}

func convertFsNotifyEvent(event fsnotify.Event) *FileEvent {

	var op Op

	// TODO I don't think you can have several events at the same time e.g. chmod + write???
	if event.Op&fsnotify.Chmod == fsnotify.Chmod {
		op |= Chmod
	}
	if event.Op&fsnotify.Create == fsnotify.Create {
		op |= Create
	}
	if event.Op&fsnotify.Write == fsnotify.Write {
		op |= Write
	}
	if event.Op&fsnotify.Remove == fsnotify.Remove {
		op |= Remove
	}
	if event.Op&fsnotify.Rename == fsnotify.Rename {
		op |= Rename
	}

	if uint32(op) != uint32(event.Op) {
		log.Fatal("fsnotify events are not parsable")
	}

	return &FileEvent{
		Name: event.Name,
		Op:   op,
	}
}
