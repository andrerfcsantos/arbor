package lib

import (
	"encoding/json"
	"fmt"
	"os/exec"
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
func CountLinesOfCode(directory string) (map[string]int, error) {
	// Try using a shell command to ensure proper working directory handling
	// This should force the shell to change to the directory and then run scc
	shellCmd := fmt.Sprintf("cd %s && scc --uloc --format json", directory)
	cmd := exec.Command("sh", "-c", shellCmd)

	output, err := cmd.Output()
	if err != nil {
		// If shell command fails, try the direct approach
		cmd = exec.Command("scc", "--uloc", "--format", "json", directory)
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
