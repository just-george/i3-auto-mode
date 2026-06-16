package config

type Root struct {
	DefaultMode string                `yaml:"default_mode" validate:"required"`
	Modes       map[string]ModeConfig `yaml:"modes" validate:"min=1,dive"`
}

type ModeConfig struct {
	ApplicationClasses []string `yaml:"application_classes" validate:"min=1"`
}
