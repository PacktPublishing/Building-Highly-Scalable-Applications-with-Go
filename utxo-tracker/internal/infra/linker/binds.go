// Package linker is where we dance with the Go linker.
package linker

// These variables will be auto assigned by the Go linker. To find where
// this magic happens search for 'ldflags' in the Mage build file: magefile.go.
// Here is a date format example: 2021-09-03 23:16:24
var (
	CommitShortHash string
	CommitLongHash  string
	CommitDate      string
	CommitSubject   string
	Committer       string
	BuildDate       string
	GoVersion       string
)