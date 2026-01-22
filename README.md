<div align="center">

# 📊 SQL Helper

**A modern Windows desktop application for converting Excel data to SQL INSERT statements**

[![CI](https://github.com/hoangtran1411/sql-helper/actions/workflows/ci.yml/badge.svg)](https://github.com/hoangtran1411/sql-helper/actions/workflows/ci.yml)
[![Release](https://github.com/hoangtran1411/sql-helper/actions/workflows/release.yml/badge.svg)](https://github.com/hoangtran1411/sql-helper/actions/workflows/release.yml)
[![Go Version](https://img.shields.io/badge/Go-1.24+-00ADD8?style=flat&logo=go)](https://go.dev/)
[![Wails](https://img.shields.io/badge/Wails-v2-DF0000?style=flat&logo=wails)](https://wails.io/)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![GitHub release](https://img.shields.io/github/v/release/hoangtran1411/sql-helper?include_prereleases)](https://github.com/hoangtran1411/sql-helper/releases)

[Features](#-features) • [Installation](#-installation) • [Usage](#-usage) • [Development](#-development) • [Contributing](#-contributing)

</div>

---

## ✨ Features

- 🚀 **Fast Excel Parsing** — Efficiently read and parse `.xlsx` files using Excelize
- 📝 **SQL Generation** — Convert Excel rows to SQL INSERT statements instantly
- 🔢 **Smart Number Detection** — Automatically detect numeric columns to avoid unnecessary quoting
- 🔄 **Find & Replace** — Transform data values before generating SQL
- 📋 **One-Click Copy** — Copy generated SQL to clipboard with a single click
- 💾 **Export to File** — Save SQL output directly to `.txt` files
- 🔄 **Auto-Update** — Built-in update checker with automatic installation for Windows
- 🎨 **Modern UI** — Clean, responsive dark-themed interface built with HTML/CSS/JS

## 📸 Screenshots

<div align="center">
<img src="docs/screenshot.png" alt="SQL Helper Screenshot" width="800"/>
</div>

## 📥 Installation

### Download Pre-built Binaries

Download the latest release for Windows from the [Releases](https://github.com/hoangtran1411/sql-helper/releases) page:

| Platform | Download |
|----------|----------|
| **Windows** (x64) | [`sql-helper-windows-amd64.exe`](https://github.com/hoangtran1411/sql-helper/releases/latest) |

*Note: macOS and Linux builds are planned. For now, you can build from source.*

### Build from Source

See the [Development](#-development) section for build instructions.

## 🚀 Usage

### Quick Start

1. **Open Excel File** — Click "Browse" to select your `.xlsx` file.
2. **Select Sheet** — Choose the worksheet containing your data from the dropdown.
3. **Configure Options** — Set the target **Table Name** and specify **Number Columns** (comma-separated indices, e.g., `1,3`).
4. **Generate SQL** — Click "**Generate SQL**" to create INSERT statements.
5. **Copy or Save** — Use the copy button for clipboard or export to a file.

### Excel File Format

Your Excel file should follow a standard table structure:
- **First Row**: Column headers (used as SQL column names).
- **Subsequent Rows**: Data rows (transformed into individual VALUE tuples).

**Example Data:**

| ID | Name | Price | Status |
|----|------|-------|--------|
| 1 | Product A | 25.50 | Active |
| 2 | Product B | 40.00 | Inactive |

**Generated SQL Output:**

```sql
INSERT INTO products (ID, Name, Price, Status) VALUES
(1, 'Product A', 25.5, 'Active'),
(2, 'Product B', 40, 'Inactive');
```

---

## 🛠️ Development

### Prerequisites

- [Go 1.24+](https://go.dev/dl/)
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

# Run in development mode (with hot reload)
make dev

# Build production binary for Windows
make build-windows
```

### Available Commands

| Command | Description |
|---------|-------------|
| `make dev` | Run in development mode with hot reload |
| `make build` | Build production binary for current platform |
| `make build-windows` | Cross-compile for Windows AMD64 |
| `make test` | Run all Go unit tests |
| `make lint` | Run golangci-lint |
| `make coverage` | Generate and view test coverage |
| `make clean` | Remove build artifacts and coverage files |

### Project Structure

```
sql-helper/
├── main.go              # Wails entry point
├── app.go               # App struct with Wails bindings
├── updater.go           # Auto-update functionality (Windows)
├── internal/
│   ├── excel/           # Excel parsing logic (Excelize)
│   └── sql/             # SQL generation & formatting logic
├── frontend/            # Vue/React/Svelte/Vanilla frontend
├── build/               # Build output & icons
└── .github/workflows/   # CI/CD (Lint, Test, Release)
```

## 🤝 Contributing

Contributions are welcome! Please check our [CONTRIBUTING.md](CONTRIBUTING.md) for guidelines.

1. Fork the Project
2. Create your Feature Branch (`git checkout -b feature/AmazingFeature`)
3. Commit your Changes (`git commit -m 'Add some AmazingFeature'`)
4. Push to the Branch (`git push origin feature/AmazingFeature`)
5. Open a Pull Request

## 📋 Roadmap

- [ ] Multiple database dialect support (PostgreSQL, SQLite, Oracle)
- [ ] Native macOS and Linux releases
- [ ] Dark/Light mode toggle (currently Dark only)
- [ ] Batch processing for multiple Excel files at once
- [ ] Custom SQL templates (INSERT IGNORE, REPLACE INTO, etc.)

## 📜 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🙏 Acknowledgments

- [Wails](https://wails.io/) — The amazing Go/Web hybrid framework.
- [Excelize](https://github.com/xuri/excelize) — Powering the Excel parsing.
- [Go-Ole](https://github.com/go-ole/go-ole) — Essential for Windows integration.

---

<div align="center">

**Created by [Hoang Tran](https://github.com/hoangtran1411)**

⭐ Star this repo if you find it useful!

</div>
