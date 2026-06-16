package cli

import (
	"flag"
	"i3-auto-mode/internal/config"
)

type CLI struct {
	ConfigPath string
}

func New() *CLI {
	var configPath string
	flag.StringVar(&configPath, "config", config.DefaultPath(), "Path to configuration file")

	flag.Parse()

	return &CLI{
		ConfigPath: configPath,
	}
}
