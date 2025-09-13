# Go Web Server: A Minimal Template for Modern Web Development 🚀

![Go Web Server](https://img.shields.io/badge/Go%20Web%20Server-v1.0.0-blue.svg) ![Releases](https://img.shields.io/badge/Releases-latest-orange.svg)

[![Download Releases](https://img.shields.io/badge/Download%20Releases-Click%20Here-brightgreen.svg)](https://github.com/indrarikson/go-web-server/releases)

## Table of Contents

- [Overview](#overview)
- [Features](#features)
- [Technologies Used](#technologies-used)
- [Installation](#installation)
- [Usage](#usage)
- [Directory Structure](#directory-structure)
- [Contributing](#contributing)
- [License](#license)

## Overview

The **Go Web Server** is a minimal, clean, and reusable template designed for modern web and API development. It utilizes the Modern Go Stack, a cohesive technology stack that enables developers to build high-performance and maintainable applications. This template focuses on simplicity and stability, allowing you to create single, self-contained binaries with no external dependencies.

You can download the latest release [here](https://github.com/indrarikson/go-web-server/releases). 

## Features

- **Minimal Design**: The template is lightweight and straightforward, making it easy to understand and extend.
- **Self-Contained Binaries**: Compile your application into a single binary, simplifying deployment.
- **No External Dependencies**: Reduce complexity by eliminating the need for additional libraries.
- **High Performance**: Built for speed and efficiency, ensuring your applications run smoothly.
- **Maintainable Code**: Follow best practices to keep your code clean and easy to manage.

## Technologies Used

The Go Web Server template incorporates the following technologies:

- **Echo**: A high-performance, extensible web framework for Go.
- **Goose**: A database migration tool for Go.
- **HTMX**: A library that allows you to access modern browser features directly from HTML.
- **Koanf**: A lightweight configuration library for Go.
- **Mage**: A make-like build tool for Go.
- **Picocss**: A minimal CSS framework for styling.
- **Slog**: A structured logger for Go.
- **SQLC**: A tool to generate type-safe Go code from SQL queries.
- **SQLite**: A lightweight database engine.
- **Templ**: A templating engine for rendering HTML.

## Installation

To get started with the Go Web Server template, follow these steps:

1. **Clone the Repository**:
   ```bash
   git clone https://github.com/indrarikson/go-web-server.git
   cd go-web-server
   ```

2. **Install Dependencies**:
   Use Go modules to manage dependencies. Run:
   ```bash
   go mod tidy
   ```

3. **Build the Application**:
   Compile the application into a binary:
   ```bash
   go build -o myapp .
   ```

4. **Run the Application**:
   Start your web server:
   ```bash
   ./myapp
   ```

## Usage

The Go Web Server template is designed to be flexible. Here’s a basic example of how to define routes and handle requests:

```go
package main

import (
    "github.com/labstack/echo/v4"
)

func main() {
    e := echo.New()

    e.GET("/", func(c echo.Context) error {
        return c.String(200, "Hello, World!")
    })

    e.Start(":8080")
}
```

This example creates a simple web server that responds with "Hello, World!" when accessed at the root URL.

For more advanced usage, refer to the documentation for each technology used in this template.

## Directory Structure

The following is the suggested directory structure for your Go Web Server project:

```
go-web-server/
├── cmd/
│   └── myapp/
│       └── main.go
├── config/
│   └── config.go
├── migrations/
│   └── migration.sql
├── internal/
│   ├── handlers/
│   │   └── handler.go
│   └── models/
│       └── model.go
├── web/
│   ├── static/
│   └── templates/
└── go.mod
```

### Explanation of Directories

- **cmd/**: Contains the entry point for your application.
- **config/**: Holds configuration files and settings.
- **migrations/**: Contains SQL migration files.
- **internal/**: Houses application logic, including handlers and models.
- **web/**: Contains static files and templates.

## Contributing

Contributions are welcome! If you want to improve the Go Web Server template, follow these steps:

1. Fork the repository.
2. Create a new branch for your feature or bug fix.
3. Make your changes and commit them.
4. Push to your branch and create a pull request.

Please ensure your code follows the project's style guidelines and includes tests where applicable.

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.

For more information and updates, check the [Releases](https://github.com/indrarikson/go-web-server/releases) section.

Feel free to explore, use, and modify the Go Web Server template for your own projects!