<template>
<div class="databases">
    <div class="toolbar">
        <div class="db-picker">
            <span class="label">Database</span>
            <select v-if="databases.length" v-model="selectedDb" @change="onDbChange">
                <option v-for="db in databases" :key="db.name" :value="db.name">
                    {{ db.name }} ({{ db.host }})
                </option>
            </select>
            <p v-else class="empty">This account has no databases configured in its settings.</p>
        </div>

        <div v-if="selectedTable" class="toolbar-actions">
            <label class="limit-label" for="row-limit">Rows</label>
            <input id="row-limit" type="number" v-model.number="rowLimit" min="1" @keyup.enter="loadTableData" />
            <button class="btn-ghost" @click="loadTableData" :disabled="loading">Apply</button>
            <button class="btn-ghost" @click="loadTableData" :disabled="loading">Refresh</button>
            <button class="btn-primary" @click="startAdd" :disabled="!tableData">+ Add row</button>
        </div>
    </div>

    <div v-if="error" class="error-msg">{{ error }}</div>
    <div v-if="notice" class="notice-msg">{{ notice }}</div>

    <div v-if="databases.length" class="layout">
        <!-- compact table list -->
        <div class="table-list">
            <div class="table-list-head">Tables <span class="count">{{ tables.length }}</span></div>
            <div v-if="loadingTables" class="muted small">Loading…</div>
            <template v-else>
                <button
                    v-for="t in tables" :key="t.name"
                    class="table-item"
                    :class="{ active: t.name === selectedTable }"
                    @click="selectTable(t)"
                    :title="t.name + ' · ' + t.engine + (t.comment ? ' · ' + t.comment : '')">
                    <span class="tname">{{ t.name }}</span>
                    <span class="trows">{{ t.rows }}</span>
                </button>
                <p v-if="tables.length === 0" class="muted small">No tables.</p>
            </template>
        </div>

        <!-- viewer / edit form -->
        <div class="viewer">
            <Spinner v-if="loading" />
            <template v-else-if="tableData && selectedTable">
                <div class="viewer-head">
                    <span class="table-title">{{ selectedTable }}</span>
                    <span class="muted small">{{ tableData.rows.length }} row(s){{ hasMoreRows ? ' (limit reached — raise Rows to see more)' : '' }}</span>
                    <span v-if="!hasPK" class="pk-warning">no primary key — rows are view-only for this table</span>
                </div>
                <div class="table-scroll">
                    <table>
                        <thead>
                            <tr>
                                <th class="actions-col"></th>
                                <th v-for="c in tableData.columns" :key="c" :class="{ pk: isPkColumn(c) }">
                                    {{ c }}<span v-if="isPkColumn(c)" class="pk-mark" title="primary key"> ·PK</span>
                                </th>
                            </tr>
                        </thead>
                        <tbody>
                            <tr v-for="(row, ri) in tableData.rows" :key="ri">
                                <td class="actions-col">
                                    <button class="btn-mini" :disabled="!hasPK" :title="hasPK ? 'Edit row' : 'Editing requires a primary key'" @click="startEdit(row)">Edit</button>
                                </td>
                                <td v-for="(val, ci) in row" :key="ci" class="cell">
                                    <span v-if="isNull(val)" class="null-val">NULL</span>
                                    <template v-else>{{ formatCell(val) }}</template>
                                </td>
                            </tr>
                            <tr v-if="tableData.rows.length === 0">
                                <td :colspan="tableData.columns.length + 1" class="empty center">No rows.</td>
                            </tr>
                        </tbody>
                    </table>
                </div>
            </template>
            <p v-else-if="!selectedTable" class="muted">Select a table to browse its data.</p>
        </div>
    </div>

    <!-- edit / add row modal -->
    <Modal v-if="formMode">
        <div class="row-form">
            <div class="form-head">
                <h3>{{ formMode === 'add' ? 'Add row to' : 'Edit row in' }} <span class="table-title">{{ selectedTable }}</span></h3>
                <button class="btn-ghost" @click="cancelForm">✕ Close</button>
            </div>
            <p class="muted small">Leave a nullable field empty to store NULL. Auto-increment columns are managed for you. Primary keys are editable when adding a row (unless auto-increment) and locked when editing.</p>
            <div v-if="formError" class="error-msg">{{ formError }}</div>

            <div class="form-fields">
                <div v-for="f in formFields" :key="f.name" class="field">
                    <label :for="'f-' + f.name">
                        {{ f.name }}
                        <span class="type">({{ f.type }}{{ f.isPk ? ', PK' : '' }}{{ f.autoIncrement ? ', auto-inc' : '' }}{{ f.nullable ? ', NULL ok' : '' }})</span>
                    </label>
                    <input
                        v-if="!isLongText(f)"
                        :id="'f-' + f.name"
                        :type="inputType(f)"
                        v-model="f.value"
                        :readonly="(formMode === 'edit' && f.isPk) || (formMode === 'add' && f.autoIncrement)" />
                    <textarea
                        v-else
                        :id="'f-' + f.name"
                        rows="4"
                        v-model="f.value"></textarea>
                </div>
                <p v-if="formFields.length === 0" class="muted">Nothing to edit — all columns are auto-managed.</p>
            </div>

            <div class="form-actions">
                <button class="btn-primary" :disabled="saving" @click="saveRow">{{ saving ? 'Saving…' : (formMode === 'add' ? 'Insert row' : 'Save changes') }}</button>
                <button class="btn-ghost" :disabled="saving" @click="cancelForm">Cancel</button>
            </div>
        </div>
    </Modal>
</div>
</template>

<script setup>
import { ref, computed } from 'vue'
import axios from 'axios'
import Spinner from '../Spinner.vue'
import Modal from '../Modal.vue'

const props = defineProps({
    account: {
        type: Object,
        default: () => ({})
    }
})

const databases = computed(() => props.account.databases || [])

const selectedDb = ref(databases.value.length ? databases.value[0].name : '')
const selectedTable = ref('')
const rowLimit = ref(100)

const tables = ref([])
const tableData = ref(null)

const loadingTables = ref(false)
const loading = ref(false)
const saving = ref(false)
const error = ref('')
const notice = ref('')

const formMode = ref(null)          // null | 'add' | 'edit'
const formFields = ref([])
const formError = ref('')
const editingPK = ref({})

const apiBase = computed(() => `/api/accounts/${encodeURIComponent(props.account.name)}/databases/${encodeURIComponent(selectedDb.value)}`)

const hasPK = computed(() => !!(tableData.value && tableData.value.primaryKey && tableData.value.primaryKey.length))
const hasMoreRows = computed(() => !!(tableData.value && tableData.value.rows.length === tableData.value.limit))

function isPkColumn(name) {
    return !!(tableData.value && tableData.value.primaryKey?.includes(name))
}

function onDbChange() {
    selectedTable.value = ''
    tables.value = []
    tableData.value = null
    error.value = ''
    notice.value = ''
    cancelForm()
    loadTables()
}

async function loadTables() {
    if (!selectedDb.value) return
    loadingTables.value = true
    error.value = ''
    try {
        const res = await axios.get(`${apiBase.value}/tables`)
        if (res.data.error) { error.value = res.data.error; return }
        tables.value = res.data.tables || []
    } catch (err) {
        error.value = apiError(err, 'Could not load tables.')
    } finally {
        loadingTables.value = false
    }
}

function selectTable(t) {
    selectedTable.value = t.name
    error.value = ''
    notice.value = ''
    cancelForm()
    loadTableData()
}

async function loadTableData() {
    if (!selectedTable.value) return
    loading.value = true
    error.value = ''
    try {
        const res = await axios.get(`${apiBase.value}/tables/${encodeURIComponent(selectedTable.value)}`, {
            params: { limit: rowLimit.value || 100 }
        })
        if (res.data.error) { error.value = res.data.error; return }
        tableData.value = res.data
    } catch (err) {
        tableData.value = null
        error.value = apiError(err, 'Could not load table data.')
    } finally {
        loading.value = false
    }
}

function apiError(err, fallback) {
    return err.response?.data?.error || err.response?.data?.message || fallback
}

// ---------- display helpers ----------

function isNull(val) {
    return val === null || val === undefined
}

function formatCell(val) {
    if (isNull(val)) return ''
    if (typeof val === 'object') {
        try { return JSON.stringify(val) } catch (e) { return String(val) }
    }
    return String(val)
}

// ---------- column input mapping ----------

const numericTypes = ['int', 'bigint', 'smallint', 'mediumint', 'tinyint', 'decimal', 'numeric', 'float', 'double', 'bit', 'year']
const longTypes = ['text', 'tinytext', 'mediumtext', 'longtext', 'json', 'blob', 'tinyblob', 'mediumblob', 'longblob']

function inputType(col) {
    const t = (col.type || '').toLowerCase()
    if (numericTypes.includes(t)) return 'number'
    if (t === 'date') return 'date'
    if (t === 'time') return 'time'
    if (t === 'datetime' || t === 'timestamp') return 'datetime-local'
    return 'text'
}

function isLongText(col) {
    return longTypes.includes((col.type || '').toLowerCase())
}

// datetime values arrive as RFC3339 (parseTime=true); <input type=datetime-local>
// wants "YYYY-MM-DDTHH:MM:SS" in local time.
function toInputValue(col, val) {
    if (isNull(val)) return ''
    const t = (col.type || '').toLowerCase()
    if (t === 'datetime' || t === 'timestamp') {
        const d = new Date(val)
        if (!isNaN(d)) {
            const pad = (n) => String(n).padStart(2, '0')
            return `${d.getFullYear()}-${pad(d.getMonth() + 1)}-${pad(d.getDate())}T${pad(d.getHours())}:${pad(d.getMinutes())}:${pad(d.getSeconds())}`
        }
    }
    return String(val)
}

// ---------- row edit / add ----------

function buildFormFields(rowValues) {
    if (!tableData.value) return []
    return tableData.value.columnInfo.map((col, i) => ({
        ...col,
        value: rowValues ? toInputValue(col, rowValues[i]) : '',
    }))
}

function startEdit(row) {
    if (!hasPK.value) {
        error.value = 'This table has no primary key, so rows can\u2019t be edited safely.'
        return
    }
    formError.value = ''
    formMode.value = 'edit'
    formFields.value = buildFormFields(row)
    editingPK.value = {}
    for (const pk of tableData.value.primaryKey) {
        const idx = tableData.value.columnInfo.findIndex((c) => c.name === pk)
        editingPK.value[pk] = idx >= 0 ? String(row[idx] ?? '') : ''
    }
}

function startAdd() {
    formError.value = ''
    formMode.value = 'add'
    formFields.value = buildFormFields(null)
    editingPK.value = {}
}

function cancelForm() {
    formMode.value = null
    formFields.value = []
    formError.value = ''
    editingPK.value = {}
}

async function saveRow() {
    if (saving.value || !tableData.value) return
    saving.value = true
    formError.value = ''
    error.value = ''
    notice.value = ''

    // Build the multipart body the server expects: col_<name> (and pk_<name>
    // on updates). Server whitelists column names against information_schema
    // and binds every value as a ? parameter.
    const body = new FormData()
    for (const f of formFields.value) {
        if (formMode.value === 'add') {
            if (f.autoIncrement) continue           // server-assigned
            if (f.isPk && f.value === '') continue  // let MySQL default/assign the PK
        } else {
            if (f.isPk) continue                    // PK identifies the row; not editable
        }
        body.append('col_' + f.name, f.value ?? '')
    }
    if (formMode.value === 'edit') {
        for (const [pk, val] of Object.entries(editingPK.value)) {
            body.append('pk_' + pk, val)
        }
    }

    try {
        const url = `${apiBase.value}/tables/${encodeURIComponent(selectedTable.value)}/rows`
        const res = formMode.value === 'add'
            ? await axios.post(url, body)
            : await axios.put(url, body)
        if (res.data.error) { formError.value = res.data.error; return }
        notice.value = res.data.message || (formMode.value === 'add' ? 'Row added.' : 'Row updated.')
        cancelForm()
        await loadTableData()
    } catch (err) {
        formError.value = apiError(err, 'Save failed — values are validated server-side, check the column types.')
    } finally {
        saving.value = false
    }
}
</script>

<style lang="scss" scoped>
@import "../../scss/variables.scss";

.databases {
    display: flex;
    flex-direction: column;
    gap: 16px;
}

.toolbar {
    display: flex;
    justify-content: space-between;
    align-items: flex-end;
    gap: 16px;
    flex-wrap: wrap;
}

.db-picker {
    display: flex;
    align-items: center;
    gap: 10px;

    .label {
        font-size: 13px;
        color: $c-od-muted;
    }

    select {
        background: $c-od-bg;
        border: 1px solid $c-od-border;
        border-radius: 10px;
        padding: 8px 12px;
        color: $c-od-fg;
        font-size: 14px;
        max-width: 340px;
    }
}

.toolbar-actions {
    display: flex;
    align-items: center;
    gap: 8px;

    .limit-label { font-size: 13px; color: $c-od-muted; }

    input {
        width: 80px;
        background: $c-od-bg;
        border: 1px solid $c-od-border;
        border-radius: 8px;
        padding: 7px 10px;
        color: $c-od-fg;
        font-size: 13px;
    }
}

.error-msg { font-size: 13px; color: $c-od-danger; }
.notice-msg { font-size: 13px; color: $c-od-success; }
.muted { color: $c-od-muted; }
.small { font-size: 12px; }
.empty { font-size: 14px; color: $c-od-muted; }
.center { text-align: center; }

.layout {
    display: flex;
    gap: 16px;
    align-items: flex-start;
}

/* ---- compact table list ---- */
.table-list {
    flex: 0 0 170px;
    display: flex;
    flex-direction: column;
    gap: 2px;
    max-height: 480px;
    overflow-y: auto;

    .table-list-head {
        font-size: 11px;
        font-weight: 600;
        text-transform: uppercase;
        letter-spacing: 0.06em;
        color: $c-od-muted;
        padding: 4px 6px;

        .count { color: $c-od-fg; }
    }

    .table-item {
        display: flex;
        justify-content: space-between;
        align-items: center;
        gap: 6px;
        padding: 5px 8px;
        font-size: 12px;
        border-radius: 6px;
        border: 1px solid transparent;
        background: transparent;
        color: $c-od-muted;
        cursor: pointer;
        text-align: left;

        .tname {
            overflow: hidden;
            text-overflow: ellipsis;
            white-space: nowrap;
        }

        .trows {
            font-size: 10px;
            color: $c-od-muted;
            background: $c-od-bg;
            border-radius: 10px;
            padding: 1px 6px;
            flex-shrink: 0;
        }

        &:hover { background: $c-od-surface-hi; color: $c-od-fg; }

        &.active {
            background: $c-od-surface-hi;
            border-color: $c-od-accent;
            color: $c-od-fg;
            font-weight: 600;
        }
    }
}

/* ---- viewer ---- */
.viewer {
    flex: 1 1 auto;
    min-width: 0;
    display: flex;
    flex-direction: column;
    gap: 8px;
}

.viewer-head {
    display: flex;
    align-items: baseline;
    gap: 12px;

    .table-title { font-size: 15px; font-weight: 600; }
    .pk-warning { font-size: 12px; color: $c-od-warn; }
}

/* horizontal scroll INSIDE the viewer, so wide tables never stretch the page */
.table-scroll {
    overflow-x: auto;
    max-height: 480px;
    overflow-y: auto;
    border: 1px solid $c-od-border;
    border-radius: 10px;

    table {
        border-collapse: collapse;
        font-size: 13px;

        th, td {
            border-bottom: 1px solid $c-od-border;
            border-right: 1px solid $c-od-border;
            padding: 6px 10px;
            white-space: nowrap;
        }

        th {
            position: sticky;
            top: 0;
            background: $c-od-surface-hi;
            font-weight: 600;
            text-align: left;
            z-index: 1;

            &.pk { color: $c-od-accent; }
            .pk-mark { font-size: 10px; color: $c-od-accent; }
        }

        td.cell {
            max-width: 340px;
            overflow: hidden;
            text-overflow: ellipsis;
            background: $c-od-surface;
            color: $c-od-fg;

            &:hover {
                white-space: normal;
                overflow: visible;
                max-width: none;
                position: relative;
                z-index: 2;
                box-shadow: 0 0 0 1px $c-od-accent;
            }
        }

        .null-val {
            color: $c-od-muted;
            font-style: italic;
            font-size: 11px;
        }

        th.actions-col, td.actions-col {
            position: sticky;
            left: 0;
            background: $c-od-surface-hi;
            z-index: 2;
            width: 1%;
        }

        td.actions-col { background: $c-od-surface; }

        tbody tr:hover td.cell { background: $c-od-surface-hi; }
    }
}

/* ---- buttons ---- */
.btn-primary {
    padding: 8px 16px;
    border-radius: 8px;
    background: linear-gradient(135deg, $c-od-accent, $c-od-accent-sub);
    color: #fff;
    font-size: 13px;
    font-weight: 600;
    border: none;
    cursor: pointer;
    white-space: nowrap;

    &:hover { opacity: 0.92; }
    &:disabled { opacity: 0.5; cursor: not-allowed; }
}

.btn-ghost {
    padding: 7px 12px;
    border-radius: 8px;
    background: transparent;
    border: 1px solid $c-od-border;
    color: $c-od-muted;
    font-size: 13px;
    cursor: pointer;

    &:hover:not(:disabled) { background: $c-od-surface-hi; color: $c-od-fg; }
    &:disabled { opacity: 0.4; cursor: not-allowed; }
}

.btn-mini {
    padding: 2px 8px;
    border-radius: 6px;
    background: transparent;
    border: 1px solid $c-od-border;
    color: $c-od-muted;
    font-size: 11px;
    cursor: pointer;

    &:hover:not(:disabled) { background: $c-od-surface-hi; color: $c-od-fg; }
    &:disabled { opacity: 0.35; cursor: not-allowed; }
}

/* ---- row edit/add form ---- */
.row-form {
    display: flex;
    flex-direction: column;
    gap: 14px;

    .form-head {
        display: flex;
        justify-content: space-between;
        align-items: center;

        h3 { font-size: 16px; font-weight: 600; }
        .table-title { color: $c-od-accent; }
    }

    .form-fields {
        display: flex;
        flex-direction: column;
        gap: 10px;
        max-height: 55vh;
        overflow-y: auto;
        padding-right: 4px;

        .field {
            display: flex;
            flex-direction: column;
            gap: 4px;

            label {
                font-size: 12px;
                font-weight: 600;
                color: $c-od-fg;

                .type { color: $c-od-muted; font-weight: 400; }
            }

            input, textarea {
                background: $c-od-bg;
                border: 1px solid $c-od-border;
                border-radius: 8px;
                padding: 8px 12px;
                color: $c-od-fg;
                font-size: 13px;
                font-family: inherit;

                &:focus { border-color: $c-od-accent; outline: none; }
                &:read-only { opacity: 0.6; }
            }
        }
    }

    .form-actions { display: flex; gap: 10px; }
}
</style>
