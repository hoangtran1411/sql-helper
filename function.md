# SQL Helper - Function Documentation for Go/Wails Migration

> **Mục đích:** Tài liệu này extract các core concepts, logic và ý tưởng từ project Next.js để chuyển đổi sang Golang + Wails desktop application.

---

## 📋 Tổng Quan Ứng Dụng

### Mô tả
**SQL Helper** là một utility tool giúp:
- Parse file Excel (`.xlsx`, `.xls`)
- Chuyển đổi dữ liệu thành SQL `INSERT` statements
- Hỗ trợ Find & Replace trước khi generate SQL
- Export kết quả ra clipboard hoặc file `.txt`

### Use Cases
1. DBA cần import data từ Excel vào database
2. Developer cần tạo seed data cho testing
3. Data migration từ spreadsheet sang SQL database

---

## 🏗️ Kiến Trúc Chuyển Đổi

### Next.js → Wails Mapping

| Next.js Component | Wails Equivalent | Ghi chú |
|-------------------|------------------|---------|
| React Component | HTML + Vanilla JS / Svelte / Vue | Frontend UI |
| Web Worker (`worker.ts`) | Go Backend Functions | Business Logic |
| `useState` / `useRef` | JavaScript Variables | State Management |
| `react-toastify` | Custom Toast hoặc Native Dialog | Notifications |
| `sweetalert2` | `runtime.MessageDialog` | Confirmation Dialogs |
| File Input | `runtime.OpenFileDialog` | File Selection |
| Clipboard API | `runtime.ClipboardSetText` | Copy to Clipboard |
| Blob Download | `os.WriteFile` | Export to File |

---

## 🔧 Core Functions (Backend - Go)

### 1. Excel Parser

```go
// ParseExcelFile đọc file Excel và trả về danh sách sheets
// Input: filePath string
// Output: []string (sheet names), error
func (a *App) ParseExcelFile(filePath string) ([]string, error) {
    // Sử dụng library: github.com/xuri/excelize/v2
    // 1. Mở file Excel
    // 2. Lấy danh sách tất cả sheet names
    // 3. Trả về cho frontend để user chọn (nếu > 1 sheet)
}
```

**Logic chi tiết:**
- Nếu file có 1 sheet → tự động process sheet đó
- Nếu file có nhiều sheets → hiển thị dialog để user chọn

---

### 2. Sheet Processor

```go
// ProcessSheet đọc dữ liệu từ một sheet cụ thể
// Input: filePath string, sheetName string
// Output: SheetData struct, error
type SheetData struct {
    Headers  []string        // Dòng đầu tiên = headers
    DataRows [][]interface{} // Các dòng còn lại = data
}

func (a *App) ProcessSheet(filePath, sheetName string) (*SheetData, error) {
    // 1. Đọc sheet theo sheetName
    // 2. Row đầu tiên → Headers
    // 3. Các row còn lại → DataRows
    // 4. Xử lý cell types: string, number, date, boolean, null
}
```

**Lưu ý xử lý:**
- Empty cells → `null` hoặc `""`
- Date cells → cần parse từ Excel serial number

---

### 3. SQL Generator (⭐ Core Logic)

```go
// GenerateSQL tạo SQL INSERT values từ data
// Input: sheetData *SheetData, numberColumns []string
// Output: string (SQL result)
func (a *App) GenerateSQL(sheetData *SheetData, numberColumns []string) string {
    // Với mỗi row trong DataRows:
    //   Với mỗi cell trong row:
    //     - Nếu column nằm trong numberColumns:
    //         - null/empty/NaN → "NULL"
    //         - có giá trị → giữ nguyên số
    //     - Nếu không phải number column:
    //         - null → "''"
    //         - Date → format "YYYY-MM-DD HH:MI:SS", wrap với quotes
    //         - String → escape single quotes ('), wrap với quotes
    //   Format: (val1, val2, val3, ...)
    // Join các rows bằng ",\n"
}
```

**SQL Formatting Rules:**

| Data Type | Example Input | SQL Output |
|-----------|---------------|------------|
| String | `Hello` | `'Hello'` |
| String với quotes | `It's OK` | `'It''s OK'` |
| Number (marked) | `123.45` | `123.45` |
| Number (unmarked) | `123` | `'123'` |
| Date | `2024-01-15 10:30:00` | `'2024-01-15 10:30:00'` |
| Null (number col) | `nil` | `NULL` |
| Null (string col) | `nil` | `''` |
| Empty string | `""` | `''` |

**Date Formatting Function:**

```go
func formatDateForSQL(t time.Time) string {
    return t.Format("2006-01-02 15:04:05")
}
```

---

### 4. Find & Replace

```go
// FindAndReplace thay thế giá trị trong data
// Input: sheetData *SheetData, findValue string, replaceWith string
// Output: *SheetData (modified)
func (a *App) FindAndReplace(sheetData *SheetData, findValue, replaceWith string) *SheetData {
    // Duyệt qua tất cả cells
    // Nếu String(cell) == findValue → thay bằng replaceWith
    // Trả về SheetData mới
}
```

**Lưu ý:**
- So sánh exact match (không phải substring)
- Sau khi replace cần regenerate SQL

---

### 5. Export Functions

```go
// CopyToClipboard copy text vào clipboard
func (a *App) CopyToClipboard(text string) error {
    return runtime.ClipboardSetText(a.ctx, text)
}

// ExportToFile lưu text ra file .txt
// Hiển thị SaveFileDialog để user chọn location
func (a *App) ExportToFile(content string) error {
    savePath, err := runtime.SaveFileDialog(a.ctx, runtime.SaveDialogOptions{
        DefaultFilename: "result.txt",
        Filters: []runtime.FileFilter{
            {DisplayName: "Text Files", Pattern: "*.txt"},
        },
    })
    if err != nil {
        return err
    }
    return os.WriteFile(savePath, []byte(content), 0644)
}
```

---

## 🎨 Frontend UI Components

### 1. File Uploader
- Button "Choose Excel File"
- Gọi `runtime.OpenFileDialog()` với filter `.xlsx, .xls`
- Hiển thị spinner khi đang processing

### 2. Sheet Selector (Optional Dialog)
- Chỉ hiển thị khi file có > 1 sheet
- Dropdown hoặc Radio buttons để chọn sheet
- Buttons: "Process" / "Cancel"

### 3. Column Selector
- Hiển thị sau khi parse thành công
- Checkbox list với tất cả column headers
- Check = đánh dấu column là NUMERIC
- Mỗi khi thay đổi → regenerate SQL

### 4. Find & Replace Panel
- Input: "Value to find"
- Input: "Replace with"
- Button: "Replace"
- Confirmation dialog trước khi replace

### 5. SQL Preview Panel
- Textarea hiển thị kết quả SQL (editable)
- Buttons:
  - 📋 **Copy** → Copy to clipboard
  - 📥 **Export** → Save to file
  - 🗑️ **Clear** → Reset tất cả

---

## 📊 State Management

### Application State

```go
type AppState struct {
    FilePath      string          // Đường dẫn file Excel hiện tại
    SheetNames    []string        // Danh sách sheets trong file
    CurrentSheet  string          // Sheet đang được xử lý
    Headers       []string        // Column headers
    DataRows      [][]interface{} // Raw data từ Excel
    NumberColumns []string        // Columns được đánh dấu là numeric
    SQLResult     string          // Kết quả SQL đã generate
    IsProcessing  bool            // Flag loading state
}
```

### State Flow

```
[User chọn file]
       ↓
[Parse Excel] → SheetNames
       ↓
[Chọn Sheet] (nếu > 1)
       ↓
[Process Sheet] → Headers, DataRows
       ↓
[Generate SQL] → SQLResult
       ↓
[User có thể:]
  ├── Toggle NumberColumns → Regenerate SQL
  ├── Find & Replace → Modify DataRows → Regenerate SQL
  ├── Edit SQL trực tiếp
  ├── Copy to Clipboard
  └── Export to File
```

---

## 🔄 Event Flow

### 1. File Upload Flow
```
1. User click "Choose File"
2. OpenFileDialog() → filePath
3. ParseExcelFile(filePath) → sheetNames
4. IF len(sheetNames) == 1:
     → ProcessSheet(filePath, sheetNames[0])
   ELSE:
     → Show Sheet Selection Dialog
     → User selects sheet
     → ProcessSheet(filePath, selectedSheet)
5. ProcessSheet() returns headers, dataRows
6. GenerateSQL(sheetData, []) → sqlResult
7. Display sqlResult in textarea
```

### 2. Column Toggle Flow
```
1. User checks/unchecks column checkbox
2. Update numberColumns array
3. GenerateSQL(sheetData, numberColumns)
4. Update sqlResult display
```

### 3. Find & Replace Flow
```
1. User enters find/replace values
2. User clicks "Replace"
3. Show confirmation dialog
4. IF confirmed:
     → FindAndReplace(sheetData, find, replace)
     → Update sheetData
     → GenerateSQL(sheetData, numberColumns)
     → Update sqlResult display
```

---

## 📦 Recommended Go Libraries

| Purpose | Library | Install |
|---------|---------|---------|
| Excel Parsing | excelize | `go get github.com/xuri/excelize/v2` |
| Wails Framework | wails | `go install github.com/wailsapp/wails/v2/cmd/wails@latest` |
| (Optional) TOML Config | toml | `go get github.com/BurntSushi/toml` |

---

## 🚀 Wails Project Structure

```
sql-helper/
├── main.go              # Entry point
├── app.go               # App struct với tất cả methods
├── wails.json           # Wails config
├── frontend/
│   ├── index.html       # Main HTML
│   ├── src/
│   │   ├── main.js      # Frontend logic
│   │   └── style.css    # Styles
│   └── dist/            # Build output
└── internal/
    ├── excel/
    │   └── parser.go    # Excel parsing logic
    └── sql/
        └── generator.go # SQL generation logic
```

---

## ⚡ Performance Considerations

### Original (Next.js)
- Sử dụng Web Worker để không block UI thread
- Heavy computation chạy trong worker

### Wails Equivalent
- Go backend tự động chạy trong goroutine riêng
- Frontend gọi Go functions không block UI
- Có thể dùng goroutines cho large file processing

---

## 🎯 Migration Checklist

- [ ] Setup Wails project
- [ ] Implement `ParseExcelFile()` với excelize
- [ ] Implement `ProcessSheet()`
- [ ] Implement `GenerateSQL()` (⭐ Critical)
- [ ] Implement `FindAndReplace()`
- [ ] Implement `CopyToClipboard()`
- [ ] Implement `ExportToFile()`
- [ ] Create HTML/CSS UI
- [ ] Bind Go functions to frontend
- [ ] Implement Sheet Selection Dialog
- [ ] Add toast/notification system
- [ ] Test with real Excel files
- [ ] Build và package

---

## 🧠 SQL Generator Logic - Deep Dive (⭐ Core Logic)

### Tổng quan

SQL Generator là **trái tim** của ứng dụng, chịu trách nhiệm chuyển đổi dữ liệu Excel thành SQL INSERT values với định dạng chính xác.

### Data Type Detection & Handling

```
┌─────────────────────────────────────────────────────────────────┐
│                    CELL VALUE INPUT                             │
└─────────────────────────────────────────────────────────────────┘
                              │
                              ▼
                    ┌─────────────────┐
                    │ Column marked   │
                    │ as NUMERIC?     │
                    └────────┬────────┘
                             │
              ┌──────────────┴──────────────┐
              │ YES                         │ NO
              ▼                             ▼
    ┌─────────────────┐           ┌─────────────────┐
    │ Xử lý Numeric   │           │ Xử lý Text      │
    └────────┬────────┘           └────────┬────────┘
             │                             │
             ▼                             ▼
    ┌─────────────────┐           ┌─────────────────┐
    │ null/empty/NaN? │           │ null/undefined? │
    └────────┬────────┘           └────────┬────────┘
             │                             │
      ┌──────┴──────┐               ┌──────┴──────┐
      │YES      │NO │               │YES      │NO │
      ▼         ▼   │               ▼         ▼   │
    NULL    Number  │             ''     Check    │
                    │                   Type      │
                    │                     │       │
                    │           ┌─────────┴───────┐
                    │           │                 │
                    │        Is Date?        Is String?
                    │           │                 │
                    │           ▼                 ▼
                    │    Format Date       Escape Quotes
                    │    YYYY-MM-DD        ' → ''
                    │    HH:MI:SS          Wrap with ''
                    │    Wrap with ''
                    │
                    ▼
              OUTPUT VALUE
```

### Processing Rules - Chi tiết

#### Rule 1: Numeric Columns
```
IF column IN numberColumns:
    IF value IS null OR undefined OR empty string:
        RETURN "NULL"
    ELSE IF NOT isNumber(value):
        RETURN "NULL"
    ELSE:
        RETURN String(value)  // Không có quotes
```

**Examples:**
| Input | Output |
|-------|--------|
| `123` | `123` |
| `45.67` | `45.67` |
| `null` | `NULL` |
| `""` | `NULL` |
| `"abc"` | `NULL` |
| `0` | `0` |
| `-10` | `-10` |

#### Rule 2: Text Columns (Non-Numeric)
```
IF column NOT IN numberColumns:
    IF value IS null OR undefined:
        RETURN "''"
    ELSE IF value IS Date:
        RETURN "'" + formatDate(value) + "'"
    ELSE IF value IS ISO_DateString:
        RETURN "'" + formatDate(parseDate(value)) + "'"
    ELSE:
        escapedValue = String(value).replace("'", "''")
        RETURN "'" + escapedValue + "'"
```

**Examples:**
| Input | Output |
|-------|--------|
| `"Hello"` | `'Hello'` |
| `"It's OK"` | `'It''s OK'` |
| `"O'Brien's"` | `'O''Brien''s'` |
| `null` | `''` |
| `Date(2024,0,15)` | `'2024-01-15 00:00:00'` |
| `"2024-01-15T10:30:00"` | `'2024-01-15 10:30:00'` |
| `123` | `'123'` (wrapped as string) |
| `true` | `'true'` |

### Date Formatting

#### Excel Date Handling

Excel lưu Date dưới dạng **serial number** (số ngày kể từ 1900-01-01).  
Khi parse với `excelize` hoặc `xlsx`, dates được chuyển thành `time.Time` hoặc `Date` object.

```go
// Go - Format date cho SQL
func formatDateForSQL(t time.Time) string {
    return t.Format("2006-01-02 15:04:05")
}
```

```javascript
// JavaScript - Format date cho SQL
function formatDateForSQL(date) {
    const pad = (n) => String(n).padStart(2, '0');
    
    return `${date.getFullYear()}-${pad(date.getMonth() + 1)}-${pad(date.getDate())} ` +
           `${pad(date.getHours())}:${pad(date.getMinutes())}:${pad(date.getSeconds())}`;
}
```

#### ISO Date Detection

```javascript
function isISODateString(str) {
    return typeof str === 'string' && 
           str.includes('T') && 
           !isNaN(Date.parse(str));
}
// "2024-01-15T10:30:00.000Z" → true
// "2024-01-15" → false (không có T)
// "Hello World" → false
```

### Complete Go Implementation

```go
package sql

import (
    "fmt"
    "strconv"
    "strings"
    "time"
)

// GenerateSQLValues tạo SQL INSERT values từ data
func GenerateSQLValues(
    headers []string,
    dataRows [][]interface{},
    numberColumns []string,
) string {
    if len(dataRows) == 0 {
        return ""
    }

    // Convert numberColumns to map for O(1) lookup
    numColSet := make(map[string]bool)
    for _, col := range numberColumns {
        numColSet[col] = true
    }

    columnCount := len(headers)
    var resultLines []string

    for _, row := range dataRows {
        var formattedValues []string

        for i := 0; i < columnCount; i++ {
            header := headers[i]
            
            // Lấy cell value, xử lý trường hợp row ngắn hơn headers
            var cellValue interface{}
            if i < len(row) {
                cellValue = row[i]
            }

            // Format cell dựa trên column type
            formatted := formatCellValue(cellValue, numColSet[header])
            formattedValues = append(formattedValues, formatted)
        }

        line := fmt.Sprintf("(%s)", strings.Join(formattedValues, ", "))
        resultLines = append(resultLines, line)
    }

    return strings.Join(resultLines, ",\n")
}

// formatCellValue format một cell value cho SQL
func formatCellValue(value interface{}, isNumeric bool) string {
    // Xử lý nil
    if value == nil {
        if isNumeric {
            return "NULL"
        }
        return "''"
    }

    // Xử lý numeric columns
    if isNumeric {
        return formatNumericValue(value)
    }

    // Xử lý text columns
    return formatTextValue(value)
}

// formatNumericValue format giá trị cho numeric column
func formatNumericValue(value interface{}) string {
    switch v := value.(type) {
    case int, int8, int16, int32, int64:
        return fmt.Sprintf("%d", v)
    case uint, uint8, uint16, uint32, uint64:
        return fmt.Sprintf("%d", v)
    case float32:
        return strconv.FormatFloat(float64(v), 'f', -1, 32)
    case float64:
        return strconv.FormatFloat(v, 'f', -1, 64)
    case string:
        // Thử parse string thành number
        s := strings.TrimSpace(v)
        if s == "" {
            return "NULL"
        }
        if _, err := strconv.ParseFloat(s, 64); err == nil {
            return s
        }
        return "NULL"
    case bool:
        if v {
            return "1"
        }
        return "0"
    default:
        return "NULL"
    }
}

// formatTextValue format giá trị cho text column
func formatTextValue(value interface{}) string {
    switch v := value.(type) {
    case time.Time:
        return fmt.Sprintf("'%s'", v.Format("2006-01-02 15:04:05"))
    case string:
        // Kiểm tra xem có phải ISO date string không
        if t, err := time.Parse(time.RFC3339, v); err == nil {
            return fmt.Sprintf("'%s'", t.Format("2006-01-02 15:04:05"))
        }
        // Escape single quotes
        escaped := strings.ReplaceAll(v, "'", "''")
        return fmt.Sprintf("'%s'", escaped)
    case bool:
        if v {
            return "'TRUE'"
        }
        return "'FALSE'"
    default:
        // Convert to string và escape
        s := fmt.Sprintf("%v", v)
        escaped := strings.ReplaceAll(s, "'", "''")
        return fmt.Sprintf("'%s'", escaped)
    }
}
```

### Edge Cases & Special Handling

| Case | Input | Expected Output | Notes |
|------|-------|-----------------|-------|
| Empty string (numeric) | `""` | `NULL` | Không thể convert |
| Empty string (text) | `""` | `''` | Keep as empty |
| Zero value | `0` | `0` | Valid number |
| Negative number | `-5` | `-5` | Valid number |
| Boolean true (numeric) | `true` | `1` | Convert to 1 |
| Boolean false (numeric) | `false` | `0` | Convert to 0 |
| Boolean true (text) | `true` | `'TRUE'` | Keep as string |
| Unicode characters | `"Việt Nam"` | `'Việt Nam'` | Keep unicode |
| Newline in string | `"Line1\nLine2"` | `'Line1\nLine2'` | Keep newlines |
| Tab in string | `"A\tB"` | `'A\tB'` | Keep tabs |
| HTML entities | `"<div>"` | `'<div>'` | Keep as-is |
| SQL injection attempt | `"'; DROP TABLE--"` | `'''; DROP TABLE--'` | Escaped! |

### Row Assembly

```
Row: ["John", 30, null, "2024-01-15T10:30:00Z", true]
Headers: ["Name", "Age", "Email", "JoinDate", "Active"]
NumberColumns: ["Age"]

Processing:
  Name (text):     "John"       → 'John'
  Age (numeric):   30           → 30
  Email (text):    null         → ''
  JoinDate (text): ISO Date     → '2024-01-15 10:30:00'
  Active (text):   true         → 'TRUE'

Output: ('John', 30, '', '2024-01-15 10:30:00', 'TRUE')
```

### Performance Considerations

```go
// Sử dụng strings.Builder cho performance tốt hơn với large datasets
func GenerateSQLValuesOptimized(
    headers []string,
    dataRows [][]interface{},
    numberColumns []string,
) string {
    if len(dataRows) == 0 {
        return ""
    }

    numColSet := make(map[string]bool)
    for _, col := range numberColumns {
        numColSet[col] = true
    }

    // Pre-allocate builder với estimated size
    estimatedSize := len(dataRows) * len(headers) * 20
    var builder strings.Builder
    builder.Grow(estimatedSize)

    columnCount := len(headers)

    for rowIdx, row := range dataRows {
        if rowIdx > 0 {
            builder.WriteString(",\n")
        }
        builder.WriteByte('(')

        for i := 0; i < columnCount; i++ {
            if i > 0 {
                builder.WriteString(", ")
            }

            var cellValue interface{}
            if i < len(row) {
                cellValue = row[i]
            }

            formatted := formatCellValue(cellValue, numColSet[headers[i]])
            builder.WriteString(formatted)
        }

        builder.WriteByte(')')
    }

    return builder.String()
}
```

### Unit Test Cases

```go
func TestFormatCellValue(t *testing.T) {
    tests := []struct {
        name      string
        value     interface{}
        isNumeric bool
        expected  string
    }{
        // Numeric column tests
        {"nil numeric", nil, true, "NULL"},
        {"int numeric", 42, true, "42"},
        {"float numeric", 3.14, true, "3.14"},
        {"empty string numeric", "", true, "NULL"},
        {"invalid string numeric", "abc", true, "NULL"},
        {"valid string number", "123", true, "123"},
        
        // Text column tests
        {"nil text", nil, false, "''"},
        {"simple string", "Hello", false, "'Hello'"},
        {"string with quote", "It's OK", false, "'It''s OK'"},
        {"number as text", 123, false, "'123'"},
        {"bool true text", true, false, "'TRUE'"},
        {"date value", time.Date(2024, 1, 15, 10, 30, 0, 0, time.UTC), false, 
            "'2024-01-15 10:30:00'"},
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            result := formatCellValue(tt.value, tt.isNumeric)
            if result != tt.expected {
                t.Errorf("formatCellValue(%v, %v) = %s; want %s",
                    tt.value, tt.isNumeric, result, tt.expected)
            }
        })
    }
}
```

---

## 📝 Sample SQL Output

**Input Excel:**
| Name | Age | JoinDate | Active |
|------|-----|----------|--------|
| John | 25 | 2024-01-15 | TRUE |
| Jane | 30 | 2024-02-20 | FALSE |
| Bob | | 2024-03-10 | TRUE |

**With Age marked as NUMERIC:**
```sql
('John', 25, '2024-01-15 00:00:00', 'TRUE'),
('Jane', 30, '2024-02-20 00:00:00', 'FALSE'),
('Bob', NULL, '2024-03-10 00:00:00', 'TRUE')
```

**Without any NUMERIC columns:**
```sql
('John', '25', '2024-01-15 00:00:00', 'TRUE'),
('Jane', '30', '2024-02-20 00:00:00', 'FALSE'),
('Bob', '', '2024-03-10 00:00:00', 'TRUE')
```

---

## 🔗 Related Resources

- [Wails Documentation](https://wails.io/docs/introduction)
- [Excelize Documentation](https://xuri.me/excelize/en/)
- [Go String Formatting](https://pkg.go.dev/fmt)

---

*Document generated: 2026-01-22*
*Source: Next.js SQL Helper Project*
