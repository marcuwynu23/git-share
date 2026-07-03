package config

import (
	"os"
	"path/filepath"
	"runtime"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Port     int    `yaml:"port"`
	ReadOnly bool   `yaml:"readonly"`
	Hostname string `yaml:"hostname"`
	Timeout  int    `yaml:"timeout"`
	Theme    string `yaml:"theme"`
}

func Default() *Config {
	return &Config{
		Port:     9720,
		ReadOnly: false,
		Hostname: "",
		Timeout:  0,
		Theme:    "dark",
	}
}

func ConfigDir() string {
	switch runtime.GOOS {
	case "windows":
		appData := os.Getenv("APPDATA")
		if appData != "" {
			return filepath.Join(appData, "git-share")
		}
		home, _ := os.UserHomeDir()
		return filepath.Join(home, "AppData", "Roaming", "git-share")
	default:
		home, _ := os.UserHomeDir()
		return filepath.Join(home, ".config", "git-share")
	}
}

func ConfigPath() string {
	return filepath.Join(ConfigDir(), "config.yaml")
}

func Load() (*Config, error) {
	cfg := Default()

	path := ConfigPath()
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return cfg, nil
		}
		return nil, err
	}

	if err := yaml.Unmarshal(data, cfg); err != nil {
		return nil, err
	}

	return cfg, nil
}

func Save(cfg *Config) error {
	dir := ConfigDir()
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	data, err := yaml.Marshal(cfg)
	if err != nil {
		return err
	}

	return os.WriteFile(ConfigPath(), data, 0644)
}

func (c *Config) Merge(c2 *Config) {
	if c2.Port != 0 && c2.Port != 9720 {
		c.Port = c2.Port
	}
	c.ReadOnly = c2.ReadOnly
	if c2.Hostname != "" {
		c.Hostname = c2.Hostname
	}
	if c2.Timeout != 0 {
		c.Timeout = c2.Timeout
	}
	if c2.Theme != "" {
		c.Theme = c2.Theme
	}
}
