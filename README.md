<p align="center">
  <img src="docs/images/go-logo.png" alt="Go Web Server Template Logo" width="400" />
</p>

<p align="center">
  <a href="https://github.com/dunamismax/go-web-server">
    <img src="https://readme-typing-svg.demolab.com/?font=Fira+Code&size=24&pause=1000&color=00ADD8&center=true&vCenter=true&width=900&lines=Modern+Go+Web+Server+Template;Radical+Simplicity+%26+Performance;Echo+%2B+Templ+%2B+HTMX+%2B+Pico.css;SQLC+%2B+SQLite+%2B+Structured+Logging;Single+Binary+Deployment;Production+Ready+Template;Zero+External+Dependencies" alt="Typing SVG" />
  </a>
</p>

<p align="center">
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

A minimal, perfect, reusable template for modern web and API development using the **Modern Go Stack** - a cohesive technology stack for building high-performance, maintainable applications with radical simplicity and stability. Creates single, self-contained binaries with no external dependencies.

**Key Features:**

- **Echo Framework**: High-performance, minimalist web framework
- **Templ**: Type-safe, component-based HTML templating in Go
- **HTMX**: Dynamic frontend interactions without JavaScript frameworks
- **Pico.css**: Minimal, semantic CSS with zero configuration
- **SQLC**: Type-safe database queries generated from SQL
- **SQLite (CGO-free)**: Pure Go database driver for single binaries
- **slog**: Structured logging throughout the application

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
| **Build**     | [go generate](https://golang.org/pkg/go/)                   | Go native code generation               |

## Quick Start

```bash
# Clone and setup
git clone https://github.com/dunamismax/go-web-server.git
cd go-web-server
go mod tidy

# Install dependencies and run
make deps
make run

# Server starts at http://localhost:8080
```

<p align="center">
  <img src="docs/images/gopher-mage.svg" alt="Gopher Mage" width="200" />
</p>

## Mage Commands

**Development:**

```bash
mage setup            # Install tools and dependencies
mage generate:all     # Generate sqlc and templ code
mage dev:server       # Start web development server
mage dev:tui          # Start TUI development tool
```

**Build & Test:**

```bash
mage build:all        # Build all applications
mage build:webapp     # Build web application only
mage build:tui        # Build TUI tool only
mage test:all         # Run all tests
mage test:webapp      # Run web application tests
mage test:tui         # Run TUI tool tests
```

**Database:**

```bash
mage database:up      # Run database migrations
mage database:down    # Rollback latest migration
mage database:reset   # Reset database (drop & recreate)
```

**Quality & Production:**

```bash
mage quality:all      # Run all quality checks
mage quality:fmt      # Format Go code
mage quality:vet      # Run go vet
mage quality:vulncheck # Check for vulnerabilities
mage clean            # Clean build artifacts
mage ci               # Full CI pipeline
```

## Applications

### Web Demo (`localhost:8080`)

Interactive web application showcasing Echo + Templ + HTMX + Pico.css + SQLC integration with user management, CRUD operations, and real-time updates.

## Development Commands

```bash
# Development
make run                # Build and run server
make dev                # Run with hot reload (requires air)
make test               # Run tests
make lint               # Run linting and security checks

# Build & Deploy
make build              # Build production binary
make clean              # Remove build artifacts
make deps               # Install development tools

# Code Generation
make generate           # Generate templ + sqlc code
make help               # Show all available commands
```

## Project Structure

```sh
go-web-server/
├── cmd/web/              # Application entry point
├── internal/
│   ├── handler/          # HTTP handlers & centralized routes
│   ├── view/             # Templ templates (.templ files)
│   ├── store/            # Database layer (SQL + generated code)
│   │   └── migrations/   # Database migrations
│   ├── config/           # Configuration management
│   └── ui/               # Static assets (embedded)
│       ├── static/       # Pico.css & HTMX
│       └── embed.go      # Go embed directive
├── bin/                  # Compiled binaries
├── Makefile             # Build automation
└── sqlc.yaml            # SQLC configuration
```

## Production Deployment

### Single Binary

```bash
make build  # Creates optimized binary in bin/server
```

The binary includes embedded Pico.css, HTMX, Templ templates, and SQLite database. **Zero external dependencies**, ~10-15MB size, instant startup.

The binary includes embedded static assets, templates, and SQLite database for easy deployment.

### Environment Variables

- `PORT`: Server port (default: 8080)
- `DATABASE_URL`: SQLite database file (default: data.db)
- `ENVIRONMENT`: Environment mode (default: development)
- `LOG_LEVEL`: Logging level (default: info)

### Systemd Service

```ini
[Unit]
Description=Go Web Server
After=network.target

[Service]
Type=simple
User=www-data
ExecStart=/usr/local/bin/server
Restart=always
RestartSec=5

[Install]
WantedBy=multi-user.target
```

### Reverse Proxy (Caddy)

```sh
yourdomain.com {
    reverse_proxy localhost:8080
}
```

## Key Features Demonstrated

**Web Application:**

- Echo framework with middleware
- Type-safe Templ templates
- HTMX dynamic interactions
- Pico.css semantic styling
- SQLC type-safe queries
- Embedded static assets
- Structured logging with slog

**Template Benefits:**

- Complete project structure
- Production-ready patterns
- Single binary deployment
- Hot reloading development
- Type safety throughout
- Modern Go practices

<p align="center">
  <a href="https://buymeacoffee.com/dunamismax" target="_blank">
    <img src="docs/images/buy-coffee-go.gif" alt="Buy Me A Coffee" style="height: 150px !important;" />
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
  <img src="docs/images/gopher-running-jumping.gif" alt="Gopher Running and Jumping" width="400" />
</p>
