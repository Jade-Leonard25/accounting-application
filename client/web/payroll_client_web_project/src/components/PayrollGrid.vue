<script setup lang="ts">
import { ref } from 'vue'
import * as XLSX from 'xlsx'
import Button from './Button.vue'
import Table from './Table.vue' // Using your custom table for consistency

export interface PayrollRow {
  id: string
  name: string
  basicCentavos: number
  allowancesCentavos: number
  sssCentavos: number
  philhealthCentavos: number
  pagibigCentavos: number
  taxCentavos: number
  netPayCentavos: number
}

const rows = ref<PayrollRow[]>([
  {
    id: '49v-82m-31k',
    name: 'Juan Dela Cruz',
    basicCentavos: 3500000, // ₱35,000.00
    allowancesCentavos: 200000, // ₱2,000.00
    sssCentavos: 0,
    philhealthCentavos: 0,
    pagibigCentavos: 0,
    taxCentavos: 0,
    netPayCentavos: 3700000,
  },
])

const emit = defineEmits<{
  (e: 'select-payslip', row: PayrollRow): void
}>()

// 1. Import existing agency Excel (.xlsx)
function handleExcelImport(event: Event) {
  const file = (event.target as HTMLInputElement).files?.[0]
  if (!file) return

  const reader = new FileReader()
  reader.onload = (e) => {
    const workbook = XLSX.read(e.target?.result, { type: 'binary' })
    const sheet = workbook.Sheets[workbook.SheetNames[0]]
    const rawData = XLSX.utils.sheet_to_json<Record<string, any>>(sheet)

    rows.value = rawData.map((r, idx) => {
      const basic = Math.round((Number(r['Basic Pay']) || 0) * 100)
      const allow = Math.round((Number(r['Allowances']) || 0) * 100)
      return {
        id: r['Employee ID'] || `EMP-${idx + 1}`,
        name: r['Name'] || 'Unnamed',
        basicCentavos: basic,
        allowancesCentavos: allow,
        sssCentavos: Math.round((Number(r['SSS']) || 0) * 100),
        philhealthCentavos: Math.round((Number(r['PhilHealth']) || 0) * 100),
        pagibigCentavos: Math.round((Number(r['Pag-IBIG']) || 0) * 100),
        taxCentavos: Math.round((Number(r['Tax']) || 0) * 100),
        netPayCentavos: basic + allow,
      }
    })
    recalculateNet()
  }
  reader.readAsBinaryString(file)
}

// 2. Auto-Compute Deductions (2026 PH Rules)
function autoComputeAll() {
  rows.value = rows.value.map((row) => {
    const basic = row.basicCentavos
    const sssBase = Math.min(Math.max(basic, 500000), 3500000)
    const sss = Math.round(sssBase * 0.05)
    const phBase = Math.min(Math.max(basic, 1000000), 10000000)
    const philhealth = Math.round((phBase * 0.05) / 2)
    const pagibig = Math.round(Math.min(basic, 1000000) * 0.02)
    const taxablePesos = Math.max(0, basic - (sss + philhealth + pagibig)) / 100
    let taxPesos = 0

    if (taxablePesos > 20833 && taxablePesos <= 33332) taxPesos = (taxablePesos - 20833) * 0.15
    else if (taxablePesos > 33332 && taxablePesos <= 66666) taxPesos = 1875 + (taxablePesos - 33333) * 0.20
    else if (taxablePesos > 66666 && taxablePesos <= 166666) taxPesos = 8541.8 + (taxablePesos - 66667) * 0.25

    const tax = Math.round(taxPesos * 100)
    const netPay = basic + row.allowancesCentavos - (sss + philhealth + pagibig + tax)

    return { ...row, sssCentavos: sss, philhealthCentavos: philhealth, pagibigCentavos: pagibig, taxCentavos: tax, netPayCentavos: netPay }
  })
}

// 3. Manual Override Recalculation
function updateCell(row: PayrollRow, field: keyof PayrollRow, pesoValue: string) {
  ;(row[field] as number) = Math.round((parseFloat(pesoValue) || 0) * 100)
  recalculateNet()
}

function recalculateNet() {
  rows.value.forEach((r) => {
    r.netPayCentavos = r.basicCentavos + r.allowancesCentavos - (r.sssCentavos + r.philhealthCentavos + r.pagibigCentavos + r.taxCentavos)
  })
}

const fmt = (centavos: number) => (centavos / 100).toFixed(2)

// Table configuration for our custom Table component
const columns = [
  { key: 'name', label: 'Employee' },
  { key: 'basicCentavos', label: 'Basic Pay' },
  { key: 'allowancesCentavos', label: 'Allowances' },
  { key: 'sssCentavos', label: 'SSS' },
  { key: 'philhealthCentavos', label: 'PhilHealth' },
  { key: 'pagibigCentavos', label: 'Pag-IBIG' },
  { key: 'taxCentavos', label: 'BIR Tax' },
  { key: 'netPayCentavos', label: 'Net Pay' },
  { key: 'actions', label: 'Payslip' },
]
</script>

<template>
  <div class="space-y-4">
    <div class="flex items-center justify-between gap-2">
      <div class="flex items-center gap-2">
        <label class="px-4 py-2 text-sm font-semibold rounded bg-slate-800 text-gray-200 border border-slate-600 hover:bg-slate-700 cursor-pointer transition-colors">
          Import Excel (.xlsx)
          <input type="file" accept=".xlsx,.xls" class="hidden" @change="handleExcelImport" />
        </label>
        <Button variant="secondary" @click="autoComputeAll">
          Auto-Compute Deductions
        </Button>
      </div>
    </div>

    <Table :columns="columns" :rows="rows" :searchable="false">
      <!-- Editable Cells -->
      <template #cell-basicCentavos="{ row }">
        <input type="number" :value="fmt(row.basicCentavos)" @change="updateCell(row, 'basicCentavos', ($event.target as HTMLInputElement).value)"
          class="w-24 px-2 py-1 text-right bg-transparent border-b border-transparent focus:border-green-400 focus:bg-slate-800 outline-none text-gray-200" />
      </template>
      
      <template #cell-allowancesCentavos="{ row }">
        <input type="number" :value="fmt(row.allowancesCentavos)" @change="updateCell(row, 'allowancesCentavos', ($event.target as HTMLInputElement).value)"
          class="w-24 px-2 py-1 text-right bg-transparent border-b border-transparent focus:border-green-400 focus:bg-slate-800 outline-none text-gray-200" />
      </template>

      <template #cell-sssCentavos="{ row }">
        <input type="number" :value="fmt(row.sssCentavos)" @change="updateCell(row, 'sssCentavos', ($event.target as HTMLInputElement).value)"
          class="w-20 px-2 py-1 text-right bg-transparent border-b border-transparent focus:border-green-400 focus:bg-slate-800 outline-none text-gray-200" />
      </template>

      <template #cell-philhealthCentavos="{ row }">
        <input type="number" :value="fmt(row.philhealthCentavos)" @change="updateCell(row, 'philhealthCentavos', ($event.target as HTMLInputElement).value)"
          class="w-20 px-2 py-1 text-right bg-transparent border-b border-transparent focus:border-green-400 focus:bg-slate-800 outline-none text-gray-200" />
      </template>

      <template #cell-pagibigCentavos="{ row }">
        <input type="number" :value="fmt(row.pagibigCentavos)" @change="updateCell(row, 'pagibigCentavos', ($event.target as HTMLInputElement).value)"
          class="w-20 px-2 py-1 text-right bg-transparent border-b border-transparent focus:border-green-400 focus:bg-slate-800 outline-none text-gray-200" />
      </template>

      <template #cell-taxCentavos="{ row }">
        <input type="number" :value="fmt(row.taxCentavos)" @change="updateCell(row, 'taxCentavos', ($event.target as HTMLInputElement).value)"
          class="w-24 px-2 py-1 text-right bg-transparent border-b border-transparent focus:border-green-400 focus:bg-slate-800 outline-none text-gray-200" />
      </template>

      <!-- Calculated Net Pay -->
      <template #cell-netPayCentavos="{ row }">
        <span class="font-semibold text-green-400">₱{{ fmt(row.netPayCentavos) }}</span>
      </template>

      <!-- Action Button -->
      <template #cell-actions="{ row }">
        <Button variant="outline" size="sm" @click="emit('select-payslip', row)">Preview Doc</Button>
      </template>
    </Table>
  </div>
</template>