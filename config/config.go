package config

import (
	"bufio"
	"log"
	"os"

	gitignore "github.com/sabhiram/go-git-ignore"
	"gopkg.in/alecthomas/kingpin.v2"
)

type Config struct {
	BaseDir           string
	ScoDir            string
	ScoConfigFile     string
	GitIgnore         *gitignore.GitIgnore // TODO rename into something git agnostiv like FileIgnore
	ScoDirPermissions os.FileMode
}

var (
	repoPath      = kingpin.Flag("path", "Path to reposiory to watch.").Default(".").String()
	scoIgnorePath = kingpin.Flag("ignoreFile", "Path to ignore file.").Default(".gitignore").String()
)

func New() *Config {
	kingpin.Parse()

	scoDir := ".sco"

	mandatoryIgnorePatterns := []string{scoDir, ".git"}
	ignorePatterns := append(mandatoryIgnorePatterns, getIgnorePatterns()...)

	log.Println("Using following ignore patterns: ", ignorePatterns)

	ignoreObject, error := gitignore.CompileIgnoreLines(ignorePatterns...)

	if error != nil {
		panic("Error when compiling ignore lines: " + error.Error())
	}

	config := Config{
		BaseDir:           *repoPath,
		ScoDir:            scoDir,
		ScoConfigFile:     "sco.json",
		GitIgnore:         ignoreObject,
		ScoDirPermissions: 0700,
	}

	// create sco internal directory
	if _, err := os.Stat(config.ScoDir); os.IsNotExist(err) {
		errr := os.MkdirAll(config.ScoDir, config.ScoDirPermissions)
		if errr != nil {
			log.Fatal(errr)
		}
	}

	return &config
}

func getIgnorePatterns() []string {
	ignoreLines := make([]string, 0, 40)

	file, err := os.Open(*scoIgnorePath)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		ignoreLines = append(ignoreLines, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	return ignoreLines
}
