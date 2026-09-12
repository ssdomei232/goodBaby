<script setup lang="ts">
/**
 * 妄想症主题的背景层：噪点、暗角、缓慢下移的扫描光带，
 * 以及像石蒜花瓣一样飘落的色片。
 *
 * 数值全部由下标推出来，刷新页面时不会跳到别的位置；
 * 小屏与「减少动效」偏好下会由 CSS 自动收敛。
 */
const PETAL_COUNT = 14

const petals = Array.from({ length: PETAL_COUNT }, (_, i) => {
  const r1 = ((i * 97) % 100) / 100
  const r2 = ((i * 53 + 17) % 100) / 100
  const r3 = ((i * 37 + 11) % 100) / 100

  return {
    id: i,
    x: `${((i * 7.3 + 3.5) % 100).toFixed(2)}%`,
    drift: `${Math.round((r1 - 0.5) * 170)}px`,
    dur: `${(18 + r2 * 16).toFixed(1)}s`,
    delay: `${(-i * 1.9).toFixed(1)}s`,
    scale: (0.55 + r3 * 1.05).toFixed(2),
    opacity: (0.32 + r1 * 0.36).toFixed(2),
    color: i % 5 === 0 ? 'var(--pa-blue)' : i % 7 === 0 ? 'var(--pa-cyan)' : 'var(--pa-crimson)',
  }
})
</script>

<template>
  <div class="pa-backdrop" aria-hidden="true">
    <div
      class="pa-orb"
      style="
        left: -6%;
        top: 12%;
        width: 42vw;
        height: 42vw;
        background: rgb(224 36 94 / 0.22);
      "
    />
    <div
      class="pa-orb"
      style="
        right: -8%;
        bottom: 4%;
        width: 38vw;
        height: 38vw;
        background: rgb(102 204 255 / 0.14);
        animation-delay: -6s;
      "
    />

    <div class="pa-backdrop__scan" />
    <div class="pa-backdrop__grain" />

    <span class="pa-ghost-text" style="right: 3%; top: 8%; font-size: 92px">妄想</span>
    <span class="pa-ghost-text" style="left: 1%; bottom: 6%; font-size: 58px">现实</span>

    <div class="pa-backdrop__petals">
      <span
        v-for="p in petals"
        :key="p.id"
        class="pa-petal"
        :style="{
          '--pa-x': p.x,
          '--pa-drift': p.drift,
          '--pa-dur': p.dur,
          '--pa-delay': p.delay,
          '--pa-scale': p.scale,
          '--pa-opacity': p.opacity,
        }"
      >
        <svg viewBox="0 0 24 24">
          <path
            d="M12 1.6c3.4 5.2 6.6 8.6 6.6 13.1 0 3.9-3 7.1-6.6 7.1S5.4 18.6 5.4 14.7C5.4 10.2 8.6 6.8 12 1.6z"
            :fill="p.color"
          />
          <path d="M12 3.6v16.4" stroke="rgb(0 0 0 / 0.22)" stroke-width="0.6" fill="none" />
        </svg>
      </span>
    </div>

    <div class="pa-backdrop__vignette" />
  </div>
</template>
