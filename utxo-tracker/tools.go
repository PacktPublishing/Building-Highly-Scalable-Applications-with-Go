//go:build mage

package main

import "github.com/magefile/mage/sh"

var toolGolangciLint = runCmdV(
	"go",
	"run",
	"github.com/golangci/golangci-lint/cmd/golangci-lint@v1.61.0",
)

var toolGoArchLint = runCmdV(
	"go",
	"run",
	"github.com/fe3dback/go-arch-lint@v1.11.6",
)

var toolOAPICodegen = runCmdV(
	"go",
	"run",
	"github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen@v2.4.1",
)

var toolKo = runCmdV(
	"go",
	"run",
	"github.com/google/ko@latest",
)

var toolMigrate = runCmdV(
	"go",
	"run",
	"-tags",
	"'postgres'",
	"github.com/golang-migrate/migrate/v4/cmd/migrate@latest",
)

// runCmdV does what sh.runCmd does but outputs to STD out
func runCmdV(cmd string, args ...string) func(args ...string) error {
	return func(args2 ...string) error {
		return sh.RunV(cmd, append(args, args2...)...)
	}
}
