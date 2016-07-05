package config

import "gopkg.in/alecthomas/kingpin.v2"

type Config struct {
	BaseDir string
}

var (
	repoPath = kingpin.Flag("path", "Path to reposiory to watch.").Default(".").String()
)

func New() *Config {
	kingpin.Parse()

	return &Config{
		BaseDir: *repoPath,
	}

}
