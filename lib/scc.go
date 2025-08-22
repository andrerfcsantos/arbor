package lib

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

// LanguageCount represents the lines of code count for a specific language
type LanguageCount struct {
	Name       string `json:"Name"`
	Lines      int    `json:"Lines"`
	Code       int    `json:"Code"`
	Comment    int    `json:"Comment"`
	Blank      int    `json:"Blank"`
	Complexity int    `json:"Complexity"`
}

// CountLinesOfCode runs the scc tool on the specified directory and returns language counts
// excludeDirs is a slice of directory names to exclude from the analysis
func CountLinesOfCode(directory string, excludeDirs []string) (map[string]int, error) {
	// Default directories to exclude if none provided
	if excludeDirs == nil {
		excludeDirs = []string{"src/test", "test"}
	}

	// Filter out non-existent directories to avoid scc failures
	var existingExcludeDirs []string
	for _, dir := range excludeDirs {
		dirPath := filepath.Join(directory, dir)
		if _, err := os.Stat(dirPath); err == nil {
			// Directory exists, add it to the list
			existingExcludeDirs = append(existingExcludeDirs, dir)
		}
		// If directory doesn't exist, simply skip it
	}

	// Build the scc command arguments
	var sccArgs []string
	sccArgs = append(sccArgs, "--uloc", "--format", "json")

	// Add exclude-dir arguments for existing directories
	// scc requires a separate --exclude-dir flag for each directory
	for _, dir := range existingExcludeDirs {
		sccArgs = append(sccArgs, "--exclude-dir", dir)
	}

	// Add the target directory
	sccArgs = append(sccArgs, directory)

	// Try using a shell command to ensure proper working directory handling
	// This should force the shell to change to the directory and then run scc
	shellCmd := fmt.Sprintf("cd %s && scc %s", directory, buildShellArgs(sccArgs))
	cmd := exec.Command("sh", "-c", shellCmd)

	output, err := cmd.Output()
	if err != nil {
		// If shell command fails, try the direct approach
		cmd = exec.Command("scc", sccArgs...)
		cmd.Dir = directory
		output, err = cmd.Output()
		if err != nil {
			return nil, fmt.Errorf("scc command failed: %w", err)
		}
	}

	var languages []LanguageCount
	err = json.Unmarshal(output, &languages)
	if err != nil {
		return nil, fmt.Errorf("failed to parse scc output: %w", err)
	}

	result := make(map[string]int)
	for _, lang := range languages {
		result[lang.Name] = lang.Lines
	}

	return result, nil
}

// buildShellArgs converts the scc arguments slice to a space-separated string for shell command
func buildShellArgs(args []string) string {
	result := ""
	for i, arg := range args {
		if i > 0 {
			result += " "
		}
		// Escape spaces in arguments if needed
		if contains(arg, " ") {
			result += fmt.Sprintf("'%s'", arg)
		} else {
			result += arg
		}
	}
	return result
}

// contains checks if a string contains a substring
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || (len(s) > len(substr) && (s[:len(substr)] == substr || s[len(s)-len(substr):] == substr || contains(s[1:], substr))))
}
