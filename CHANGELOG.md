# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

---

## [Unreleased]

### Refactored

- **Unified SQL Batch Streaming (`internal/sql`)**:
  - Introduced `sql.BatchWriter` with `NewBatchWriter`, `WriteRow`, and `Close` to handle SQL prefix generation, column formatting, batch delimiters, and statement boundaries into any `io.Writer`.
  - Refactored `GenerateBatchSQL` to stream directly into `strings.Builder` using `BatchWriter`.
  - Decoupled `ExportSQLStream` from desktop UI dialogs in `app.go`, enabling direct `O(1)` memory streaming to files or buffers without GUI mocks.
- **Consolidated Tabular Ingestion & Transformations (`internal/excel`)**:
  - Relocated cell transformation (`FindAndReplace`) from `internal/sql` to `internal/excel` to respect domain boundaries and layer responsibilities.
  - Unified empty-string-to-`nil` cell normalization and column bounds padding into `internal/excel.normalizeRow`.
  - Implemented `excel.StreamRows` streaming normalized `[]any` rows directly to any consumer callback (`bw.WriteRow`), eliminating manual row conversions and duplicate replacement loops from `app.go`.
  - Simplified `ExportSQLStream` from 44 lines down to 22 lines.
  - Provided backward-compatible type alias `type Replacement = excel.Replacement` in `app.go`.
- **Buffer Ownership Clarification (`app.go`)**:
  - Preserved explicit buffer ownership: `GenerateAndSaveSQL` manages `bufio.Writer` allocation and explicitly calls `Flush()`, removing ad-hoc flush type assertions from `BatchWriter.Close()`.

### Added

- **Comprehensive Deterministic Testing**:
  - Added unit tests for `ExportSQLStream` in `app_test.go` covering standard batch exports, multi-batch boundary chunking, real-time find-and-replace transformations, `ValuesOnly` mode, column filtering, error scenarios, and end-to-end file exports.
  - Added table-driven tests for `sql.BatchWriter` in `internal/sql/generator_test.go` covering nil writer validation, batch sizes (unlimited, 1, custom), `ValuesOnly` mode, column subsets, and error writers.
  - Added unit tests for `normalizeRow`, `FindAndReplace`, and `StreamRows` in `internal/excel/parser_test.go`.
  - Added `TestAppFindAndReplace_ExcelDelegation` in `app_test.go` validating backend-to-Excel delegation.
- **Governance & Policies**:
  - Added `CODE_OF_CONDUCT.md` and `SECURITY.md`.
  - Added automated Dependabot merge workflows.

### Removed

- Evacuated `FindAndReplace` from `internal/sql/generator.go` and its associated test cases from `internal/sql/generator_test.go`.
- Deleted duplicate string replacement loops and manual empty-string-to-`nil` row mapping in `app.go`.
- Removed legacy dummy file writing test (`TestExportToFileLogic`) in favor of direct end-to-end streaming tests (`TestExportSQLStream_ToFile`).

### Changed

- Elevated test statement coverage to **97.0%** in `internal/excel` and **96.9%** in `internal/sql` (well above the 70% threshold).
- Upgraded `github.com/xuri/excelize/v2` to `v2.11.0`.
- Upgraded `github.com/wailsapp/wails/v3` dependency.

---

## [1.3.0] - 2026-09-11

### Added

- Implemented `sql.ResolveColumns` for robust column name mapping, validation, and numeric set lookup.
- Added streaming iterator `excel.IterateSheet` for `O(1)` memory consumption on massive Excel spreadsheets.
- Added table-driven tests for column resolution and streaming sheet iteration.

### Changed

- Replaced `interface{}` with `any` across SQL generation functions and tests.
- Improved version comparison and semver parsing in auto-updater.

---

## [1.2.0] - 2026-09-09

### Fixed

- Corrected checkout action versions across CI and release workflows.
- Adjusted Wails v3 bindings generation step ordering in GitHub Actions.

---

## [1.1.1] - 2026-09-08

### Added

- Modern dashboard layout revamp with Bootstrap 5 and fixed toolbars.
- Responsive update notification badges and progress toasts for desktop auto-updater.

---

## [1.1.0] - 2026-09-07

### Added

- Excel file reading and sheet selection modal for multi-sheet workbooks.
- SQL `INSERT INTO` and raw values generation with column selector cards.
- Real-time find-and-replace tool for preview data rows.
- Native clipboard copying and save file dialog integration via Wails v3 runtime.

---

## [1.0.0] - 2026-09-06

### Added

- Initial release of SQL Helper desktop application.
- Single-executable Windows desktop build powered by Wails v3 and Go.
