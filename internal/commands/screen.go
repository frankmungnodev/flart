package commands

import (
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

func CreateScreen(screenName string) error {
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	// Add required dependencies
	dependencies := []string{
		"flutter_bloc",
		"equatable",
	}

	// Add Freezed dependencies if enabled
	if cfg.Screens.UseFreezed {
		if err := utils.AddFreezedDependencies(cfg.ProjectDir); err != nil {
			return err
		}
	}

	// Add dependencies
	for _, dep := range dependencies {
		if err := utils.AddDependency(dep, cfg.ProjectDir); err != nil {
			return fmt.Errorf("failed to add dependency %s: %w", dep, err)
		}
	}

	// Convert screen name to proper case
	caser := cases.Title(language.English)
	screenName = caser.String(screenName)

	// Create directory structure
	screenDir := filepath.Join(cfg.ProjectDir, "lib/screens", strings.ToLower(screenName))
	var stateDir string
	if cfg.Screens.UseCubit {
		stateDir = filepath.Join(screenDir, "cubit")
	} else {
		stateDir = filepath.Join(screenDir, "bloc")
	}

	dirs := []string{screenDir, stateDir}
	for _, dir := range dirs {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("failed to create directory %s: %w", dir, err)
		}
	}

	// Create files with templates
	var files map[string]string
	if cfg.Screens.UseCubit {
		files = map[string]string{
			filepath.Join(screenDir, strings.ToLower(screenName)+"_screen.dart"): templates.GenerateScreen(screenName, true),
			filepath.Join(stateDir, strings.ToLower(screenName)+"_cubit.dart"):   templates.GenerateCubit(screenName),
			filepath.Join(stateDir, strings.ToLower(screenName)+"_state.dart"):   templates.GenerateState(screenName, true, cfg.Screens.UseFreezed),
		}
	} else {
		files = map[string]string{
			filepath.Join(screenDir, strings.ToLower(screenName)+"_screen.dart"): templates.GenerateScreen(screenName, false),
			filepath.Join(stateDir, strings.ToLower(screenName)+"_bloc.dart"):    templates.GenerateBloc(screenName),
			filepath.Join(stateDir, strings.ToLower(screenName)+"_event.dart"):   templates.GenerateEvent(screenName),
			filepath.Join(stateDir, strings.ToLower(screenName)+"_state.dart"):   templates.GenerateState(screenName, false, cfg.Screens.UseFreezed),
		}
	}

	for filePath, content := range files {
		if err := os.WriteFile(filePath, []byte(content), 0644); err != nil {
			return fmt.Errorf("failed to write file %s: %w", filePath, err)
		}

		// Format the created file using dart format
		cmd := exec.Command("dart", "format", filePath)
		cmd.Dir = cfg.ProjectDir
		if err := cmd.Run(); err != nil {
			return fmt.Errorf("failed to format file %s: %w", filePath, err)
		}
	}

	// Run build_runner if freezed is enabled
	if cfg.Screens.UseFreezed {
		cmd := exec.Command("dart", "run", "build_runner", "build", "--delete-conflicting-outputs")
		cmd.Dir = cfg.ProjectDir
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			return fmt.Errorf("failed to run build_runner: %w", err)
		}
	}

	return nil
}
