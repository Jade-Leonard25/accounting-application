<script setup lang="ts">
import { ref, onMounted, onBeforeUnmount } from 'vue'
import type { HTMLAttributes } from 'vue'
import { motion, AnimatePresence, MotionConfig } from 'motion-v'
import { clsx, type ClassValue } from 'clsx'
import { twMerge } from 'tailwind-merge'
import { cva, type VariantProps } from 'class-variance-authority'

const cn = (...inputs: ClassValue[]) => twMerge(clsx(inputs))

const menuVariants = cva(
  'absolute z-50 min-w-48 rounded-md border border-slate-700 bg-slate-800 p-1 text-sm text-gray-200 shadow-lg',
  {
    variants: {
      align: {
        start: 'left-0 origin-top-left',
        end: 'right-0 origin-top-right',
      },
      side: {
        bottom: 'top-full mt-2',
        top: 'bottom-full mb-2',
      },
    },
    defaultVariants: {
      align: 'start',
      side: 'bottom',
    },
  },
)

type MenuVariants = VariantProps<typeof menuVariants>

interface Props {
  align?: MenuVariants['align']
  side?: MenuVariants['side']
  class?: HTMLAttributes['class']
}

const props = defineProps<Props>()

const open = defineModel<boolean>('open', { default: false })
const root = ref<HTMLElement | null>(null)

const toggle = () => (open.value = !open.value)
const close = () => (open.value = false)

const onPointerDown = (e: PointerEvent) => {
  if (open.value && root.value && !root.value.contains(e.target as Node)) close()
}
const onKeyDown = (e: KeyboardEvent) => {
  if (e.key === 'Escape') close()
}

onMounted(() => {
  document.addEventListener('pointerdown', onPointerDown)
  document.addEventListener('keydown', onKeyDown)
})
onBeforeUnmount(() => {
  document.removeEventListener('pointerdown', onPointerDown)
  document.removeEventListener('keydown', onKeyDown)
})
</script>

<template>
  <MotionConfig reduced-motion="user">
    <div ref="root" class="relative inline-block self-start">
      <slot name="trigger" :open="open" :toggle="toggle" />

      <AnimatePresence>
        <motion.div
          v-if="open"
          role="menu"
          :initial="{ opacity: 0, scale: 0.95, y: -6 }"
          :animate="{ opacity: 1, scale: 1, y: 0 }"
          :exit="{ opacity: 0, scale: 0.95, y: -6 }"
          :transition="{ type: 'spring', stiffness: 500, damping: 32 }"
          :class="cn(menuVariants({ align: props.align, side: props.side }), props.class)"
        >
          <slot :close="close" />
        </motion.div>
      </AnimatePresence>
    </div>
  </MotionConfig>
</template>