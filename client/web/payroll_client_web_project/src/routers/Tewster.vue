<script setup lang="ts">
import Dropdown from '../components/Dropdown.vue'
import Button from '../components/Button.vue'
</script>

<template>
  <div class="flex min-h-screen items-start justify-end bg-slate-900 p-6 text-gray-200">
    <Dropdown align="end">
      <template #trigger="{ toggle, open }">
        <Button variant="ghost" size="sm" :aria-expanded="open" @click="toggle">Menu ▾</Button>
      </template>

      <template #default="{ close }">
        <button
          role="menuitem"
          class="w-full cursor-pointer rounded px-3 py-2 text-left hover:bg-slate-700"
          @click="close()"
        >
          Profile
        </button>
      </template>
    </Dropdown>
  </div>
</template>