<script setup lang="ts">
import { computed, onMounted, reactive, ref, watch } from 'vue'
import { useEditor, EditorContent } from '@tiptap/vue-3'
import StarterKit from '@tiptap/starter-kit'
import { get, set } from 'idb-keyval'
import Button from './Button.vue'

const props = defineProps<{
  selectedEmployee: {
    name: string
    basicCentavos: number
    sssCentavos: number
    philhealthCentavos: number
    pagibigCentavos: number
    taxCentavos: number
    netPayCentavos: number
  } | null
  employeeId?: string
  accountLast4?: string
  period?: string // e.g. "November 2026"
}>()

// ---- Design state (persisted) -----------------------------------------
const DEFAULT_DESIGN = {
  companyName: 'ACME Corporation',
  companyAddress: '123 Business Ave, Makati City',
  footerNote: 'This is a system-generated document. No signature required.',
  showLogo: true,
  showNotes: true,
  logo: null as string | null,
  font: 'sans' as 'sans' | 'serif' | 'mono',
  accent: '#16a34a',
  headerBg: '#f8fafc',
  headerText: '#111827',
  paperBg: '#ffffff',
  ink: '#374151',
  netBg: '#f3f4f6',
}

type Design = typeof DEFAULT_DESIGN
const design = reactive<Design>({ ...DEFAULT_DESIGN })
const ready = ref(false)

const DESIGN_KEY = 'payslip:design'

onMounted(async () => {
  try {
    const saved = await get<Partial<Design>>(DESIGN_KEY)
    if (saved) Object.assign(design, saved)
  } catch {
    /* storage unavailable */
  } finally {
    ready.value = true
  }
})

watch(
  design,
  (v) => {
    if (!ready.value) return
    set(DESIGN_KEY, JSON.parse(JSON.stringify(v))).catch(() => {})
  },
  { deep: true },
)

const reset = () => Object.assign(design, DEFAULT_DESIGN)

// ---- Colors and presets -----------------------------------------------
const colorFields = [
  { key: 'accent', label: 'Accent' },
  { key: 'headerBg', label: 'Header' },
  { key: 'headerText', label: 'Header text' },
  { key: 'paperBg', label: 'Paper' },
  { key: 'ink', label: 'Body text' },
  { key: 'netBg', label: 'Net pay box' },
] as const

const presets = [
  { name: 'Classic', accent: '#111827', headerBg: '#f8fafc', headerText: '#111827', paperBg: '#ffffff', ink: '#374151', netBg: '#f3f4f6' },
  { name: 'Ocean', accent: '#2563eb', headerBg: '#eff6ff', headerText: '#1e3a8a', paperBg: '#ffffff', ink: '#334155', netBg: '#dbeafe' },
  { name: 'Forest', accent: '#16a34a', headerBg: '#f0fdf4', headerText: '#14532d', paperBg: '#ffffff', ink: '#374151', netBg: '#dcfce7' },
  { name: 'Sunset', accent: '#ea580c', headerBg: '#fff7ed', headerText: '#7c2d12', paperBg: '#ffffff', ink: '#44403c', netBg: '#ffedd5' },
  { name: 'Bold', accent: '#facc15', headerBg: '#111827', headerText: '#ffffff', paperBg: '#ffffff', ink: '#1f2937', netBg: '#fef9c3' },
]

const applyPreset = (p: (typeof presets)[number]) => {
  const { name, ...colors } = p
  Object.assign(design, colors)
}

// Pick readable text color on top of any background
const textOn = (hex: string) => {
  const h = hex.replace('#', '')
  const full = h.length === 3 ? h.split('').map((c) => c + c).join('') : h
  const n = parseInt(full, 16)
  const lum = 0.299 * ((n >> 16) & 255) + 0.587 * ((n >> 8) & 255) + 0.114 * (n & 255)
  return lum > 150 ? '#111827' : '#ffffff'
}

const fontStack = computed(() => ({
  sans: 'ui-sans-serif, system-ui, sans-serif',
  serif: 'ui-serif, Georgia, "Times New Roman", serif',
  mono: 'ui-monospace, "SFMono-Regular", Menlo, monospace',
})[design.font])

const paperVars = computed(() => ({
  '--accent': design.accent,
  '--on-accent': textOn(design.accent),
  '--hdr': design.headerBg,
  '--hdr-text': design.headerText,
  '--paper': design.paperBg,
  '--ink': design.ink,
  '--net': design.netBg,
  '--line': `${design.ink}26`, // body text at ~15% opacity
  fontFamily: fontStack.value,
}))

// ---- Logo upload (use your own company logo only) ---------------------
const logoError = ref('')
const onLogo = (e: Event) => {
  logoError.value = ''
  const input = e.target as HTMLInputElement
  const file = input.files?.[0]
  input.value = ''
  if (!file) return
  if (!/^image\/(png|jpeg|svg\+xml|webp)$/.test(file.type)) {
    logoError.value = 'Use a PNG, JPG, SVG, or WebP image.'
    return
  }
  if (file.size > 500 * 1024) {
    logoError.value = 'Logo must be under 500 KB.'
    return
  }
  const reader = new FileReader()
  reader.onload = () => {
    design.logo = String(reader.result)
    design.showLogo = true
  }
  reader.readAsDataURL(file)
}

// ---- Inline editing on the paper --------------------------------------
const onInline = (e: Event, key: 'companyName' | 'companyAddress' | 'footerNote') => {
  design[key] = (e.target as HTMLElement).innerText.trim()
}

// ---- Notes editor -----------------------------------------------------
const notesEditor = useEditor({
  extensions: [StarterKit],
  content: '<p>Thank you for your hard work this month!</p>',
  editorProps: {
    attributes: { class: 'focus:outline-none min-h-[100px] text-sm text-slate-300' },
  },
})
const notesHtml = computed(() => notesEditor.value?.getHTML() || '')
const hasNotes = computed(() => design.showNotes && notesHtml.value !== '<p></p>')

// ---- Data -------------------------------------------------------------
const fmt = (c: number) => (c / 100).toLocaleString('en-PH', { minimumFractionDigits: 2 })

const previewData = computed(
  () =>
    props.selectedEmployee ?? {
      name: 'Juan Dela Cruz',
      basicCentavos: 3500000,
      sssCentavos: 175000,
      philhealthCentavos: 87500,
      pagibigCentavos: 20000,
      taxCentavos: 250000,
      netPayCentavos: 2967500,
    },
)

const totalDeductions = computed(() => {
  const d = previewData.value
  return d.sssCentavos + d.philhealthCentavos + d.pagibigCentavos + d.taxCentavos
})

const periodLabel = computed(
  () => props.period ?? new Date().toLocaleDateString('en-PH', { month: 'long', year: 'numeric' }),
)

const print = () => window.print()

const inputClass =
  'w-full rounded border border-slate-600 bg-slate-900 px-3 py-2 text-sm text-white focus:outline-none focus:ring-2 focus:ring-green-400'
</script>

<template>
  <div class="grid grid-cols-1 gap-6 lg:grid-cols-12">
    <!-- LEFT: settings -->
    <div class="flex flex-col gap-4 lg:col-span-4">
      <div class="rounded-lg border border-slate-700 bg-slate-800 p-5">
        <h3 class="mb-4 text-lg font-bold text-white">Company Details</h3>
        <div class="flex flex-col gap-3">
          <div>
            <label class="mb-1 block text-xs text-slate-400">Company Name</label>
            <input v-model="design.companyName" type="text" :class="inputClass" />
          </div>
          <div>
            <label class="mb-1 block text-xs text-slate-400">Address</label>
            <input v-model="design.companyAddress" type="text" :class="inputClass" />
          </div>
          <div>
            <label class="mb-1 block text-xs text-slate-400">Footer Note</label>
            <input v-model="design.footerNote" type="text" :class="inputClass" />
          </div>
          <label class="mt-1 flex cursor-pointer items-center gap-2">
            <input v-model="design.showNotes" type="checkbox" class="rounded border-slate-600" />
            <span class="text-sm text-slate-300">Show notes section</span>
          </label>
        </div>
      </div>

      <div class="flex flex-1 flex-col rounded-lg border border-slate-700 bg-slate-800 p-5">
        <h3 class="mb-2 text-lg font-bold text-white">Additional Notes</h3>
        <p class="mb-3 text-xs text-slate-400">Printed at the bottom of the payslip.</p>
        <div class="flex-1 rounded border border-slate-600 bg-slate-900 p-3">
          <EditorContent :editor="notesEditor" />
        </div>
      </div>
    </div>

    <!-- RIGHT: design toolbar + paper -->
    <div class="flex flex-col gap-4 lg:col-span-8">
      <!-- Design toolbar -->
      <div class="flex flex-col gap-4 rounded-lg border border-slate-700 bg-slate-800 p-4">
        <div class="flex flex-wrap items-center gap-2">
          <span class="mr-1 text-xs font-medium uppercase tracking-wide text-slate-400">Presets</span>
          <button
            v-for="p in presets"
            :key="p.name"
            type="button"
            class="flex cursor-pointer items-center gap-2 rounded border border-slate-600 px-2 py-1 text-xs text-slate-200 transition-colors hover:bg-slate-700"
            @click="applyPreset(p)"
          >
            <span class="h-3 w-3 rounded-full" :style="{ background: p.accent }" />
            {{ p.name }}
          </button>
        </div>

        <div class="grid grid-cols-2 gap-3 sm:grid-cols-3 lg:grid-cols-6">
          <label v-for="f in colorFields" :key="f.key" class="flex flex-col gap-1">
            <span class="text-xs text-slate-400">{{ f.label }}</span>
            <span class="flex items-center gap-2 rounded border border-slate-600 bg-slate-900 p-1">
              <input
                v-model="design[f.key]"
                type="color"
                class="h-7 w-7 shrink-0 cursor-pointer rounded border-0 bg-transparent p-0"
                :aria-label="f.label"
              />
              <span class="truncate font-mono text-xs text-slate-300">{{ design[f.key] }}</span>
            </span>
          </label>
        </div>

        <div class="flex flex-wrap items-end justify-between gap-3">
          <div class="flex flex-wrap items-end gap-3">
            <label class="flex flex-col gap-1">
              <span class="text-xs text-slate-400">Font</span>
              <select v-model="design.font" :class="[inputClass, 'h-9 py-0']">
                <option value="sans">Sans</option>
                <option value="serif">Serif</option>
                <option value="mono">Mono</option>
              </select>
            </label>

            <div class="flex flex-col gap-1">
              <span class="text-xs text-slate-400">Logo</span>
              <div class="flex items-center gap-2">
                <label
                  class="inline-flex h-9 cursor-pointer items-center rounded border border-slate-600 bg-slate-700 px-3 text-xs font-medium text-white hover:bg-slate-600"
                >
                  Upload
                  <input type="file" accept="image/png,image/jpeg,image/svg+xml,image/webp" class="sr-only" @change="onLogo" />
                </label>
                <Button v-if="design.logo" variant="ghost" size="sm" type="button" @click="design.logo = null">Remove</Button>
                <label class="flex cursor-pointer items-center gap-2 text-xs text-slate-300">
                  <input v-model="design.showLogo" type="checkbox" class="rounded border-slate-600" />
                  Show
                </label>
              </div>
            </div>
          </div>

          <div class="flex gap-2">
            <Button variant="secondary" size="sm" type="button" @click="reset">Reset design</Button>
            <Button size="sm" type="button" @click="print">Print / Save PDF</Button>
          </div>
        </div>
        <p v-if="logoError" class="text-xs text-red-400">{{ logoError }}</p>
        <p class="text-xs text-slate-500">Tip: click the company name, address, or footer on the paper to edit it directly.</p>
      </div>

      <!-- The paper -->
      <div class="flex justify-center">
        <div
          id="payslip-paper"
          :style="paperVars"
          class="relative flex w-full max-w-[800px] flex-col overflow-hidden rounded-sm bg-[var(--paper)] text-[var(--ink)] shadow-2xl"
        >
          <!-- Header -->
          <div class="flex items-start justify-between border-b border-[var(--line)] bg-[var(--hdr)] p-8">
            <div class="flex items-center gap-4">
              <template v-if="design.showLogo">
                <img v-if="design.logo" :src="design.logo" alt="Company logo" class="h-12 max-w-[120px] object-contain" />
                <div
                  v-else
                  class="flex h-12 w-12 items-center justify-center rounded bg-[var(--accent)] text-xl font-bold text-[var(--on-accent)]"
                >
                  {{ (design.companyName || '?').charAt(0) }}
                </div>
              </template>
              <div>
                <h1
                  contenteditable="plaintext-only"
                  spellcheck="false"
                  class="rounded text-2xl font-extrabold text-[var(--hdr-text)] outline-none focus:ring-2 focus:ring-[var(--accent)]"
                  @keydown.enter.prevent
                  @blur="onInline($event, 'companyName')"
                >{{ design.companyName }}</h1>
                <p
                  contenteditable="plaintext-only"
                  spellcheck="false"
                  class="rounded text-sm text-[var(--hdr-text)] opacity-70 outline-none focus:ring-2 focus:ring-[var(--accent)]"
                  @keydown.enter.prevent
                  @blur="onInline($event, 'companyAddress')"
                >{{ design.companyAddress }}</p>
              </div>
            </div>
            <div class="text-right">
              <p class="text-xs font-bold uppercase tracking-wider text-[var(--hdr-text)] opacity-60">Payslip</p>
              <p class="text-sm font-medium text-[var(--hdr-text)]">{{ periodLabel }}</p>
            </div>
          </div>

          <!-- Employee info -->
          <div class="grid grid-cols-2 gap-6 border-b border-[var(--line)] p-8">
            <div>
              <p class="mb-1 text-xs uppercase tracking-wide opacity-60">Employee Name</p>
              <p class="font-semibold">{{ previewData.name }}</p>
            </div>
            <div v-if="employeeId">
              <p class="mb-1 text-xs uppercase tracking-wide opacity-60">Employee ID</p>
              <p class="font-semibold">{{ employeeId }}</p>
            </div>
          </div>

          <!-- Earnings and deductions -->
          <div class="flex-1 p-8">
            <table class="w-full text-sm">
              <thead>
                <tr class="border-b border-[var(--line)]">
                  <th class="pb-2 text-left font-semibold">Earnings</th>
                  <th class="pb-2 text-right font-semibold">Amount</th>
                  <th class="w-8"></th>
                  <th class="pb-2 text-left font-semibold">Deductions</th>
                  <th class="pb-2 text-right font-semibold">Amount</th>
                </tr>
              </thead>
              <tbody>
                <tr>
                  <td class="py-2">Basic Pay</td>
                  <td class="py-2 text-right">₱{{ fmt(previewData.basicCentavos) }}</td>
                  <td></td>
                  <td class="py-2">SSS Contribution</td>
                  <td class="py-2 text-right">₱{{ fmt(previewData.sssCentavos) }}</td>
                </tr>
                <tr>
                  <td class="py-2">Allowances</td>
                  <td class="py-2 text-right">₱0.00</td>
                  <td></td>
                  <td class="py-2">PhilHealth</td>
                  <td class="py-2 text-right">₱{{ fmt(previewData.philhealthCentavos) }}</td>
                </tr>
                <tr>
                  <td></td><td></td><td></td>
                  <td class="py-2">Pag-IBIG</td>
                  <td class="py-2 text-right">₱{{ fmt(previewData.pagibigCentavos) }}</td>
                </tr>
                <tr>
                  <td></td><td></td><td></td>
                  <td class="py-2">Withholding Tax</td>
                  <td class="py-2 text-right">₱{{ fmt(previewData.taxCentavos) }}</td>
                </tr>
              </tbody>
              <tfoot>
                <tr class="border-t border-[var(--line)] font-semibold">
                  <td class="py-3">Total Earnings</td>
                  <td class="py-3 text-right">₱{{ fmt(previewData.basicCentavos) }}</td>
                  <td></td>
                  <td class="py-3">Total Deductions</td>
                  <td class="py-3 text-right">₱{{ fmt(totalDeductions) }}</td>
                </tr>
              </tfoot>
            </table>
          </div>

          <!-- Net pay -->
          <div class="flex items-center justify-between border-t border-[var(--line)] bg-[var(--net)] p-8">
            <div>
              <p class="text-sm font-medium">Net Take-Home Pay</p>
              <p v-if="accountLast4" class="mt-1 text-xs opacity-60">Deposited to account ending in {{ accountLast4 }}</p>
            </div>
            <p class="text-3xl font-extrabold text-[var(--accent)]">₱{{ fmt(previewData.netPayCentavos) }}</p>
          </div>

          <!-- Notes -->
          <div v-if="hasNotes" class="border-t border-[var(--line)] p-8 pt-0">
            <p class="mb-2 mt-4 text-xs font-semibold uppercase tracking-wider opacity-60">Notes</p>
            <div class="text-sm [&_p]:mb-2" v-html="notesHtml" />
          </div>

          <!-- Footer -->
          <div class="mt-auto border-t border-[var(--line)] p-4 text-center">
            <p
              contenteditable="plaintext-only"
              spellcheck="false"
              class="rounded text-[10px] opacity-60 outline-none focus:ring-2 focus:ring-[var(--accent)]"
              @keydown.enter.prevent
              @blur="onInline($event, 'footerNote')"
            >{{ design.footerNote }}</p>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<style>
/* Print only the paper, with its colors */
@media print {
  @page { margin: 12mm; }
  body * { visibility: hidden !important; }
  #payslip-paper, #payslip-paper * { visibility: visible !important; }
  #payslip-paper {
    position: absolute;
    inset: 0 auto auto 0;
    width: 100%;
    max-width: none;
    box-shadow: none;
    -webkit-print-color-adjust: exact;
    print-color-adjust: exact;
  }
}
</style>