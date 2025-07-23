//go:build mage

package main

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"
)

const (
	packageName = "github.com/dunamismax/go-web-server"
	binaryName  = "server"
	buildDir    = "bin"
	tmpDir      = "tmp"
)

// Default target to run when none is specified
var Default = Build

// Build generates code and builds the server binary
func Build() error {
	mg.SerialDeps(Generate, buildServer)
	return nil
}

func buildServer() error {
	fmt.Println("🔨 Building server...")
	
	if err := sh.Run("mkdir", "-p", buildDir); err != nil {
		return fmt.Errorf("failed to create build directory: %w", err)
	}

	ldflags := "-s -w -X main.version=1.0.0 -X main.buildTime=" + getCurrentTime()
	binaryPath := filepath.Join(buildDir, binaryName)
	
	// Add .exe extension on Windows
	if runtime.GOOS == "windows" {
		binaryPath += ".exe"
	}

	return sh.RunV("go", "build", "-ldflags="+ldflags, "-o", binaryPath, "./cmd/web")
}

func getCurrentTime() string {
	output, err := sh.Output("date", "-u", "+%Y-%m-%dT%H:%M:%SZ")
	if err != nil {
		return "unknown"
	}
	return output
}

// Generate runs all code generation
func Generate() error {
	fmt.Println("⚡ Generating code...")
	mg.Deps(generateSqlc, generateTempl)
	return nil
}

func generateSqlc() error {
	fmt.Println("  📊 Generating sqlc code...")
	return sh.RunV("sqlc", "generate")
}

func generateTempl() error {
	fmt.Println("  🎨 Generating templ code...")
	return sh.RunV("templ", "generate")
}

// Test runs all tests with coverage
func Test() error {
	fmt.Println("🧪 Running tests...")
	return sh.RunV("go", "test", "-race", "-coverprofile=coverage.out", "-covermode=atomic", "./...")
}

// TestVerbose runs tests with verbose output
func TestVerbose() error {
	fmt.Println("🧪 Running tests (verbose)...")
	return sh.RunV("go", "test", "-race", "-v", "-coverprofile=coverage.out", "-covermode=atomic", "./...")
}

// Coverage shows test coverage report
func Coverage() error {
	mg.Deps(Test)
	fmt.Println("📈 Generating coverage report...")
	return sh.RunV("go", "tool", "cover", "-html=coverage.out", "-o", "coverage.html")
}

// Fmt formats and tidies code
func Fmt() error {
	fmt.Println("✨ Formatting and tidying...")
	
	if err := sh.RunV("go", "mod", "tidy"); err != nil {
		return fmt.Errorf("failed to tidy modules: %w", err)
	}
	
	if err := sh.RunV("go", "fmt", "./..."); err != nil {
		return fmt.Errorf("failed to format code: %w", err)
	}
	
	// Format templ files if templ is available
	if err := sh.Run("which", "templ"); err == nil {
		fmt.Println("  🎨 Formatting templ files...")
		if err := sh.RunV("templ", "fmt", "."); err != nil {
			fmt.Printf("Warning: failed to format templ files: %v\n", err)
		}
	}
	
	return nil
}

// Vet analyzes code for common errors
func Vet() error {
	fmt.Println("🔍 Running go vet...")
	return sh.RunV("go", "vet", "./...")
}

// VulnCheck scans for known vulnerabilities
func VulnCheck() error {
	fmt.Println("🛡️  Running vulnerability check...")
	return sh.RunV("govulncheck", "./...")
}

// StaticCheck runs staticcheck linter
func StaticCheck() error {
	fmt.Println("🔬 Running staticcheck...")
	
	// Try to find staticcheck in GOPATH/bin
	gopath := os.Getenv("GOPATH")
	if gopath == "" {
		// Default GOPATH
		if home := os.Getenv("HOME"); home != "" {
			gopath = filepath.Join(home, "go")
		}
	}
	
	staticcheckPath := filepath.Join(gopath, "bin", "staticcheck")
	
	// Check if staticcheck exists, install if not
	if _, err := os.Stat(staticcheckPath); os.IsNotExist(err) {
		fmt.Println("Installing staticcheck...")
		if err := sh.RunV("go", "install", "honnef.co/go/tools/cmd/staticcheck@latest"); err != nil {
			return fmt.Errorf("failed to install staticcheck: %w", err)
		}
	}
	
	return sh.RunV(staticcheckPath, "./...")
}

// Run builds and runs the server
func Run() error {
	mg.SerialDeps(Build)
	fmt.Println("🚀 Starting server...")
	
	binaryPath := filepath.Join(buildDir, binaryName)
	if runtime.GOOS == "windows" {
		binaryPath += ".exe"
	}
	
	return sh.RunV(binaryPath)
}

// Dev starts development server with hot reload
func Dev() error {
	fmt.Println("🔥 Starting development server with hot reload...")
	
	// Ensure air is available
	if err := sh.Run("which", "air"); err != nil {
		fmt.Println("Installing air...")
		if err := sh.RunV("go", "install", "github.com/air-verse/air@latest"); err != nil {
			return fmt.Errorf("failed to install air: %w", err)
		}
	}
	
	return sh.RunV("air")
}

// Clean removes built binaries and generated files
func Clean() error {
	fmt.Println("🧹 Cleaning up...")
	
	// Remove build directory
	if err := sh.Rm(buildDir); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("failed to remove build directory: %w", err)
	}
	
	// Remove tmp directory
	if err := sh.Rm(tmpDir); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("failed to remove tmp directory: %w", err)
	}
	
	// Remove coverage files
	sh.Rm("coverage.out")
	sh.Rm("coverage.html")
	
	fmt.Println("✅ Clean complete!")
	return nil
}

// Setup installs required development tools
func Setup() error {
	fmt.Println("🚀 Setting up development environment...")
	
	tools := map[string]string{
		"templ":        "github.com/a-h/templ/cmd/templ@latest",
		"sqlc":         "github.com/sqlc-dev/sqlc/cmd/sqlc@latest",
		"govulncheck":  "golang.org/x/vuln/cmd/govulncheck@latest",
		"air":          "github.com/air-verse/air@latest",
		"staticcheck":  "honnef.co/go/tools/cmd/staticcheck@latest",
	}
	
	for tool, pkg := range tools {
		fmt.Printf("  📦 Installing %s...\n", tool)
		if err := sh.RunV("go", "install", pkg); err != nil {
			return fmt.Errorf("failed to install %s: %w", tool, err)
		}
	}
	
	// Download module dependencies
	fmt.Println("📥 Downloading dependencies...")
	if err := sh.RunV("go", "mod", "download"); err != nil {
		return fmt.Errorf("failed to download dependencies: %w", err)
	}
	
	fmt.Println("✅ Setup complete!")
	fmt.Println("💡 Next steps:")
	fmt.Println("   • Run 'mage dev' to start development with hot reload")
	fmt.Println("   • Run 'mage test' to run tests")
	fmt.Println("   • Run 'mage build' to create production binary")
	
	return nil
}

// Lint runs all linters and checks
func Lint() error {
	fmt.Println("🔍 Running all linters...")
	mg.Deps(Vet, StaticCheck, VulnCheck)
	return nil
}

// CI runs all checks suitable for continuous integration
func CI() error {
	fmt.Println("🏗️  Running CI pipeline...")
	mg.SerialDeps(Generate, Fmt, Lint, Test, Build)
	
	// Show build info
	if err := showBuildInfo(); err != nil {
		fmt.Printf("Warning: failed to show build info: %v\n", err)
	}
	
	fmt.Println("✅ CI pipeline completed successfully!")
	return nil
}

// Docker builds a Docker image (optional)
func Docker() error {
	fmt.Println("🐳 Building Docker image...")
	return sh.RunV("docker", "build", "-t", "go-web-server", ".")
}

// Help prints a help message with available commands
func Help() {
	fmt.Println(`
✨ Go Web Server Magefile ✨

Available commands:

Development:
  mage setup (s)        Install tools and dependencies
  mage generate (g)     Generate sqlc and templ code
  mage dev (d)          Start development server with hot reload
  mage run (r)          Build and run server

Build & Test:
  mage build (b)        Build optimized server binary (default)
  mage test (t)         Run all tests with coverage
  mage testverbose (tv) Run tests with verbose output
  mage coverage (co)    Generate HTML coverage report
  mage clean (c)        Clean build artifacts and coverage files

Quality & Production:
  mage fmt (f)          Format and tidy Go code
  mage vet (v)          Run go vet static analysis
  mage vulncheck (vc)   Check for security vulnerabilities
  mage staticcheck (sc) Run advanced static analysis
  mage lint (l)         Run all linters (vet + staticcheck + vulncheck)
  mage ci               Complete CI pipeline with build info
  mage docker           Build a Docker image (optional)

Other:
  mage help (h)         Show this help message
	`)
}

// showBuildInfo displays information about the built binary
func showBuildInfo() error {
	binaryPath := filepath.Join(buildDir, binaryName)
	if runtime.GOOS == "windows" {
		binaryPath += ".exe"
	}

	if _, err := os.Stat(binaryPath); os.IsNotExist(err) {
		return fmt.Errorf("binary not found: %s", binaryPath)
	}

	fmt.Println("\n📦 Build Information:")

	// Show binary size
	if info, err := os.Stat(binaryPath); err == nil {
		size := info.Size()
		fmt.Printf("   Binary size: %.2f MB\n", float64(size)/1024/1024)
	}

	// Show Go version
	if version, err := sh.Output("go", "version"); err == nil {
		fmt.Printf("   Go version: %s\n", version)
	}

	return nil
}

// Aliases for common commands
var Aliases = map[string]interface{}{
	"b":  Build,
	"g":  Generate,
	"t":  Test,
	"tv": TestVerbose,
	"co": Coverage,
	"f":  Fmt,
	"v":  Vet,
	"vc": VulnCheck,
	"sc": StaticCheck,
	"r":  Run,
	"d":  Dev,
	"c":  Clean,
	"s":  Setup,
	"l":  Lint,
	"h":  Help,
}
