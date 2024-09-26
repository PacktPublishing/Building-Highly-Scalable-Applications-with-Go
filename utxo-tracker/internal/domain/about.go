package domain

// ServiceVersion holds information used to identify the deployed service version
type ServiceVersion struct {
	CommitDate      string
	Committer       string
	CommitShortHash string
	CommitLongHash  string
	CommitSubject   string
	BuildDate       string
	GoVersion       string
}
