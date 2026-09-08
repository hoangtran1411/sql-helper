/**
 * SQL Helper - Wails v3 Frontend JavaScript
 * 
 * Frontend logic with Wails v3 ES module bindings.
 */

import * as App from "./bindings/github.com/hoangtran1411/sql-helper/app.js";
import { Events } from "@wailsio/runtime";

// ===================================
// Application State
// ===================================
const AppState = {
    filePath: '',
    sheetNames: [],
    currentSheet: '',
    headers: [],
    dataRows: [],
    selectedColumns: [],
    numberColumns: [],
    tableName: '',
    batchSize: 1000,
    valuesOnly: false,
    replacements: [],
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
    btnSelectAllCols: document.getElementById('btnSelectAllCols'),
    btnDeselectAllCols: document.getElementById('btnDeselectAllCols'),

    // Inputs
    valueToFind: document.getElementById('valueToFind'),
    replacementValue: document.getElementById('replacementValue'),
    sqlResult: document.getElementById('sqlResult'),
    sheetSelect: document.getElementById('sheetSelect'),
    formatInsert: document.getElementById('formatInsert'),
    formatValues: document.getElementById('formatValues'),
    tableNameInput: document.getElementById('tableNameInput'),
    batchSizeInput: document.getElementById('batchSizeInput'),

    // Containers & Groups
    tableNameGroup: document.getElementById('tableNameGroup'),
    batchSizeGroup: document.getElementById('batchSizeGroup'),
    columnSelectorCard: document.getElementById('columnSelectorCard'),
    columnTableBody: document.getElementById('columnTableBody'),
    loadingOverlay: document.getElementById('loadingOverlay'),
    toastContainer: document.getElementById('toastContainer'),
    confirmMessage: document.getElementById('confirmMessage'),

    // Modals
    sheetModal: null,
    confirmModal: null,

    // Footer
    appVersion: document.getElementById('appVersion'),
    updateBadge: document.getElementById('updateBadge'),
    newVersionLabel: document.getElementById('newVersionLabel')
};

// ===================================
// Helper functions
// ===================================
function escapeHtml(str) {
    if (str === null || str === undefined) return '';
    return String(str)
        .replace(/&/g, '&amp;')
        .replace(/</g, '&lt;')
        .replace(/>/g, '&gt;')
        .replace(/"/g, '&quot;')
        .replace(/'/g, '&#039;');
}

// ===================================
// Initialize
// ===================================
document.addEventListener('DOMContentLoaded', () => {
    // Initialize Bootstrap modals
    elements.sheetModal = new bootstrap.Modal(document.getElementById('sheetModal'));
    elements.confirmModal = new bootstrap.Modal(document.getElementById('confirmModal'));

    // Bind event listeners
    bindEventListeners();

    // Listen for backend progress events
    if (Events && Events.On) {
        Events.On('updateProgress', (event) => {
            const msg = typeof event === 'string' ? event : (event?.data || 'Updating...');
            showToast(msg, 'info');
        });
    }

    // Display version and check for updates
    displayCurrentVersion();
    checkUpdate();

    console.log('SQL Helper v3 initialized');
});

// ===================================
// Auto Update
// ===================================
async function displayCurrentVersion() {
    try {
        const ver = await App.GetCurrentVersion();
        if (ver) {
            elements.appVersion.textContent = ver;
        }
    } catch (err) {
        console.error('Failed to get version:', err);
    }
}

async function checkUpdate() {
    try {
        const info = await App.CheckForUpdate();

        if (info && info.available) {
            elements.newVersionLabel.textContent = info.latestVersion;
            elements.updateBadge.style.display = 'inline-flex';
            elements.updateBadge.onclick = () => handleUpdate(info.downloadUrl);
        }
    } catch (err) {
        console.error('Failed to check for updates:', err);
    }
}

// Update handler
async function handleUpdate(url) {
    if (!confirm('Download and install new update? The app will restart.')) {
        return;
    }

    showLoading(true);
    try {
        await App.PerformUpdate(url);
    } catch (err) {
        showLoading(false);
        showToast('Update failed: ' + err, 'error');
    }
}

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

    // SQL Format & Config
    if (elements.formatInsert && elements.formatValues) {
        elements.formatInsert.addEventListener('change', handleFormatChange);
        elements.formatValues.addEventListener('change', handleFormatChange);
    }

    if (elements.tableNameInput) {
        elements.tableNameInput.addEventListener('input', (e) => {
            AppState.tableName = e.target.value.trim();
            debouncedGenerateSQL(150);
        });
    }

    if (elements.batchSizeInput) {
        elements.batchSizeInput.addEventListener('input', (e) => {
            const val = parseInt(e.target.value, 10);
            AppState.batchSize = isNaN(val) || val < 0 ? 0 : val;
            debouncedGenerateSQL(150);
        });
    }

    // Column selection shortcuts
    if (elements.btnSelectAllCols) {
        elements.btnSelectAllCols.addEventListener('click', handleSelectAllColumns);
    }
    if (elements.btnDeselectAllCols) {
        elements.btnDeselectAllCols.addEventListener('click', handleDeselectAllColumns);
    }

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
        const result = await App.OpenExcelFile();

        const filePath = result.filePath || result.FilePath;
        if (!result || !filePath) {
            showLoading(false);
            return; // User cancelled
        }

        AppState.filePath = filePath;
        AppState.sheetNames = result.sheetNames || result.SheetNames || [];
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

function handleFormatChange() {
    AppState.valuesOnly = elements.formatValues.checked;
    if (elements.tableNameGroup && elements.batchSizeGroup) {
        elements.tableNameGroup.style.opacity = AppState.valuesOnly ? '0.5' : '1';
        elements.batchSizeGroup.style.opacity = AppState.valuesOnly ? '0.5' : '1';
        elements.tableNameInput.disabled = AppState.valuesOnly;
        elements.batchSizeInput.disabled = AppState.valuesOnly;
    }
    generateSQL();
}

async function processSheet(sheetName) {
    showLoading(true);

    try {
        const result = await App.ProcessSheet(AppState.filePath, sheetName);

        if (!result) {
            showLoading(false);
            showToast('Failed to read sheet data', 'error');
            return;
        }

        AppState.headers = result.headers || result.Headers || [];
        AppState.dataRows = result.dataRows || result.DataRows || [];
        AppState.selectedColumns = [...AppState.headers]; // default: all columns selected
        AppState.numberColumns = [];
        AppState.replacements = [];
        AppState.tableName = sheetName;

        if (elements.tableNameInput) {
            elements.tableNameInput.value = sheetName;
        }
        if (elements.batchSizeInput && !elements.batchSizeInput.value) {
            elements.batchSizeInput.value = AppState.batchSize;
        }

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
    elements.columnTableBody.innerHTML = '';

    AppState.headers.forEach((header, index) => {
        const tr = document.createElement('tr');
        const isSelected = AppState.selectedColumns.includes(header);
        const isNumeric = AppState.numberColumns.includes(header);

        tr.innerHTML = `
            <td class="text-center text-muted fw-semibold">${index + 1}</td>
            <td class="fw-medium">${escapeHtml(header)}</td>
            <td class="text-center">
                <div class="form-check d-inline-block">
                    <input class="form-check-input col-include-check" type="checkbox" 
                           id="col-inc-${index}" data-header="${escapeHtml(header)}"
                           ${isSelected ? 'checked' : ''}>
                </div>
            </td>
            <td class="text-center">
                <div class="form-check d-inline-block">
                    <input class="form-check-input col-numeric-check" type="checkbox" 
                           id="col-num-${index}" data-header="${escapeHtml(header)}"
                           ${isNumeric ? 'checked' : ''} ${isSelected ? '' : 'disabled'}>
                </div>
            </td>
        `;
        elements.columnTableBody.appendChild(tr);
    });

    // Bind checkbox events
    elements.columnTableBody.querySelectorAll('.col-include-check').forEach(chk => {
        chk.addEventListener('change', handleColumnIncludeToggle);
    });
    elements.columnTableBody.querySelectorAll('.col-numeric-check').forEach(chk => {
        chk.addEventListener('change', handleColumnNumericToggle);
    });
}

async function handleColumnIncludeToggle(event) {
    const header = event.target.dataset.header;
    const isChecked = event.target.checked;
    const tr = event.target.closest('tr');
    const numCheck = tr ? tr.querySelector('.col-numeric-check') : null;

    if (isChecked) {
        if (!AppState.selectedColumns.includes(header)) {
            AppState.selectedColumns.push(header);
        }
        if (numCheck) numCheck.disabled = false;
    } else {
        AppState.selectedColumns = AppState.selectedColumns.filter(h => h !== header);
        if (numCheck) numCheck.disabled = true;
    }

    await generateSQL();
}

async function handleColumnNumericToggle(event) {
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

let debounceTimer = null;
function debouncedGenerateSQL(delay = 150) {
    if (debounceTimer) clearTimeout(debounceTimer);
    debounceTimer = setTimeout(() => {
        generateSQL();
    }, delay);
}

function updateColumnCheckboxes() {
    if (!elements.columnTableBody) return;
    elements.columnTableBody.querySelectorAll('tr').forEach(tr => {
        const incCheck = tr.querySelector('.col-include-check');
        const numCheck = tr.querySelector('.col-numeric-check');
        if (incCheck) {
            const h = incCheck.dataset.header;
            const isSelected = AppState.selectedColumns.includes(h);
            incCheck.checked = isSelected;
            if (numCheck) {
                numCheck.disabled = !isSelected;
                numCheck.checked = AppState.numberColumns.includes(h);
            }
        }
    });
}

async function handleSelectAllColumns() {
    AppState.selectedColumns = [...AppState.headers];
    updateColumnCheckboxes();
    await generateSQL();
}

async function handleDeselectAllColumns() {
    AppState.selectedColumns = [];
    updateColumnCheckboxes();
    await generateSQL();
}

// ===================================
// SQL Generation
// ===================================
async function generateSQL() {
    if (AppState.dataRows.length === 0 || AppState.selectedColumns.length === 0) {
        AppState.sqlResult = '';
        elements.sqlResult.value = '';
        updateButtonStates();
        return;
    }

    try {
        const options = {
            tableName: AppState.tableName || AppState.currentSheet || 'my_table',
            selectedColumns: AppState.selectedColumns,
            numberColumns: AppState.numberColumns,
            batchSize: typeof AppState.batchSize === 'number' ? AppState.batchSize : 1000,
            valuesOnly: AppState.valuesOnly
        };

        const result = await App.GenerateSQL(
            AppState.headers,
            AppState.dataRows,
            options
        );

        AppState.sqlResult = result || '';
        elements.sqlResult.value = AppState.sqlResult;
        updateButtonStates();
    } catch (err) {
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
        AppState.dataRows = await App.FindAndReplace(
            AppState.dataRows,
            findValue,
            replaceValue
        );

        AppState.replacements.push({
            find: findValue,
            replace: replaceValue
        });

        await generateSQL();

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
        if (navigator.clipboard && navigator.clipboard.writeText) {
            await navigator.clipboard.writeText(AppState.sqlResult);
        } else {
            await App.CopyToClipboard(AppState.sqlResult);
        }
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
        const options = {
            tableName: AppState.tableName || AppState.currentSheet || 'my_table',
            selectedColumns: AppState.selectedColumns,
            numberColumns: AppState.numberColumns,
            batchSize: typeof AppState.batchSize === 'number' ? AppState.batchSize : 1000,
            valuesOnly: AppState.valuesOnly
        };

        const saved = await App.GenerateAndSaveSQL(
            AppState.filePath,
            AppState.currentSheet,
            AppState.headers,
            options,
            AppState.replacements
        );
        if (saved) {
            showToast('File exported successfully!', 'success');
        }
    } catch (err) {
        if (!String(err).toLowerCase().includes('cancel')) {
            showToast('Failed to export: ' + err, 'error');
        }
    }
}

function handleClear() {
    AppState.headers = [];
    AppState.dataRows = [];
    AppState.selectedColumns = [];
    AppState.numberColumns = [];
    AppState.replacements = [];
    AppState.sqlResult = '';
    AppState.currentSheet = '';
    AppState.tableName = '';

    elements.sqlResult.value = '';
    if (elements.tableNameInput) elements.tableNameInput.value = '';
    elements.columnSelectorCard.style.display = 'none';
    if (elements.columnTableBody) elements.columnTableBody.innerHTML = '';

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

    toastElement.addEventListener('hidden.bs.toast', () => {
        toastElement.remove();
    });
}

// Expose for debugging
window.AppState = AppState;
