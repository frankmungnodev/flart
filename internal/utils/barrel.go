package utils

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

func UpdateBarrelFile(dir, name, barrelFileName string) error {
	barrelPath := filepath.Join(dir, barrelFileName)
	fileName := ToSnakeCase(name)
	exportLine := fmt.Sprintf("export '%s.dart';", fileName)

	// Create barrel file if it doesn't exist
	if !FileExists(barrelPath) {
		// Scan directory for existing files
		entries, err := os.ReadDir(dir)
		if err != nil {
			return fmt.Errorf("failed to read directory: %w", err)
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
				name != barrelFileName {
				exports = append(exports, fmt.Sprintf("export '%s';", name))
			}
		}

		// Only add the new file if it's not already in the exports list
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

func UpdateScreenBarrelFile(dir, name, barrelFileName string) error {
	barrelPath := filepath.Join(dir, barrelFileName)
	screenName := ToSnakeCase(name)
	exportLine := fmt.Sprintf("export '%s/%s.dart';", screenName, screenName)

	// Create barrel file if it doesn't exist
	if !FileExists(barrelPath) {
		// Scan directory for screen subdirectories
		entries, err := os.ReadDir(dir)
		if err != nil {
			return fmt.Errorf("failed to read directory: %w", err)
		}

		var exports []string
		for _, entry := range entries {
			if !entry.IsDir() {
				continue // Skip non-directories
			}

			screenDir := entry.Name()
			screenFile := filepath.Join(screenDir, screenDir+".dart")
			if FileExists(filepath.Join(dir, screenFile)) {
				exports = append(exports, fmt.Sprintf("export '%s';", screenFile))
			}
		}

		// Only add the new file if it's not already in the exports list
		if !containsExport(exports, exportLine) {
			exports = append(exports, exportLine)
		}

		// Sort exports for consistency
		sort.Strings(exports)

		// Write barrel file with all exports
		content := strings.Join(exports, "\n") + "\n"
		return os.WriteFile(barrelPath, []byte(content), 0644)
	}

	// Rest of the function remains the same
	content, err := os.ReadFile(barrelPath)
	if err != nil {
		return fmt.Errorf("failed to read barrel file: %w", err)
	}

	if strings.Contains(string(content), exportLine) {
		return nil
	}

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
