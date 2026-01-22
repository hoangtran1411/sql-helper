# Contributing to SQL Helper

First off, thank you for considering contributing to SQL Helper! It's people like you that make this project better for everyone.

## 📋 Table of Contents

- [Code of Conduct](#code-of-conduct)
- [Getting Started](#getting-started)
- [Development Setup](#development-setup)
- [How to Contribute](#how-to-contribute)
- [Pull Request Process](#pull-request-process)
- [Coding Standards](#coding-standards)
- [Commit Messages](#commit-messages)
- [Issue Guidelines](#issue-guidelines)

## 📜 Code of Conduct

This project and everyone participating in it is governed by our commitment to providing a welcoming and inclusive environment. By participating, you are expected to:

- Use welcoming and inclusive language
- Be respectful of differing viewpoints and experiences
- Gracefully accept constructive criticism
- Focus on what is best for the community
- Show empathy towards other community members

## 🚀 Getting Started

### Prerequisites

Before you begin, ensure you have the following installed:

- [Go 1.23+](https://go.dev/dl/)
- [Node.js 20+](https://nodejs.org/)
- [Git](https://git-scm.com/)
- [Wails CLI](https://wails.io/docs/gettingstarted/installation)

```bash
go install github.com/wailsapp/wails/v2/cmd/wails@latest
```

### Development Setup

1. **Fork the repository** on GitHub

2. **Clone your fork**
   ```bash
   git clone https://github.com/YOUR_USERNAME/sql-helper.git
   cd sql-helper
   ```

3. **Add upstream remote**
   ```bash
   git remote add upstream https://github.com/hoangtran1411/sql-helper.git
   ```

4. **Install dependencies**
   ```bash
   go mod download
   ```

5. **Install development tools**
   ```bash
   make install-tools
   ```

6. **Run in development mode**
   ```bash
   make dev
   ```

## 🔧 How to Contribute

### Reporting Bugs

Before creating a bug report, please check existing issues to avoid duplicates.

When filing an issue, include:

- **Clear title** describing the issue
- **Steps to reproduce** the behavior
- **Expected behavior** vs actual behavior
- **Screenshots** if applicable
- **Environment details** (OS, Go version, etc.)

### Suggesting Features

Feature requests are welcome! Please provide:

- Clear use case description
- Expected behavior
- Any alternatives you've considered
- Additional context or mockups

### Contributing Code

1. Check existing issues or create a new one
2. Comment on the issue to let others know you're working on it
3. Fork and create a feature branch
4. Make your changes
5. Submit a pull request

## 🔀 Pull Request Process

### Before Submitting

1. **Sync with upstream**
   ```bash
   git fetch upstream
   git rebase upstream/main
   ```

2. **Create a feature branch**
   ```bash
   git checkout -b feature/your-feature-name
   ```

3. **Make your changes** following our coding standards

4. **Run tests**
   ```bash
   make test
   ```

5. **Run linter**
   ```bash
   make lint
   ```

6. **Format code**
   ```bash
   make fmt
   ```

### Submitting

1. Push your branch to your fork
   ```bash
   git push origin feature/your-feature-name
   ```

2. Open a Pull Request against `main` branch

3. Fill in the PR template with:
   - Description of changes
   - Related issue(s)
   - Screenshots (if UI changes)
   - Testing done

### Review Process

- All PRs require at least one approving review
- CI checks must pass (lint, test, build)
- Maintain or improve test coverage
- Address all review comments

## 📏 Coding Standards

### Go Code

- Follow [Effective Go](https://go.dev/doc/effective_go) guidelines
- Use `gofmt` for formatting
- Run `golangci-lint` before committing
- Write tests for new functionality
- Maintain 70%+ code coverage

### File Organization

```
internal/           # Private packages
├── excel/          # Excel parsing logic
└── sql/            # SQL generation logic
```

### Naming Conventions

| Type | Convention | Example |
|------|------------|---------|
| Files | `snake_case.go` | `parser_test.go` |
| Packages | `lowercase` | `excel`, `sql` |
| Functions | `PascalCase` (exported) | `ParseExcelFile` |
| Functions | `camelCase` (unexported) | `parseRow` |
| Variables | `camelCase` | `rowCount` |
| Constants | `PascalCase` | `MaxRowLimit` |

### Documentation

- All exported functions must have doc comments
- Use complete sentences in comments
- Document edge cases and design decisions

```go
// ParseExcelFile reads an Excel file and returns a list of sheet names.
// Returns an error if the file cannot be opened or is not a valid Excel file.
func ParseExcelFile(filePath string) ([]string, error) {
    // ...
}
```

### Error Handling

- Always wrap errors with context
- Use `fmt.Errorf` with `%w` for wrapping
- Handle all errors explicitly

```go
if err != nil {
    return fmt.Errorf("failed to parse row %d: %w", rowNum, err)
}
```

## 💬 Commit Messages

We follow [Conventional Commits](https://www.conventionalcommits.org/):

```
<type>(<scope>): <subject>

[optional body]

[optional footer]
```

### Types

| Type | Description |
|------|-------------|
| `feat` | New feature |
| `fix` | Bug fix |
| `docs` | Documentation only |
| `style` | Formatting, no code change |
| `refactor` | Code change without fix/feature |
| `perf` | Performance improvement |
| `test` | Adding tests |
| `chore` | Maintenance tasks |
| `ci` | CI/CD changes |

### Examples

```
feat(sql): add PostgreSQL dialect support

fix(excel): handle empty cells correctly
Closes #42

docs: update README with new usage examples

chore(deps): update Excelize to v2.9.0
```

## 🐛 Issue Guidelines

### Bug Reports

Use the bug report template and include:

```markdown
**Describe the bug**
A clear description of what the bug is.

**To Reproduce**
Steps to reproduce:
1. Go to '...'
2. Click on '...'
3. See error

**Expected behavior**
What you expected to happen.

**Screenshots**
If applicable, add screenshots.

**Environment:**
 - OS: [e.g., Windows 11]
 - Version: [e.g., v1.2.0]
 - Go Version: [e.g., 1.23]
```

### Feature Requests

```markdown
**Is your feature request related to a problem?**
A clear description of the problem.

**Describe the solution you'd like**
What you want to happen.

**Describe alternatives you've considered**
Any alternative solutions or features.

**Additional context**
Any other context or screenshots.
```

---

## 🙏 Thank You!

Your contributions make this project better. Whether it's:

- 🐛 Reporting bugs
- 💡 Suggesting features
- 📝 Improving documentation
- 🔧 Submitting code

Every contribution is valued and appreciated!

---

<div align="center">

**Questions?** Feel free to open an issue or reach out!

</div>
