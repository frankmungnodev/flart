package commands

import (
	"bufio"
	"flart/internal/config"
	"flart/internal/templates"
	"flart/internal/utils"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
)

func CreateModel(modelName string) error {
	// Load config
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	// Validate critical config values
	if cfg.ProjectDir == nil {
		return fmt.Errorf("project directory not configured")
	}

	// Use config's UseFreezed value, which should have a default set in the config
	useFreezed := cfg.Models.UseFreezed != nil && *cfg.Models.UseFreezed

	// Prepare paths using config's project directory
	projectDir := *cfg.ProjectDir
	modelDir := filepath.Join(projectDir, "lib", "models")
	testDir := filepath.Join(projectDir, "test", "models")

	// Convert to snake case for file names
	snakeCase := utils.ToSnakeCase(modelName)
	modelFile := filepath.Join(modelDir, snakeCase+".dart")
	testFile := filepath.Join(testDir, snakeCase+"_test.dart")

	// Check existing files with user confirmation
	existingFiles := []string{}
	if utils.FileExists(modelFile) {
		existingFiles = append(existingFiles, modelFile)
	}
	if utils.FileExists(testFile) {
		existingFiles = append(existingFiles, testFile)
	}

	if len(existingFiles) > 0 {
		fmt.Println("Warning: The following files already exist:")
		for _, file := range existingFiles {
			fmt.Printf("- %s\n", file)
		}

		fmt.Print("Do you want to overwrite these files? (y/N): ")
		reader := bufio.NewReader(os.Stdin)
		response, err := reader.ReadString('\n')
		if err != nil {
			return fmt.Errorf("failed to read user input: %w", err)
		}

		response = strings.ToLower(strings.TrimSpace(response))
		if response != "y" && response != "yes" {
			return fmt.Errorf("operation cancelled by user")
		}
	}

	// Ensure directories exist
	dirsToCreate := []string{modelDir, testDir}
	for _, dir := range dirsToCreate {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("failed to create directory %s: %w", dir, err)
		}
	}

	// Add dependencies
	if err := utils.AddDependency("equatable", projectDir); err != nil {
		return fmt.Errorf("failed to add equatable dependency: %w", err)
	}

	// Add Freezed dependencies if enabled in config
	if useFreezed {
		if err := utils.AddFreezedDependencies(projectDir); err != nil {
			return fmt.Errorf("failed to add freezed dependencies: %w", err)
		}
	}

	// Prepare files to create
	files := map[string]string{
		modelFile: templates.GenerateModel(modelName, useFreezed),
		testFile:  templates.GenerateModelTest(modelName, useFreezed, projectDir),
	}

	// Write and format files
	for filePath, content := range files {
		if err := writeAndFormatFile(filePath, content, projectDir); err != nil {
			return err
		}
	}

	// Run build_runner if freezed is enabled
	if useFreezed {
		cmd := exec.Command("dart", "run", "build_runner", "build", "--delete-conflicting-outputs")
		cmd.Dir = projectDir
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			return fmt.Errorf("failed to run build_runner: %w", err)
		}
	}

	// Update barrel file
	if err := updateBarrelFile(modelDir, modelName); err != nil {
		return fmt.Errorf("failed to update barrel file: %w", err)
	}

	return nil
}

// Helper function to write and format file
func writeAndFormatFile(filePath, content, projectDir string) error {
	if err := os.WriteFile(filePath, []byte(content), 0644); err != nil {
		return fmt.Errorf("failed to write file %s: %w", filePath, err)
	}

	cmd := exec.Command("dart", "format", filePath)
	cmd.Dir = projectDir
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to format file %s: %w", filePath, err)
	}

	return nil
}

func updateBarrelFile(modelDir, modelName string) error {
	barrelPath := filepath.Join(modelDir, "models.dart")
	modelFileName := utils.ToSnakeCase(modelName)
	exportLine := fmt.Sprintf("export '%s.dart';", modelFileName)

	// Create barrel file if it doesn't exist
	if !utils.FileExists(barrelPath) {
		// Scan directory for existing model files
		entries, err := os.ReadDir(modelDir)
		if err != nil {
			return fmt.Errorf("failed to read models directory: %w", err)
		}

		var exports []string
		for _, entry := range entries {
			if entry.IsDir() {
				continue // Skip directories
			}

			name := entry.Name()
			// Only include .dart files, exclude .g.dart and .freezed.dart
			if strings.HasSuffix(name, ".dart") &&
				!strings.HasSuffix(name, ".g.dart") &&
				!strings.HasSuffix(name, ".freezed.dart") &&
				name != "models.dart" {
				exports = append(exports, fmt.Sprintf("export '%s';", name))
			}
		}

		// Only add the new model if it's not already in the exports list
		if !containsExport(exports, exportLine) {
			exports = append(exports, exportLine)
		}

		// Sort exports for consistency
		sort.Strings(exports)

		// Write barrel file with all exports
		content := strings.Join(exports, "\n") + "\n"
		return os.WriteFile(barrelPath, []byte(content), 0644)
	}

	// Read existing barrel file
	content, err := os.ReadFile(barrelPath)
	if err != nil {
		return fmt.Errorf("failed to read barrel file: %w", err)
	}

	// Check if export already exists
	if strings.Contains(string(content), exportLine) {
		return nil
	}

	// Append new export
	f, err := os.OpenFile(barrelPath, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("failed to open barrel file: %w", err)
	}
	defer f.Close()

	if _, err := f.WriteString(exportLine + "\n"); err != nil {
		return fmt.Errorf("failed to update barrel file: %w", err)
	}

	return nil
}

// Helper function to check if an export line already exists in the slice
func containsExport(exports []string, target string) bool {
	for _, export := range exports {
		if export == target {
			return true
		}
	}
	return false
}
