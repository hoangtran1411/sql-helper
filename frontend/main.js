/**
 * SQL Helper - Wails Frontend JavaScript
 * 
 * Frontend logic với Wails bindings để gọi Go backend.
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
async function handleChooseFile() {
    showLoading(true);
    
    try {
        const result = await window.go.main.App.OpenExcelFile();
        
        if (!result || !result.filePath) {
            showLoading(false);
            return; // User cancelled
        }
        
        AppState.filePath = result.filePath;
        AppState.sheetNames = result.sheetNames;
        handleSheetSelection();
    } catch (err) {
        showLoading(false);
        showToast('Error opening file: ' + err, 'error');
    }
}

function handleSheetSelection() {
    showLoading(false);
    
    if (AppState.sheetNames.length === 0) {
        showToast('No sheets found in the Excel file.', 'error');
        return;
    }
    
    if (AppState.sheetNames.length === 1) {
        // Auto process if only 1 sheet
        AppState.currentSheet = AppState.sheetNames[0];
        processSheet(AppState.currentSheet);
    } else {
        // Show sheet selection modal
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

async function processSheet(sheetName) {
    showLoading(true);
    
    try {
        const result = await window.go.main.App.ProcessSheet(AppState.filePath, sheetName);
        
        AppState.headers = result.headers;
        AppState.dataRows = result.dataRows;
        AppState.numberColumns = [];
        
        renderColumnSelector();
        await generateSQL();
        showLoading(false);
        showToast('Sheet processed successfully!', 'success');
    } catch (err) {
        showLoading(false);
        showToast('Error processing sheet: ' + err, 'error');
    }
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

async function handleColumnToggle(event) {
    const header = event.target.dataset.header;
    const isChecked = event.target.checked;
    
    if (isChecked) {
        if (!AppState.numberColumns.includes(header)) {
            AppState.numberColumns.push(header);
        }
    } else {
        AppState.numberColumns = AppState.numberColumns.filter(h => h !== header);
    }
    
    await generateSQL();
}

// ===================================
// SQL Generation
// ===================================
async function generateSQL() {
    if (AppState.dataRows.length === 0) {
        AppState.sqlResult = '';
        elements.sqlResult.value = '';
        updateButtonStates();
        return;
    }
    
    showLoading(true);
    
    try {
        const result = await window.go.main.App.GenerateSQL(
            AppState.headers,
            AppState.dataRows,
            AppState.numberColumns
        );
        
        AppState.sqlResult = result;
        elements.sqlResult.value = result;
        updateButtonStates();
        showLoading(false);
    } catch (err) {
        showLoading(false);
        showToast('Error generating SQL: ' + err, 'error');
    }
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
    
    // Update message and show confirmation modal
    elements.confirmMessage.textContent = 
        `This will replace all occurrences of "${findValue}" with "${replaceValue}". Continue?`;
    elements.confirmModal.show();
}

async function handleConfirmReplace() {
    const findValue = elements.valueToFind.value;
    const replaceValue = elements.replacementValue.value;
    
    elements.confirmModal.hide();
    showLoading(true);
    
    try {
        // Call Go backend for replacement
        AppState.dataRows = await window.go.main.App.FindAndReplace(
            AppState.dataRows,
            findValue,
            replaceValue
        );
        
        // Regenerate SQL
        await generateSQL();
        
        // Clear inputs
        elements.valueToFind.value = '';
        elements.replacementValue.value = '';
        
        showToast('Values replaced successfully!', 'success');
    } catch (err) {
        showLoading(false);
        showToast('Error replacing values: ' + err, 'error');
    }
}

// ===================================
// SQL Actions
// ===================================
async function handleCopy() {
    if (!AppState.sqlResult) {
        showToast('Nothing to copy.', 'info');
        return;
    }
    
    try {
        await window.go.main.App.CopyToClipboard(AppState.sqlResult);
        showToast('Copied to clipboard!', 'success');
    } catch (err) {
        showToast('Failed to copy: ' + err, 'error');
    }
}

async function handleExport() {
    if (!AppState.sqlResult) {
        showToast('Nothing to export.', 'info');
        return;
    }
    
    try {
        await window.go.main.App.ExportToFile(AppState.sqlResult);
        showToast('File exported successfully!', 'success');
    } catch (err) {
        showToast('Failed to export: ' + err, 'error');
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
 * Show toast notification
 * @param {string} message - Message content
 * @param {string} type - Type: 'success', 'error', 'warning', 'info'
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
    
    // Remove toast element after hidden
    toastElement.addEventListener('hidden.bs.toast', () => {
        toastElement.remove();
    });
}

// ===================================
// Expose for debugging (optional)
// ===================================
window.AppState = AppState;
window.generateSQL = generateSQL;
