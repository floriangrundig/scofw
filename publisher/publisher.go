package publisher

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"

	"github.com/floriangrundig/scofw/config"
	"github.com/floriangrundig/scofw/fw"
	gitconfig "github.com/floriangrundig/scofw/git/config"
	"github.com/floriangrundig/scofw/util"
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
			publisher.log(msg)
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
