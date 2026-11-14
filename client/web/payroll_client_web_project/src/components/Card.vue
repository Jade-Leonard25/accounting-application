<script setup lang="ts">
import type { HTMLAttributes } from 'vue'
import { clsx, type ClassValue } from 'clsx'
import { twMerge } from 'tailwind-merge'
import { cva, type VariantProps } from 'class-variance-authority'

const cn = (...inputs: ClassValue[]) => twMerge(clsx(inputs))

const cardVariants = cva('rounded-lg border text-gray-200 shadow-sm', {
  variants: {
    variant: {
      default: 'bg-slate-900 border-slate-700',
      muted: 'bg-slate-800 border-slate-700',
      outline: 'bg-transparent border-slate-600',
      destructive: 'bg-slate-900 border-red-400',
    },
    padding: {
      none: '',
      sm: 'p-3',
      default: 'p-6',
      lg: 'p-8',
    },
  },
  defaultVariants: {
    variant: 'default',
    padding: 'default',
  },
})

type CardVariants = VariantProps<typeof cardVariants>

interface Props {
  variant?: CardVariants['variant']
  padding?: CardVariants['padding']
  as?: string
  class?: HTMLAttributes['class']
}

const props = withDefaults(defineProps<Props>(), {
  as: 'div',
})
</script>

<template>
  <component
    :is="props.as"
    :class="cn(cardVariants({ variant: props.variant, padding: props.padding }), 'flex flex-col gap-4', props.class)"
  >
    <div v-if="$slots.title || $slots.description" class="flex flex-col gap-1.5">
      <h3 v-if="$slots.title" class="text-2xl font-bold leading-none tracking-tight">
        <slot name="title" />
      </h3>
      <p v-if="$slots.description" class="text-sm text-slate-400">
        <slot name="description" />
      </p>
    </div>

    <div v-if="$slots.default">
      <slot />
    </div>

    <div v-if="$slots.footer" class="flex items-center gap-2">
      <slot name="footer" />
    </div>
  </component>
</template>