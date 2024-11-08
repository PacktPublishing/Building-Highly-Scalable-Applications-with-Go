package domain

import "time"

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

// Service instance uniquely identifies a running service instance
type ServiceInstance struct {
	Name        string
	Version     ServiceVersion
	PID         int
	StartupTime time.Time
}
