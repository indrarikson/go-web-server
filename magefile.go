//go:build mage

package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"
)

// Default target to run when none is specified
// If not set, running mage will list available targets
// mg:default Build
var Default = Build.All

type Build mg.Namespace

// All generates code and builds all applications
func (Build) All() error {
	mg.SerialDeps(Generate.All, Build.Webapp, Build.TUI)
	return nil
}

// Webapp builds the webapp binary
func (Build) Webapp() error {
	fmt.Println("Building webapp...")
	return sh.RunV("go", "build", "-o", "bin/webapp", "./cmd/webapp")
}

// TUI builds the TUI tool binary
func (Build) TUI() error {
	fmt.Println("Building TUI tool...")
	return sh.RunV("go", "build", "-o", "bin/tui-tool", "./cmd/tui-tool")
}

type Generate mg.Namespace

// All runs all code generation
func (Generate) All() error {
	mg.Deps(Generate.Sqlc, Generate.Templ)
	return nil
}

// Sqlc generates type-safe Go code from SQL
func (Generate) Sqlc() error {
	fmt.Println("Generating sqlc code...")
	return sh.RunV("sqlc", "generate", "-f", "storage/sqlc.yaml")
}

// Templ generates Go code from templ files
func (Generate) Templ() error {
	fmt.Println("Generating templ code...")
	return sh.RunV("templ", "generate")
}

type Test mg.Namespace

// All runs all tests
func (Test) All() error {
	fmt.Println("Running all tests...")
	return sh.RunV("go", "test", "-v", "./...")
}

// Webapp runs webapp specific tests
func (Test) Webapp() error {
	fmt.Println("Running webapp tests...")
	return sh.RunV("go", "test", "-v", "./cmd/webapp/...", "./internal/webapp/...")
}

// TUI runs TUI tool specific tests
func (Test) TUI() error {
	fmt.Println("Running TUI tool tests...")
	return sh.RunV("go", "test", "-v", "./cmd/tui-tool/...", "./internal/tui-tool/...")
}

type Quality mg.Namespace

// All runs all quality checks
func (Quality) All() error {
	mg.Deps(Quality.Fmt, Quality.Vet, Quality.VulnCheck)
	return nil
}

// Fmt formats all Go code
func (Quality) Fmt() error {
	fmt.Println("Formatting Go code...")
	return sh.RunV("go", "fmt", "./...")
}

// Vet analyzes code for common errors
func (Quality) Vet() error {
	fmt.Println("Running go vet...")
	return sh.RunV("go", "vet", "./...")
}

// VulnCheck scans for known vulnerabilities
func (Quality) VulnCheck() error {
	fmt.Println("Running vulnerability check...")
	return sh.RunV("govulncheck", "./...")
}

type Database mg.Namespace

// Up runs database migrations
func (Database) Up() error {
	fmt.Println("Running database migrations...")
	// Use sqlite3 command directly since migrate CLI may not have sqlite support
	if err := sh.RunV("sqlite3", "./app.db", ".read ./storage/migrations/001_create_tasks.up.sql"); err != nil {
		return err
	}
	return sh.RunV("sqlite3", "./app.db", ".read ./storage/migrations/002_add_indexes.up.sql")
}

// Down rolls back the latest migration
func (Database) Down() error {
	fmt.Println("Rolling back migration...")
	return sh.RunV("migrate", "-path", "./storage/migrations", "-database", "file://./app.db", "down", "1")
}

// Reset drops all data and re-runs migrations
func (Database) Reset() error {
	fmt.Println("Resetting database...")
	if err := sh.Rm("./app.db"); err != nil && !os.IsNotExist(err) {
		return err
	}
	return Database{}.Up()
}

type Dev mg.Namespace

// Server starts the development server with live reloading
func (Dev) Server() error {
	mg.SerialDeps(Generate.All, Build.Webapp)
	fmt.Println("Starting development server...")
	return sh.RunV("./bin/webapp")
}

// TUI starts the TUI tool for development
func (Dev) TUI() error {
	mg.SerialDeps(Generate.All, Build.TUI)
	fmt.Println("Starting TUI tool...")
	return sh.RunV("./bin/tui-tool")
}

// Watch watches for file changes and rebuilds (requires entr or similar)
func (Dev) Watch() error {
	fmt.Println("Starting file watcher (requires 'entr' command)...")
	// Use shell to properly pipe commands
	return sh.RunV("sh", "-c", "find . -name '*.go' -o -name '*.templ' -o -name '*.sql' | entr -r mage dev:server")
}

// Clean removes built binaries and generated files
func Clean() error {
	fmt.Println("Cleaning up...")
	
	// Remove binaries
	if err := sh.Rm("bin"); err != nil && !os.IsNotExist(err) {
		return err
	}
	
	// Remove generated templ files
	templFiles, err := filepath.Glob("internal/**/*_templ.go")
	if err != nil {
		return err
	}
	
	for _, file := range templFiles {
		if err := sh.Rm(file); err != nil && !os.IsNotExist(err) {
			return err
		}
	}
	
	// Remove database
	if err := sh.Rm("app.db"); err != nil && !os.IsNotExist(err) {
		return err
	}
	
	fmt.Println("Clean complete!")
	return nil
}

// Setup installs required tools and dependencies
func Setup() error {
	fmt.Println("Setting up development environment...")
	
	tools := map[string]string{
		"templ":        "github.com/a-h/templ/cmd/templ@latest",
		"sqlc":         "github.com/sqlc-dev/sqlc/cmd/sqlc@latest",
		"migrate":      "github.com/golang-migrate/migrate/v4/cmd/migrate@latest",
		"govulncheck":  "golang.org/x/vuln/cmd/govulncheck@latest",
	}
	
	for tool, pkg := range tools {
		fmt.Printf("Installing %s...\n", tool)
		if err := sh.RunV("go", "install", pkg); err != nil {
			return fmt.Errorf("failed to install %s: %w", tool, err)
		}
	}
	
	// Initialize module dependencies
	fmt.Println("Downloading dependencies...")
	if err := sh.RunV("go", "mod", "download"); err != nil {
		return err
	}
	
	fmt.Println("Setup complete! Run 'mage dev:server' to start development")
	return nil
}

// CI runs all checks suitable for continuous integration
func CI() error {
	mg.SerialDeps(Generate.All, Quality.All, Test.All, Build.All)
	fmt.Println("CI pipeline completed successfully!")
	return nil
}