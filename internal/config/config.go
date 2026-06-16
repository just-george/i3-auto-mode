package config

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/go-playground/validator/v10"
	"go.yaml.in/yaml/v3"
)

const (
	appDirName     = "i3-auto-mode"
	configFileName = "config.yaml"
)

// DefaultPath returns the default config file location:
// $XDG_CONFIG_HOME/i3-auto-mode/config.yaml, falling back to
// ~/.config/i3-auto-mode/config.yaml if XDG_CONFIG_HOME is unset.
func DefaultPath() string {
	configHome := os.Getenv("XDG_CONFIG_HOME")
	if configHome == "" {
		homeDir, homeErr := os.UserHomeDir()
		if homeErr != nil {
			log.Printf(
				"could not determine home directory, falling back to %q relative to the current directory: %v",
				filepath.Join(appDirName, configFileName),
				homeErr,
			)
			return filepath.Join(appDirName, configFileName)
		}
		configHome = filepath.Join(homeDir, ".config")
	}
	return filepath.Join(configHome, appDirName, configFileName)
}

func FromFile(path string, validate *validator.Validate) (*Root, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("reading config file %q: %w", path, err)
	}

	var config Root
	if err = yaml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("parsing config file %q: %w", path, err)
	}

	if err = validate.Struct(&config); err != nil {
		return nil, fmt.Errorf("validating config file %q: %w", path, err)
	}

	return &config, nil
}
