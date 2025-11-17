// In file: kafka-finder/main.go
package main

import (
	"fmt"
	"log"
	"path"

	// Import all our local packages
	"bug_analyzer/kafka-finder/pkg/analysis"
	"bug_analyzer/kafka-finder/pkg/exporter"
	"bug_analyzer/kafka-finder/pkg/git"
	"bug_analyzer/kafka-finder/pkg/jira"
	"bug_analyzer/kafka-finder/pkg/types"
)

// (All consts are GONE)

func main() {
	log.Println("Starting Bug Finder...")

	// 1. Load Configuration
	config, err := LoadConfig("config.json")
	if err != nil {
		log.Fatalf("Failed to load config.json: %v", err)
	}

	log.Printf("Configuration loaded for project: %s", config.ProjectName)
	log.Printf("Analyzing local repo: %s", config.RepoPath)

	// 2. Get Bugs from JIRA
	jiraClient := jira.NewClient()
	log.Printf("Searching for %d bugs...", config.MaxBugs)

	bugKeys, err := jiraClient.SearchBugs(config.ProjectName, config.MaxBugs)
	if err != nil {
		log.Fatalf("Error searching JIRA: %v", err)
	}
	log.Printf("Found %d bug keys. Starting analysis...\n", len(bugKeys))

	var finalResults []types.BugResult

	// 3. Process each bug
	for _, bugID := range bugKeys {
		// Step 3a: Find the commit in Git
		commitHash, err := git.FindCommit(config.RepoPath, bugID)
		if err != nil {
			log.Printf("  - %s: ERROR searching git: %v\n", bugID, err)
			continue
		}
		if commitHash == "" {
			continue
		}

		if len(config.AnalysisKeywords) == 0 {
			log.Printf("  - %s: MATCH FOUND (no keywords) -> %s\n", bugID, commitHash)
			result := types.BugResult{
				BugID:      bugID,
				CommitHash: commitHash,
				JiraURL:    path.Join(config.JiraURL, bugID),
				CommitURL:  path.Join(config.RepoURL, commitHash),
			}
			finalResults = append(finalResults, result)
			continue // Skip to the next bug
		}

		// Step 3b: Get the diff *only for test files*
		diff, err := git.GetCommitDiff(config.RepoPath, commitHash)
		if err != nil {
			log.Printf("  - %s: ERROR getting diff for %s: %v\n", bugID, commitHash, err)
			continue
		}
		if diff == "" {
			continue
		}

		// Step 3c: Analyze the diff (using config keywords)
		foundMatch, err := analysis.AnalyzeDiff(diff, config.AnalysisKeywords)
		if err != nil {
			log.Printf("  - %s: ERROR analyzing diff for %s: %v\n", bugID, commitHash, err)
			continue
		}

		if foundMatch {
			log.Printf("  - %s: MATCH FOUND! -> %s\n", bugID, commitHash)

			result := types.BugResult{
				BugID:      bugID,
				CommitHash: commitHash,
				JiraURL:    path.Join(config.JiraURL, bugID),
				CommitURL:  path.Join(config.RepoURL, commitHash),
			}
			finalResults = append(finalResults, result)
		}
	}

	// 4. Print summary and export to CSV
	log.Println("--------------------------------------------------")
	log.Printf("Analysis complete. Found %d matching bugs.\n", len(finalResults))
	log.Println("--------------------------------------------------")

	if len(finalResults) == 0 {
		log.Println("No bugs matched all criteria.")
	} else {
		for _, result := range finalResults {
			fmt.Printf("Bug: %s -> %s\n", result.BugID, result.CommitURL)
		}

		log.Printf("\nSaving results to %s...", config.OutputFile)
		if err := exporter.WriteToCSV(config.OutputFile, finalResults); err != nil {
			log.Fatalf("Error writing to CSV: %v", err)
		}
		log.Printf("Successfully saved results to %s\n", config.OutputFile)
	}
}
