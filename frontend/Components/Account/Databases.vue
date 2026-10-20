<template>
<div class="databases">
    <div class="toolbar">
        <div class="db-picker">
            <span class="label">Database</span>
            <select v-if="databases.length" v-model="selectedDb">
                <option v-for="db in databases" :key="db.name" :value="db.name">
                    {{ db.name }} ({{ db.host }})
                </option>
            </select>
            <p v-else class="empty">This account has no databases configured in its settings.</p>
        </div>

        <div v-if="selectedTable" class="toolbar-actions">
            <label class="limit-label" for="row-limit">Rows</label>
            <input id="row-limit" type="number" v-model.number="rowLimit" min="1" @keyup.enter="applyLimit" />
            <button class="btn-ghost" @click="applyLimit" :disabled="loading">Apply</button>
            <button class="btn-ghost" @click="loadTableData" :disabled="loading">Refresh</button>
            <button class="btn-primary" @click="startAdd" :disabled="!tableData">+ Add row</button>
        </div>
    </div>

    <div v-if="error" class="error-msg">{{ error }}</div>

    <transition name="toast-fade">
        <div v-if="toast" class="toast" :class="'toast-' + toast.kind" role="status">
            <span class="toast-icon">{{ toast.kind === 'success' ? '✓' : '!' }}</span>{{ toast.text }}
        </div>
    </transition>

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
                    <span class="muted small">{{ rangeText }}</span>
                    <span v-if="!hasPK" class="pk-warning">no primary key — rows are view-only for this table</span>
                    <button
                        v-if="selectedCount > 0"
                        class="btn-danger-solid btn-tight"
                        :disabled="deleting"
                        @click="startBulkDelete">Delete selected ({{ selectedCount }})</button>
                </div>
                <div class="table-scroll">
                    <table>
                        <thead>
                            <tr>
                                <th class="check-col" @click.stop>
                                    <input
                                        type="checkbox"
                                        :checked="allOnPageSelected"
                                        :disabled="!hasPK || selectableOnPage === 0"
                                        @change="toggleSelectAll"
                                        title="Select all rows on this page" />
                                </th>
                                <th class="actions-col"></th>
                                <th
                                    v-for="c in tableData.columns" :key="c"
                                    class="sortable"
                                    :class="{ sorted: sortBy === c }"
                                    :title="'Sort by ' + c"
                                    @click="toggleSort(c)">
                                    {{ c }}<span v-if="isPkColumn(c)" class="pk-mark" title="primary key"> ·PK</span><button
                                        v-if="sortBy === c"
                                        class="sort-dir"
                                        :title="sortDir === 'asc' ? 'Reverse sort' : 'Sort ascending'"
                                        @click.stop="reverseSort">{{ sortDir === 'asc' ? '▲' : '▼' }}</button>
                                </th>
                            </tr>
                        </thead>
                        <tbody>
                            <tr v-for="(row, ri) in tableData.rows" :key="ri"
                                :class="{ expanded: expandedRow === ri }"
                                @click="toggleRowExpand(ri)">
                                <td class="check-col" @click.stop>
                                    <input
                                        type="checkbox"
                                        :checked="isSelected(row)"
                                        :disabled="!isSelectableRow(row)"
                                        @change="toggleSelect(row)" />
                                </td>
                                <td class="actions-col" @click.stop>
                                    <button class="btn-mini" :disabled="!hasPK" :title="hasPK ? 'Edit row' : 'Editing requires a primary key'" @click="startEdit(row)">Edit</button>
                                    <button class="btn-mini btn-danger" :disabled="!hasPK" :title="hasPK ? 'Delete row' : 'Deleting requires a primary key'" @click="startDelete(row)">Del</button>
                                </td>
                                <td v-for="(val, ci) in row" :key="ci" class="cell" :title="expandedRow === ri ? '' : 'Click row to expand'">
                                    <span v-if="isNull(val)" class="null-val">NULL</span>
                                    <template v-else>{{ formatCell(val) }}</template>
                                </td>
                            </tr>
                            <tr v-if="tableData.rows.length === 0">
                                <td :colspan="tableData.columns.length + 2" class="empty center">No rows.</td>
                            </tr>
                        </tbody>
                    </table>
                </div>
                <div v-if="pageCount > 1" class="pager">
                    <button class="btn-ghost" :disabled="page <= 1" @click="gotoPage(1)">«</button>
                    <button class="btn-ghost" :disabled="page <= 1" @click="gotoPage(page - 1)">‹ Prev</button>
                    <span class="muted small">page {{ page }} of {{ pageCount }}</span>
                    <button class="btn-ghost" :disabled="page >= pageCount" @click="gotoPage(page + 1)">Next ›</button>
                    <button class="btn-ghost" :disabled="page >= pageCount" @click="gotoPage(pageCount)">»</button>
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
                <div v-for="f in formFields" :key="f.name" class="field" :class="{ 'field-invalid': fieldErrors[f.name] }">
                    <label :for="'f-' + f.name">
                        {{ f.name }}
                        <span class="type">({{ f.columnType || f.type }}{{ f.isPk ? ', PK' : '' }}{{ f.autoIncrement ? ', auto-inc' : '' }}{{ f.nullable ? ', NULL ok' : '' }})</span>
                    </label>
                    <input
                        v-if="!isLongText(f)"
                        :id="'f-' + f.name"
                        :type="inputType()"
                        :placeholder="inputPlaceholder(f)"
                        v-model="f.value"
                        :readonly="(formMode === 'edit' && f.isPk) || (formMode === 'add' && f.autoIncrement)" />
                    <textarea
                        v-else
                        :id="'f-' + f.name"
                        rows="4"
                        v-model="f.value"></textarea>
                    <p v-if="fieldErrors[f.name]" class="field-error">{{ fieldErrors[f.name] }}</p>
                </div>
                <p v-if="formFields.length === 0" class="muted">Nothing to edit — all columns are auto-managed.</p>
            </div>

            <div class="form-actions">
                <button class="btn-primary" :disabled="saving" @click="saveRow">{{ saving ? 'Saving…' : (formMode === 'add' ? 'Insert row' : 'Save changes') }}</button>
                <button class="btn-ghost" :disabled="saving" @click="cancelForm">Cancel</button>
            </div>
        </div>
    </Modal>

    <!-- delete confirmation modal (compact: it only names the target rows) -->
    <Modal v-if="deleteTarget" size="compact">
        <div class="row-form">
            <div class="form-head">
                <h3>Delete row from <span class="table-title">{{ selectedTable }}</span>?</h3>
            </div>
            <p class="muted small">This permanently removes the row identified by:</p>
            <ul class="pk-list">
                <li v-for="[name, val] in deleteTarget.pkValues" :key="name">
                    <strong>{{ name }}</strong> = <span class="pk-val">{{ isNull(val) ? 'NULL' : val }}</span>
                </li>
            </ul>
            <div v-if="deleteError" class="error-msg">{{ deleteError }}</div>
            <div class="form-actions">
                <button class="btn-danger-solid" :disabled="deleting" @click="confirmDelete">{{ deleting ? 'Deleting…' : 'Delete row' }}</button>
                <button class="btn-ghost" :disabled="deleting" @click="cancelDelete">Cancel</button>
            </div>
        </div>
    </Modal>
    <!-- bulk delete confirmation modal (compact: a row-list prompt should
         not take over the whole screen) -->
    <Modal v-if="bulkDeleteTargets" size="compact">
        <div class="row-form">
            <div class="form-head">
                <h3>Delete {{ bulkDeleteTargets.length }} row(s) from <span class="table-title">{{ selectedTable }}</span>?</h3>
            </div>
            <p class="muted small">This permanently removes every row listed below. This cannot be undone.</p>
            <ul class="pk-list pk-list-bulk">
                <li v-for="(pkValues, i) in bulkDeleteTargets" :key="i">{{
                    pkValues.map(([name, val]) => `${name}=${isNull(val) ? 'NULL' : val}`).join(', ')
                }}</li>
            </ul>
            <div v-if="deleteError" class="error-msg">{{ deleteError }}</div>
            <div class="form-actions">
                <button class="btn-danger-solid" :disabled="deleting" @click="confirmBulkDelete">{{ deleting ? 'Deleting…' : `Delete ${bulkDeleteTargets.length} row(s)` }}</button>
                <button class="btn-ghost" :disabled="deleting" @click="cancelBulkDelete">Cancel</button>
            </div>
        </div>
    </Modal>
</div>
</template>

<script setup>
import { ref, computed, onMounted, watch } from 'vue'
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
const toast = ref(null)               // { text, kind } — transient confirmation banner
let toastTimer = null

function showToast(text, kind = 'success') {
    if (toastTimer) clearTimeout(toastTimer)
    toast.value = { text, kind }
    toastTimer = setTimeout(() => { toast.value = null; toastTimer = null }, 4000)
}

const formMode = ref(null)          // null | 'add' | 'edit'
const formFields = ref([])
const formError = ref('')
const fieldErrors = ref({})           // { [columnName]: 'message' } — per-field validation
const editingPK = ref({})

const expandedRow = ref(null)       // row index whose cells are expanded
const deleting = ref(false)
const deleteTarget = ref(null)      // { pkValues: [[name, value], ...] } | null
const deleteError = ref('')

// ---------- pagination ----------
// Rows are fetched one page at a time (server LIMIT/OFFSET) so tables larger
// than the Rows limit stay reachable. `page` is 1-based; offset is derived.
const page = ref(1)

const totalRows = computed(() => (tableData.value ? Number(tableData.value.totalRows ?? tableData.value.rows.length) : 0))
const pageCount = computed(() => {
    const size = Number(tableData.value?.limit) || rowLimit.value || 100
    return Math.max(1, Math.ceil(totalRows.value / size))
})
const rangeText = computed(() => {
    if (!tableData.value) return ''
    const total = totalRows.value
    const shown = tableData.value.rows.length
    if (total === 0) return '0 rows'
    const first = (page.value - 1) * (Number(tableData.value.limit) || shown) + 1
    const last = first + shown - 1
    return `rows ${first}–${last} of ${total}`
})

// ---------- column sorting ----------
// Sorting is server-side (ORDER BY over the WHOLE table, not just this page);
// the server validates the column against information_schema before use.
const sortColumn = ref('')          // '' = default order (server uses PK ASC)
const sortDir = ref('asc')          // 'asc' | 'desc'
const sortBy = computed(() => sortColumn.value)

function toggleSort(col) {
    if (sortColumn.value === col) {
        // clicking the already-sorted column reverses it (matches the
        // reverse-sort button; both are the user's "flip direction" gesture)
        reverseSort()
        return
    }
    sortColumn.value = col
    sortDir.value = 'asc'
    page.value = 1
    loadTableData()
}

function reverseSort() {
    if (!sortColumn.value) return
    sortDir.value = sortDir.value === 'asc' ? 'desc' : 'asc'
    page.value = 1
    loadTableData()
}

function gotoPage(p) {
    const target = Math.min(Math.max(1, p), pageCount.value)
    if (target === page.value) return
    page.value = target
    loadTableData()
}

function applyLimit() {
    page.value = 1
    loadTableData()
}

// ---------- row selection ----------
// Selection is keyed by the row's primary-key values (serialized) so it
// survives re-renders and stays precise across pages. Rows from other pages
// stay selected (they carry their PK), but "select all" only ever touches the
// rows currently displayed. Tables without a PK cannot be selected: there is
// no way to address their rows for deletion.
const selectedKeys = ref(new Set())

function pkKeyOf(row) {
    if (!tableData.value?.primaryKey?.length) return null
    return JSON.stringify(tableData.value.primaryKey.map((pk) => {
        const idx = tableData.value.columnInfo.findIndex((c) => c.name === pk)
        return idx >= 0 ? row[idx] : null
    }))
}

function isSelectableRow(row) {
    if (!hasPK.value) return false
    const key = pkKeyOf(row)
    if (key === null) return false
    return JSON.parse(key).every((v) => !isNull(v))   // NULL PK can't be targeted
}

function isSelected(row) {
    const key = pkKeyOf(row)
    return key !== null && selectedKeys.value.has(key)
}

const selectedCount = computed(() => selectedKeys.value.size)

const selectableOnPage = computed(() =>
    (tableData.value?.rows || []).filter(isSelectableRow).length)

const allOnPageSelected = computed(() =>
    selectableOnPage.value > 0 &&
    (tableData.value?.rows || []).filter(isSelectableRow).every(isSelected))

function toggleSelect(row) {
    const key = pkKeyOf(row)
    if (key === null) return
    const next = new Set(selectedKeys.value)
    if (next.has(key)) next.delete(key)
    else next.add(key)
    selectedKeys.value = next
}

// Deliberately page-scoped: only the rows currently displayed are selected,
// never the whole table — a 10M-row "select all" would be a loaded gun.
function toggleSelectAll() {
    const next = new Set(selectedKeys.value)
    const target = !allOnPageSelected.value
    for (const row of tableData.value?.rows || []) {
        if (!isSelectableRow(row)) continue
        const key = pkKeyOf(row)
        if (key === null) continue
        if (target) next.add(key)
        else next.delete(key)
    }
    selectedKeys.value = next
}

function clearSelection() {
    selectedKeys.value = new Set()
}

// ---------- bulk delete ----------
const bulkDeleteTargets = ref(null) // [[ [pkName, pkValue], ... ], ...] | null

function startBulkDelete() {
    if (!tableData.value || selectedCount.value === 0) return
    deleteError.value = ''
    const pks = tableData.value.primaryKey || []
    const targets = []
    let skipped = 0
    for (const keyJson of selectedKeys.value) {
        const raw = JSON.parse(keyJson)
        if (!pks.length || raw.length !== pks.length || raw.some((v) => isNull(v))) {
            skipped += 1
            continue
        }
        targets.push(pks.map((pk, i) => [pk, raw[i]]))
    }
    if (targets.length === 0) {
        // nothing is addressable — surface it where the user clicked rather
        // than opening an empty confirmation sheet
        error.value = 'None of the selected rows can be deleted safely (each row needs a complete, non-NULL primary key).'
        return
    }
    if (skipped > 0) {
        deleteError.value = `${skipped} selected row(s) lack a complete primary key and will be skipped.`
    }
    bulkDeleteTargets.value = targets
}

function cancelBulkDelete() {
    bulkDeleteTargets.value = null
    deleteError.value = ''
}

function pkKeyOfPkValues(pkValues) {
    return JSON.stringify(pkValues.map(([, v]) => v))
}

async function confirmBulkDelete() {
    if (deleting.value || !bulkDeleteTargets.value) return
    deleting.value = true
    deleteError.value = ''
    const url = `${apiBase.value}/tables/${encodeURIComponent(selectedTable.value)}/rows`
    const failed = []
    let deleted = 0
    try {
        // One request per row: the DELETE endpoint addresses a single row by
        // its full PK, and doing them one at a time means a FK-blocked or
        // already-deleted row cannot roll back or mask the rest. Results are
        // summarized in one toast.
        for (const pkValues of bulkDeleteTargets.value) {
            const body = new FormData()
            for (const [name, val] of pkValues) body.append('pk_' + name, val)
            try {
                const res = await axios.delete(url, { data: body })
                if (res.data.error) failed.push(pkValues.map(([n, v]) => `${n}=${v}`).join(', '))
                else deleted += 1
            } catch (err) {
                failed.push(pkValues.map(([n, v]) => `${n}=${v}`).join(', '))
            }
        }
        if (deleted > 0) {
            showToast(`Deleted ${deleted} row(s) from ${selectedTable.value}.`)
        }
        if (failed.length > 0) {
            deleteError.value = `${failed.length} row(s) could not be deleted: ${failed.join('; ')}`
            // keep the modal open with only the failures selected for retry
            const keep = new Set()
            for (const pkValues of bulkDeleteTargets.value) {
                const desc = pkValues.map(([n, v]) => `${n}=${v}`).join(', ')
                if (failed.includes(desc)) keep.add(pkKeyOfPkValues(pkValues))
            }
            selectedKeys.value = keep
        } else {
            cancelBulkDelete()
            clearSelection()
        }
        await loadTableData()
    } finally {
        deleting.value = false
    }
}

function toggleRowExpand(ri) {
    expandedRow.value = expandedRow.value === ri ? null : ri
}

const apiBase = computed(() => `/api/accounts/${encodeURIComponent(props.account.name)}/databases/${encodeURIComponent(selectedDb.value)}`)

const hasPK = computed(() => !!(tableData.value && tableData.value.primaryKey && tableData.value.primaryKey.length))

function isPkColumn(name) {
    return !!(tableData.value && tableData.value.primaryKey?.includes(name))
}

function onDbChange() {
    selectedTable.value = ''
    tables.value = []
    tableData.value = null
    error.value = ''
    page.value = 1
    sortColumn.value = ''
    sortDir.value = 'asc'
    cancelForm()
    cancelBulkDelete()
    clearSelection()
    loadTables()
}

// load on mount too: the select is pre-filled with the first database, so a
// @change handler would never fire until the user manually switches dbs.
watch(selectedDb, onDbChange, { immediate: true })

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
    expandedRow.value = null
    page.value = 1
    // Sort is per-table state: a sort column from another table would 400.
    sortColumn.value = ''
    sortDir.value = 'asc'
    cancelForm()
    cancelBulkDelete()
    clearSelection()
    loadTableData()
}

async function loadTableData() {
    if (!selectedTable.value) return
    loading.value = true
    error.value = ''
    expandedRow.value = null
    try {
        const params = {
            limit: rowLimit.value || 100,
            offset: (page.value - 1) * (rowLimit.value || 100),
        }
        if (sortColumn.value) {
            params.sort = sortColumn.value
            params.dir = sortDir.value
        }
        const res = await axios.get(`${apiBase.value}/tables/${encodeURIComponent(selectedTable.value)}`, { params })
        if (res.data.error) { error.value = res.data.error; return }
        tableData.value = res.data
        // Default sort is the first PK column ASC — which is exactly the
        // server's default order, so adopting it here costs no extra request
        // and gives the header the sorted highlight + reverse button out of
        // the box (previously a hard-coded highlight on the PK header).
        if (!sortColumn.value && tableData.value.primaryKey?.length) {
            sortColumn.value = tableData.value.primaryKey[0]
            sortDir.value = 'asc'
        }
        // A delete (here or elsewhere) can empty the current page; clamp back
        // into range instead of stranding the user on "page 7 of 6".
        if (tableData.value.rows.length === 0 && page.value > 1 && totalRows.value > 0 && page.value > pageCount.value) {
            page.value = pageCount.value
            loading.value = false
            return loadTableData()
        }
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

const longTypes = ['text', 'tinytext', 'mediumtext', 'longtext', 'json', 'blob', 'tinyblob', 'mediumblob', 'longblob']

// Every column renders as type="text" and is validated by validateValue()
// against the declared column type. Native type="number"/"date"/"time"/
// "datetime-local" are deliberately NOT used: on an invalid entry the browser
// reports an empty string, so the server sees "" and maps it to NULL for
// nullable columns — the user's typo silently wipes the value with no error
// (verified: "25:99" in a TIME column and "abc" in an INT column both saved
// NULL). type="text" keeps what was typed so it can be rejected with a real
// message, and the placeholder shows the expected format.
function inputType() {
    return 'text'
}

function inputPlaceholder(col) {
    const t = (col.type || '').toLowerCase()
    if (t === 'date') return 'YYYY-MM-DD'
    if (t === 'time') return 'HH:MM:SS'
    if (t === 'datetime' || t === 'timestamp') return 'YYYY-MM-DD HH:MM:SS'
    if (t === 'enum') {
        const allowed = [...(col.columnType || '').matchAll(/'([^']*)'/g)].map((m) => m[1])
        if (allowed.length) return allowed.join(' | ')
    }
    if (t === 'set') return 'comma-separated: x, y, z'
    return ''
}

function isLongText(col) {
    return longTypes.includes((col.type || '').toLowerCase())
}

// The server returns DATE/DATETIME/TIME as MySQL's own canonical text (see
// openDBForAccount: parseTime is off), so these are treated as plain strings.
// Never new Date() them: that applies timezone math to values that have no
// timezone and silently shifts a date on save.
function toInputValue(col, val) {
    if (isNull(val)) return ''
    const t = (col.type || '').toLowerCase()
    if (t === 'datetime' || t === 'timestamp') {
        // MySQL canonical "YYYY-MM-DD HH:MM:SS[.ffffff]"; drop any timezone
        // suffix the driver may have appended and keep the value as written.
        return String(val).trim().slice(0, 26)
    }
    if (t === 'date') {
        return String(val).trim().slice(0, 10)
    }
    if (t === 'time') {
        // keep HH:MM:SS (or fractional), dropping any trailing junk
        const m = String(val).trim().match(/^([+-]?\d{1,3}:\d{2}(?::\d{2}(?:\.\d+)?)?)/)
        return m ? m[1] : String(val).trim()
    }
    return String(val)
}

// ---------- client-side validation ----------
//
// The server is the authority (it re-checks every identifier against
// information_schema and binds all values), but MySQL in non-strict mode
// silently coerces junk — "abc" into an INT becomes 0, and the browser's
// type="number" quietly turns an invalid entry into "" which the server then
// maps to NULL. That wipes data with no feedback, so we reject obvious type
// errors here using the declared column type from information_schema, while
// still letting MySQL be the final say on anything subtle.

function parseColumnType(col) {
    const raw = (col.columnType || col.type || '').toLowerCase()
    const base = raw.split('(')[0].trim()
    const unsigned = raw.includes('unsigned')
    const out = { base, unsigned, raw }
    const lenMatch = raw.match(/\((\d+)\)/)
    const decMatch = raw.match(/\((\d+)\s*,\s*(\d+)\)/)
    if (decMatch) { out.precision = parseInt(decMatch[1], 10); out.scale = parseInt(decMatch[2], 10) }
    else if (lenMatch) { out.length = parseInt(lenMatch[1], 10) }
    if (base === 'enum' || base === 'set') {
        out.allowed = [...raw.matchAll(/'([^']*)'/g)].map((m) => m[1])
    }
    return out
}

const intRanges = {
    tinyint:   { min: -128,        max: 127 },
    smallint:  { min: -32768,      max: 32767 },
    mediumint: { min: -8388608,    max: 8388607 },
    int:       { min: -2147483648, max: 2147483647 },
    integer:   { min: -2147483648, max: 2147483647 },
    bigint:    { min: -9223372036854775808n, max: 9223372036854775807n },
    serial:    { min: 0,           max: 4294967295 },
}

const EMPTY_OK = Symbol('empty-ok')

function validateValue(col, value) {
    const raw = value === null || value === undefined ? '' : String(value)
    const trimmed = raw.trim()

    if (trimmed === '') {
        // empty means NULL only where the column allows it; otherwise MySQL
        // would silently insert 0 / '' and destroy the previous value.
        if (col.nullable) return EMPTY_OK
        if (col.autoIncrement) return EMPTY_OK
        return 'This column does not allow NULL — enter a value.'
    }

    const t = parseColumnType(col)
    const isInt = t.base in intRanges || t.base === 'bit'

    if (isInt) {
        if (!/^[+-]?\d+$/.test(trimmed)) return `Must be a whole number (${t.raw}) — "${raw}" is not valid.`
        if (t.base === 'bit') {
            const bits = t.length || 1
            const max = 2 ** bits - 1
            if (+trimmed < 0 || +trimmed > max) return `bit(${bits}) accepts 0 to ${max}.`
            return null
        }
        if (t.base === 'year') {
            const y = +trimmed
            if (!(y === 0 || (y >= 1901 && y <= 2155))) return 'YEAR accepts 1901–2155 (or 0).'
            return null
        }
        const r = intRanges[t.base]
        if (t.unsigned && trimmed.startsWith('-')) return `${t.raw} cannot be negative (unsigned).`
        // BigInt throughout: Number() would lose precision near the bigint bounds.
        // Unsigned max is 2*signedMax+1 (e.g. tinyint: 127 -> 255), NOT the signed
        // max — using the signed bound would block legitimate high values.
        let big, min, max
        try {
            big = BigInt(trimmed)
            const signedMax = BigInt(r.max)
            min = BigInt(t.unsigned ? 0 : r.min)
            max = t.unsigned ? signedMax * 2n + 1n : signedMax
        } catch (e) { return `Must be a whole number (${t.raw}).` }
        if (big < min || big > max) return `Out of range for ${t.raw} (${min} to ${max}).`
        return null
    }

    if (t.base === 'decimal' || t.base === 'numeric') {
        if (!/^[+-]?(\d+(\.\d*)?|\.\d+)([eE][+-]?\d+)?$/.test(trimmed)) return `Must be a number (${t.raw}) — "${raw}" is not valid.`
        if (t.scale !== undefined) {
            const frac = (trimmed.split('.')[1] || '').length
            if (frac > t.scale) return `${t.raw} allows at most ${t.scale} decimal place(s); you entered ${frac}.`
        }
        if (t.precision !== undefined) {
            const digits = trimmed.replace(/^[+-]?0+/, '').replace('.', '').replace(/e.*$/i, '').length
            if (digits > t.precision) return `${t.raw} allows at most ${t.precision} total digits.`
        }
        return null
    }

    if (t.base === 'float' || t.base === 'double' || t.base === 'real') {
        if (!/^[+-]?(\d+(\.\d*)?|\.\d+)([eE][+-]?\d+)?$/.test(trimmed) || !isFinite(Number(trimmed))) {
            return `Must be a number (${t.raw}) — "${raw}" is not valid.`
        }
        return null
    }

    if (t.base === 'char' || t.base === 'varchar') {
        if (t.length !== undefined && raw.length > t.length) return `Too long: ${t.raw} holds ${t.length} characters, you entered ${raw.length}.`
        return null
    }

    if (t.base === 'binary' || t.base === 'varbinary') {
        if (t.length !== undefined && raw.length > t.length) return `Too long: ${t.raw} holds ${t.length} bytes.`
        return null
    }

    if (t.base === 'enum') {
        if (t.allowed && !t.allowed.includes(raw)) return `Must be one of: ${t.allowed.join(', ')}.`
        return null
    }

    if (t.base === 'set') {
        if (t.allowed) {
            const parts = raw.split(',').map((s) => s.trim()).filter((s) => s !== '')
            const bad = parts.filter((p) => !t.allowed.includes(p))
            if (bad.length) return `Not allowed in set: ${bad.join(', ')} (options: ${t.allowed.join(', ')}).`
        }
        return null
    }

    if (t.base === 'date') {
        if (!/^\d{4}-\d{2}-\d{2}$/.test(trimmed) || isNaN(Date.parse(trimmed))) return 'Use a valid date, YYYY-MM-DD.'
        return null
    }

    if (t.base === 'time') {
        // MySQL TIME spans -838:59:59 to 838:59:59, so hours may exceed 24 but
        // minutes and seconds must still be 0-59 (digit-count alone let "25:99" through).
        const m = trimmed.match(/^([+-])?(\d{1,3}):(\d{1,2})(:(\d{1,2})(\.\d+)?)?$/)
        if (!m) return 'Use a valid time, HH:MM or HH:MM:SS.'
        const [, , hh, mi, , ss] = m
        if (+mi > 59) return `TIME minutes must be 0-59 (got ${mi}).`
        if (ss !== undefined && +ss > 59) return `TIME seconds must be 0-59 (got ${ss}).`
        if (+hh > 838) return 'TIME hours must be 0-838.'
        return null
    }

    if (t.base === 'datetime' || t.base === 'timestamp') {
        if (!/^\d{4}-\d{2}-\d{2}([T ]\d{2}:\d{2}(:\d{2}(\.\d+)?)?)?$/.test(trimmed) || isNaN(Date.parse(trimmed))) {
            return 'Use a valid date/time, YYYY-MM-DD HH:MM:SS.'
        }
        return null
    }

    if (t.base === 'json') {
        try { JSON.parse(trimmed) } catch (e) { return 'Not valid JSON.' }
        return null
    }

    // text/blob and anything unusual: MySQL decides.
    return null
}

function validateFields() {
    const errs = {}
    const editing = formMode.value === 'edit'
    for (const f of formFields.value) {
        // server-managed columns are not validated
        if (formMode.value === 'add' && f.autoIncrement) continue
        if (editing && f.isPk) continue
        // On update, only the fields that will actually be sent are validated.
        // An untouched column is omitted from the SET clause entirely, so a
        // stored value the UI cannot represent (a legacy '0000-00-00' date, a
        // binary blob) must not block an unrelated edit to some other column.
        if (editing && String(f.value ?? '') === String(f.original ?? '')) continue
        const msg = validateValue(f, f.value)
        if (msg && msg !== EMPTY_OK) errs[f.name] = msg
    }
    fieldErrors.value = errs
    return Object.keys(errs).length === 0
}


function buildFormFields(rowValues) {
    if (!tableData.value) return []
    return tableData.value.columnInfo.map((col, i) => ({
        ...col,
        value: rowValues ? toInputValue(col, rowValues[i]) : '',
        // Baseline the row was opened with, so saveRow() can tell which fields
        // the user actually changed and send ONLY those.
        original: rowValues ? toInputValue(col, rowValues[i]) : '',
    }))
}

function startEdit(row) {
    if (!hasPK.value) {
        error.value = 'This table has no primary key, so rows can\u2019t be edited safely.'
        return
    }
    formError.value = ''
    fieldErrors.value = {}
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
    fieldErrors.value = {}
    formMode.value = 'add'
    formFields.value = buildFormFields(null)
    editingPK.value = {}
}

function cancelForm() {
    formMode.value = null
    formFields.value = []
    formError.value = ''
    fieldErrors.value = {}
    editingPK.value = {}
}

async function saveRow() {
    if (saving.value || !tableData.value) return

    // Reject obvious type violations before anything is sent, so a typo can
    // never silently wipe a value (see validateValue for why this matters).
    if (!validateFields()) {
        formError.value = 'Fix the highlighted field(s) before saving — nothing was changed.'
        return
    }

    saving.value = true
    formError.value = ''
    error.value = ''

    // Build the multipart body the server expects: col_<name> (and pk_<name>
    // on updates). Server whitelists column names against information_schema
    // and binds every value as a ? parameter.
    //
    // On UPDATE we send ONLY the fields the user actually changed. The server
    // builds its SET clause purely from the col_* fields that arrive, so any
    // column we omit is left completely untouched in the row — including
    // columns the grid may have rendered lossily (binary, long text) and
    // ON UPDATE CURRENT_TIMESTAMP / generated columns, which a blind
    // full-row rewrite would clobber or churn. startEdit baselines `original`
    // from the row it just loaded, so an unmodified field is provably equal to
    // what is stored and is safe to omit. INSERT still sends every field.
    const body = new FormData()
    const changed = []
    for (const f of formFields.value) {
        if (formMode.value === 'add') {
            if (f.autoIncrement) continue           // server-assigned
            if (f.isPk && f.value === '') continue  // let MySQL default/assign the PK
            body.append('col_' + f.name, f.value ?? '')
        } else {
            if (f.isPk) continue                    // PK identifies the row; not editable
            if (String(f.value ?? '') === String(f.original ?? '')) continue
            body.append('col_' + f.name, f.value ?? '')
            changed.push(f.name)
        }
    }
    if (formMode.value === 'edit') {
        if (changed.length === 0) {
            saving.value = false
            formError.value = 'Nothing to save — no fields were changed.'
            return
        }
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
        const verb = formMode.value === 'add' ? 'Row added to' : 'Row updated in'
        showToast(`${verb} ${selectedTable.value} — ${res.data.message || 'saved'}.`)
        cancelForm()
        await loadTableData()
    } catch (err) {
        formError.value = apiError(err, 'Save failed — values are validated server-side, check the column types.')
    } finally {
        saving.value = false
    }
}

// ---------- row delete ----------

function startDelete(row) {
    if (!hasPK.value) {
        error.value = 'This table has no primary key, so rows can\u2019t be deleted safely.'
        return
    }
    deleteError.value = ''
    const pkValues = []
    for (const pk of tableData.value.primaryKey) {
        const idx = tableData.value.columnInfo.findIndex((c) => c.name === pk)
        pkValues.push([pk, idx >= 0 ? row[idx] : null])
    }
    deleteTarget.value = { pkValues }
}

function cancelDelete() {
    deleteTarget.value = null
    deleteError.value = ''
}

async function confirmDelete() {
    if (deleting.value || !deleteTarget.value) return
    deleting.value = true
    deleteError.value = ''
    error.value = ''

    const body = new FormData()
    for (const [name, val] of deleteTarget.value.pkValues) {
        if (isNull(val)) {
            deleteError.value = 'Rows with a NULL primary key can\u2019t be targeted safely.'
            deleting.value = false
            return
        }
        body.append('pk_' + name, val)
    }

    try {
        const url = `${apiBase.value}/tables/${encodeURIComponent(selectedTable.value)}/rows`
        const res = await axios.delete(url, { data: body })
        if (res.data.error) { deleteError.value = res.data.error; return }
        const idDesc = deleteTarget.value.pkValues.map(([n, v]) => `${n}=${isNull(v) ? 'NULL' : v}`).join(', ')
        showToast(`Deleted row from ${selectedTable.value} (${idDesc}).`)
        // drop it from any selection so a re-inserted row with the same PK
        // does not appear pre-selected
        const key = JSON.stringify(deleteTarget.value.pkValues.map(([, v]) => v))
        if (selectedKeys.value.has(key)) {
            const next = new Set(selectedKeys.value)
            next.delete(key)
            selectedKeys.value = next
        }
        cancelDelete()
        await loadTableData()
    } catch (err) {
        deleteError.value = apiError(err, 'Delete failed — the row is targeted by its primary key, check the values.')
    } finally {
        deleting.value = false
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

/* ---- confirmation toast ---- */
.toast {
    position: fixed;
    right: 24px;
    bottom: 24px;
    z-index: 1200;
    display: flex;
    align-items: center;
    gap: 10px;
    max-width: 420px;
    padding: 12px 16px;
    border-radius: 10px;
    font-size: 13px;
    font-weight: 500;
    box-shadow: 0 8px 28px rgba(0, 0, 0, 0.35);
    border: 1px solid transparent;

    .toast-icon {
        display: inline-flex;
        align-items: center;
        justify-content: center;
        width: 20px;
        height: 20px;
        border-radius: 50%;
        font-size: 12px;
        font-weight: 700;
        flex-shrink: 0;
        color: #fff;
    }

    &.toast-success {
        background: $c-od-surface-hi;
        border-color: $c-od-success;
        color: $c-od-fg;
        .toast-icon { background: $c-od-success; }
    }

    &.toast-error {
        background: $c-od-surface-hi;
        border-color: $c-od-danger;
        color: $c-od-fg;
        .toast-icon { background: $c-od-danger; }
    }
}

.toast-fade-enter-active, .toast-fade-leave-active { transition: opacity 0.25s ease, transform 0.25s ease; }
.toast-fade-enter-from, .toast-fade-leave-to { opacity: 0; transform: translateY(8px); }
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

            .pk-mark { font-size: 10px; color: $c-od-accent; }
        }

        th.sortable {
            cursor: pointer;
            user-select: none;

            &:hover { color: $c-od-fg; }

            // the column the rows are currently sorted by keeps the accent
            // highlight (the same treatment the PK header used to hard-code)
            &.sorted {
                color: $c-od-accent;
                box-shadow: inset 0 -2px 0 var(--c-od-accent);
            }

            .sort-dir {
                margin-left: 4px;
                padding: 0 3px;
                border: none;
                background: transparent;
                color: $c-od-accent;
                font-size: 10px;
                line-height: 1.4;
                cursor: pointer;
                border-radius: 4px;

                &:hover { background: rgba(255, 255, 255, 0.08); }
            }
        }

        td.cell {
            max-width: 340px;
            overflow: hidden;
            text-overflow: ellipsis;
            background: $c-od-surface;
            color: $c-od-fg;
            cursor: pointer;
        }

        /* cells expand only after the row is clicked (no hover-expand) */
        tbody tr.expanded td.cell {
            white-space: normal;
            overflow: visible;
            max-width: none;
            position: relative;
            z-index: 2;
            box-shadow: 0 0 0 1px $c-od-accent;
        }

        .null-val {
            color: $c-od-muted;
            font-style: italic;
            font-size: 11px;
        }

        th.actions-col, td.actions-col,
        th.check-col, td.check-col {
            position: sticky;
            left: 0;
            background: $c-od-surface-hi;
            z-index: 2;
            width: 1%;
        }

        td.check-col { background: $c-od-surface; }

        th.check-col + th.actions-col { left: 34px; }

        td.check-col input[type='checkbox'],
        th.check-col input[type='checkbox'] {
            accent-color: var(--c-od-accent);
            cursor: pointer;
            margin: 0;

            &:disabled { cursor: not-allowed; opacity: 0.4; }
        }

        tbody tr:hover td.cell { background: $c-od-surface-hi; }
    }
}

/* ---- pagination ---- */
.pager {
    display: flex;
    align-items: center;
    gap: 8px;
    justify-content: flex-end;
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

    &.btn-danger:hover:not(:disabled) { border-color: $c-od-danger; color: $c-od-danger; }
}

.btn-danger-solid {
    padding: 8px 16px;
    border-radius: 8px;
    background: $c-od-danger;
    color: #fff;
    font-size: 13px;
    font-weight: 600;
    border: none;
    cursor: pointer;
    white-space: nowrap;

    &:hover { opacity: 0.9; }
    &:disabled { opacity: 0.5; cursor: not-allowed; }

    &.btn-tight { padding: 4px 12px; font-size: 12px; margin-left: auto; }
}

.pk-list {
    list-style: none;
    padding: 0;
    margin: 0;
    display: flex;
    flex-direction: column;
    gap: 4px;
    font-size: 13px;

    .pk-val { color: $c-od-accent; font-family: monospace; }
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

            .field-error {
                margin: 0;
                font-size: 11px;
                color: $c-od-danger;
            }
        }

        .field-invalid {
            input, textarea {
                border-color: $c-od-danger;
                background: rgba(255, 90, 90, 0.06);
            }
        }
    }

    .form-actions { display: flex; gap: 10px; }
}
</style>
