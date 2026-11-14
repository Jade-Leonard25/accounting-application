<script setup lang="ts">
import { computed, reactive, ref } from 'vue'
import { useRouter } from 'vue-router'
import { z } from 'zod'
import { signUp } from '../lib/auth-client'
import Button from '../components/Button.vue'
import Input from '../components/Input.vue'
import Select from '../components/Select.vue'

const router = useRouter()

const roleOptions = [
  { label: 'Admin', value: 'admin' },
  { label: 'User', value: 'user' },
]

const schema = z.object({
  name: z.string().min(1, 'Name is required'),
  email: z.string().min(1, 'Email is required').pipe(z.email('Enter a valid email')),
  password: z.string().min(8, 'Password must be at least 8 characters'),
  role: z.enum(['admin', 'user'], { error: 'Choose a role' }),
})

type Field = keyof z.infer<typeof schema>

const form = reactive({ name: '', email: '', password: '', role: 'user' })
const errors = ref<Partial<Record<Field, string>>>({})
const loading = ref(false)
const serverError = ref('')

// Exactly what is sent to the server
const payload = computed(() => ({
  name: form.name,
  email: form.email,
  password: form.password,
  role: form.role,
}))

// Dev-only preview, password hidden
const isDev = import.meta.env.DEV
const preview = computed(() => JSON.stringify({ ...payload.value, password: '••••••••' }, null, 2))

const handleSubmit = async () => {
  serverError.value = ''
  errors.value = {}

  const result = schema.safeParse(form)
  if (!result.success) {
    for (const issue of result.error.issues) {
      const key = issue.path[0] as Field
      if (!errors.value[key]) errors.value[key] = issue.message
    }
    return
  }

  loading.value = true
  try {
    // `role` is an extra field, so it needs the server and client config described below
    const response = await signUp.email(payload.value as any)

    if (response?.error) {
      serverError.value = response.error.message || 'Failed to create account.'
      return
    }

    router.push({ name: 'dashboard' })
  } catch (e: any) {
    serverError.value = e?.message || 'An unexpected error occurred.'
  } finally {
    loading.value = false
  }
}
</script>

<template>
  <div class="grid min-h-screen grid-cols-1 bg-zinc-950 text-zinc-100 lg:grid-cols-2">
    <!-- Left: visual panel -->
    <div class="relative hidden flex-col justify-between overflow-hidden border-r border-zinc-800/80 bg-zinc-900 p-10 lg:flex">
      <div class="z-10 flex items-center gap-2 text-lg font-semibold tracking-tight">
        <span class="h-3 w-3 rounded-full bg-emerald-500"></span>
        Payroll
      </div>

      <div class="z-10 my-auto flex justify-center">
        <img
          src="/auth.jpeg"
          alt="Payroll"
          class="w-full max-w-md rounded-lg border border-zinc-800 shadow-2xl brightness-90 contrast-105"
        />
      </div>

      <p class="z-10 text-xs text-zinc-500">Internal Payroll & Employee Management</p>
    </div>

    <!-- Right: form -->
    <div class="flex items-center justify-center p-6 sm:p-12">
      <div class="w-full max-w-sm space-y-6">
        <div class="space-y-1.5">
          <h1 class="text-2xl font-semibold tracking-tight">Create an account</h1>
          <p class="text-sm text-zinc-400">Enter your details to get started</p>
        </div>

        <form class="space-y-4" novalidate @submit.prevent="handleSubmit">
          <div class="space-y-1.5">
            <label for="name" class="text-xs font-medium text-zinc-300">Full Name</label>
            <Input
              id="name"
              v-model="form.name"
              placeholder="John Doe"
              :variant="errors.name ? 'error' : 'default'"
            />
            <p v-if="errors.name" class="text-xs text-red-400">{{ errors.name }}</p>
          </div>

          <div class="space-y-1.5">
            <label for="email" class="text-xs font-medium text-zinc-300">Email</label>
            <Input
              id="email"
              v-model="form.email"
              type="email"
              placeholder="name@company.com"
              :variant="errors.email ? 'error' : 'default'"
            />
            <p v-if="errors.email" class="text-xs text-red-400">{{ errors.email }}</p>
          </div>

          <div class="space-y-1.5">
            <label for="password" class="text-xs font-medium text-zinc-300">Password</label>
            <Input
              id="password"
              v-model="form.password"
              type="password"
              placeholder="••••••••"
              :variant="errors.password ? 'error' : 'default'"
            />
            <p v-if="errors.password" class="text-xs text-red-400">{{ errors.password }}</p>
          </div>

          <div class="space-y-1.5">
            <span class="text-xs font-medium text-zinc-300">Role</span>
            <Select
              v-model="form.role"
              :options="roleOptions"
              placeholder="Select a role"
              :variant="errors.role ? 'error' : 'default'"
            />
            <p v-if="errors.role" class="text-xs text-red-400">{{ errors.role }}</p>
          </div>

          <p v-if="serverError" class="text-xs text-red-400">{{ serverError }}</p>

          <Button type="submit" class="w-full" :loading="loading">
            {{ loading ? 'Creating account...' : 'Create account' }}
          </Button>
        </form>

        <!-- Dev only: what will be sent -->
        <details v-if="isDev" class="rounded border border-zinc-800 p-3 text-xs text-zinc-400">
          <summary class="cursor-pointer select-none">Request preview (dev only)</summary>
          <pre class="mt-2 overflow-x-auto">POST /api/auth/sign-up/email
{{ preview }}</pre>
        </details>

        <p class="text-center text-xs text-zinc-500">
          Already have an account?
          <router-link to="/login" class="text-zinc-200 hover:underline">Sign in</router-link>
        </p>
      </div>
    </div>
  </div>
</template>