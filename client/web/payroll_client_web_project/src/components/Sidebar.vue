<script setup lang="ts">
import { ref } from 'vue'
import type { HTMLAttributes } from 'vue'
import { clsx, type ClassValue } from 'clsx'
import { twMerge } from 'tailwind-merge'
import { cva } from 'class-variance-authority'
import Button from './Button.vue'
const cn = (...inputs: ClassValue[]) => twMerge(clsx(inputs))

const itemVariants = cva(
  'flex w-full cursor-pointer items-center justify-between rounded px-3 py-2 text-left text-sm font-medium transition-colors focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-green-400',
  {
    variants: {
      active: {
        true: 'bg-slate-800 text-white',
        false: 'text-slate-400 hover:bg-slate-800/60 hover:text-white',
      },
    },
    defaultVariants: { active: false },
  },
)

interface Item {
  key: string
  label: string
  badge?: number
}

const props = defineProps<{
  items: Item[]
  class?: HTMLAttributes['class']
}>()

const active = defineModel<string>({ required: true })

// State for minimize/close
const isCollapsed = ref(false)
const isClosed = ref(false)
</script>

<template>
  <!-- If closed, show a small button to reopen it -->
  <div v-if="isClosed" class="p-3 bg-slate-950 border-r-2 border-slate-700">
    <Button
      @click="isClosed = false"
      class="p-2 text-slate-400 hover:text-white transition-colors rounded hover:bg-slate-800"
      title="Open Sidebar"
    >
      <svg xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="m9 18 6-6-6-6"/></svg>
    </Button>
  </div>

  <!-- Sidebar -->
  <aside 
    v-else
    :class="cn(
      'flex shrink-0 flex-col gap-1 border-r-2 border-slate-700 bg-slate-950 p-3 transition-all duration-300 ease-in-out', 
      isCollapsed ? 'w-16' : 'w-56',
      props.class
    )"
  >
    <!-- Header with Minimize and Close -->
    <div :class="cn('flex items-center mb-2 pb-2 border-b border-slate-800', isCollapsed ? 'justify-center' : 'justify-between')">
      <span v-if="!isCollapsed" class="text-xs font-bold tracking-widest text-slate-500 uppercase pl-1">
        Menu
      </span>
      <div class="flex items-center gap-1">
        <!-- Minimize Button -->
        <Button
          @click="isCollapsed = !isCollapsed"
          class="p-1.5 text-slate-400 hover:text-white transition-colors rounded hover:bg-slate-800"
          :title="isCollapsed ? 'Expand Sidebar' : 'Minimize Sidebar'"
        >
          <svg v-if="!isCollapsed" xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="m15 18-6-6 6-6"/></svg>
          <svg v-else xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="m9 18 6-6-6-6"/></svg>
        </Button>
        
        <!-- Close Button (only show when expanded) -->
      
      </div>
    </div>

    <!-- Navigation Items -->
    <Button
      v-for="item in items"
      :key="item.key"
      type="button"
      :aria-current="active === item.key ? 'page' : undefined"
      :class="itemVariants({ active: active === item.key })"
      @click="active = item.key"
      :title="isCollapsed ? item.label : undefined"
    >
      <!-- If collapsed, just show the first letter as an avatar/icon placeholder -->
      <span v-if="isCollapsed" class="mx-auto flex size-6 items-center justify-center rounded bg-slate-800 text-xs">
        {{ item.label.charAt(0) }}
      </span>

      <!-- If expanded, show the full text -->
      <template v-else>
        <span>{{ item.label }}</span>
        <span
          v-if="item.badge"
          class="rounded-full bg-green-500 px-2 py-0.5 text-xs font-semibold text-slate-950"
        >
          {{ item.badge }}
        </span>
      </template>
    </Button>

    <!-- Bottom-left: account area -->
    <div v-if="$slots.footer" :class="cn('mt-auto flex flex-col gap-3 border-t border-slate-700 pt-3', isCollapsed && 'items-center')">
      <slot name="footer" />
    </div>
  </aside>
</template>