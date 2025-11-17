package main

import (
	"encoding/json"
	"fmt"
	"os"
)

// Config struct matches our config.json
type Config struct {
	ProjectName      string   `json:"project_name"`
	RepoPath         string   `json:"local_repo_path"`
	MaxBugs          int      `json:"max_bugs_to_find"`
	OutputFile       string   `json:"output_csv_file"`
	JiraURL          string   `json:"jira_base_url"`
	RepoURL          string   `json:"repo_commit_url"`
	AnalysisKeywords []string `json:"analysis_keywords"`
}

// LoadConfig reads and parses the config.json file
func LoadConfig(filename string) (*Config, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("could not open config file %s: %w", filename, err)
	}
	defer file.Close()

	var config Config
	if err := json.NewDecoder(file).Decode(&config); err != nil {
		return nil, fmt.Errorf("could not parse config file %s: %w", config, err)
	}

	return &config, nil
}
