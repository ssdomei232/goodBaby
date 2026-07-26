<script setup lang="ts">
// 品牌标志：信号波纹，寓意「定期发出的平安信号」
defineProps<{ size?: number }>()
</script>

<template>
  <svg
    :width="size ?? 22"
    :height="size ?? 22"
    viewBox="0 0 24 24"
    fill="none"
    stroke="currentColor"
    stroke-width="2"
    stroke-linecap="round"
  >
    <circle cx="12" cy="12" r="2.2" fill="currentColor" stroke="none" />
    <path d="M8.6 15.4a4.8 4.8 0 0 1 0-6.8" />
    <path d="M15.4 8.6a4.8 4.8 0 0 1 0 6.8" />
    <path d="M5.8 18.2a8.8 8.8 0 0 1 0-12.4" />
    <path d="M18.2 5.8a8.8 8.8 0 0 1 0 12.4" />
  </svg>
</template>
