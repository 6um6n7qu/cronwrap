// Package config provides configuration loading and validation for cronwrap.
package config

import (
	"errors"
	"os"
	"time"

	"gopkg.in/yaml.v3"
)

// Config holds the full cronwrap configuration.
type Config struct {
	Job      JobConfig      `yaml:"job"`
	Notifier NotifierConfig `yaml:"notifier"`
	Logger   LoggerConfig   `yaml:"logger"`
}

// JobConfig describes the command to run and its constraints.
type JobConfig struct {
	Name    string        `yaml:"name"`
	Command string        `yaml:"command"`
	Args    []string      `yaml:"args"`
	Timeout time.Duration `yaml:"timeout"`
}

// NotifierConfig holds SMTP alerting settings.
type NotifierConfig struct {
	SMTPHost   string   `yaml:"smtp_host"`
	SMTPPort   int      `yaml:"smtp_port"`
	From       string   `yaml:"from"`
	Recipients []string `yaml:"recipients"`
}

// LoggerConfig controls log output behaviour.
type LoggerConfig struct {
	Level  string `yaml:"level"`
	Output string `yaml:"output"` // "stdout" or a file path
}

// Load reads and parses a YAML config file at the given path.
func Load(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}

	if err := cfg.validate(); err != nil {
		return nil, err
	}

	return &cfg, nil
}

// validate checks that required fields are present and values are sane.
func (c *Config) validate() error {
	if c.Job.Name == "" {
		return errors.New("config: job.name is required")
	}
	if c.Job.Command == "" {
		return errors.New("config: job.command is required")
	}
	if c.Job.Timeout < 0 {
		return errors.New("config: job.timeout must not be negative")
	}
	if c.NotifierConfig().SMTPPort < 0 || c.NotifierConfig().SMTPPort > 65535 {
		return errors.New("config: notifier.smtp_port must be between 0 and 65535")
	}
	return nil
}

// NotifierConfig returns the notifier sub-config (helper for readability).
func (c *Config) NotifierConfig() NotifierConfig {
	return c.Notifier
}
