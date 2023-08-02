package model

// Config represents values in the config file.
type Config struct {
	Source Database `yaml:"source"`
	Target Database `yaml:"target"`
}
