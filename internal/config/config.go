package config

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

// Tab represents a single tab within a cmux workspace.
type Tab struct {
	Name    string `yaml:"name"`
	Dir     string `yaml:"dir,omitempty"`
	Command string `yaml:"command,omitempty"`
	Agent   string `yaml:"agent,omitempty"`
	Color   string `yaml:"color,omitempty"`
}

// Workspace represents a cmux workspace configuration.
type Workspace struct {
	Name  string `yaml:"name"`
	Color string `yaml:"color,omitempty"`
	Icon  string `yaml:"icon,omitempty"`
	Label string `yaml:"label,omitempty"`
	Tabs  []Tab  `yaml:"tabs"`
}

// Config is the top-level configuration structure.
type Config struct {
	Workspaces []Workspace `yaml:"workspaces"`
}

// DefaultConfigDir returns the default config directory path.
func DefaultConfigDir() string {
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".config", "cmux-kiro")
}

// DefaultConfigPath returns the default config file path.
func DefaultConfigPath() string {
	return filepath.Join(DefaultConfigDir(), "workspace.yaml")
}

// Load reads the config file and returns a Config.
func Load(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read config: %w", err)
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("failed to parse YAML: %w", err)
	}

	return &cfg, nil
}

// Save writes the Config to a YAML file.
func Save(path string, cfg *Config) error {
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return err
	}

	data, err := yaml.Marshal(cfg)
	if err != nil {
		return err
	}

	return os.WriteFile(path, data, 0644)
}

// DefaultConfig returns an example configuration.
func DefaultConfig() *Config {
	return &Config{
		Workspaces: []Workspace{
			{
				Name:  "devops",
				Color: "#4CAF50",
				Tabs: []Tab{
					{
						Name:    "infra",
						Dir:     "~/projects/terraform",
						Command: "kiro-cli chat",
					},
					{
						Name:    "query",
						Dir:     "~",
						Command: "kiro-cli chat",
						Agent:   "query",
					},
				},
			},
		},
	}
}
