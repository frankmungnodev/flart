package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type ModelConfig struct {
	UseFreezed bool `json:"useFreezed"`
}

type ScreenConfig struct {
	UseCubit   bool `json:"useCubit"`
	UseFreezed bool `json:"useFreezed"`
}

type Config struct {
	ProjectDir string       `json:"projectDir"`
	Models     ModelConfig  `json:"models"`
	Screens    ScreenConfig `json:"screens"`
}

const (
	configFileName = "flart_config.json"
	defaultUseFreezed = false
	defaultUseCubit = false
)

func Load() (*Config, error) {
	cfg, err := loadConfigFromFile()
	if err != nil {
		if os.IsNotExist(err) {
			return createDefaultConfig()
		}
		return nil, fmt.Errorf("config loading failed: %w", err)
	}

	if err := validateAndSanitizeConfig(cfg); err != nil {
		return nil, err
	}

	return cfg, nil
}

func loadConfigFromFile() (*Config, error) {
	configPath, err := getConfigPath()
	if err != nil {
		return nil, err
	}

	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, err
	}

	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("invalid config format: %w", err)
	}

	return &cfg, nil
}

func getConfigPath() (string, error) {
	configDir, err := os.UserConfigDir()
	if err != nil {
		return "", fmt.Errorf("failed to get config directory: %w", err)
	}
	return filepath.Join(configDir, "flart", configFileName), nil
}

func createDefaultConfig() (*Config, error) {
	currentDir, err := os.Getwd()
	if err != nil {
		return nil, fmt.Errorf("failed to get current directory: %w", err)
	}

	return &Config{
		ProjectDir: currentDir,
		Models:     ModelConfig{UseFreezed: defaultUseFreezed},
		Screens:    ScreenConfig{
			UseCubit:   defaultUseCubit,
			UseFreezed: defaultUseFreezed,
		},
	}, nil
}

func validateAndSanitizeConfig(cfg *Config) error {
	// Expand tilde and relative paths
	expandedPath, err := expandPath(cfg.ProjectDir)
	if err != nil {
		return fmt.Errorf("invalid project directory: %w", err)
	}
	cfg.ProjectDir = expandedPath

	// Verify directory exists
	if _, err := os.Stat(cfg.ProjectDir); err != nil {
		return fmt.Errorf("project directory validation failed: %w", err)
	}

	return nil
}

func expandPath(path string) (string, error) {
	if strings.HasPrefix(path, "~") {
		home, err := os.UserHomeDir()
		if err != nil {
			return "", err
		}
		path = filepath.Join(home, path[1:])
	}

	absPath, err := filepath.Abs(path)
	if err != nil {
		return "", fmt.Errorf("path resolution failed: %w", err)
	}

	return absPath, nil
}

func Save(cfg *Config) error {
	if err := validateAndSanitizeConfig(cfg); err != nil {
		return fmt.Errorf("cannot save invalid config: %w", err)
	}

	configPath, err := getConfigPath()
	if err != nil {
		return err
	}

	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return fmt.Errorf("config serialization failed: %w", err)
	}

	if err := os.MkdirAll(filepath.Dir(configPath), 0755); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	if err := os.WriteFile(configPath, data, 0644); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}

	return nil
}