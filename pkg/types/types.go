// In file: pkg/types/types.go
package types

// BugResult holds the final data for a single matching bug
type BugResult struct {
	BugID      string
	CommitHash string
	JiraURL    string
	CommitURL  string
}
