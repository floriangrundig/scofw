package config

import (
	"bufio"
	log_ "log"
	"os"

	gitignore "github.com/sabhiram/go-git-ignore"
)

type Config struct {
	ProjectDir        string
	ScoDir            string
	ScoConfigFile     string
	VerboseOutput     bool
	Logger            *log_.Logger
	GitIgnore         *gitignore.GitIgnore // TODO rename into something git agnostiv like FileIgnore
	ScoDirPermissions os.FileMode
}

var (
	log *log_.Logger
)

func GetVerboseLoggingFlag() bool {
	log_.Println("verbose flag: ", *Verbose)
	return *Verbose
}

func New(globalConfig *GlobalConfig, projectDefinition *ProjectConfig, logger *log_.Logger) *Config {
	log = logger

	mandatoryIgnorePatterns := []string{".sco", ".git", "*___jb_*"}
	ignorePatterns := append(mandatoryIgnorePatterns, getIgnorePatterns(globalConfig, projectDefinition.IgnoreFiles)...)
	log.Println("Using following ignore patterns: ", ignorePatterns)

	ignoreObject, error := gitignore.CompileIgnoreLines(ignorePatterns...)

	if error != nil {
		log.Fatal("Error when compiling ignore lines: " + error.Error())
	}

	config := Config{
		ProjectDir:        projectDefinition.Path,
		ScoDir:            projectDefinition.ScoDir,
		ScoConfigFile:     "sco.json",
		GitIgnore:         ignoreObject,
		ScoDirPermissions: 0700,
		VerboseOutput:     *Verbose,
		Logger:            log,
	}

	return &config
}

func getIgnorePatterns(globalConfig *GlobalConfig, ignoreFiles []string) []string {
	ignoreLines := make([]string, 0, 4)

	for _, ignoreFile := range ignoreFiles {

		// file, err := os.Open(*scoIgnorePath)
		file, err := os.Open(ignoreFile)
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
	}

	for _, ignorePattern := range globalConfig.IgnorePatterns {
		ignoreLines = append(ignoreLines, ignorePattern)
	}

	return ignoreLines
}
