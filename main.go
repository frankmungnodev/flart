package main

import (
	"flart/internal/commands"
	"fmt"
	"os"

	"github.com/AlecAivazis/survey/v2"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func run() error {
	// Handle command line arguments
	if len(os.Args) > 2 {
		return handleCommandLine(os.Args[1], os.Args[2])
	} else if len(os.Args) > 1 {
		return handleBuildCommand(os.Args[1])
	}

	return handleInteractive()
}

func handleBuildCommand(command string) error {
	switch command {
	case "build:runner":
		if err := commands.BuildRunner(); err != nil {
			return err
		}
		fmt.Printf("Freeze build_runner done!")

	default:
		fmt.Printf("Unknown command: %s\n", command)
	}

	return nil
}

func handleCommandLine(command, name string) error {
	switch command {
	case "make:model":
		if err := commands.CreateModel(name); err != nil {
			return fmt.Errorf("failed to create model: %w", err)
		}
		fmt.Printf("Model %s created successfully!\n", name)

	case "make:screen":
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
	var choice string
	options := []string{"New Project", "New Screen", "New Model"}

	prompt := &survey.Select{
		Message: "Choose an option:",
		Options: options,
	}

	if err := survey.AskOne(prompt, &choice); err != nil {
		return fmt.Errorf("failed to get user choice: %w", err)
	}

	switch choice {
	case "New Project":
		return fmt.Errorf("project creation not implemented yet")

	case "New Screen":
		var screenName string
		if err := survey.AskOne(&survey.Input{
			Message: "Enter screen name:",
		}, &screenName); err != nil {
			return fmt.Errorf("failed to get screen name: %w", err)
		}

		if err := commands.CreateScreen(screenName); err != nil {
			return fmt.Errorf("failed to create screen: %w", err)
		}
		fmt.Printf("Screen %s created successfully!\n", screenName)

	case "New Model":
		var modelName string
		if err := survey.AskOne(&survey.Input{
			Message: "Enter model name:",
		}, &modelName); err != nil {
			return fmt.Errorf("failed to get model name: %w", err)
		}

		if err := commands.CreateModel(modelName); err != nil {
			return fmt.Errorf("failed to create model: %w", err)
		}
		fmt.Printf("Model %s created successfully!\n", modelName)
	}

	return nil
}
