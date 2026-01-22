/**
 * SQL Helper - Wails Frontend JavaScript
 * 
 * Đây là template JavaScript cho Wails frontend.
 * Các hàm Go sẽ được bind và gọi thông qua window.go.main.App
 * 
 * Khi tích hợp với Wails, uncomment các dòng gọi Go backend.
 */

// ===================================
// Application State
// ===================================
const AppState = {
    filePath: '',
    sheetNames: [],
    currentSheet: '',
    headers: [],
    dataRows: [],
    numberColumns: [],
    sqlResult: '',
    isProcessing: false
};

// ===================================
// DOM Elements
// ===================================
const elements = {
    // Buttons
    btnChooseFile: document.getElementById('btnChooseFile'),
    btnReplace: document.getElementById('btnReplace'),
    btnCopy: document.getElementById('btnCopy'),
    btnExport: document.getElementById('btnExport'),
    btnClear: document.getElementById('btnClear'),
    btnProcessSheet: document.getElementById('btnProcessSheet'),
    btnConfirmReplace: document.getElementById('btnConfirmReplace'),
    
    // Inputs
    fileInput: document.getElementById('fileInput'),
    valueToFind: document.getElementById('valueToFind'),
    replacementValue: document.getElementById('replacementValue'),
    sqlResult: document.getElementById('sqlResult'),
    sheetSelect: document.getElementById('sheetSelect'),
    
    // Containers
    columnSelectorCard: document.getElementById('columnSelectorCard'),
    columnCheckboxes: document.getElementById('columnCheckboxes'),
    loadingOverlay: document.getElementById('loadingOverlay'),
    toastContainer: document.getElementById('toastContainer'),
    confirmMessage: document.getElementById('confirmMessage'),
    
    // Modals
    sheetModal: null,
    confirmModal: null
};

// ===================================
// Initialize
// ===================================
document.addEventListener('DOMContentLoaded', () => {
    // Initialize Bootstrap modals
    elements.sheetModal = new bootstrap.Modal(document.getElementById('sheetModal'));
    elements.confirmModal = new bootstrap.Modal(document.getElementById('confirmModal'));
    
    // Bind event listeners
    bindEventListeners();
    
    console.log('SQL Helper initialized');
});

// ===================================
// Event Listeners
// ===================================
function bindEventListeners() {
    // File selection
    elements.btnChooseFile.addEventListener('click', handleChooseFile);
    elements.fileInput.addEventListener('change', handleFileSelected);
    
    // Find & Replace
    elements.btnReplace.addEventListener('click', handleReplaceClick);
    elements.btnConfirmReplace.addEventListener('click', handleConfirmReplace);
    
    // SQL Actions
    elements.btnCopy.addEventListener('click', handleCopy);
    elements.btnExport.addEventListener('click', handleExport);
    elements.btnClear.addEventListener('click', handleClear);
    
    // Sheet selection
    elements.btnProcessSheet.addEventListener('click', handleProcessSheet);
    
    // SQL textarea change
    elements.sqlResult.addEventListener('input', (e) => {
        AppState.sqlResult = e.target.value;
        updateButtonStates();
    });
}

// ===================================
// File Handling
// ===================================
function handleChooseFile() {
    // Trong Wails, sử dụng runtime.OpenFileDialog thay vì input file
    // Uncomment khi tích hợp Wails:
    /*
    window.go.main.App.OpenExcelFile().then(result => {
        if (result.filePath) {
            AppState.filePath = result.filePath;
            AppState.sheetNames = result.sheetNames;
            handleSheetSelection();
        }
    }).catch(err => {
        showToast('Error opening file: ' + err, 'error');
    });
    */
    
    // Fallback cho browser testing:
    elements.fileInput.click();
}

function handleFileSelected(event) {
    const file = event.target.files[0];
    if (!file) return;
    
    showLoading(true);
    
    // Trong Wails, file sẽ được xử lý bởi Go backend
    // Đây là placeholder cho browser testing
    
    setTimeout(() => {
        // Giả lập response từ Go backend
        AppState.sheetNames = ['Sheet1', 'Sheet2']; // Mock data
        handleSheetSelection();
        showLoading(false);
    }, 500);
    
    // Reset file input
    event.target.value = '';
}

function handleSheetSelection() {
    if (AppState.sheetNames.length === 0) {
        showToast('No sheets found in the Excel file.', 'error');
        return;
    }
    
    if (AppState.sheetNames.length === 1) {
        // Tự động process nếu chỉ có 1 sheet
        AppState.currentSheet = AppState.sheetNames[0];
        processSheet(AppState.currentSheet);
    } else {
        // Hiển thị modal chọn sheet
        populateSheetSelect();
        elements.sheetModal.show();
    }
}

function populateSheetSelect() {
    elements.sheetSelect.innerHTML = '';
    AppState.sheetNames.forEach(name => {
        const option = document.createElement('option');
        option.value = name;
        option.textContent = name;
        elements.sheetSelect.appendChild(option);
    });
}

function handleProcessSheet() {
    const selectedSheet = elements.sheetSelect.value;
    if (!selectedSheet) {
        showToast('Please select a sheet.', 'warning');
        return;
    }
    
    AppState.currentSheet = selectedSheet;
    elements.sheetModal.hide();
    processSheet(selectedSheet);
}

function processSheet(sheetName) {
    showLoading(true);
    
    // Uncomment khi tích hợp Wails:
    /*
    window.go.main.App.ProcessSheet(AppState.filePath, sheetName).then(result => {
        AppState.headers = result.headers;
        AppState.dataRows = result.dataRows;
        AppState.numberColumns = [];
        
        renderColumnSelector();
        generateSQL();
        showLoading(false);
        showToast('Sheet processed successfully!', 'success');
    }).catch(err => {
        showLoading(false);
        showToast('Error processing sheet: ' + err, 'error');
    });
    */
    
    // Mock data cho browser testing
    setTimeout(() => {
        AppState.headers = ['Name', 'Age', 'Email', 'Active', 'JoinDate'];
        AppState.dataRows = [
            ['John Doe', 30, 'john@example.com', 'TRUE', '2024-01-15'],
            ['Jane Smith', 25, 'jane@example.com', 'FALSE', '2024-02-20'],
            ['Bob Wilson', null, 'bob@example.com', 'TRUE', '2024-03-10']
        ];
        AppState.numberColumns = [];
        
        renderColumnSelector();
        generateSQL();
        showLoading(false);
        showToast('Sheet processed successfully!', 'success');
    }, 500);
}

// ===================================
// Column Selector
// ===================================
function renderColumnSelector() {
    if (AppState.headers.length === 0) {
        elements.columnSelectorCard.style.display = 'none';
        return;
    }
    
    elements.columnSelectorCard.style.display = 'block';
    elements.columnCheckboxes.innerHTML = '';
    
    AppState.headers.forEach((header, index) => {
        const div = document.createElement('div');
        div.className = 'form-check';
        div.innerHTML = `
            <input class="form-check-input" type="checkbox" 
                   id="col-${index}" 
                   data-header="${header}"
                   ${AppState.numberColumns.includes(header) ? 'checked' : ''}>
            <label class="form-check-label" for="col-${index}">${header}</label>
        `;
        elements.columnCheckboxes.appendChild(div);
    });
    
    // Bind checkbox events
    elements.columnCheckboxes.querySelectorAll('input[type="checkbox"]').forEach(checkbox => {
        checkbox.addEventListener('change', handleColumnToggle);
    });
}

function handleColumnToggle(event) {
    const header = event.target.dataset.header;
    const isChecked = event.target.checked;
    
    if (isChecked) {
        if (!AppState.numberColumns.includes(header)) {
            AppState.numberColumns.push(header);
        }
    } else {
        AppState.numberColumns = AppState.numberColumns.filter(h => h !== header);
    }
    
    generateSQL();
}

// ===================================
// SQL Generation (⭐ Core Logic)
// ===================================
function generateSQL() {
    if (AppState.dataRows.length === 0) {
        AppState.sqlResult = '';
        elements.sqlResult.value = '';
        updateButtonStates();
        return;
    }
    
    showLoading(true);
    
    // Uncomment khi tích hợp Wails (Go sẽ xử lý):
    /*
    window.go.main.App.GenerateSQL(
        AppState.headers,
        AppState.dataRows,
        AppState.numberColumns
    ).then(result => {
        AppState.sqlResult = result;
        elements.sqlResult.value = result;
        updateButtonStates();
        showLoading(false);
    }).catch(err => {
        showLoading(false);
        showToast('Error generating SQL: ' + err, 'error');
    });
    */
    
    // JavaScript implementation (cho browser testing)
    setTimeout(() => {
        const sqlLines = AppState.dataRows.map(row => {
            const formattedValues = row.map((cell, index) => {
                const header = AppState.headers[index];
                return formatCellValue(cell, AppState.numberColumns.includes(header));
            });
            return `(${formattedValues.join(', ')})`;
        });
        
        AppState.sqlResult = sqlLines.join(',\n');
        elements.sqlResult.value = AppState.sqlResult;
        updateButtonStates();
        showLoading(false);
    }, 100);
}

/**
 * Format cell value cho SQL
 * @param {*} value - Giá trị cell
 * @param {boolean} isNumeric - Column có phải numeric không
 * @returns {string} - Giá trị đã format
 */
function formatCellValue(value, isNumeric) {
    // Xử lý null/undefined
    if (value === null || value === undefined) {
        return isNumeric ? 'NULL' : "''";
    }
    
    // Xử lý empty string
    if (typeof value === 'string' && value.trim() === '') {
        return isNumeric ? 'NULL' : "''";
    }
    
    // Xử lý numeric column
    if (isNumeric) {
        const num = Number(value);
        if (isNaN(num)) {
            return 'NULL';
        }
        return String(num);
    }
    
    // Xử lý Date
    if (value instanceof Date) {
        return `'${formatDateForSQL(value)}'`;
    }
    
    // Xử lý string có thể là ISO date
    if (typeof value === 'string' && isISODateString(value)) {
        const date = new Date(value);
        if (!isNaN(date.getTime())) {
            return `'${formatDateForSQL(date)}'`;
        }
    }
    
    // Xử lý string thông thường - escape single quotes
    const stringValue = String(value).replace(/'/g, "''");
    return `'${stringValue}'`;
}

/**
 * Format Date thành SQL datetime string
 * @param {Date} date 
 * @returns {string} - YYYY-MM-DD HH:MI:SS
 */
function formatDateForSQL(date) {
    const pad = (n) => String(n).padStart(2, '0');
    
    const year = date.getFullYear();
    const month = pad(date.getMonth() + 1);
    const day = pad(date.getDate());
    const hours = pad(date.getHours());
    const minutes = pad(date.getMinutes());
    const seconds = pad(date.getSeconds());
    
    return `${year}-${month}-${day} ${hours}:${minutes}:${seconds}`;
}

/**
 * Kiểm tra string có phải ISO date format không
 * @param {string} str 
 * @returns {boolean}
 */
function isISODateString(str) {
    // Kiểm tra ISO 8601 format: YYYY-MM-DDTHH:MM:SS
    return typeof str === 'string' && 
           str.includes('T') && 
           !isNaN(Date.parse(str));
}

// ===================================
// Find & Replace
// ===================================
function handleReplaceClick() {
    const findValue = elements.valueToFind.value;
    const replaceValue = elements.replacementValue.value;
    
    if (AppState.dataRows.length === 0) {
        showToast('Nothing to replace. Upload a file first.', 'info');
        return;
    }
    
    if (!findValue) {
        showToast('Please enter a "Value to find".', 'warning');
        return;
    }
    
    // Cập nhật message và hiển thị modal xác nhận
    elements.confirmMessage.textContent = 
        `This will replace all occurrences of "${findValue}" with "${replaceValue}". Continue?`;
    elements.confirmModal.show();
}

function handleConfirmReplace() {
    const findValue = elements.valueToFind.value;
    const replaceValue = elements.replacementValue.value;
    
    elements.confirmModal.hide();
    showLoading(true);
    
    // Thực hiện replace
    AppState.dataRows = AppState.dataRows.map(row =>
        row.map(cell => {
            if (String(cell) === findValue) {
                return replaceValue;
            }
            return cell;
        })
    );
    
    // Regenerate SQL
    generateSQL();
    
    // Clear inputs
    elements.valueToFind.value = '';
    elements.replacementValue.value = '';
    
    showToast('Values replaced successfully!', 'success');
}

// ===================================
// SQL Actions
// ===================================
function handleCopy() {
    if (!AppState.sqlResult) {
        showToast('Nothing to copy.', 'info');
        return;
    }
    
    // Uncomment khi tích hợp Wails:
    /*
    window.go.main.App.CopyToClipboard(AppState.sqlResult).then(() => {
        showToast('Copied to clipboard!', 'success');
    }).catch(err => {
        showToast('Failed to copy: ' + err, 'error');
    });
    */
    
    // Browser fallback
    navigator.clipboard.writeText(AppState.sqlResult).then(() => {
        showToast('Copied to clipboard!', 'success');
    }).catch(err => {
        console.error('Copy failed:', err);
        showToast('Failed to copy.', 'error');
    });
}

function handleExport() {
    if (!AppState.sqlResult) {
        showToast('Nothing to export.', 'info');
        return;
    }
    
    // Uncomment khi tích hợp Wails:
    /*
    window.go.main.App.ExportToFile(AppState.sqlResult).then(() => {
        showToast('File exported successfully!', 'success');
    }).catch(err => {
        showToast('Failed to export: ' + err, 'error');
    });
    */
    
    // Browser fallback - download file
    try {
        const blob = new Blob([AppState.sqlResult], { type: 'text/plain;charset=utf-8' });
        const url = URL.createObjectURL(blob);
        const link = document.createElement('a');
        link.href = url;
        link.download = 'result.txt';
        document.body.appendChild(link);
        link.click();
        document.body.removeChild(link);
        URL.revokeObjectURL(url);
        showToast('File exported successfully!', 'success');
    } catch (err) {
        console.error('Export failed:', err);
        showToast('Failed to export file.', 'error');
    }
}

function handleClear() {
    AppState.headers = [];
    AppState.dataRows = [];
    AppState.numberColumns = [];
    AppState.sqlResult = '';
    AppState.currentSheet = '';
    
    elements.sqlResult.value = '';
    elements.columnSelectorCard.style.display = 'none';
    elements.columnCheckboxes.innerHTML = '';
    
    updateButtonStates();
    showToast('Result cleared.', 'info');
}

// ===================================
// UI Utilities
// ===================================
function updateButtonStates() {
    const hasResult = AppState.sqlResult.length > 0;
    elements.btnCopy.disabled = !hasResult;
    elements.btnExport.disabled = !hasResult;
    elements.btnClear.disabled = !hasResult;
}

function showLoading(show) {
    AppState.isProcessing = show;
    elements.loadingOverlay.style.display = show ? 'flex' : 'none';
    elements.btnChooseFile.disabled = show;
    elements.btnReplace.disabled = show;
}

/**
 * Hiển thị toast notification
 * @param {string} message - Nội dung thông báo
 * @param {string} type - Loại: 'success', 'error', 'warning', 'info'
 */
function showToast(message, type = 'info') {
    const toastId = `toast-${Date.now()}`;
    const iconMap = {
        success: 'fas fa-check-circle text-success',
        error: 'fas fa-times-circle text-danger',
        warning: 'fas fa-exclamation-triangle text-warning',
        info: 'fas fa-info-circle text-info'
    };
    
    const toastHTML = `
        <div id="${toastId}" class="toast toast-${type}" role="alert" aria-live="assertive">
            <div class="toast-header">
                <i class="${iconMap[type]} me-2"></i>
                <strong class="me-auto">${type.charAt(0).toUpperCase() + type.slice(1)}</strong>
                <button type="button" class="btn-close" data-bs-dismiss="toast"></button>
            </div>
            <div class="toast-body">
                ${message}
            </div>
        </div>
    `;
    
    elements.toastContainer.insertAdjacentHTML('beforeend', toastHTML);
    
    const toastElement = document.getElementById(toastId);
    const toast = new bootstrap.Toast(toastElement, { autohide: true, delay: 3000 });
    toast.show();
    
    // Tự động xóa toast element sau khi ẩn
    toastElement.addEventListener('hidden.bs.toast', () => {
        toastElement.remove();
    });
}

// ===================================
// Expose for debugging (optional)
// ===================================
window.AppState = AppState;
window.generateSQL = generateSQL;
