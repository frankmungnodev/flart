package config

import (
	"encoding/json"
	"fmt"
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

func Load() (*Config, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}

	configPath := filepath.Join(home, "Projects", "flart", "flart_config.json")
	data, err := os.ReadFile(configPath)
	if err != nil {
		currentDir, err := os.Getwd()
		if err != nil {
			return nil, fmt.Errorf("failed to get current directory: %w", err)
		}
		return &Config{ProjectDir: &currentDir}, nil
	}

	var config Config
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to parse config file: %w", err)
	}

	if config.ProjectDir != nil && (*config.ProjectDir)[0] == '~' {
		*config.ProjectDir = filepath.Join(home, (*config.ProjectDir)[1:])
	}

	// Set default project directory if not specified
	if config.ProjectDir == nil {
		currentDir, err := os.Getwd()
		if err != nil {
			return nil, fmt.Errorf("failed to get current directory: %w", err)
		}
		config.ProjectDir = &currentDir
	}

	// Handle tilde in project directory path
	if strings.HasPrefix(*config.ProjectDir, "~") {
		expandedPath := filepath.Join(home, (*config.ProjectDir)[1:])
		config.ProjectDir = &expandedPath
	}

	// Convert to absolute path
	absPath, err := filepath.Abs(*config.ProjectDir)
	if err != nil {
		return nil, fmt.Errorf("failed to get absolute path: %w", err)
	}
	config.ProjectDir = &absPath

	// Check if directory exists
	if _, err := os.Stat(*config.ProjectDir); os.IsNotExist(err) {
		return nil, fmt.Errorf("project directory does not exist: %s", *config.ProjectDir)
	}

	// Initialize configs if nil
	if config.Models == nil {
		config.Models = &ModelConfig{}
	}
	if config.Screens == nil {
		config.Screens = &ScreenConfig{}
	}

	// Set default model config
	if config.Models.UseFreezed == nil {
		defaultVal := false
		config.Models.UseFreezed = &defaultVal
	}

	// Set default screen config
	if config.Screens.UseCubit == nil {
		defaultVal := false
		config.Screens.UseCubit = &defaultVal
	}
	if config.Screens.UseFreezed == nil {
		defaultVal := false
		config.Screens.UseFreezed = &defaultVal
	}

	return &config, nil
}

func Save(cfg *Config) error {
	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	if err := os.WriteFile("flutter_artisan_config.json", data, 0644); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}

	return nil
}
