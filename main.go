package main

import (
	"flart/internal/commands"
	"flart/internal/config"
	"fmt"
	"os"
	"os/exec"

	"github.com/AlecAivazis/survey/v2"
)

const (
	cmdNewScreen     = "New Screen"
	cmdNewModel      = "New Model"
	cmdBuildRunner   = "Build Runner"
	cmdWatchRunner   = "Watch Runner"
	cmdMakeModel     = "make:model"
	cmdMakeScreen    = "make:screen"
	cmdBuildRunnerCL = "build:runner"
	cmdWatchRunnerCL = "watch:runner"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func run() error {
	if len(os.Args) > 2 {
		return handleMakeCommand(os.Args[1], os.Args[2])
	} else if len(os.Args) > 1 {
		return handleBuildCommand(os.Args[1])
	}
	return handleInteractive()
}

func handleBuildCommand(command string) error {
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	switch command {
	case cmdBuildRunnerCL:
		return runBuildRunner(*cfg.ProjectDir, false)
	case cmdWatchRunnerCL:
		return runBuildRunner(*cfg.ProjectDir, true)
	default:
		return fmt.Errorf("unknown command: %s", command)
	}
}

func runBuildRunner(projectDir string, watch bool) error {
	args := []string{"run", "build_runner"}
	if watch {
		args = append(args, "watch", "-d")
	} else {
		args = append(args, "build", "--delete-conflicting-outputs")
	}

	cmd := exec.Command("dart", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Dir = projectDir

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to run build_runner: %w", err)
	}

	action := "Build"
	if watch {
		action = "Watch"
	}
	fmt.Printf("%s build_runner completed successfully!\n", action)
	return nil
}

func handleMakeCommand(command, name string) error {
	switch command {
	case cmdMakeModel:
		if err := commands.CreateModel(name); err != nil {
			return fmt.Errorf("failed to create model: %w", err)
		}
		fmt.Printf("Model %s created successfully!\n", name)

	case cmdMakeScreen:
		if err := commands.CreateScreen(name); err != nil {
			return fmt.Errorf("failed to create screen: %w", err)
		}
		fmt.Printf("Screen %s created successfully!\n", name)

	default:
		return fmt.Errorf("unknown command: %s", command)
	}
	return nil
}

func handleInteractive() error {
	options := []string{
		cmdNewScreen,
		cmdNewModel,
		cmdBuildRunner,
		cmdWatchRunner,
	}

	var choice string
	prompt := &survey.Select{
		Message: "Choose an option:",
		Options: options,
	}

	if err := survey.AskOne(prompt, &choice); err != nil {
		return fmt.Errorf("failed to get user choice: %w", err)
	}

	switch choice {
	case cmdNewScreen:
		return handleNamePrompt("screen", commands.CreateScreen)

	case cmdNewModel:
		return handleNamePrompt("model", commands.CreateModel)

	case cmdBuildRunner:
		cfg, err := config.Load()
		if err != nil {
			return fmt.Errorf("failed to load config: %w", err)
		}
		return runBuildRunner(*cfg.ProjectDir, false)

	case cmdWatchRunner:
		cfg, err := config.Load()
		if err != nil {
			return fmt.Errorf("failed to load config: %w", err)
		}
		return runBuildRunner(*cfg.ProjectDir, true)
	}

	return nil
}

func handleNamePrompt(itemType string, createFn func(string) error) error {
	var name string
	if err := survey.AskOne(&survey.Input{
		Message: fmt.Sprintf("Enter %s name:", itemType),
	}, &name); err != nil {
		return fmt.Errorf("failed to get %s name: %w", itemType, err)
	}

	if err := createFn(name); err != nil {
		return fmt.Errorf("failed to create %s: %w", itemType, err)
	}
	fmt.Printf("%s %s created successfully!\n", itemType, name)
	return nil
}
