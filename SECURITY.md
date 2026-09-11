# Security Policy

The SQL Helper team takes the security of our application and the privacy of our users seriously. This document outlines our security policies, supported versions, and how to report vulnerabilities responsibly.

## Supported Versions

We provide security updates and bug fixes for the latest release versions:

| Version | Supported          |
| ------- | ------------------ |
| 1.3.x   | :white_check_mark: |
| 1.2.x   | :x:                |
| < 1.2.0 | :x:                |

We strongly encourage all users to stay updated with the latest release available from [GitHub Releases](https://github.com/hoangtran1411/sql-helper/releases).

## Local Data Privacy & Security Model

SQL Helper is architected as a local desktop utility:

- **100% Offline Processing**: All Excel reading, data parsing, and SQL generation occur entirely on your local machine.
- **No Remote Telemetry or Data Uploads**: Your spreadsheet contents, table names, and generated SQL files are never transmitted to any external server or third-party service.
- **Auto-Update Safety**: Application update checks query the public GitHub Releases API (`https://api.github.com/repos/hoangtran1411/sql-helper/releases/latest`) solely to compare version strings and download official release binaries.

## Reporting a Vulnerability

If you discover a potential security vulnerability in SQL Helper, please follow responsible disclosure practices:

1. **Do NOT report security vulnerabilities via public GitHub Issues or discussions.**
2. Report vulnerabilities privately via one of the following methods:
   - **GitHub Private Vulnerability Reporting**: Submit a private advisory through the [Security Advisories](https://github.com/hoangtran1411/sql-helper/security/advisories) page.
   - **Email**: Send a detailed report to [chi3xitin2010@gmail.com](mailto:chi3xitin2010@gmail.com) with the subject line `[Security] SQL Helper Vulnerability Report`.

### What to Include in Your Report

Please provide as much relevant information as possible to help us reproduce and resolve the issue quickly:

- Type of issue (e.g., buffer overflow, arbitrary code execution, unsafe file write, path traversal)
- Affected version(s) of SQL Helper
- Step-by-step instructions or proof-of-concept (PoC) to reproduce the vulnerability
- Any sample files (e.g., specially crafted `.xlsx` files) if applicable
- Potential impact of the issue

## Response Timeline

- **Initial Acknowledgment**: Within 48 hours of receiving the report.
- **Status Assessment**: Within 5 business days, including confirmation of the issue and preliminary assessment of severity.
- **Resolution & Release**: A patch will be prioritized and published as a security release via GitHub Releases as soon as verified.
- **Public Disclosure**: Coordinated disclosure will occur after a patched release is made available to users.
