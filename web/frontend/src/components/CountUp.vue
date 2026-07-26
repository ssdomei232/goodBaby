<script setup lang="ts">
import { onUnmounted, ref, watch } from 'vue'

const props = withDefaults(defineProps<{ value?: number | null; duration?: number }>(), {
  duration: 700,
})

const display = ref(0)
let frame: number | undefined
let fallback: number | undefined

/** 缓出曲线，末尾减速，观感比线性自然 */
const easeOut = (t: number) => 1 - Math.pow(1 - t, 3)

function clearTimers() {
  if (frame !== undefined) cancelAnimationFrame(frame)
  if (fallback !== undefined) clearTimeout(fallback)
  frame = undefined
  fallback = undefined
}

function animateTo(target: number) {
  clearTimers()

  // 这些情况下不做动画，直接显示真实数字：
  // 系统要求减少动效、页面在后台(rAF 会被冻结)、环境没有 rAF
  const reduceMotion = window.matchMedia('(prefers-reduced-motion: reduce)').matches
  if (reduceMotion || document.hidden || typeof requestAnimationFrame !== 'function') {
    display.value = target
    return
  }

  const from = display.value
  const delta = target - from
  if (delta === 0) return

  const start = performance.now()
  const step = (now: number) => {
    const progress = Math.min(1, (now - start) / props.duration)
    display.value = Math.round(from + delta * easeOut(progress))
    if (progress < 1) {
      frame = requestAnimationFrame(step)
    } else {
      clearTimers()
    }
  }
  frame = requestAnimationFrame(step)

  // 安全网：rAF 被节流或根本不回调时，兜底把真实值显示出来。
  // 数字展示的正确性优先于动画效果。
  fallback = window.setTimeout(() => {
    display.value = target
    clearTimers()
  }, props.duration + 300)
}

watch(
  () => props.value,
  (v) => animateTo(typeof v === 'number' ? v : 0),
  { immediate: true },
)

onUnmounted(clearTimers)
</script>

<template>
  <span>{{ props.value === null || props.value === undefined ? '-' : display }}</span>
</template>
