// In file: pkg/exporter/csv.go
package exporter

import (
	"encoding/csv"
	"fmt"
	"os"

	"bug_analyzer/kafka-finder/pkg/types" // Import our shared struct
)

// WriteToCSV saves the list of BugResult structs to a CSV file.
func WriteToCSV(filename string, results []types.BugResult) error {
	// 1. Create or truncate the file
	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("could not create file: %w", err)
	}
	defer file.Close()

	// 2. Create a CSV writer
	writer := csv.NewWriter(file)
	defer writer.Flush() // Ensure all data is written

	// 3. Write the header row
	header := []string{"BugID", "JiraURL", "CommitHash", "CommitURL"}
	if err := writer.Write(header); err != nil {
		return fmt.Errorf("could not write header: %w", err)
	}

	// 4. Write all the data rows
	for _, res := range results {
		row := []string{
			res.BugID,
			res.JiraURL,
			res.CommitHash,
			res.CommitURL,
		}
		if err := writer.Write(row); err != nil {
			return fmt.Errorf("could not write row for %s: %w", res.BugID, err)
		}
	}

	return nil
}
