package config

import (
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

// Config holds user preferences
type Config struct {
	DefaultBaseBranch string `yaml:"default_base_branch"`
	FeaturePrefix     string `yaml:"feature_prefix"`
	HotfixPrefix      string `yaml:"hotfix_prefix"`
}

// DefaultConfig returns default configuration
func DefaultConfig() Config {
	return Config{
		DefaultBaseBranch: "dev",
		FeaturePrefix:     "feat",
		HotfixPrefix:      "hotfix",
	}
}

// Load reads config from file
func Load() (Config, error) {
	configPath, err := getConfigPath()
	if err != nil {
		return DefaultConfig(), err
	}

	data, err := os.ReadFile(configPath)
	if err != nil {
		if os.IsNotExist(err) {
			return DefaultConfig(), nil
		}
		return DefaultConfig(), err
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return DefaultConfig(), err
	}

	return cfg, nil
}

// Save writes config to file
func Save(cfg Config) error {
	configPath, err := getConfigPath()
	if err != nil {
		return err
	}

	// Ensure directory exists
	dir := filepath.Dir(configPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	data, err := yaml.Marshal(cfg)
	if err != nil {
		return err
	}

	return os.WriteFile(configPath, data, 0644)
}

func getConfigPath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	return filepath.Join(home, ".config", "gh-flow", "config.yaml"), nil
}
