<div align="center">

# 📊 SQL Helper

**A modern desktop application for converting Excel data to SQL INSERT statements**

[![CI](https://github.com/hoangtran1411/sql-helper/actions/workflows/ci.yml/badge.svg)](https://github.com/hoangtran1411/sql-helper/actions/workflows/ci.yml)
[![Release](https://github.com/hoangtran1411/sql-helper/actions/workflows/release.yml/badge.svg)](https://github.com/hoangtran1411/sql-helper/actions/workflows/release.yml)
[![Go Version](https://img.shields.io/badge/Go-1.23+-00ADD8?style=flat&logo=go)](https://go.dev/)
[![Wails](https://img.shields.io/badge/Wails-v2-DF0000?style=flat&logo=wails)](https://wails.io/)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![GitHub release](https://img.shields.io/github/v/release/hoangtran1411/sql-helper?include_prereleases)](https://github.com/hoangtran1411/sql-helper/releases)
[![Downloads](https://img.shields.io/github/downloads/hoangtran1411/sql-helper/total)](https://github.com/hoangtran1411/sql-helper/releases)

[Features](#-features) • [Installation](#-installation) • [Usage](#-usage) • [Development](#-development) • [Contributing](#-contributing)

</div>

---

## ✨ Features

- 🚀 **Fast Excel Parsing** — Efficiently read and parse `.xlsx` files using Excelize
- 📝 **SQL Generation** — Convert Excel rows to SQL INSERT statements instantly
- 🔢 **Smart Number Detection** — Automatically detect numeric columns to avoid quoting
- 🔄 **Find & Replace** — Transform data values before generating SQL
- 📋 **One-Click Copy** — Copy generated SQL to clipboard with a single click
- 💾 **Export to File** — Save SQL output directly to `.sql` files
- 🔄 **Auto-Update** — Built-in update checker with one-click installation
- 🎨 **Modern UI** — Clean, responsive interface built with HTML/CSS/JS
- 🖥️ **Cross-Platform** — Available for Windows, macOS, and Linux

## 📸 Screenshots

<div align="center">
<img src="docs/screenshot.png" alt="SQL Helper Screenshot" width="800"/>
</div>

## 📥 Installation

### Download Pre-built Binaries

Download the latest release for your platform from the [Releases](https://github.com/hoangtran1411/sql-helper/releases) page:

| Platform | Download |
|----------|----------|
| Windows | `sql-helper-windows-amd64.zip` |
| macOS | `sql-helper-darwin-universal.tar.gz` |
| Linux | `sql-helper-linux-amd64.tar.gz` |

### Build from Source

See the [Development](#-development) section for build instructions.

## 🚀 Usage

### Quick Start

1. **Open Excel File** — Click "Browse" to select your `.xlsx` file
2. **Select Sheet** — Choose the worksheet containing your data
3. **Configure Options** — Set table name and number columns
4. **Generate SQL** — Click "Generate" to create INSERT statements
5. **Copy or Save** — Copy to clipboard or export to file

### Excel File Format

Your Excel file should have:
- **Row 1**: Column headers (these become SQL column names)
- **Row 2+**: Data rows (each row becomes an INSERT statement)

**Example:**

| ID | Name | Price | Quantity |
|----|------|-------|----------|
| 1 | Widget A | 29.99 | 100 |
| 2 | Widget B | 49.99 | 50 |

### Generated SQL

```sql
INSERT INTO products (ID, Name, Price, Quantity) VALUES
(1, 'Widget A', 29.99, 100),
(2, 'Widget B', 49.99, 50);
```

### Number Columns

Specify which columns contain numeric values (e.g., `1,3,4` for columns 1, 3, and 4). These values won't be wrapped in quotes.

### Find & Replace

Transform data before SQL generation:
- **Find**: Text pattern to search for
- **Replace**: Replacement text

Useful for escaping special characters or standardizing data formats.

## 🛠️ Development

### Prerequisites

- [Go 1.23+](https://go.dev/dl/)
- [Node.js 20+](https://nodejs.org/)
- [Wails CLI](https://wails.io/docs/gettingstarted/installation)

```bash
# Install Wails CLI
go install github.com/wailsapp/wails/v2/cmd/wails@latest
```

### Clone & Build

```bash
# Clone repository
git clone https://github.com/hoangtran1411/sql-helper.git
cd sql-helper

# Install dependencies
go mod download

# Run in development mode
make dev
# or
wails dev

# Build production binary
make build
# or
wails build
```

### Available Commands

| Command | Description |
|---------|-------------|
| `make dev` | Run in development mode with hot reload |
| `make build` | Build production binary |
| `make test` | Run all tests |
| `make lint` | Run golangci-lint |
| `make coverage` | Generate test coverage report |
| `make fmt` | Format code with gofmt |
| `make clean` | Remove build artifacts |

### Project Structure

```
sql-helper/
├── main.go              # Wails entry point
├── app.go               # App struct with Wails bindings
├── updater.go           # Auto-update functionality
├── internal/
│   ├── excel/           # Excel parsing logic
│   │   ├── parser.go
│   │   └── parser_test.go
│   └── sql/             # SQL generation logic
│       ├── generator.go
│       ├── generator_test.go
│       ├── formatter.go
│       └── formatter_test.go
├── frontend/            # Web frontend (HTML/CSS/JS)
├── build/               # Build output
└── .github/workflows/   # CI/CD pipelines
```

### Running Tests

```bash
# Run all tests
make test

# Run tests with coverage
make coverage

# View coverage report in browser
go tool cover -html=coverage.out
```

### Linting

```bash
# Install golangci-lint
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

# Run linter
make lint
```

## 🤝 Contributing

Contributions are welcome! Please read our [Contributing Guide](CONTRIBUTING.md) for details on:

- Code of Conduct
- Development workflow
- Submitting pull requests
- Coding standards

## 📋 Roadmap

- [ ] Support for `.xls` (legacy Excel) files
- [ ] Multiple database dialect support (MySQL, PostgreSQL, SQLite)
- [ ] Batch processing for multiple sheets
- [ ] Custom SQL templates
- [ ] Dark mode theme
- [ ] Localization support

## 📜 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🙏 Acknowledgments

- [Wails](https://wails.io/) — Build desktop apps with Go and web technologies
- [Excelize](https://github.com/xuri/excelize) — Go library for reading/writing Excel files
- [golangci-lint](https://golangci-lint.run/) — Fast Go linters runner

---

<div align="center">

**Made with ❤️ by [Hoang Tran](https://github.com/hoangtran1411)**

⭐ Star this repo if you find it useful!

</div>
