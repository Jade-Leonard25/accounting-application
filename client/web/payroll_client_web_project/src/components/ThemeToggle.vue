<script setup lang="ts">
import { computed } from 'vue'
import Dropdown from './Dropdown.vue'
import Button from './Button.vue'
import { useTheme, type ThemeMode } from '../lib/theme'

const { mode, setMode } = useTheme()

const options: { value: ThemeMode; label: string; icon: string }[] = [
  { value: 'light', label: 'Light', icon: '☀️' },
  { value: 'system', label: 'System', icon: '🖥️' },
  { value: 'dark', label: 'Dark', icon: '🌙' },
]

const current = computed(() => options.find((o) => o.value === mode.value)!)
</script>

<template>
  <Dropdown align="end">
    <template #trigger="{ toggle, open }">
      <Button variant="ghost" size="sm" :aria-expanded="open" @click="toggle">
        {{ current.icon }} {{ current.label }}
      </Button>
    </template>

    <template #default="{ close }">
      <button
        v-for="o in options"
        :key="o.value"
        role="menuitemradio"
        :aria-checked="mode === o.value"
        class="flex w-full cursor-pointer items-center justify-between gap-3 rounded px-3 py-2 text-left hover:bg-slate-700"
        @click="setMode(o.value); close()"
      >
        <span>{{ o.icon }} {{ o.label }}</span>
        <span v-if="mode === o.value" class="text-green-400">✓</span>
      </button>
    </template>
  </Dropdown>
</template>