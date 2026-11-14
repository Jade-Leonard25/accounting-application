<script setup lang="ts">
import PayslipEditor from '../components/PayslipEditor.vue'
</script>

<template>
  <div class="max-w-4xl">
    <h1 class="text-2xl font-bold text-white mb-2">Payslip Template</h1>
    <p class="text-slate-400 text-sm mb-6">Design the template that will be used when printing payslips.</p>
    
    <PayslipEditor :selected-employee="null" />
  </div>
</template>