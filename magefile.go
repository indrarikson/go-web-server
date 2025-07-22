//go:build mage

package main

import (
	"fmt"
	"os"

	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"
)

// Default target to run when none is specified
var Default = Build

// Build generates code and builds the server binary
func Build() error {
	mg.SerialDeps(Generate, buildServer)
	return nil
}

func buildServer() error {
	fmt.Println("Building server...")
	if err := sh.Run("mkdir", "-p", "bin"); err != nil {
		return err
	}
	return sh.RunV("go", "build", "-ldflags=-s -w", "-o", "./bin/server", "./cmd/web")
}

// Generate runs all code generation
func Generate() error {
	mg.Deps(generateSqlc, generateTempl)
	return nil
}

func generateSqlc() error {
	fmt.Println("Generating sqlc code...")
	return sh.RunV("sqlc", "generate")
}

func generateTempl() error {
	fmt.Println("Generating templ code...")
	return sh.RunV("go", "generate", "./...")
}

// Test runs all tests
func Test() error {
	fmt.Println("Running tests...")
	return sh.RunV("go", "test", "-race", "-cover", "./...")
}

// Fmt formats and tidies code
func Fmt() error {
	fmt.Println("Tidying and formatting...")
	if err := sh.RunV("go", "mod", "tidy"); err != nil {
		return err
	}
	return sh.RunV("go", "fmt", "./...")
}

// Vet analyzes code for common errors
func Vet() error {
	fmt.Println("Running go vet...")
	return sh.RunV("go", "vet", "./...")
}

// VulnCheck scans for known vulnerabilities
func VulnCheck() error {
	fmt.Println("Running vulnerability check...")
	return sh.RunV("govulncheck", "./...")
}

// Run builds and runs the server
func Run() error {
	mg.SerialDeps(Build)
	fmt.Println("Starting server...")
	return sh.RunV("./bin/server")
}

// Dev starts development server with hot reload
func Dev() error {
	fmt.Println("Starting development server with hot reload...")
	
	// Check if air is installed, install if not
	if err := sh.Run("which", "air"); err != nil {
		fmt.Println("Installing air...")
		if err := sh.RunV("go", "install", "github.com/air-verse/air@latest"); err != nil {
			return err
		}
	}
	
	return sh.RunV("air")
}

// Clean removes built binaries and generated files
func Clean() error {
	fmt.Println("Cleaning up...")
	
	// Remove binaries
	if err := sh.Rm("bin"); err != nil && !os.IsNotExist(err) {
		return err
	}
	
	// Remove tmp directory
	if err := sh.Rm("tmp"); err != nil && !os.IsNotExist(err) {
		return err
	}
	
	fmt.Println("Clean complete!")
	return nil
}

// Setup installs required development tools
func Setup() error {
	fmt.Println("Setting up development environment...")
	
	tools := map[string]string{
		"templ":        "github.com/a-h/templ/cmd/templ@latest",
		"sqlc":         "github.com/sqlc-dev/sqlc/cmd/sqlc@latest",
		"govulncheck":  "golang.org/x/vuln/cmd/govulncheck@latest",
		"air":          "github.com/air-verse/air@latest",
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
	
	fmt.Println("Setup complete! Run 'mage dev' to start development with hot reload")
	return nil
}

// Lint runs all linters (vet + security scan)
func Lint() error {
	mg.Deps(Vet, VulnCheck)
	return nil
}

// CI runs all checks suitable for continuous integration
func CI() error {
	mg.SerialDeps(Generate, Fmt, Vet, VulnCheck, Test, Build)
	fmt.Println("CI pipeline completed successfully!")
	return nil
}