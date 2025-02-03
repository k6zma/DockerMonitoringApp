package flags

import (
	"flag"
)

type AppFlags struct {
	ConfigFilePath string
}

func ParseFlags() *AppFlags {
	configFile := flag.String("config", "lol.json", "Path to json configuration file")

	flag.Parse()

	return &AppFlags{
		ConfigFilePath: *configFile,
	}
}
