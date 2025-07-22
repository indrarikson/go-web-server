<p align="center">
  <img src="https://github.com/dunamismax/images/blob/main/go-logo.png" alt="Go Web Server Template Logo" width="400" />
</p>

<p align="center">
  <a href="https://github.com/dunamismax/go-web-server">
    <img src="https://readme-typing-svg.demolab.com/?font=Fira+Code&size=24&pause=1000&color=00ADD8&center=true&vCenter=true&width=900&lines=Modern+Go+Web+Server+Template;Radical+Simplicity+%26+Performance;Echo+%2B+Templ+%2B+HTMX+%2B+Pico.css;SQLC+%2B+SQLite+%2B+Structured+Logging;Single+Binary+Deployment;Production+Ready+Template;Zero+External+Dependencies" alt="Typing SVG" />
  </a>
</p>

<p align="center">
  <a href="https://github.com/dunamismax/go-web-server/actions/workflows/ci.yml"><img src="https://github.com/dunamismax/go-web-server/workflows/CI%20Pipeline/badge.svg" alt="CI Status"></a>
  <a href="https://golang.org/"><img src="https://img.shields.io/badge/Go-1.24+-00ADD8.svg?logo=go" alt="Go Version"></a>
  <a href="https://echo.labstack.com/"><img src="https://img.shields.io/badge/Framework-Echo-00ADD8.svg?logo=go" alt="Echo Framework"></a>
  <a href="https://templ.guide/"><img src="https://img.shields.io/badge/Templates-Templ-00ADD8.svg?logo=go" alt="Templ"></a>
  <a href="https://htmx.org/"><img src="https://img.shields.io/badge/Frontend-HTMX-3D72D7.svg?logo=htmx" alt="HTMX"></a>
  <a href="https://picocss.com/"><img src="https://img.shields.io/badge/CSS-Pico.css-13795B.svg" alt="Pico.css"></a>
  <a href="https://sqlc.dev/"><img src="https://img.shields.io/badge/Queries-SQLC-00ADD8.svg?logo=go" alt="SQLC"></a>
  <a href="https://www.sqlite.org/"><img src="https://img.shields.io/badge/Database-SQLite-003B57.svg?logo=sqlite" alt="SQLite"></a>
  <a href="https://pkg.go.dev/log/slog"><img src="https://img.shields.io/badge/Logging-slog-00ADD8.svg?logo=go" alt="Go slog"></a>
  <a href="https://golang.org/pkg/go/"><img src="https://img.shields.io/badge/Build-go%20generate-purple.svg?logo=go" alt="Go Generate"></a>
  <a href="https://opensource.org/licenses/MIT"><img src="https://img.shields.io/badge/License-MIT-green.svg" alt="MIT License"></a>
</p>

---

## About

A minimal, clean, reusable template for modern web and API development using the **Modern Go Stack** - a cohesive technology stack for building high-performance, maintainable applications with radical simplicity and stability. Creates single, self-contained binaries with no external dependencies.

**Key Features:**

- **Echo Framework**: High-performance, minimalist web framework
- **Templ**: Type-safe, component-based HTML templating in Go
- **HTMX**: Dynamic frontend interactions without JavaScript frameworks
- **Pico.css**: Minimal, semantic CSS with zero configuration
- **SQLC**: Type-safe database queries generated from SQL
- **SQLite (CGO-free)**: Pure Go database driver for single binaries
- **slog**: Structured logging throughout the application
- **golang-migrate**: Professional database migration system
- **Comprehensive testing**: Example tests with coverage reporting
- **Advanced middleware**: Security, logging, rate limiting, and error handling

## Tech Stack

| Layer         | Technology                                                  | Purpose                                 |
| ------------- | ----------------------------------------------------------- | --------------------------------------- |
| **Language**  | [Go 1.24](https://go.dev/doc/)                              | Latest performance & features           |
| **Framework** | [Echo](https://echo.labstack.com/)                          | High-performance web framework          |
| **Templates** | [Templ](https://templ.guide/)                               | Type-safe Go HTML components            |
| **Logging**   | [slog](https://pkg.go.dev/log/slog)                         | Structured logging with JSON output     |
| **Database**  | [SQLite](https://www.sqlite.org/)                           | Self-contained, serverless database     |
| **Queries**   | [SQLC](https://sqlc.dev/)                                   | Generate type-safe Go from SQL          |
| **DB Driver** | [modernc.org/sqlite](https://pkg.go.dev/modernc.org/sqlite) | Pure Go, CGO-free SQLite driver         |
| **Frontend**  | [HTMX](https://htmx.org/)                                   | Dynamic interactions without JavaScript |
| **CSS**       | [Pico.css](https://picocss.com/)                            | Minimal, semantic CSS framework         |
| **Assets**    | [Go Embed](https://pkg.go.dev/embed)                        | Single binary with embedded resources   |
| **Config**    | Standard Library                                            | Environment-based configuration         |
| **Migrations** | [golang-migrate](https://github.com/golang-migrate/migrate) | Database migration management           |
| **Build**     | [Mage](https://magefile.org/)                               | Go-based build automation               |
| **Testing**   | [testify](https://github.com/stretchr/testify)              | Enhanced testing with assertions        |

<p align="center">
  <img src="https://github.com/dunamismax/images/blob/main/gopher-mage.svg" alt="Gopher Mage" width="200" />
</p>

## Quick Start

```bash
# Clone and setup
git clone https://github.com/dunamismax/go-web-server.git
cd go-web-server
go mod tidy

# Install dependencies and run
mage setup
mage run

# Server starts at http://localhost:8080
```

## Mage Commands

**Development:**

```bash
mage setup            # Install tools and dependencies (templ, sqlc, air, etc.)
mage generate         # Generate sqlc and templ code
mage dev              # Start development server with hot reload
mage run              # Build and run server
```

**Build & Test:**

```bash
mage build            # Build optimized server binary (default target)
mage test             # Run all tests with coverage
mage testverbose      # Run tests with verbose output
mage coverage         # Generate HTML coverage report
mage clean            # Clean build artifacts and coverage files
```

**Quality & Production:**

```bash
mage fmt              # Format and tidy Go code (includes templ formatting)
mage vet              # Run go vet static analysis
mage vulncheck        # Check for security vulnerabilities
mage staticcheck      # Run advanced static analysis
mage lint             # Run all linters (vet + staticcheck + vulncheck)
mage ci               # Complete CI pipeline with build info
```

## Applications

### Web Demo (`localhost:8080`)

Interactive web application showcasing Echo + Templ + HTMX + Pico.css + SQLC integration with user management, CRUD operations, and real-time updates.

## Project Structure

```sh
go-web-server/
├── .air.toml             # Hot reload configuration
├── .github/workflows/    # CI/CD pipeline
├── cmd/web/              # Application entry point
├── internal/
│   ├── handler/          # HTTP handlers & centralized routes + tests
│   ├── view/             # Templ templates (.templ files)
│   ├── store/            # Database layer (SQL + generated code) + tests
│   │   └── migrations/   # Database migrations (golang-migrate)
│   ├── config/           # Configuration management
│   └── ui/               # Static assets (embedded)
│       ├── static/       # Pico.css & HTMX
│       └── embed.go      # Go embed directive
├── bin/                  # Compiled binaries
├── magefile.go          # Mage build automation
├── sqlc.yaml            # SQLC configuration
└── coverage.out/html     # Test coverage reports
```

## Production Deployment

### Single Binary

```bash
mage build  # Creates optimized binary in bin/server (~10MB)
```

The binary includes embedded Pico.css, HTMX, Templ templates, and SQLite database. **Zero external dependencies**, single file deployment with instant startup.

### Environment Variables

- `PORT`: Server port (default: 8080)
- `HOST`: Server host (default: "")
- `DATABASE_URL`: SQLite database file (default: data.db)
- `ENVIRONMENT`: Environment mode (default: development)
- `LOG_LEVEL`: Logging level - debug, info, warn, error (default: info)
- `LOG_FORMAT`: Log format - text or json (default: text)
- `DEBUG`: Enable debug mode (default: false)
- `RUN_MIGRATIONS`: Auto-run database migrations (default: true)
- `ENABLE_CORS`: Enable CORS middleware (default: true)

## Key Features Demonstrated

**Web Application:**

- Echo framework with comprehensive middleware stack
- Type-safe Templ templates with component architecture
- HTMX dynamic interactions with theme switching
- Pico.css semantic styling with dark/light mode
- SQLC type-safe queries with context support
- golang-migrate database migrations
- Embedded static assets (Pico.css, HTMX)
- Structured logging with slog throughout
- Comprehensive health check endpoint

**Development Experience:**

- Hot reloading with Air configuration
- Complete test suite with coverage reporting
- Advanced static analysis (staticcheck, go vet, govulncheck)
- Automated formatting for Go and Templ files
- Single-command CI pipeline (mage ci)
- Cross-platform build support

**Production Ready:**

- Security middleware (CORS, rate limiting, secure headers)
- Graceful shutdown handling
- Request ID tracing and structured logging
- Environment-based configuration
- Single binary deployment (~10MB)
- Zero external dependencies

<p align="center">
  <a href="https://buymeacoffee.com/dunamismax" target="_blank">
    <img src="https://github.com/dunamismax/images/blob/main/buy-coffee-go.gif" alt="Buy Me A Coffee" style="height: 150px !important;" />
  </a>
</p>

<p align="center">
  <a href="https://twitter.com/dunamismax" target="_blank"><img src="https://img.shields.io/badge/Twitter-%231DA1F2.svg?&style=for-the-badge&logo=twitter&logoColor=white" alt="Twitter"></a>
  <a href="https://bsky.app/profile/dunamismax.bsky.social" target="_blank"><img src="https://img.shields.io/badge/Bluesky-blue?style=for-the-badge&logo=bluesky&logoColor=white" alt="Bluesky"></a>
  <a href="https://reddit.com/user/dunamismax" target="_blank"><img src="https://img.shields.io/badge/Reddit-%23FF4500.svg?&style=for-the-badge&logo=reddit&logoColor=white" alt="Reddit"></a>
  <a href="https://discord.com/users/dunamismax" target="_blank"><img src="https://img.shields.io/badge/Discord-dunamismax-7289DA.svg?style=for-the-badge&logo=discord&logoColor=white" alt="Discord"></a>
  <a href="https://signal.me/#p/+dunamismax.66" target="_blank"><img src="https://img.shields.io/badge/Signal-dunamismax.66-3A76F0.svg?style=for-the-badge&logo=signal&logoColor=white" alt="Signal"></a>
</p>

## License

This project is licensed under the **MIT License** - see the [LICENSE](LICENSE) file for details.

---

<p align="center">
  <strong>The Modern Go Stack</strong><br>
  <sub>Echo • Templ • HTMX • Pico.css • SQLC • SQLite • slog • Single Binary</sub>
</p>

<p align="center">
  <img src="https://github.com/dunamismax/images/blob/main/gopher-running-jumping.gif" alt="Gopher Running and Jumping" width="400" />
</p>
