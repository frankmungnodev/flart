package utils

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

type dependencyType string

const (
	regular dependencyType = "dependencies"
	dev     dependencyType = "dev_dependencies"
)

func AddDependency(dependency, projectDir string) error {
	return addDependencyHelper(dependency, projectDir, regular)
}

func AddDevDependency(dependency, projectDir string) error {
	return addDependencyHelper(dependency, projectDir, dev)
}

func addDependencyHelper(dependency string, projectDir string, depType dependencyType) error {
	args := []string{"pub", "add"}

	if depType == dev {
		args = append(args, fmt.Sprintf("dev:%s", dependency))
	} else {
		args = append(args, dependency)
	}

	cmd := exec.Command("flutter", args...)
	cmd.Dir = projectDir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to add %s dependency %s: %w", depType, dependency, err)
	}

	return nil
}

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

func GetFlutterPackageName(projectDir string) (string, error) {
	pubspecPath := filepath.Join(projectDir, "pubspec.yaml")
	content, err := os.ReadFile(pubspecPath)
	if err != nil {
		return "", fmt.Errorf("failed to read pubspec.yaml: %w", err)
	}

	lines := strings.Split(string(content), "\n")
	for _, line := range lines {
		if strings.HasPrefix(strings.TrimSpace(line), "name:") {
			return strings.TrimSpace(strings.TrimPrefix(strings.TrimSpace(line), "name:")), nil
		}
	}

	return "", fmt.Errorf("package name not found in pubspec.yaml")
}
