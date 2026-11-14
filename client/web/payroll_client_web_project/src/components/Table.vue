<script setup lang="ts" generic="T extends Record<string, any>">
import { computed, ref, watch } from 'vue'
import type { HTMLAttributes } from 'vue'
import { clsx, type ClassValue } from 'clsx'
import { twMerge } from 'tailwind-merge'
import { cva, type VariantProps } from 'class-variance-authority'
import Input from './Input.vue'
import Button from './Button.vue'
import Drawer from './Drawer.vue'

const cn = (...inputs: ClassValue[]) => twMerge(clsx(inputs))

const tableVariants = cva('overflow-hidden rounded-lg border text-gray-200', {
  variants: {
    variant: {
      default: 'border-slate-700 bg-slate-900',
      muted: 'border-slate-700 bg-slate-800/40',
    },
  },
  defaultVariants: { variant: 'muted' },
})

const cellVariants = cva('text-left', {
  variants: {
    density: {
      default: 'p-3',
      compact: 'px-3 py-1.5 text-sm',
    },
  },
  defaultVariants: { density: 'default' },
})

type TableVariants = VariantProps<typeof tableVariants>
type CellVariants = VariantProps<typeof cellVariants>

interface Column {
  key: string
  label: string
  sortable?: boolean
  drawerOnly?: boolean   // shown in the drawer, hidden in the table
  hideInDrawer?: boolean // shown in the table, hidden in the drawer
  editable?: boolean     // set false to make a field read-only (default true)
  type?: 'text' | 'number' | 'email'
  options?: string[]     // renders a <select> when editing
  class?: string
}

interface Props {
  columns: Column[]
  rows: T[]
  rowKey?: string
  searchable?: boolean
  searchKeys?: string[]
  searchPlaceholder?: string
  loading?: boolean
  emptyText?: string
  rowDrawer?: boolean
  editable?: boolean // adds the Edit button to the drawer
  onSave?: (row: T, values: Record<string, any>) => void | Promise<void>
  drawerSide?: 'right' | 'left' | 'top' | 'bottom'   // CHANGED: all four sides
  drawerSize?: 'sm' | 'md' | 'lg' | 'full'           // NEW
  variant?: TableVariants['variant']
  density?: CellVariants['density']
  class?: HTMLAttributes['class']
}

const props = withDefaults(defineProps<Props>(), {
  rowKey: 'id',
  searchable: true,
  searchPlaceholder: 'Search...',
  emptyText: 'No results.',
  rowDrawer: false,
  editable: false,
  drawerSide: 'bottom',   // CHANGED: was 'right'
  drawerSize: 'lg',       // NEW
})

const search = defineModel<string>('search', { default: '' })
const sortKey = ref<string | null>(null)
const sortDir = ref<'asc' | 'desc'>('asc')

const tableColumns = computed(() => props.columns.filter((c) => !c.drawerOnly))
const drawerColumns = computed(() => props.columns.filter((c) => !c.hideInDrawer))
const fields = computed(() => props.searchKeys ?? props.columns.map((c) => c.key))

const filtered = computed(() => {
  const q = search.value.trim().toLowerCase()
  if (!q) return props.rows
  return props.rows.filter((row) =>
    fields.value.some((f) => String(row[f] ?? '').toLowerCase().includes(q)),
  )
})

const visible = computed(() => {
  if (!sortKey.value) return filtered.value
  const key = sortKey.value
  const dir = sortDir.value === 'asc' ? 1 : -1
  return [...filtered.value].sort((a, b) => {
    const x = a[key]
    const y = b[key]
    if (typeof x === 'number' && typeof y === 'number') return (x - y) * dir
    return String(x ?? '').localeCompare(String(y ?? '')) * dir
  })
})

const toggleSort = (col: Column) => {
  if (!col.sortable) return
  if (sortKey.value === col.key) {
    sortDir.value = sortDir.value === 'asc' ? 'desc' : 'asc'
  } else {
    sortKey.value = col.key
    sortDir.value = 'asc'
  }
}

// ---- Row drawer --------------------------------------------------------
const selectedKey = ref<unknown>(null)

const selectedIndex = computed(() =>
  selectedKey.value === null
    ? -1
    : visible.value.findIndex((r) => r[props.rowKey] === selectedKey.value),
)
const selectedRow = computed(() => (selectedIndex.value >= 0 ? visible.value[selectedIndex.value] : null))

const drawerOpen = computed({
  get: () => selectedIndex.value >= 0,
  set: (v: boolean) => {
    if (!v) selectedKey.value = null
  },
})

const onRowClick = (e: MouseEvent, row: T) => {
  if (!props.rowDrawer) return
  if ((e.target as HTMLElement).closest('button, a, input, select, label')) return
  selectedKey.value = row[props.rowKey]
}

const step = (delta: number) => {
  const next = visible.value[selectedIndex.value + delta]
  if (next) selectedKey.value = next[props.rowKey]
}

// ---- Editing -----------------------------------------------------------
const editing = ref(false)
const saving = ref(false)
const editError = ref('')
const draft = ref<Record<string, any>>({})

const isEditable = (c: Column) => props.editable && c.editable !== false && !c.key.startsWith('_')

const startEdit = () => {
  if (!selectedRow.value) return
  draft.value = Object.fromEntries(
    drawerColumns.value.filter(isEditable).map((c) => [c.key, selectedRow.value![c.key] ?? '']),
  )
  editError.value = ''
  editing.value = true
}

const cancelEdit = () => {
  editing.value = false
  editError.value = ''
}

const saveEdit = async () => {
  if (!selectedRow.value || !props.onSave) return
  const values = Object.fromEntries(
    Object.entries(draft.value).map(([k, v]) => {
      const col = props.columns.find((c) => c.key === k)
      return [k, col?.type === 'number' ? Number(v) : v]
    }),
  )
  saving.value = true
  editError.value = ''
  try {
    await props.onSave(selectedRow.value, values)
    editing.value = false
  } catch (e: any) {
    editError.value = e?.message || 'Could not save changes.'
  } finally {
    saving.value = false
  }
}

// Leave edit mode when another row is opened or the drawer closes
watch(selectedKey, cancelEdit)

const fieldClass =
  'h-10 w-full rounded border border-slate-600 bg-slate-800 px-3 text-sm text-gray-200 focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-green-400'
</script>

<template>
  <div class="flex flex-col gap-3">
    <div v-if="searchable" class="flex items-center justify-between gap-4">
      <Input
        class="max-w-xs"
        type="search"
        :placeholder="searchPlaceholder"
        :model-value="search"
        @update:model-value="search = String($event ?? '')"
      />
      <span class="text-xs text-slate-400">{{ visible.length }} of {{ rows.length }}</span>
    </div>

    <div :class="cn(tableVariants({ variant }), props.class)">
      <div class="overflow-x-auto">
        <table class="w-full border-collapse">
          <thead>
            <tr class="border-b border-slate-700 bg-slate-800 text-sm text-slate-300">
              <th
                v-for="col in tableColumns"
                :key="col.key"
                :class="cn(cellVariants({ density }), col.sortable && 'cursor-pointer select-none hover:text-white text-sm font-bold', col.class)"
                :aria-sort="sortKey === col.key ? (sortDir === 'asc' ? 'ascending' : 'descending') : undefined"
                @click="toggleSort(col)"
              >
                {{ col.label }}
                <span v-if="sortKey === col.key" class="ml-1 text-green-400">
                  {{ sortDir === 'asc' ? '↑' : '↓' }}
                </span>
              </th>
            </tr>
          </thead>

          <tbody>
            <tr v-if="loading">
              <td :colspan="tableColumns.length" class="p-6 text-center text-slate-400 text-sm">Loading records...</td>
            </tr>
            <tr v-else-if="visible.length === 0">
              <td :colspan="tableColumns.length" class="p-6 text-center text-slate-400 text-sm">
                <slot name="empty">
                  {{ search ? `No results for "${search}".` : emptyText }}
                </slot>
              </td>
            </tr>
            <template v-else>
              <tr
                v-for="row in visible"
                :key="row[rowKey]"
                :class="[
                  'border-b border-slate-700/50 transition-colors hover:bg-slate-800/60 text-sm',
                  rowDrawer && 'cursor-pointer',
                ]"
                :tabindex="rowDrawer ? 0 : undefined"
                @click="onRowClick($event, row)"
                @keydown.enter.self="rowDrawer && (selectedKey = row[rowKey])"
              >
                <td v-for="col in tableColumns" :key="col.key" :class="cn(cellVariants({ density }), col.class)">
                  <slot :name="`cell-${col.key}`" :row="row" :value="row[col.key]">
                    {{ row[col.key] }}
                  </slot>
                </td>
              </tr>
            </template>
          </tbody>
        </table>
      </div>
    </div>

    <!-- CHANGED: passes drawerSize through -->
    <Drawer
      v-if="rowDrawer"
      v-model:open="drawerOpen"
      :side="drawerSide"
      :size="drawerSize"
      :close-on-overlay="!editing"
    >
      <template #title>
        <slot name="drawer-title" :row="selectedRow">{{ editing ? 'Edit details' : 'Row details' }}</slot>
      </template>
      <template #description>
        {{ selectedIndex + 1 }} of {{ visible.length }}
      </template>

      <slot name="drawer" :row="selectedRow">
        <dl v-if="selectedRow" class="grid grid-cols-1 gap-x-6 gap-y-2 sm:grid-cols-2 lg:grid-cols-3">
          <div v-for="col in drawerColumns" :key="col.key" class="flex flex-col gap-1 py-3">
            <dt class="text-xs font-medium uppercase tracking-wide text-slate-400">{{ col.label }}</dt>
            <dd class="break-words text-sm text-white">
              <!-- Edit mode: inputs -->
              <slot v-if="editing && isEditable(col)" :name="`edit-${col.key}`" :draft="draft" :col="col">
                <select v-if="col.options" v-model="draft[col.key]" :class="fieldClass">
                  <option v-for="o in col.options" :key="o" :value="o">{{ o }}</option>
                </select>
                <Input v-else v-model="draft[col.key]" :type="col.type ?? 'text'" />
              </slot>
              <!-- View mode: values -->
              <slot v-else :name="`cell-${col.key}`" :row="selectedRow" :value="selectedRow[col.key]">
                {{ selectedRow[col.key] === '' || selectedRow[col.key] == null ? '—' : selectedRow[col.key] }}
              </slot>
            </dd>
          </div>
        </dl>
        <p v-if="editError" class="mt-3 text-sm text-red-400">{{ editError }}</p>
      </slot>

      <template #footer="{ close }">
        <template v-if="editing">
          <Button variant="secondary" size="sm" :disabled="saving" @click="cancelEdit">Cancel</Button>
          <Button size="sm" :loading="saving" @click="saveEdit">Save changes</Button>
        </template>
        <template v-else>
          <Button variant="secondary" size="sm" :disabled="selectedIndex <= 0" @click="step(-1)">← Prev</Button>
          <Button variant="secondary" size="sm" :disabled="selectedIndex >= visible.length - 1" @click="step(1)">Next →</Button>
          <Button v-if="editable" variant="secondary" size="sm" @click="startEdit">Edit</Button>
          <Button size="sm" @click="close()">Close</Button>
        </template>
      </template>
    </Drawer>
  </div>
</template>