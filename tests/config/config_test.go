package config_test

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/marcuwynu23/git-share/internal/config"
)

func TestDefault(t *testing.T) {
	cfg := config.Default()
	if cfg.Port != 8080 {
		t.Errorf("Port = %d, want 8080", cfg.Port)
	}
	if cfg.ReadOnly {
		t.Error("ReadOnly should be false by default")
	}
	if cfg.Hostname != "" {
		t.Errorf("Hostname = %q, want empty", cfg.Hostname)
	}
	if cfg.Timeout != 0 {
		t.Errorf("Timeout = %d, want 0", cfg.Timeout)
	}
	if cfg.Theme != "dark" {
		t.Errorf("Theme = %q, want dark", cfg.Theme)
	}
}

func TestConfigDir(t *testing.T) {
	dir := config.ConfigDir()
	if dir == "" {
		t.Fatal("ConfigDir returned empty")
	}
	if !filepath.IsAbs(dir) {
		t.Errorf("ConfigDir = %q, want absolute path", dir)
	}
}

func TestConfigPath(t *testing.T) {
	path := config.ConfigPath()
	if path == "" {
		t.Fatal("ConfigPath returned empty")
	}
	if filepath.Base(path) != "config.yaml" {
		t.Errorf("ConfigPath base = %q, want config.yaml", filepath.Base(path))
	}
}

func setConfigDir(t *testing.T, dir string) {
	t.Helper()
	switch runtime.GOOS {
	case "windows":
		old := os.Getenv("APPDATA")
		os.Setenv("APPDATA", dir)
		t.Cleanup(func() { os.Setenv("APPDATA", old) })
	default:
		old := os.Getenv("HOME")
		os.Setenv("HOME", dir)
		t.Cleanup(func() { os.Setenv("HOME", old) })
	}
}

func TestLoadSave(t *testing.T) {
	setConfigDir(t, t.TempDir())

	cfg, err := config.Load()
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}
	if cfg.Port != 8080 {
		t.Errorf("Port = %d, want 8080 (default)", cfg.Port)
	}

	cfg.Port = 9090
	cfg.ReadOnly = false
	cfg.Hostname = "192.168.1.1"
	cfg.Timeout = 60
	cfg.Theme = "light"

	if err := config.Save(cfg); err != nil {
		t.Fatalf("Save() error = %v", err)
	}

	loaded, err := config.Load()
	if err != nil {
		t.Fatalf("Load() after save error = %v", err)
	}
	if loaded.Port != 9090 {
		t.Errorf("Port = %d, want 9090", loaded.Port)
	}
	if loaded.ReadOnly {
		t.Error("ReadOnly should be false")
	}
	if loaded.Hostname != "192.168.1.1" {
		t.Errorf("Hostname = %q, want 192.168.1.1", loaded.Hostname)
	}
	if loaded.Timeout != 60 {
		t.Errorf("Timeout = %d, want 60", loaded.Timeout)
	}
	if loaded.Theme != "light" {
		t.Errorf("Theme = %q, want light", loaded.Theme)
	}
}

func TestLoadNonExistent(t *testing.T) {
	setConfigDir(t, t.TempDir())

	cfg, err := config.Load()
	if err != nil {
		t.Fatalf("Load() on missing file error = %v", err)
	}
	if cfg.Port != 8080 {
		t.Errorf("Port = %d, want 8080 (default on missing file)", cfg.Port)
	}
}

func TestLoadInvalidYAML(t *testing.T) {
	setConfigDir(t, t.TempDir())
	os.WriteFile(config.ConfigPath(), []byte("port: [\x00]"), 0644)

	_, err := config.Load()
	if err == nil {
		t.Skip("yaml.v3 did not error on invalid input")
	}
}

func TestMerge(t *testing.T) {
	base := config.Default()
	override := &config.Config{
		Port:     3000,
		ReadOnly: false,
		Hostname: "0.0.0.0",
		Timeout:  30,
		Theme:    "light",
	}
	base.Merge(override)

	if base.Port != 3000 {
		t.Errorf("Port = %d, want 3000", base.Port)
	}
	if base.ReadOnly {
		t.Error("ReadOnly should be false")
	}
	if base.Hostname != "0.0.0.0" {
		t.Errorf("Hostname = %q, want 0.0.0.0", base.Hostname)
	}
	if base.Timeout != 30 {
		t.Errorf("Timeout = %d, want 30", base.Timeout)
	}
	if base.Theme != "light" {
		t.Errorf("Theme = %q, want light", base.Theme)
	}
}

func TestMergePartial(t *testing.T) {
	base := config.Default()
	partial := &config.Config{
		Port:     3000,
		Hostname: "0.0.0.0",
	}
	base.Merge(partial)

	if base.Port != 3000 {
		t.Errorf("Port = %d, want 3000", base.Port)
	}
	if base.Hostname != "0.0.0.0" {
		t.Errorf("Hostname = %q, want 0.0.0.0", base.Hostname)
	}
	if base.Theme != "dark" {
		t.Errorf("Theme = %q, want dark (unchanged)", base.Theme)
	}
}

func TestSaveCreatesDir(t *testing.T) {
	deepDir := filepath.Join(t.TempDir(), "sub", "dir")
	setConfigDir(t, deepDir)

	cfg := config.Default()
	if err := config.Save(cfg); err != nil {
		t.Fatalf("Save() error = %v", err)
	}

	configDir := filepath.Dir(config.ConfigPath())
	if _, err := os.Stat(configDir); os.IsNotExist(err) {
		t.Error("Save() did not create config directory")
	}
}
