import { computed, ref, watchEffect } from 'vue'

export type ThemeMode = 'light' | 'system' | 'dark'

const KEY = 'payroll:theme'

const read = (): ThemeMode => {
  try {
    const v = localStorage.getItem(KEY)
    if (v === 'light' || v === 'dark' || v === 'system') return v
  } catch {
    /* storage unavailable */
  }
  return 'system'
}

const mode = ref<ThemeMode>(read())
const media = window.matchMedia('(prefers-color-scheme: dark)')
const systemDark = ref(media.matches)
media.addEventListener('change', (e) => (systemDark.value = e.matches))

const isDark = computed(() => mode.value === 'dark' || (mode.value === 'system' && systemDark.value))

watchEffect(() => {
  document.documentElement.classList.toggle('dark', isDark.value)
})

export function useTheme() {
  const setMode = (m: ThemeMode) => {
    mode.value = m
    try {
      localStorage.setItem(KEY, m)
    } catch {
      /* ignore */
    }
  }
  return { mode, isDark, setMode }
}