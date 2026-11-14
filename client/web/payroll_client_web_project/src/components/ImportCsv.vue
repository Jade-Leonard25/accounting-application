<script setup lang="ts">
import { computed, onMounted, ref, watch } from 'vue'
import Papa from 'papaparse'
import { get, set, del } from 'idb-keyval'
import axios from 'axios'
import { z } from 'zod'
import Button from './Button.vue'
import Table from './Table.vue'

const emit = defineEmits<{ saved: [] }>()
const count = defineModel<number>('count', { default: 0 })

const DRAFT_KEY = 'payroll:csv-draft'

type DraftRow = { rowId: string } & Record<string, string>

const rows = ref<DraftRow[]>([])
const headers = ref<string[]>([])
const fileName = ref('')
const parseError = ref('')
const dragging = ref(false)
const saving = ref(false)
const progress = ref(0)
const serverErrors = ref<Record<string, string>>({})
const ready = ref(false)
const editMode = ref(false)

const schema = z.object({
  first_name: z.string().min(1, 'First name is required'),
  last_name: z.string().min(1, 'Last name is required'),
  email: z.email('Invalid email'),
  base_salary_cents: z
    .number('Salary must be a number')
    .int('Salary must be whole cents')
    .min(0, 'Salary cannot be negative'),
  tax_id: z.string().optional(),
  employment_status: z.enum(['Full-time', 'Part-time', 'Contract']),
  pay_frequency: z.enum(['Semi-monthly', 'Monthly', 'Weekly']),
})

const REQUIRED = ['first_name', 'last_name', 'email', 'base_salary_cents', 'employment_status', 'pay_frequency']

const API = axios.create({
  baseURL: 'http://localhost:4002/api',
  withCredentials: true,
})

// ---- Draft persistence (IndexedDB) ------------------------------------
onMounted(async () => {
  try {
    const draft = await get<{ fileName: string; headers?: string[]; rows: DraftRow[] }>(DRAFT_KEY)
    if (draft) {
      rows.value = draft.rows
      fileName.value = draft.fileName
      headers.value =
        draft.headers ?? Object.keys(draft.rows[0] ?? {}).filter((k) => k !== 'rowId')
    }
  } catch {
    /* IndexedDB unavailable, work in memory only */
  } finally {
    ready.value = true
  }
})

watch(
  rows,
  async (value) => {
    count.value = value.length
    if (!ready.value) return
    try {
      if (value.length === 0) await del(DRAFT_KEY)
      else
        await set(
          DRAFT_KEY,
          JSON.parse(JSON.stringify({ fileName: fileName.value, headers: headers.value, rows: value })),
        )
    } catch {
      /* ignore storage errors */
    }
  },
  { deep: true, immediate: true },
)

// ---- Parsing ----------------------------------------------------------
const toText = (v: unknown) =>
  v === null || v === undefined ? '' : typeof v === 'object' ? JSON.stringify(v) : String(v)

const load = (name: string, records: Record<string, unknown>[], fields: string[]) => {
  if (!records.length || !fields.length) {
    parseError.value = 'The file has no rows.'
    return
  }
  fileName.value = name
  serverErrors.value = {}
  editMode.value = false
  headers.value = fields
  rows.value = records.map((r, i) => ({
    rowId: `${Date.now()}-${i}`,
    ...Object.fromEntries(fields.map((f) => [f, toText(r[f]).trim()])),
  }))
}

const handleFile = async (file: File | undefined) => {
  parseError.value = ''
  if (!file) return
  const lower = file.name.toLowerCase()

  if (lower.endsWith('.json')) {
    try {
      let data = JSON.parse(await file.text())
      // accept { "employees": [...] } as well as a plain array
      if (!Array.isArray(data) && data && typeof data === 'object') {
        data = Object.values(data).find(Array.isArray) ?? [data]
      }
      if (!Array.isArray(data)) throw new Error('JSON must be an array of objects.')
      const records = data.filter((d: unknown) => d && typeof d === 'object') as Record<string, unknown>[]
      const fields = [...new Set(records.flatMap((r) => Object.keys(r)))]
      load(file.name, records, fields)
    } catch (e: any) {
      parseError.value = e?.message || 'Could not read this JSON file.'
    }
    return
  }

  if (!lower.endsWith('.csv')) {
    parseError.value = 'Please choose a .csv or .json file.'
    return
  }

  Papa.parse<Record<string, string>>(file, {
    header: true,
    skipEmptyLines: true,
    transformHeader: (h) => h.trim().toLowerCase(),
    complete: (result) => {
      const fields = (result.meta.fields ?? []).filter(Boolean)
      load(file.name, result.data, fields)
    },
    error: (e) => (parseError.value = e.message),
  })
}

const onPick = (e: Event) => {
  const input = e.target as HTMLInputElement
  handleFile(input.files?.[0])
  input.value = '' // allow re-selecting the same file
}
const onDrop = (e: DragEvent) => {
  dragging.value = false
  handleFile(e.dataTransfer?.files?.[0])
}

// ---- Validation -------------------------------------------------------
const isEmployeeData = computed(() => REQUIRED.every((h) => headers.value.includes(h)))

const toEmployee = (r: DraftRow) => ({
  first_name: r.first_name,
  last_name: r.last_name,
  email: r.email,
  base_salary_cents: Number(r.base_salary_cents),
  tax_id: r.tax_id ?? '',
  employment_status: r.employment_status,
  pay_frequency: r.pay_frequency,
})

const rowErrors = computed(() => {
  const map: Record<string, string> = {}
  if (!isEmployeeData.value) return map
  for (const r of rows.value) {
    const res = schema.safeParse(toEmployee(r))
    if (!res.success) map[r.rowId] = res.error.issues[0]?.message ?? 'Invalid row'
  }
  return map
})

const tableRows = computed(() =>
  rows.value.map((r) => ({
    ...r,
    _check: serverErrors.value[r.rowId] ?? rowErrors.value[r.rowId] ?? 'Ready',
    _ok: !serverErrors.value[r.rowId] && !rowErrors.value[r.rowId],
  })),
)

const validCount = computed(() => tableRows.value.filter((r) => r._ok).length)
const invalidCount = computed(() => rows.value.length - validCount.value)

const columns = computed(() => [
  ...headers.value.map((h) => ({ key: h, label: h, sortable: true })),
  ...(isEmployeeData.value ? [{ key: '_check', label: 'Check', sortable: true }] : []),
  ...(editMode.value ? [{ key: '_actions', label: '', hideInDrawer: true, class: 'text-right' }] : []),
])

// ---- Editing ----------------------------------------------------------
// Used by the drawer's Save button
const saveRow = (row: any, values: Record<string, any>) => {
  const i = rows.value.findIndex((r) => r.rowId === row.rowId)
  if (i < 0) return
  const next = { ...rows.value[i] }
  for (const [k, v] of Object.entries(values)) next[k] = String(v)
  rows.value[i] = next
  delete serverErrors.value[row.rowId]
}

// Used by the inline cell inputs
const cellSlot = (h: string) => `cell-${h}`

const onCellChange = (e: Event, rowId: string, key: string) => {
  const r = rows.value.find((x) => x.rowId === rowId)
  if (!r) return
  r[key] = (e.target as HTMLInputElement).value
  delete serverErrors.value[rowId]
}

const addRow = () => {
  rows.value.unshift({
    rowId: crypto.randomUUID(),
    ...Object.fromEntries(headers.value.map((h) => [h, ''])),
  })
  editMode.value = true
}

const deleteRow = (rowId: string) => {
  rows.value = rows.value.filter((r) => r.rowId !== rowId)
  delete serverErrors.value[rowId]
}

const exportCsv = () => {
  const csv = Papa.unparse({
    fields: headers.value,
    data: rows.value.map((r) => headers.value.map((h) => r[h] ?? '')),
  })
  const url = URL.createObjectURL(new Blob([csv], { type: 'text/csv;charset=utf-8' }))
  const a = document.createElement('a')
  a.href = url
  a.download = `${fileName.value.replace(/\.(csv|json)$/i, '') || 'export'}-edited.csv`
  a.click()
  URL.revokeObjectURL(url)
}

// ---- Save to backend --------------------------------------------------
const generateEmployeeId = () => {
  const block = () => Math.random().toString(36).slice(2, 5).padEnd(3, '0')
  return `${block()}-${block()}-${block()}`
}

const saveAll = async () => {
  saving.value = true
  progress.value = 0
  serverErrors.value = {}
  const saved = new Set<string>()
  const todo = rows.value.filter((r) => !rowErrors.value[r.rowId])

  for (const r of todo) {
    try {
      await API.post('/employees', { ...toEmployee(r), id: generateEmployeeId() })
      saved.add(r.rowId)
    } catch (e: any) {
      serverErrors.value[r.rowId] =
        e.response?.data?.error || e.response?.data?.message || 'Server rejected this row'
    }
    progress.value++
  }

  // Keep only rows that were not saved, so a retry never creates duplicates
  rows.value = rows.value.filter((r) => !saved.has(r.rowId))
  saving.value = false
  if (saved.size) emit('saved')
}

const discard = () => {
  rows.value = []
  headers.value = []
  fileName.value = ''
  serverErrors.value = {}
  parseError.value = ''
  editMode.value = false
}

const cellInputClass =
  'w-full min-w-24 rounded border border-slate-600 bg-slate-800 px-2 py-1 text-sm text-gray-200 focus:outline-none focus:ring-2 focus:ring-green-400'
</script>

<template>
  <div class="flex flex-col gap-6">
    <div>
      <h1 class="text-2xl font-bold text-white">Import Data</h1>
      <p class="mt-1 text-sm text-slate-400">
        Upload a CSV or JSON file, edit it, then save. Your draft is kept in this browser until you save or discard it.
      </p>
    </div>

    <label
      class="flex cursor-pointer flex-col items-center gap-2 rounded-lg border-2 border-dashed p-8 text-center transition-colors"
      :class="dragging ? 'border-green-400 bg-slate-800' : 'border-slate-600 hover:border-slate-500'"
      @dragover.prevent="dragging = true"
      @dragleave="dragging = false"
      @drop.prevent="onDrop"
    >
      <input type="file" accept=".csv,.json,text/csv,application/json" class="sr-only" @change="onPick" />
      <span class="font-medium text-white">Drop a CSV or JSON file here or click to browse</span>
      <span class="text-xs text-slate-400">
        Any columns work for preview. To save to the database, include: first_name, last_name, email,
        base_salary_cents, employment_status, pay_frequency (tax_id optional)
      </span>
    </label>

    <p v-if="parseError" class="text-sm text-red-400">{{ parseError }}</p>

    <template v-if="rows.length">
      <div class="flex flex-wrap items-center justify-between gap-3">
        <p class="text-sm text-slate-300">
          <span class="font-medium text-white">{{ fileName }}</span>
          · {{ rows.length }} rows · {{ headers.length }} columns
          <template v-if="isEmployeeData">
            · <span class="text-green-400">{{ validCount }} ready</span>
            <template v-if="invalidCount">
              · <span class="text-red-400">{{ invalidCount }} with problems</span>
            </template>
          </template>
        </p>
        <div class="flex flex-wrap gap-2">
          <Button :variant="editMode ? 'default' : 'secondary'" @click="editMode = !editMode">
            {{ editMode ? 'Done editing' : 'Edit cells' }}
          </Button>
          <Button variant="secondary" @click="addRow">+ Add row</Button>
          <Button variant="secondary" @click="exportCsv">Export CSV</Button>
          <Button variant="secondary" :disabled="saving" @click="discard">Discard draft</Button>
          <Button
            :loading="saving"
            :disabled="!isEmployeeData || validCount === 0"
            :title="isEmployeeData ? '' : 'Only employee-shaped files can be saved to the database'"
            @click="saveAll"
          >
            {{ saving ? `Saving ${progress}/${validCount + invalidCount}` : `Save ${validCount} to database` }}
          </Button>
        </div>
      </div>

      <Table
        :columns="columns"
        :rows="tableRows"
        row-key="rowId"
        :row-drawer="!editMode"
        editable
        :on-save="saveRow"
        search-placeholder="Search draft..."
      >
        <!-- One slot per column: inputs in edit mode, plain text otherwise -->
        <template v-for="h in headers" :key="h" #[cellSlot(h)]="{ row }">
          <input
            v-if="editMode"
            :value="row[h]"
            :class="cellInputClass"
            @change="onCellChange($event, row.rowId, h)"
          />
          <template v-else>{{ row[h] }}</template>
        </template>

        <template #cell-_check="{ row }">
          <span :class="row._ok ? 'text-green-400' : 'text-red-400'" class="text-xs">{{ row._check }}</span>
        </template>

        <template #cell-_actions="{ row }">
          <Button variant="destructive" size="sm" @click="deleteRow(row.rowId)">✕</Button>
        </template>
      </Table>
    </template>
  </div>
</template>