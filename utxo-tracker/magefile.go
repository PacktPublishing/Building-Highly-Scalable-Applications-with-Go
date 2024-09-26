//go:build mage

package main

import (
	"fmt"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/hannesdejager/utxo-tracker/internal/domain"
	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"
)

// The path to the generated Go code for the account server REST API.
const RestAPIFile = "internal/infra/api/restv1/api.gen.go"

// Run executes the account-service
func Run() error {
	mg.Deps(Generate)
	return sh.RunV(
		"go",
		"run",
		"./cmd/account-service/main.go",
	)
}

// Fmt formats the Go code in the project
func Fmt() error {
	e := sh.RunV("go", "fmt", "./...")
	if e != nil {
		return e
	}
	return sh.RunV("go", "fmt", "magefile.go")
}

// Check does static analysis
func Check() error {
	mg.Deps(Generate)
	fmt.Println("=== Architecture Linter ===\n")
	err := toolGoArchLint(
		"check",
		"--arch-file",
		"arch.yaml",
	)
	if err != nil {
		return err
	}

	fmt.Println("=== Code Linter ===\n")
	return toolGolangciLint(
		"run",
		"-c",
		"golangci.yaml",
	)
}

// Build compiles the binaries
func Build() error {
	mg.Deps(Generate)
	v, e := versionInfo()
	if e != nil {
		return e
	}
	return buildCmd("account-service", v)
}

// Clean removes build artifacts
func Clean() error {
	e := sh.RunV("rm", "-f", "account-service")
	if e != nil {
		return e
	}
	return sh.RunV("rm", "-f", RestAPIFile)
}

// Generate does code generation
func Generate() error {
	return sh.RunV(
		"go", "generate", "./...",
	)
}

// module returns the Go module name
func module() string {
	m, _ := sh.Output("go", "list", "-m")
	return m
}

// versionInfo retrieves versioning information from Git
func versionInfo() (r domain.ServiceVersion, e error) {
	gitRP := sh.OutCmd("git", "rev-parse")
	r.CommitShortHash, e = gitRP("--short", "HEAD")
	if e != nil {
		return
	}
	r.CommitLongHash, e = gitRP("HEAD")
	if e != nil {
		return
	}
	gitShow := sh.OutCmd("git", "show", "-s")
	r.CommitDate, e = gitShow("--format=%ci",
		r.CommitShortHash)
	if e != nil {
		return
	}
	r.CommitSubject, e = gitShow("--format=%s",
		r.CommitShortHash)
	if e != nil {
		return
	}
	r.Committer, e = gitShow("--format='%cn <%ce>'",
		r.CommitShortHash)
	if e != nil {
		return
	}
	r.GoVersion = runtime.Version()
	r.BuildDate = time.Now().Format("2006-01-02 15:04:05")
	return
}

func buildCmd(service string, v domain.ServiceVersion) error {
	bindPath := fmt.Sprintf("%s/internal/infra/linker", module())
	return sh.RunV("go", "build", "-o", service,
		"-ldflags", fmt.Sprintf(
			"-X '%[1]s.CommitShortHash=%[2]s' "+
				"-X '%[1]s.CommitLongHash=%[3]s' "+
				"-X '%[1]s.CommitDate=%[4]s' "+
				"-X '%[1]s.CommitSubject=%[5]s' "+
				"-X '%[1]s.Committer=%[6]s' "+
				"-X '%[1]s.BuildDate=%[7]s' "+
				"-X '%[1]s.GoVersion=%[8]s'",
			bindPath,
			v.CommitShortHash,
			v.CommitLongHash,
			v.CommitDate,
			strings.Replace(v.CommitSubject, "'", "`", -1),
			strings.Replace(v.Committer, "'", "", -1),
			v.BuildDate,
			v.GoVersion),
		filepath.Join("cmd", service, "main.go"))
}
