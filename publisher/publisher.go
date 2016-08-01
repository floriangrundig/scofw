package publisher

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/FlorianGrundig/scofw/config"
	"github.com/FlorianGrundig/scofw/fw"
	gitconfig "github.com/FlorianGrundig/scofw/git/config"
	"github.com/FlorianGrundig/scofw/util"
)

type Message struct {
	FileEvent *fw.FileEvent
	Patch     *string
}

type Publisher struct {
	config                    *config.Config
	gitConfig                 *gitconfig.Config
	util                      *util.Util
	fileChangedMessageChannel chan *Message
}

func New(config *config.Config, gitConfig *gitconfig.Config, util *util.Util, fileChangedMessageChannel chan *Message) *Publisher {
	fmt.Println("Creating Publisher...")

	util.CreateScoFolder("logs")

	return &Publisher{
		config:    config,
		gitConfig: gitConfig,
		util:      util,
		fileChangedMessageChannel: fileChangedMessageChannel,
	}
}

func (publisher *Publisher) Start() {
	go func() {
		for {
			// wait for incoming messages
			msg, ok := <-publisher.fileChangedMessageChannel

			if !ok {
				log.Println("Shutting down publisher")
				break
			}

			if *msg.Patch != "" {
				publisher.log(msg)

				publisher.logInGourceFormat(msg)
			}
		}
	}()
}

func (publisher *Publisher) log(msg *Message) {
	// TODO to open a file and creating a logger each call is insufficient -> store loggers per session in a map to reuse them
	filename := fmt.Sprintf("%s.log", publisher.gitConfig.CurrentScoSession)

	file, err := os.OpenFile(filepath.Join(publisher.config.ScoDir, "logs", filename), os.O_CREATE|os.O_WRONLY|os.O_APPEND, publisher.config.ScoDirPermissions)
	if err != nil {
		log.Fatalln("Failed to open log file", filename, ":", err)
	}

	defer file.Close()

	multi := io.MultiWriter(file, os.Stdout)

	mylog := log.New(multi, "", log.Ldate|log.Ltime)
	mylog.Println(*msg.Patch)

}

func (publisher *Publisher) logInGourceFormat(msg *Message) {
	// TODO to open a file and creating a logger each call is insufficient -> store loggers per session in a map to reuse them
	filename := fmt.Sprintf("%s.gource.log", publisher.gitConfig.CurrentScoSession)

	file, err := os.OpenFile(filepath.Join(publisher.config.ScoDir, "logs", filename), os.O_CREATE|os.O_WRONLY|os.O_APPEND, publisher.config.ScoDirPermissions)
	if err != nil {
		log.Fatalln("Failed to open log file", filename, ":", err)
	}

	defer file.Close()

	text := fmt.Sprintf("%v|flg|%s|%s\n", int32(time.Now().Unix()), fileEventToString(msg.FileEvent), string(msg.FileEvent.Name))

	if _, err = file.WriteString(text); err != nil {
		panic(err)
	}

}

func fileEventToString(event *fw.FileEvent) string {

	if event.Op&fw.Chmod == fw.Chmod {
		return "M"
	}
	if event.Op&fw.Create == fw.Create {
		return "A"
	}
	if event.Op&fw.Write == fw.Write {
		return "M"
	}
	if event.Op&fw.Remove == fw.Remove {
		return "D"
	}

	return "M" // rename
}
