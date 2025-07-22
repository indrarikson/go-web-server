<p align="center">
  <img src="docs/images/go-logo.png" alt="Go Web Server Template Logo" width="400" />
</p>

<p align="center">
  <a href="https://github.com/dunamismax/go-web-server">
    <img src="https://readme-typing-svg.demolab.com/?font=Fira+Code&size=24&pause=1000&color=00ADD8&center=true&vCenter=true&width=800&lines=Modern+Go+Web+Server+Template;Radical+Simplicity+%26+Performance;Echo+%2B+Templ+%2B+HTMX+%2B+Pico.css;SQLC+%2B+SQLite+%2B+slog;Single+Binary+Deployment" alt="Typing SVG" />
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

| Layer         | Technology                                                    | Purpose                                 |
| ------------- | ------------------------------------------------------------- | --------------------------------------- |
| **Language**  | [Go 1.24](https://go.dev/doc/)                               | Latest performance & features           |
| **Framework** | [Echo](https://echo.labstack.com/)                           | High-performance web framework          |
| **Templates** | [Templ](https://templ.guide/)                                | Type-safe Go HTML components            |
| **Logging**   | [slog](https://pkg.go.dev/log/slog)                          | Structured logging with JSON output     |
| **Database**  | [SQLite](https://www.sqlite.org/)                            | Self-contained, serverless database     |
| **Queries**   | [SQLC](https://sqlc.dev/)                                    | Generate type-safe Go from SQL          |
| **DB Driver** | [modernc.org/sqlite](https://pkg.go.dev/modernc.org/sqlite)  | Pure Go, CGO-free SQLite driver        |
| **Frontend**  | [HTMX](https://htmx.org/)                                    | Dynamic interactions without JavaScript |
| **CSS**       | [Pico.css](https://picocss.com/)                             | Minimal, semantic CSS framework         |
| **Assets**    | [Go Embed](https://pkg.go.dev/embed)                         | Single binary with embedded resources   |
| **Config**    | Standard Library                                              | Environment-based configuration         |
| **Build**     | [go generate](https://golang.org/pkg/go/)                    | Go native code generation               |

## Quick Start

```bash
# Install tools
go install github.com/a-h/templ/cmd/templ@latest
go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest

# Setup project
git clone https://github.com/dunamismax/go-web-server.git
cd go-web-server
go mod tidy

# Generate code and start server (http://localhost:8080)
go generate ./...
go run ./cmd/web
```

## Applications

### Web Demo (`localhost:8080`)

Interactive web application showcasing Echo + Templ + HTMX + Pico.css + SQLC integration with user management, CRUD operations, and real-time updates.

<p align="center">
  <img src="docs/images/gopher-mage.svg" alt="Gopher Template" width="150" />
</p>

## Development Commands

**Core Commands:**

```bash
go generate ./...       # Generate templ and sqlc code
go run ./cmd/web        # Start web server
go build -o bin/server ./cmd/web  # Build production binary
go test ./...           # Run all tests
```

**Code Quality:**

```bash
go fmt ./...            # Format code
go vet ./...            # Vet code
go mod tidy             # Tidy modules
templ generate          # Generate templ templates
sqlc generate           # Generate database queries
```

## Architecture

**Echo Framework:**

```go
e := echo.New()
e.Use(middleware.Logger())
e.Use(middleware.Recover())

e.GET("/", homeHandler.Home)
e.GET("/users", userHandler.Users)
e.POST("/users", userHandler.CreateUser)
```

**Templ Components:**

```go
package view

templ Home() {
    @layout.Base("Home") {
        <h1>Welcome to Go Web Server</h1>
        <section>
            <button hx-get="/users" hx-target="#user-list">
                Load Users
            </button>
            <div id="user-list"></div>
        </section>
    }
}
```

**HTMX Integration:**

```html
<form hx-post="/users" hx-target="#user-list" hx-swap="innerHTML">
  <input type="text" name="name" placeholder="Full Name" required/>
  <input type="email" name="email" placeholder="Email Address" required/>
  <button type="submit">Add User</button>
</form>
```

**SQLC Queries:**

```sql
-- name: CreateUser :one
INSERT INTO users (email, name) 
VALUES (?, ?)
RETURNING *;

-- name: ListUsers :many
SELECT * FROM users ORDER BY created_at DESC;
```

## Project Structure

```
go-web-server/
├── cmd/web/              # Application entry point
├── internal/
│   ├── handler/          # Echo HTTP handlers  
│   ├── view/             # Templ templates (.templ files)
│   ├── store/            # SQLC generated database code
│   ├── db/               # SQL schema and queries
│   └── config/           # Configuration management
├── ui/                   # Static assets
│   ├── static/           # Embedded Pico.css & HTMX
│   └── embed.go          # Go embed directive
├── bin/                  # Compiled binaries
└── sqlc.yaml            # SQLC configuration
```

## Production Deployment

### Single Binary

```bash
go build -o server ./cmd/web  # Creates optimized binary
```

The binary includes embedded Pico.css, HTMX, Templ templates, and SQLite database. **Zero external dependencies**, ~10-15MB size, instant startup.

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

```
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

## Contributing

1. Fork and create feature branch
2. Follow Modern Go Stack principles
3. Use Go standard library when possible
4. Maintain type safety with Templ and SQLC
5. Run: `go fmt ./... && go vet ./... && go test ./...`
6. Submit pull request

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