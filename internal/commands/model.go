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
