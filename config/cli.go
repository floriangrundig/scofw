package config

import "gopkg.in/alecthomas/kingpin.v2"

var (
	// TODO we should add one flag to add a list of ignore patterns from command line e.g. -ignore="out target"
	// scoIgnorePath = kingpin.Flag("ignoreFile", "Path to ignore file.").Default(".gitignore").String()
	ConfigFile = kingpin.Flag("configFile", "Configuration File").Short('c').String()
	Verbose    = kingpin.Flag("verbose", "Verbose mode.").Short('v').Bool()
)

func ParseCLI() {
	kingpin.Parse()
}
