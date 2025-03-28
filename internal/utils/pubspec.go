package utils

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

type DependencyType string

const (
	Regular DependencyType = "dependencies"
	Dev     DependencyType = "dev_dependencies"
)

type PubspecYaml struct {
	Dependencies    map[string]interface{} `yaml:"dependencies"`
	DevDependencies map[string]interface{} `yaml:"dev_dependencies"`
}

func AddDependency(dependency, projectDir string) error {
	return addDependencyHelper(dependency, projectDir, Regular)
}

func AddDevDependency(dependency, projectDir string) error {
	return addDependencyHelper(dependency, projectDir, Dev)
}

func isDependencyExists(dependency, projectDir string, depType DependencyType) (bool, error) {
	// Read pubspec.yaml
	pubspecPath := filepath.Join(projectDir, "pubspec.yaml")
	content, err := os.ReadFile(pubspecPath)
	if err != nil {
		return false, fmt.Errorf("failed to read pubspec.yaml: %w", err)
	}

	// Parse YAML content
	var pubspec PubspecYaml
	if err := yaml.Unmarshal(content, &pubspec); err != nil {
		return false, fmt.Errorf("failed to parse pubspec.yaml: %w", err)
	}

	// Check dependency existence based on type
	switch depType {
	case Dev:
		_, exists := pubspec.DevDependencies[dependency]
		return exists, nil
	default:
		_, exists := pubspec.Dependencies[dependency]
		return exists, nil
	}
}

// addDependencyHelper is a shared method to add dependencies
func addDependencyHelper(dependency, projectDir string, depType DependencyType) error {
	// Check if dependency already exists
	exists, err := isDependencyExists(dependency, projectDir, depType)
	if err != nil {
		return fmt.Errorf("failed to check dependency existence: %w", err)
	}

	// Skip if dependency exists
	if exists {
		fmt.Printf("Dependency %s already exists in %s\n", dependency, depType)
		return nil
	}

	// Prepare command arguments
	args := []string{"pub", "add"}
	if depType == Dev {
		args = append(args, fmt.Sprintf("dev:%s", dependency))
	} else {
		args = append(args, dependency)
	}

	// Execute Flutter pub add command
	cmd := exec.Command("flutter", args...)
	cmd.Dir = projectDir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to add %s dependency %s: %w", depType, dependency, err)
	}

	return nil
}

// AddFreezedDependencies adds Freezed-related dependencies to the project
func AddFreezedDependencies(projectDir string) error {
	// Add regular dependencies
	if err := AddDependency("freezed_annotation", projectDir); err != nil {
		return fmt.Errorf("failed to add freezed_annotation dependency: %w", err)
	}

	// Add dev dependencies
	devDependencies := []string{
		"freezed",
		"build_runner",
		"json_serializable",
	}

	for _, dep := range devDependencies {
		if err := AddDevDependency(dep, projectDir); err != nil {
			return fmt.Errorf("failed to add %s dependency: %w", dep, err)
		}
	}

	return nil
}

// GetFlutterPackageName retrieves the package name from pubspec.yaml
func GetFlutterPackageName(projectDir string) (string, error) {
	// Read pubspec.yaml
	pubspecPath := filepath.Join(projectDir, "pubspec.yaml")
	content, err := os.ReadFile(pubspecPath)
	if err != nil {
		return "", fmt.Errorf("failed to read pubspec.yaml: %w", err)
	}

	// Parse lines to find package name
	for _, line := range strings.Split(string(content), "\n") {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "name:") {
			return strings.TrimSpace(strings.TrimPrefix(line, "name:")), nil
		}
	}

	return "", fmt.Errorf("package name not found in pubspec.yaml")
}
