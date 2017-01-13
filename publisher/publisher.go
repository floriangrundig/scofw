package publisher

import (
	"bytes"
	"fmt"
	"io"
	log_ "log"
	"os"
	"path/filepath"
	"time"

	"github.com/floriangrundig/scofw/config"
	"github.com/floriangrundig/scofw/fw"
	gitconfig "github.com/floriangrundig/scofw/git/config"
	"github.com/floriangrundig/scofw/util"
)

var (
	log *log_.Logger
)

type Message struct {
	FileEvent *fw.FileEvent
	Patch     *string
}

type ServerMessage struct {
	FileChanges       *string
	CurrentScoSession *string
	ProjectName       *string
}

type Publisher struct {
	config                    *config.Config
	gitConfig                 *gitconfig.Config
	util                      *util.Util
	fileChangedMessageChannel <-chan *Message
	serverChannel             chan<- *ServerMessage
}

func New(config *config.Config, gitConfig *gitconfig.Config, util *util.Util, fileChangedMessageChannel <-chan *Message, serverChannel chan<- *ServerMessage) *Publisher {
	log = config.Logger
	log.Println("Creating Publisher...")

	util.CreateScoFolder("logs")

	return &Publisher{
		config:    config,
		gitConfig: gitConfig,
		util:      util,
		fileChangedMessageChannel: fileChangedMessageChannel,
		serverChannel:             serverChannel,
	}
}

func (publisher *Publisher) Start() {
	go func() {
		for {
			// wait for incoming messages
			msg, ok := <-publisher.fileChangedMessageChannel
			if !ok {
				log.Println("Shutting Down Publisher")
				break
			}

			if *msg.Patch != "" {
				publisher.publishToServer(msg)
				log.Println("publish into log file")
				publisher.log(msg)
				log.Println("publish gource file")
				publisher.logInGourceFormat(msg)
			}
		}
	}()
}

func (publisher *Publisher) publishToServer(msg *Message) {

	var buf bytes.Buffer
	logger := log_.New(&buf, "", log_.Ldate|log_.Ltime)
	logger.Println(*msg.Patch)

	fileChanges := fmt.Sprint(&buf)

	transformedMsg := &ServerMessage{
		FileChanges:       &fileChanges,
		CurrentScoSession: &publisher.gitConfig.CurrentScoSession,
		ProjectName:       &publisher.config.ProjectName,
	}

	log.Print("Publish to server... ")
	publisher.serverChannel <- transformedMsg
	log.Println("[DONE]")
}

func (publisher *Publisher) log(msg *Message) {
	// TODO toopen a file and creating a logger each call is insufficient -> store loggers per session in a map to reuse them
	filename := fmt.Sprintf("%s.log", publisher.gitConfig.CurrentScoSession)

	file, err := os.OpenFile(filepath.Join(publisher.config.ScoDir, "logs", filename), os.O_CREATE|os.O_WRONLY|os.O_APPEND, publisher.config.ScoDirPermissions)
	if err != nil {
		log.Fatalln("Failed to open log file", filename, ":", err)
	}

	defer file.Close()

	var writer io.Writer

	if publisher.config.VerboseOutput {
		writer = io.MultiWriter(file, os.Stdout)
	} else {

		writer = io.Writer(file)
	}

	mylog := log_.New(writer, "", log_.Ldate|log_.Ltime)
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

	text := fmt.Sprintf("%v|Trifork|%s|%s\n", int32(time.Now().Unix()), fileEventToString(msg.FileEvent), publisher.util.ToProjectRelativePath(msg.FileEvent.Name))

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
