<script setup lang="ts">
import type { HTMLAttributes } from 'vue'
import { clsx, type ClassValue } from 'clsx'
import { twMerge } from 'tailwind-merge'
import { cva, type VariantProps } from 'class-variance-authority'

const cn = (...inputs: ClassValue[]) => twMerge(clsx(inputs))

const buttonVariants = cva(
  'inline-flex items-center justify-center gap-2 whitespace-nowrap rounded font-semibold transition-colors cursor-pointer focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-green-400 focus-visible:ring-offset-2 focus-visible:ring-offset-slate-900 disabled:pointer-events-none disabled:opacity-50 cursor-pointer',
  {
    variants: {
      variant: {
        default: 'bg-green-500 text-slate-900 hover:bg-green-400',
        destructive: 'bg-red-500 text-white hover:bg-red-400',
        outline: 'border border-slate-600 bg-transparent text-gray-200 hover:bg-slate-800',
        secondary: 'bg-slate-700 text-gray-200 hover:bg-slate-600',
        ghost: 'bg-transparent text-gray-200 hover:bg-slate-800',
        link: 'text-green-400 underline-offset-4 hover:underline',
      },
      size: {
        default: 'h-10 px-4 py-2 text-sm',
        sm: 'h-8 px-3 text-xs',
        lg: 'h-12 px-6 text-base',
        icon: 'h-10 w-10',
      },
    },
    defaultVariants: {
      variant: 'default',
      size: 'default',
    },
  },
)

type ButtonVariants = VariantProps<typeof buttonVariants>

interface Props {
  variant?: ButtonVariants['variant']
  size?: ButtonVariants['size']
  as?: string
  type?: 'button' | 'submit' | 'reset'
  disabled?: boolean
  loading?: boolean
  class?: HTMLAttributes['class']
}

const props = withDefaults(defineProps<Props>(), {
  as: 'button',
  type: 'button',
})
</script>

<template>
  <component
    :is="props.as"
    :type="props.as === 'button' ? props.type : undefined"
    :disabled="props.disabled || props.loading"
    :class="cn(buttonVariants({ variant: props.variant, size: props.size }), props.class)"
  >
    <span
      v-if="props.loading"
      class="size-4 animate-spin rounded-full border-2 border-current border-t-transparent"
    />
    <slot />
  </component>
</template>