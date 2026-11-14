<script setup lang="ts">
import { computed, nextTick, onBeforeUnmount, onMounted, ref, watch } from 'vue'
import type { HTMLAttributes } from 'vue'
import { motion, AnimatePresence, MotionConfig } from 'motion-v'
import { clsx, type ClassValue } from 'clsx'
import { twMerge } from 'tailwind-merge'
import { cva, type VariantProps } from 'class-variance-authority'
import Button from './Button.vue'

const cn = (...inputs: ClassValue[]) => twMerge(clsx(inputs))

const drawerVariants = cva(
  'fixed z-50 flex flex-col border-slate-700 bg-slate-900 text-gray-200 shadow-2xl outline-none',
  {
    variants: {
      side: {
        right: 'right-0 top-0 h-full border-l',
        left: 'left-0 top-0 h-full border-r',
        top: 'left-0 top-0 w-full border-b rounded-b-xl',
        bottom: 'bottom-0 left-0 w-full border-t rounded-t-xl',
      },
      size: {
        sm: '',
        md: '',
        lg: '',
        full: '',
      },
    },
    compoundVariants: [
      { side: ['left', 'right'], size: 'sm', class: 'w-full max-w-xs' },
      { side: ['left', 'right'], size: 'md', class: 'w-full max-w-md' },
      { side: ['left', 'right'], size: 'lg', class: 'w-full max-w-2xl' },
      { side: ['left', 'right'], size: 'full', class: 'w-full' },
      { side: ['top', 'bottom'], size: 'sm', class: 'max-h-[30vh]' },
      { side: ['top', 'bottom'], size: 'md', class: 'max-h-[50vh]' },
      { side: ['top', 'bottom'], size: 'lg', class: 'max-h-[80vh]' },
      { side: ['top', 'bottom'], size: 'full', class: 'h-full rounded-none' },
    ],
    defaultVariants: { side: 'bottom', size: 'lg' },
  },
)

type DrawerVariants = VariantProps<typeof drawerVariants>

interface Props {
  side?: NonNullable<DrawerVariants['side']>
  size?: DrawerVariants['size']
  closeOnOverlay?: boolean
  class?: HTMLAttributes['class']
}

const props = withDefaults(defineProps<Props>(), {
  side: 'bottom',
  size: 'lg',
  closeOnOverlay: true,
})

const open = defineModel<boolean>('open', { default: false })
const panel = ref<HTMLElement | null>(null)

const close = () => (open.value = false)

// Where the panel starts and ends its slide
const offscreen = computed(() => {
  switch (props.side) {
    case 'left': return { x: '-100%' }
    case 'top': return { y: '-100%' }
    case 'bottom': return { y: '100%' }
    default: return { x: '100%' }
  }
})
const onscreen = computed(() => ({ x: 0, y: 0 }))

const onKeyDown = (e: KeyboardEvent) => {
  if (open.value && e.key === 'Escape') close()
}

onMounted(() => document.addEventListener('keydown', onKeyDown))
onBeforeUnmount(() => {
  document.removeEventListener('keydown', onKeyDown)
  document.body.style.overflow = ''
})

// Lock page scroll while open, and move focus into the panel
watch(open, async (isOpen) => {
  document.body.style.overflow = isOpen ? 'hidden' : ''
  if (isOpen) {
    await nextTick()
    panel.value?.focus()
  }
})
</script>

<template>
  <Teleport to="body">
    <MotionConfig reduced-motion="user">
      <AnimatePresence>
        <motion.div
          v-if="open"
          key="drawer-overlay"
          class="fixed inset-0 z-40 bg-black/60"
          :initial="{ opacity: 0 }"
          :animate="{ opacity: 1 }"
          :exit="{ opacity: 0 }"
          :transition="{ duration: 0.2 }"
          aria-hidden="true"
          @click="closeOnOverlay && close()"
        />

        <motion.div
          v-if="open"
          key="drawer-panel"
          ref="panel"
          role="dialog"
          aria-modal="true"
          tabindex="-1"
          :initial="offscreen"
          :animate="onscreen"
          :exit="offscreen"
          :transition="{ type: 'spring', stiffness: 380, damping: 36 }"
          :class="cn(drawerVariants({ side: props.side, size: props.size }), props.class)"
        >
          <div
            v-if="$slots.title || $slots.description"
            class="flex items-start justify-between gap-4 border-b border-slate-700 p-6"
          >
            <div class="flex flex-col gap-1.5">
              <h2 v-if="$slots.title" class="text-xl font-bold leading-none tracking-tight text-white">
                <slot name="title" />
              </h2>
              <p v-if="$slots.description" class="text-sm text-slate-400">
                <slot name="description" />
              </p>
            </div>
            <Button variant="ghost" size="icon" aria-label="Close" class="-mr-2 -mt-2 shrink-0" @click="close">
              ✕
            </Button>
          </div>

          <div class="flex-1 overflow-y-auto p-6">
            <slot :close="close" />
          </div>

          <div v-if="$slots.footer" class="flex items-center justify-end gap-2 border-t border-slate-700 p-6">
            <slot name="footer" :close="close" />
          </div>
        </motion.div>
      </AnimatePresence>
    </MotionConfig>
  </Teleport>
</template>