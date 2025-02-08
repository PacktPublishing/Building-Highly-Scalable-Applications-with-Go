//go:build mage

package main

import (
	"context"
	"fmt"
	"path"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/google/go-containerregistry/pkg/name"
	"github.com/google/go-containerregistry/pkg/v1/remote"
	"github.com/google/ko/pkg/build"
	"github.com/google/ko/pkg/publish"
	"github.com/hannesdejager/utxo-tracker/internal/domain"
	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"

	_ "github.com/golang-migrate/migrate/v4/database/postgres"
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

type Image mg.Namespace

// Account_service creates a Docker image for the account service
func (Image) Account_service() error {
	mg.Deps(Generate)
	v, e := versionInfo()
	if e != nil {
		return e
	}
	name, e := koImg("account-service", v)
	if e != nil {
		return fmt.Errorf("could not build image: %w", e)
	}
	fmt.Println(name)
	return nil
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

func ldflags(v domain.ServiceVersion) string {
	bindPath := fmt.Sprintf("%s/internal/infra/linker",
		module())
	return fmt.Sprintf(
		"-X '%[1]s.CommitShortHash=%[2]s' "+
			"-X '%[1]s.CommitLongHash=%[3]s' "+
			"-X '%[1]s.CommitDate=%[4]s' "+
			"-X '%[1]s.CommitSubject=%[5]s' "+
			"-X '%[1]s.Committer=%[6]s' "+
			"-X '%[1]s.BuildDate=%[7]s' "+
			"-X '%[1]s.GoVersion=%[8]s' "+
			"-s -w",
		bindPath,
		v.CommitShortHash,
		v.CommitLongHash,
		v.CommitDate,
		strings.Replace(v.CommitSubject, "'", "`", -1),
		strings.Replace(v.Committer, "'", "", -1),
		v.BuildDate,
		v.GoVersion)
}

func buildCmd(service string, v domain.ServiceVersion) error {
	env := map[string]string{"CGO_ENABLED": "0"}
	return sh.RunWithV(env, "go", "build", "-o", service,
		"-ldflags", ldflags(v),
		filepath.Join("cmd", service, "main.go"))
}

func imageNamer(base, importpath string) string {
	return path.Join(base, path.Base(importpath))
}

func koImg(service string, v domain.ServiceVersion) (
	name.Reference, error) {
	ctx := context.Background()

	b, err := build.NewGo(ctx, ".",
		build.WithPlatforms("linux/amd64"),
		build.WithDefaultLdflags(
			strings.Split(ldflags(v), " "),
		),
		build.WithBaseImages(func(
			ctx context.Context,
			_ string) (name.Reference, build.Result, error) {
			ref := name.MustParseReference(
				"cgr.dev/chainguard/static:latest")
			base, err := remote.Index(
				ref, remote.WithContext(ctx))
			return ref, base, err
		}),
	)
	if err != nil {
		return nil,
			fmt.Errorf("could not create image builder: %w", err)
	}

	importPath, err := b.QualifyImport("./cmd/" + service)
	if err != nil {
		return nil,
			fmt.Errorf("failed to qualify import path: %w", err)
	}

	r, err := b.Build(ctx, importPath)
	if err != nil {
		return nil,
			fmt.Errorf("failed to build image: %w", err)
	}

	p, err := publish.NewDaemon(
		imageNamer,
		[]string{"latest"},
		publish.WithLocalDomain("utxo-tracker"),
	)
	if err != nil {
		return nil,
			fmt.Errorf("failed to create publisher: %w", err)
	}
	defer p.Close()

	ref, err := p.Publish(ctx, r, importPath)
	if err != nil {
		return nil,
			fmt.Errorf("failed to publish: %w", err)
	}
	return ref, nil
}

type Account_service mg.Namespace

func (Account_service) Db_up() error {
	mg.Deps(Generate)
	return toolMigrate(
		"-source",
		"file://internal/infra/aspostgres/schema",
		"-database",
		"postgresql://as:as@localhost:5432/as?sslmode=disable",
		"up",
	)
}

func (Account_service) Db_down() error {
	mg.Deps(Generate)
	return toolMigrate(
		"-source",
		"file://internal/infra/aspostgres/schema",
		"-database",
		"postgresql://as:as@localhost:5432/as?sslmode=disable",
		"down",
	)
}
