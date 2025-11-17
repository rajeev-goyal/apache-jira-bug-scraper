// In file: pkg/analysis/analyzer.go
package analysis

import (
	"bufio"
	"strings"
)

// (The kafkaKeywords variable is GONE)

// AnalyzeDiff scans a git diff string for a list of keywords.
//
// It only checks *added lines* (starting with '+').
//
// Returns true if a keyword is found, false otherwise.
func AnalyzeDiff(diffText string, keywords []string) (bool, error) {
	// Check if the user forgot to provide keywords
	if len(keywords) == 0 {
		return false, nil // No keywords, so no match
	}

	scanner := bufio.NewScanner(strings.NewReader(diffText))

	for scanner.Scan() {
		line := scanner.Text()

		// 1. We only care about lines that were *added*
		if !strings.HasPrefix(line, "+") {
			continue
		}

		// 2. Check if this new line contains any of our keywords
		for _, keyword := range keywords {
			if strings.Contains(line, keyword) {
				// Found it! This diff is a match.
				return true, nil
			}
		}
	}

	if err := scanner.Err(); err != nil {
		return false, err
	}

	// Scanned the whole diff, no keywords were found
	return false, nil
}
