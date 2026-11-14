<script setup lang="ts">
import { computed, ref } from 'vue'
import type { HTMLAttributes } from 'vue'
import { clsx, type ClassValue } from 'clsx'
import { twMerge } from 'tailwind-merge'
import { cva, type VariantProps } from 'class-variance-authority'
import Dropdown from './Dropdown.vue'

const cn = (...inputs: ClassValue[]) => twMerge(clsx(inputs))

const selectVariants = cva(
  'flex w-full items-center justify-between rounded border bg-slate-800 text-gray-200 transition-colors focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-offset-2 focus-visible:ring-offset-slate-900 disabled:cursor-not-allowed disabled:opacity-50',
  {
    variants: {
      variant: {
        default: 'border-slate-600 focus-visible:ring-green-400',
        error: 'border-red-400 focus-visible:ring-red-400',
      },
      size: {
        default: 'h-10 px-3 py-2 text-sm',
        sm: 'h-8 px-2 text-xs',
        lg: 'h-12 px-4 text-base',
      },
    },
    defaultVariants: { variant: 'default', size: 'default' },
  },
)

type SelectVariants = VariantProps<typeof selectVariants>

interface Option {
  label: string
  value: string | number
  disabled?: boolean
}

interface Props {
  options?: (Option | string)[]
  variant?: SelectVariants['variant']
  size?: SelectVariants['size']
  disabled?: boolean
  class?: HTMLAttributes['class']
  placeholder?: string
}

const props = withDefaults(defineProps<Props>(), {
  variant: 'default',
  size: 'default',
  placeholder: 'Select an option...',
})

const model = defineModel<string | number>()
const open = ref(false)
const highlightedIndex = ref(-1)

const normalizedOptions = computed<Option[]>(() =>
  (props.options ?? []).map((opt) => (typeof opt === 'string' ? { label: opt, value: opt } : opt)),
)

const selectedOption = computed(() => normalizedOptions.value.find((opt) => opt.value === model.value))

const selectOption = (option: Option) => {
  if (option.disabled) return
  model.value = option.value
  open.value = false
}

// Attached to the trigger button only, so it never steals keys from other inputs
const onKeyDown = (e: KeyboardEvent) => {
  if (props.disabled) return
  const count = normalizedOptions.value.length
  if (!count) return

  if (!open.value) {
    if (e.key === 'Enter' || e.key === ' ' || e.key === 'ArrowDown') {
      e.preventDefault()
      open.value = true
      const i = normalizedOptions.value.findIndex((o) => o.value === model.value)
      highlightedIndex.value = i === -1 ? 0 : i
    }
    return
  }

  if (e.key === 'ArrowDown') {
    e.preventDefault()
    highlightedIndex.value = (highlightedIndex.value + 1) % count
  } else if (e.key === 'ArrowUp') {
    e.preventDefault()
    highlightedIndex.value = (highlightedIndex.value - 1 + count) % count
  } else if (e.key === 'Enter' || e.key === ' ') {
    e.preventDefault()
    const opt = normalizedOptions.value[highlightedIndex.value]
    if (opt) selectOption(opt)
  } else if (e.key === 'Tab') {
    open.value = false
  }
}
</script>

<template>
  <Dropdown v-model:open="open" align="start" root-class="block w-full" class="w-full min-w-full">
    <template #trigger="{ toggle }">
      <button
        type="button"
        :disabled="props.disabled"
        :class="cn(selectVariants({ variant: props.variant, size: props.size }), props.class)"
        aria-haspopup="listbox"
        :aria-expanded="open"
        @click="toggle"
        @keydown="onKeyDown"
      >
        <span :class="!selectedOption && 'text-slate-500'">
          {{ selectedOption ? selectedOption.label : props.placeholder }}
        </span>
        <svg
          xmlns="http://www.w3.org/2000/svg"
          viewBox="0 0 24 24"
          fill="none"
          stroke="currentColor"
          stroke-width="2"
          stroke-linecap="round"
          stroke-linejoin="round"
          :class="cn('size-4 text-slate-400 transition-transform duration-200', open && 'rotate-180')"
        >
          <path d="m6 9 6 6 6-6" />
        </svg>
      </button>
    </template>

    <div class="flex flex-col py-1" role="listbox">
      <button
        v-for="(opt, index) in normalizedOptions"
        :key="opt.value"
        type="button"
        role="option"
        tabindex="-1"
        :aria-selected="model === opt.value"
        :disabled="opt.disabled"
        :class="cn(
          'flex w-full items-center px-3 py-2 text-left text-sm transition-colors',
          opt.disabled ? 'cursor-not-allowed opacity-50' : 'cursor-pointer',
          model === opt.value ? 'bg-green-500/10 font-medium text-green-400' : 'text-gray-200 hover:bg-slate-700/50',
          highlightedIndex === index && 'bg-slate-700/50',
        )"
        @click="selectOption(opt)"
        @mouseenter="highlightedIndex = index"
      >
        <span class="flex-1">{{ opt.label }}</span>
        <svg
          v-if="model === opt.value"
          xmlns="http://www.w3.org/2000/svg"
          viewBox="0 0 24 24"
          fill="none"
          stroke="currentColor"
          stroke-width="2"
          stroke-linecap="round"
          stroke-linejoin="round"
          class="size-4"
        >
          <path d="M20 6 9 17l-5-5" />
        </svg>
      </button>
    </div>
  </Dropdown>
</template>