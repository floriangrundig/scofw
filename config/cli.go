package config

import "gopkg.in/alecthomas/kingpin.v2"

var (
	ConfigFile = kingpin.Flag("configFile", "Configuration File").Short('c').String() // TODO make it mandatory or have a sane default
	Verbose    = kingpin.Flag("verbose", "Verbose mode.").Short('v').Bool()
)

func ParseCLI() {
	kingpin.Parse()
}
