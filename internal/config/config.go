package config

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
)

type ModelConfig struct {
	UseFreezed *bool `json:"useFreezed"`
}

type ScreenConfig struct {
	UseCubit   *bool `json:"useCubit"`
	UseFreezed *bool `json:"useFreezed"`
}

type Config struct {
	ProjectDir *string       `json:"projectDir"`
	Models     *ModelConfig  `json:"models"`
	Screens    *ScreenConfig `json:"screens"`
}

// configFileName is consistent across save and load operations
const configFileName = "flart_config.json"

func Load() (*Config, error) {
	// Create a default config with explicit default values
	cfg := &Config{
		ProjectDir: new(string),
		Models: &ModelConfig{
			UseFreezed: new(bool),
		},
		Screens: &ScreenConfig{
			UseCubit:   new(bool),
			UseFreezed: new(bool),
		},
	}

	// Set default values explicitly
	*cfg.ProjectDir = "."
	*cfg.Models.UseFreezed = false
	*cfg.Screens.UseCubit = false
	*cfg.Screens.UseFreezed = false

	// Determine the config file path
	currentDir, err := os.Getwd()
	if err != nil {
		return nil, fmt.Errorf("failed to get current working directory: %w", err)
	}

	configPath := filepath.Join(currentDir, configFileName)

	// Check if config file exists
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Printf("Config file not found at %s, using defaults", configPath)
		return cfg, nil
	}

	// Read config file
	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file %s: %w", configPath, err)
	}

	// Unmarshal config
	if err := json.Unmarshal(data, cfg); err != nil {
		return nil, fmt.Errorf("failed to parse config file %s: %w", configPath, err)
	}

	// Resolve project directory
	if cfg.ProjectDir != nil {
		// Handle home directory expansion
		if strings.HasPrefix(*cfg.ProjectDir, "~/") {
			homeDir, err := os.UserHomeDir()
			if err != nil {
				return nil, fmt.Errorf("failed to get home directory: %w", err)
			}
			*cfg.ProjectDir = filepath.Join(homeDir, (*cfg.ProjectDir)[2:])
		}

		// Convert to absolute path
		absPath, err := filepath.Abs(*cfg.ProjectDir)
		if err != nil {
			return nil, fmt.Errorf("failed to resolve project directory: %w", err)
		}
		*cfg.ProjectDir = absPath
	}

	return cfg, nil
}

func Save(cfg *Config) error {
	// Ensure the config directory exists
	currentDir, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get current working directory: %w", err)
	}

	configPath := filepath.Join(currentDir, configFileName)

	// Marshal config with indentation
	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	// Write config file with more permissive permissions
	// 0644 allows read/write for owner, read for others
	if err := os.WriteFile(configPath, data, 0644); err != nil {
		return fmt.Errorf("failed to write config file %s: %w", configPath, err)
	}

	log.Printf("Config saved to %s", configPath)
	return nil
}
