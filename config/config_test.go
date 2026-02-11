package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestEffectsOrderPersistence(t *testing.T) {
	// Create a temporary directory for test files
	tempDir, err := os.MkdirTemp("", "chroma_config_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	configPath := filepath.Join(tempDir, "test_config.toml")

	// Create config with effects order
	config := DefaultConfig()
	config.EffectsOrder = []string{"granular", "filter", "delay"}

	// Save config
	err = Save(config, configPath)
	if err != nil {
		t.Fatalf("Failed to save config: %v", err)
	}

	// Load config
	loaded, err := LoadPath(configPath)
	if err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}

	if len(loaded.EffectsOrder) != 3 ||
		loaded.EffectsOrder[0] != "granular" ||
		loaded.EffectsOrder[1] != "filter" ||
		loaded.EffectsOrder[2] != "delay" {
		t.Errorf("Expected effects order [granular, filter, delay], got %v", loaded.EffectsOrder)
	}
}
