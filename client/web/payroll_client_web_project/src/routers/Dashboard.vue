<script setup lang="ts">
import { ref, reactive, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import axios from 'axios'
import { z } from 'zod'
import { signOut, useSession } from '../lib/auth-client'
import Button from '../components/Button.vue'
import Input from '../components/Input.vue'
import Card from '../components/Card.vue'
import Dropdown from '../components/Dropdown.vue'
import Table from '../components/Table.vue'
import Sidebar from '../components/Sidebar.vue'
import ImportCsv from '../components/ImportCsv.vue'
import ThemeToggle from '../components/ThemeToggle.vue'

// Import the new payroll components
import PayrollGrid from '../components/PayrollGrid.vue'
import PayslipEditor from '../components/PayslipEditor.vue'

const router = useRouter()
const session = useSession()

// ---- Sidebar -----------------------------------------------------------
const section = ref('employees')
const draftCount = ref(0)

const sidebarItems = computed(() => [
  { key: 'employees', label: 'Employees' },
  { key: 'import', label: 'Import CSV', badge: draftCount.value },
  { key: 'payroll', label: 'Process Payroll' }, // 👈 Added
  { key: 'templates', label: 'Payslip Template' }, // 👈 Added
])

// ---- Types and state ---------------------------------------------------
interface Employee {
  id: string
  user_id?: string | null
  first_name: string
  last_name: string
  email: string
  base_salary_cents: number
  tax_id?: string
  employment_status: 'Full-time' | 'Part-time' | 'Contract'
  pay_frequency: 'Semi-monthly' | 'Monthly' | 'Weekly'
}

const employees = ref<Employee[]>([])
const loading = ref(false)
const error = ref('')
const submitting = ref(false)
const showModal = ref(false)
const isEditing = ref(false)
const currentId = ref<string | null>(null)

// ---- Table -------------------------------------------------------------
const columns = [
  { key: 'name', label: 'Name', sortable: true, editable: false, hideInDrawer: true },
  { key: 'first_name', label: 'First Name', drawerOnly: true },
  { key: 'last_name', label: 'Last Name', drawerOnly: true },
  { key: 'email', label: 'Email', sortable: true, type: 'email' as const },
  { key: 'tax_id', label: 'Tax ID', drawerOnly: true },
  { key: 'employment_status', label: 'Status', sortable: true, options: ['Full-time', 'Part-time', 'Contract'] },
  { key: 'pay_frequency', label: 'Pay Frequency', drawerOnly: true, options: ['Semi-monthly', 'Monthly', 'Weekly'] },
  { key: 'base_salary_cents', label: 'Base Salary', sortable: true, type: 'number' as const },
  { key: 'actions', label: 'Actions', class: 'text-right', hideInDrawer: true, editable: false },
]

const rows = computed(() =>
  employees.value.map((e) => ({ ...e, name: `${e.first_name} ${e.last_name}` })),
)

// ---- Form --------------------------------------------------------------
const emptyForm = (): Omit<Employee, 'id'> => ({
  first_name: '',
  last_name: '',
  email: '',
  base_salary_cents: 0,
  tax_id: '',
  employment_status: 'Full-time',
  pay_frequency: 'Semi-monthly',
})
const saveEmployee = async (row: any, values: Record<string, any>) => {
  const payload = {
    first_name: row.first_name,
    last_name: row.last_name,
    email: row.email,
    base_salary_cents: row.base_salary_cents,
    tax_id: row.tax_id ?? '',
    employment_status: row.employment_status,
    pay_frequency: row.pay_frequency,
    ...values,
  }
  const parsed = schema.safeParse(payload)
  if (!parsed.success) throw new Error(parsed.error.issues[0]?.message ?? 'Invalid data')
  try {
    await API.put(`/employees/${row.id}`, parsed.data)
  } catch (e: any) {
    throw new Error(e.response?.data?.error || e.response?.data?.message || 'Update failed')
  }
  await fetchEmployees()
}
const form = reactive(emptyForm())

const schema = z.object({
  first_name: z.string().min(1, 'First name is required'),
  last_name: z.string().min(1, 'Last name is required'),
  email: z.email({ error: 'Enter a valid email' }),
  base_salary_cents: z
    .number({ error: 'Salary must be a number' })
    .int('Salary must be a whole number')
    .min(0, 'Salary cannot be negative'),
  tax_id: z.string().optional(),
  employment_status: z.enum(['Full-time', 'Part-time', 'Contract']),
  pay_frequency: z.enum(['Semi-monthly', 'Monthly', 'Weekly']),
})

// ---- API ---------------------------------------------------------------
const API = axios.create({
  baseURL: 'http://localhost:4002/api',
  withCredentials: true,
})

const fetchEmployees = async () => {
  loading.value = true
  error.value = ''
  try {
    const { data } = await API.get<{ employees: Employee[] } | Employee[]>('/employees')
    employees.value = Array.isArray(data) ? data : data.employees ?? []
  } catch (e: any) {
    if (e.response?.status === 401) {
      router.push({ name: 'login' })
      return
    }
    error.value = e.response?.data?.error || 'Failed to load employees.'
    employees.value = []
  } finally {
    loading.value = false
  }
}

onMounted(fetchEmployees)

const onImported = () => {
  fetchEmployees()
  section.value = 'employees'
}

// ---- Modal handlers ----------------------------------------------------
const openCreateModal = () => {
  isEditing.value = false
  currentId.value = null
  Object.assign(form, emptyForm())
  error.value = ''
  showModal.value = true
}

const openEditModal = (emp: any) => {
  isEditing.value = true
  currentId.value = emp.id
  Object.assign(form, {
    first_name: emp.first_name,
    last_name: emp.last_name,
    email: emp.email,
    base_salary_cents: emp.base_salary_cents,
    tax_id: emp.tax_id ?? '',
    employment_status: emp.employment_status ?? 'Full-time',
    pay_frequency: emp.pay_frequency ?? 'Semi-monthly',
  })
  error.value = ''
  showModal.value = true
}

const generateEmployeeId = (): string => {
  const block = () => Math.random().toString(36).slice(2, 5).padEnd(3, '0')
  return `${block()}-${block()}-${block()}`
}

const handleSubmit = async () => {
  const parsed = schema.safeParse(form)
  if (!parsed.success) {
    error.value = parsed.error.issues[0]?.message ?? 'Invalid form data'
    return
  }

  submitting.value = true
  error.value = ''
  try {
    if (isEditing.value && currentId.value) {
      await API.put(`/employees/${currentId.value}`, parsed.data)
    } else {
      await API.post('/employees', { ...parsed.data, id: generateEmployeeId() })
    }
    showModal.value = false
    await fetchEmployees()
  } catch (e: any) {
    error.value = e.response?.data?.error || e.response?.data?.message || 'Operation failed'
  } finally {
    submitting.value = false
  }
}

const handleLogout = async () => {
  try {
    await signOut()
  } catch {
    /* ignore */
  } finally {
    router.push({ name: 'login' })
  }
}

const selectClass =
  'h-10 w-full rounded border border-slate-600 bg-slate-800 px-3 text-sm text-gray-200 focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-green-400'

// ---- Payroll state -----------------------------------------------------
const selectedPayrollRow = ref<any>(null)
</script>

<template>
  <div class="flex min-h-screen flex-col bg-slate-900 text-gray-200">
    <!-- ROOF -->
    <nav class="flex items-center justify-between border-b border-slate-700 bg-slate-950 px-6 py-4">
      <span class="text-xl font-bold text-white">Payroll Dashboard</span>
      <ThemeToggle />
    </nav>

    <!-- HOUSE -->
    <div class="flex flex-1">
      <!-- HALLWAY -->
      <Sidebar v-model="section" :items="sidebarItems">
        <template #footer>
          <div v-if="session.data?.user" class="flex items-center gap-3 px-1">
            <div class="flex h-9 w-9 shrink-0 items-center justify-center rounded-full bg-green-500 text-sm font-semibold text-slate-950">
              {{ (session.data.user.name || session.data.user.email || '?').charAt(0).toUpperCase() }}
            </div>
            <div class="min-w-0">
              <p v-if="session.data.user.name" class="truncate text-sm font-medium text-white">
                {{ session.data.user.name }}
              </p>
              <p class="truncate text-xs text-slate-400">{{ session.data.user.email }}</p>
            </div>
          </div>

          <Dropdown side="top" align="start">
            <template #trigger="{ toggle, open }">
              <Button variant="destructive" size="sm" :aria-expanded="open" @click="toggle">
                Logout
              </Button>
            </template>

            <template #default="{ close }">
              <div class="flex flex-col gap-3 p-2">
                <p class="text-sm text-slate-300">Are you sure you want to log out?</p>
                <div class="flex justify-end gap-2">
                  <Button variant="secondary" size="sm" @click="close()">Cancel</Button>
                  <Button variant="destructive" size="sm" @click="handleLogout">Logout</Button>
                </div>
              </div>
            </template>
          </Dropdown>
        </template>
      </Sidebar>

      <!-- ROOM -->
      <main class="mx-auto min-w-0 max-w-7xl flex-1 p-6">
        
        <!-- Door 1: Employees room -->
        <div v-if="section === 'employees'">
          <div class="mb-6 flex items-center justify-between">
            <h1 class="text-2xl font-bold text-white">Employee Management</h1>
            <Button @click="openCreateModal">+ Add Employee</Button>
          </div>

          <p v-if="error && !showModal" class="mb-4 text-red-400">{{ error }}</p>

          <Table
            :columns="columns"
            :rows="rows"
            :loading="loading"
            row-drawer
            editable
            :on-save="saveEmployee"
            :search-keys="['name', 'email', 'employment_status', 'tax_id']"
            search-placeholder="Search employees..."
          >
            <template #cell-name="{ value }">
              <span class="font-medium text-white">{{ value }}</span>
            </template>

            <template #cell-employment_status="{ value }">
              <span class="rounded border border-blue-500/30 bg-blue-500/20 px-2 py-1 text-xs text-blue-400">
                {{ value }}
              </span>
            </template>

            <template #cell-actions="{ row }">
              <Button variant="secondary" size="sm" @click="openEditModal(row as Employee)">Edit</Button>
            </template>

            <template #empty>
              No data yet. Click
              <span class="font-semibold text-green-400">"+ Add Employee"</span>
              to create your first record.
            </template>
          </Table>
        </div>

        <!-- Door 2: Import CSV room -->
        <div v-if="section === 'import'">
          <ImportCsv v-model:count="draftCount" @saved="onImported" />
        </div>

        <!-- Door 3: Process Payroll room -->
        <div v-if="section === 'payroll'">
          <div class="mb-6">
            <h1 class="text-2xl font-bold text-white">Process Payroll</h1>
            <p class="text-slate-400 text-sm mt-1">Upload an Excel file or auto-compute deductions.</p>
          </div>
          <Card padding="none">
            <PayrollGrid @select-payslip="selectedPayrollRow = $event" />
          </Card>
        </div>

        <!-- Door 4: Payslip Template room -->
        <div v-if="section === 'templates'">
          <div class="mb-6">
            <h1 class="text-2xl font-bold text-white">Payslip Template</h1>
            <p class="text-slate-400 text-sm mt-1">Design the template that will be used when printing payslips.</p>
          </div>
          <Card padding="none">
            <!-- If a payroll row is selected, pass it. Otherwise, show the live demo. -->
            <PayslipEditor :selected-employee="selectedPayrollRow" />
          </Card>
        </div>

      </main>
    </div>

    <!-- MODAL: the pop-up for add/edit employee -->
    <div v-if="showModal" class="fixed inset-0 flex items-center justify-center bg-black/70 p-4">
      <Card class="w-full max-w-md">
        <template #title>{{ isEditing ? 'Edit Employee' : 'Add New Employee' }}</template>
        <template #description>
          {{ isEditing ? 'Update this employee’s details.' : 'Fill in the details for the new employee.' }}
        </template>

        <form id="employee-form" class="flex flex-col gap-1.5" @submit.prevent="handleSubmit">
          <label for="first_name" class="text-sm text-slate-300">First Name</label>
          <Input id="first_name" v-model="form.first_name" required />

          <label for="last_name" class="mt-1.5 text-sm text-slate-300">Last Name</label>
          <Input id="last_name" v-model="form.last_name" required />

          <label for="email" class="mt-1.5 text-sm text-slate-300">Email</label>
          <Input id="email" v-model="form.email" type="email" required />

          <label for="salary" class="mt-1.5 text-sm text-slate-300">Base Salary (in Cents)</label>
          <Input id="salary" v-model.number="form.base_salary_cents" type="number" min="0" required />

          <label for="tax_id" class="mt-1.5 text-sm text-slate-300">Tax ID</label>
          <Input id="tax_id" v-model="form.tax_id" />

          <label for="status" class="mt-1.5 text-sm text-slate-300">Employment Status</label>
          <select id="status" v-model="form.employment_status" :class="selectClass">
            <option>Full-time</option>
            <option>Part-time</option>
            <option>Contract</option>
          </select>

          <label for="frequency" class="mt-1.5 text-sm text-slate-300">Pay Frequency</label>
          <select id="frequency" v-model="form.pay_frequency" :class="selectClass">
            <option>Semi-monthly</option>
            <option>Monthly</option>
            <option>Weekly</option>
          </select>

          <p v-if="error" class="mt-2 text-sm text-red-400">{{ error }}</p>
        </form>

        <template #footer>
          <div class="flex w-full justify-end gap-2">
            <Button variant="secondary" @click="showModal = false">Cancel</Button>
            <Button type="submit" form="employee-form" :loading="submitting">
              {{ submitting ? 'Saving...' : isEditing ? 'Save Changes' : 'Create Employee' }}
            </Button>
          </div>
        </template>
      </Card>
    </div>
  </div>
</template>