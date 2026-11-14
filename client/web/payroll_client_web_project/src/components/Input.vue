<script setup lang="ts">
import type { HTMLAttributes } from 'vue'
import { clsx, type ClassValue } from 'clsx'
import { twMerge } from 'tailwind-merge'
import { cva, type VariantProps } from 'class-variance-authority'

const cn = (...inputs: ClassValue[]) => twMerge(clsx(inputs))

const inputVariants = cva(
  'flex w-full rounded border bg-slate-800 text-gray-200 placeholder:text-slate-500 transition-colors focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-offset-2 focus-visible:ring-offset-slate-900 disabled:cursor-not-allowed disabled:opacity-50',
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
    defaultVariants: {
      variant: 'default',
      size: 'default',
    },
  },
)

type InputVariants = VariantProps<typeof inputVariants>

interface Props {
  variant?: InputVariants['variant']
  size?: InputVariants['size']
  type?: string
  disabled?: boolean
  class?: HTMLAttributes['class']
}

const props = withDefaults(defineProps<Props>(), {
  type: 'text',
})

const model = defineModel<string>()
</script>

<template>
  <input
    v-model="model"
    :type="props.type"
    :disabled="props.disabled"
    :class="cn(inputVariants({ variant: props.variant, size: props.size }), props.class)"
  />
</template>