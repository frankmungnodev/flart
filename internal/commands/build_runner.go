package commands

import (
	"flart/internal/config"
	"flart/internal/utils"
	"fmt"
	"os"
	"os/exec"
)

func BuildRunner() error {
	// Load config
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	// Add freezed and build_runner as dev dependencies
	if err := utils.AddDependency("freezed_annotation", *cfg.ProjectDir); err != nil {
		return fmt.Errorf("failed to add freezed_annotation dependency: %w", err)
	}
	if err := utils.AddDevDependency("freezed", *cfg.ProjectDir); err != nil {
		return fmt.Errorf("failed to add freezed dependency: %w", err)
	}
	if err := utils.AddDevDependency("build_runner", *cfg.ProjectDir); err != nil {
		return fmt.Errorf("failed to add build_runner dependency: %w", err)
	}

	// Run build_runner
	cmd := exec.Command("dart", "run", "build_runner", "build", "--delete-conflicting-outputs")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Dir = *cfg.ProjectDir
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to run build_runner: %w", err)
	}

	return nil
}
