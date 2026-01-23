package prd

import (
	"fmt"
	"os"
)

// LoadPRD loads and parses a PRD file (markdown format)
func LoadPRD(filename string) (*PRD, error) {
	// Read file
	data, err := os.ReadFile(filename)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("PRD file not found: %s", filename)
		}
		return nil, fmt.Errorf("failed to read PRD file: %w", err)
	}

	// Parse markdown content
	prd, err := ParseMarkdown(string(data))
	if err != nil {
		return nil, fmt.Errorf("failed to parse PRD file: %w", err)
	}

	return prd, nil
}
