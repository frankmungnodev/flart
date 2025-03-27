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

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

func CreateModel(modelName string) error {
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	useFreezed :=false
	if(cfg.Models != nil && cfg.Models.UseFreezed != nil){
		useFreezed = *cfg.Models.UseFreezed;
	}

	projectDir := "."
	if(cfg.ProjectDir != nil){
		projectDir = *cfg.ProjectDir
	}
	// Convert model name to proper case
	caser := cases.Title(language.English)
	modelName = caser.String(modelName)

	// Check if files already exist
	modelDir := filepath.Join(projectDir, "lib/models")
	testDir := filepath.Join(projectDir, "test/models")

	modelFile := filepath.Join(modelDir, strings.ToLower(modelName)+".dart")
	testFile := filepath.Join(testDir, strings.ToLower(modelName)+"_test.dart")

	if utils.FileExists(modelFile) || utils.FileExists(testFile) {
		fmt.Printf("Warning: One or more files already exist:\n")
		if utils.FileExists(modelFile) {
			fmt.Printf("- %s\n", modelFile)
		}
		if utils.FileExists(testFile) {
			fmt.Printf("- %s\n", testFile)
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

	// Add required dependencies
	if err := utils.AddDependency("equatable", projectDir); err != nil {
		return fmt.Errorf("failed to add equatable dependency: %w", err)
	}

	// Add Freezed dependencies if enabled
	if useFreezed {
		if err := utils.AddFreezedDependencies(projectDir); err != nil {
			return err
		}
	}

	// Create directory structure
	modelDir = filepath.Join(projectDir, "lib/models")
	testDir = filepath.Join(projectDir, "test/models")

	dirs := []string{modelDir, testDir}
	for _, dir := range dirs {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("failed to create directory %s: %w", dir, err)
		}
	}

	// Create files with templates
	files := map[string]string{
		filepath.Join(modelDir, strings.ToLower(modelName)+".dart"):     templates.GenerateModel(modelName, useFreezed),
		filepath.Join(testDir, strings.ToLower(modelName)+"_test.dart"): templates.GenerateModelTest(modelName, useFreezed, projectDir),
	}

	for filePath, content := range files {
		if err := os.WriteFile(filePath, []byte(content), 0644); err != nil {
			return fmt.Errorf("failed to write file %s: %w", filePath, err)
		}

		// Format the created file using dart format
		cmd := exec.Command("dart", "format", filePath)
		cmd.Dir = projectDir
		if err := cmd.Run(); err != nil {
			return fmt.Errorf("failed to format file %s: %w", filePath, err)
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
