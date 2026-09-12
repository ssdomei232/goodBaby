<script setup lang="ts">
import { computed, onBeforeUnmount, onMounted, ref } from 'vue'
import { Check, MagicStick } from '@element-plus/icons-vue'
import { skinOptions, useTheme, type ThemeSkin } from '@/composables/useTheme'

/**
 * 外观主题选择器。
 *
 * 用自绘的浮层而不是 el-popover：这样预览小样可以完全跟着皮肤走，
 * 也不必和 teleport 出来的 popper 抢样式。
 */
const { skin, isDark, setSkin } = useTheme()

const open = ref(false)
const root = ref<HTMLElement | null>(null)

const currentName = computed(() => skinOptions.find((o) => o.id === skin.value)?.name ?? '')

/** 每个皮肤的迷你预览配色，跟随当前明暗模式 */
function previewOf(id: ThemeSkin) {
  const paranoia = id === 'paranoia'
  if (isDark.value) {
    return paranoia
      ? { bg: '#0d0a11', side: '#1d0f1b', card: '#150f19', accent: '#e0245e', line: '#3a2130' }
      : { bg: '#101216', side: '#17191f', card: '#1a1d23', accent: '#66ccff', line: '#2b303a' }
  }
  return paranoia
    ? { bg: '#f5f1ee', side: '#241a22', card: '#fffdfc', accent: '#b3163c', line: '#e4d3d8' }
    : { bg: '#f4f6f8', side: '#232630', card: '#ffffff', accent: '#3fa9e0', line: '#dde3e9' }
}

async function pick(id: ThemeSkin, event: MouseEvent) {
  open.value = false
  if (id === skin.value) return
  await setSkin(id, event)
}

function onDocumentClick(event: MouseEvent) {
  if (root.value && !root.value.contains(event.target as Node)) {
    open.value = false
  }
}

function onKeydown(event: KeyboardEvent) {
  if (event.key === 'Escape') open.value = false
}

onMounted(() => {
  document.addEventListener('click', onDocumentClick, true)
  document.addEventListener('keydown', onKeydown)
})

onBeforeUnmount(() => {
  document.removeEventListener('click', onDocumentClick, true)
  document.removeEventListener('keydown', onKeydown)
})
</script>

<template>
  <div ref="root" class="skin-picker">
    <button
      class="icon-btn skin-btn"
      type="button"
      :aria-expanded="open"
      aria-haspopup="dialog"
      :title="`外观主题：${currentName}`"
      @click.stop="open = !open"
    >
      <el-icon :size="17"><MagicStick /></el-icon>
      <span class="skin-btn__pip" :class="skin" />
    </button>

    <Transition name="skin-pop">
      <div v-if="open" class="skin-panel" role="dialog" aria-label="外观主题">
        <div class="skin-panel__head">
          <span class="skin-panel__title">外观主题</span>
          <span class="skin-panel__hint">Appearance</span>
        </div>

        <button
          v-for="option in skinOptions"
          :key="option.id"
          type="button"
          class="skin-option"
          :class="{ 'is-active': option.id === skin }"
          @click.stop="pick(option.id, $event)"
        >
          <span
            class="skin-option__preview"
            :style="{
              background: previewOf(option.id).bg,
              borderColor: previewOf(option.id).line,
            }"
          >
            <i class="skin-option__side" :style="{ background: previewOf(option.id).side }" />
            <i class="skin-option__bar" :style="{ background: previewOf(option.id).accent }" />
            <i
              class="skin-option__card"
              :style="{
                background: previewOf(option.id).card,
                borderColor: previewOf(option.id).line,
              }"
            />
            <i
              class="skin-option__card skin-option__card--alt"
              :style="{
                background: previewOf(option.id).card,
                borderColor: previewOf(option.id).line,
              }"
            />
          </span>

          <span class="skin-option__meta">
            <span class="skin-option__name">
              {{ option.name }}
              <el-icon v-if="option.id === skin" :size="14" class="skin-option__check">
                <Check />
              </el-icon>
            </span>
            <span class="skin-option__sub">{{ option.sub }}</span>
            <span class="skin-option__swatches">
              <i v-for="color in option.swatches" :key="color" :style="{ background: color }" />
            </span>
          </span>
        </button>

        <div class="skin-panel__foot">
          <span>暗色模式由顶栏的月亮按钮切换</span>
        </div>
      </div>
    </Transition>
  </div>
</template>

<style scoped>
.skin-picker {
  position: relative;
  display: flex;
}

.skin-btn {
  position: relative;
  width: 34px;
  height: 34px;
  padding: 0;
  border: none;
  border-radius: 10px;
  background: transparent;
  color: var(--el-text-color-regular);
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  transition: background 0.2s var(--gb-ease);
}

.skin-btn:hover {
  background: var(--el-fill-color);
}

.skin-btn__pip {
  position: absolute;
  right: 5px;
  bottom: 5px;
  width: 7px;
  height: 7px;
  border-radius: 2px;
  background: var(--gb-primary);
  box-shadow: 0 0 0 2px var(--gb-card);
}

.skin-btn__pip.paranoia {
  background: #e0245e;
}

/* ---------- 浮层 ---------- */

.skin-panel {
  position: absolute;
  top: calc(100% + 10px);
  right: 0;
  z-index: 1600;
  width: 296px;
  padding: 14px;
  box-sizing: border-box;
  background: var(--el-bg-color-overlay);
  border: 1px solid var(--el-border-color-light);
  border-radius: 12px;
  box-shadow: 0 18px 44px rgb(16 12 24 / 0.22);
}

html.paranoia .skin-panel {
  border-radius: 6px;
  border-color: var(--pa-line);
  box-shadow:
    0 18px 44px rgb(0 0 0 / 0.4),
    0 0 0 1px rgb(224 36 94 / 0.12);
}

.skin-panel__head {
  display: flex;
  align-items: baseline;
  justify-content: space-between;
  gap: 8px;
  padding: 0 2px 10px;
}

.skin-panel__title {
  font-size: 14px;
  font-weight: 700;
  color: var(--el-text-color-primary);
}

html.paranoia .skin-panel__title {
  font-family: var(--pa-serif);
  letter-spacing: 0.06em;
}

.skin-panel__hint {
  font-family: ui-monospace, SFMono-Regular, Consolas, monospace;
  font-size: 10px;
  letter-spacing: 0.16em;
  text-transform: uppercase;
  color: var(--el-text-color-secondary);
}

.skin-option {
  display: flex;
  align-items: center;
  gap: 12px;
  width: 100%;
  padding: 9px;
  margin-bottom: 6px;
  border: 1px solid transparent;
  border-radius: 9px;
  background: transparent;
  text-align: left;
  cursor: pointer;
  transition:
    background 0.2s var(--gb-ease),
    border-color 0.2s var(--gb-ease);
}

html.paranoia .skin-option {
  border-radius: 4px;
}

.skin-option:hover {
  background: var(--el-fill-color-light);
}

.skin-option.is-active {
  border-color: var(--gb-primary);
  background: var(--el-fill-color-light);
}

/* ---------- 迷你预览 ---------- */

.skin-option__preview {
  position: relative;
  flex: none;
  width: 62px;
  height: 44px;
  border: 1px solid;
  border-radius: 5px;
  overflow: hidden;
}

html.paranoia .skin-option__preview {
  border-radius: 3px;
}

.skin-option__side {
  position: absolute;
  left: 0;
  top: 0;
  bottom: 0;
  width: 15px;
}

.skin-option__bar {
  position: absolute;
  left: 15px;
  right: 0;
  top: 0;
  height: 7px;
  opacity: 0.85;
}

.skin-option__card {
  position: absolute;
  left: 22px;
  top: 13px;
  width: 32px;
  height: 11px;
  border: 1px solid;
  border-radius: 2px;
}

.skin-option__card--alt {
  top: 28px;
  width: 22px;
}

.skin-option__meta {
  display: flex;
  flex-direction: column;
  gap: 3px;
  min-width: 0;
}

.skin-option__name {
  display: flex;
  align-items: center;
  gap: 5px;
  font-size: 13px;
  font-weight: 600;
  color: var(--el-text-color-primary);
}

.skin-option__check {
  color: var(--gb-primary);
}

.skin-option__sub {
  font-size: 11px;
  color: var(--el-text-color-secondary);
}

.skin-option__swatches {
  display: flex;
  gap: 4px;
  margin-top: 3px;
}

.skin-option__swatches i {
  width: 14px;
  height: 4px;
  border-radius: 2px;
}

.skin-panel__foot {
  padding: 6px 2px 0;
  border-top: 1px solid var(--el-border-color-extra-light);
  margin-top: 4px;
  font-size: 11px;
  line-height: 1.7;
  color: var(--el-text-color-secondary);
}

/* ---------- 浮层进出场 ---------- */

.skin-pop-enter-active,
.skin-pop-leave-active {
  transition:
    opacity 0.18s var(--gb-ease),
    transform 0.22s var(--gb-ease);
}

.skin-pop-enter-from,
.skin-pop-leave-to {
  opacity: 0;
  transform: translateY(-6px) scale(0.97);
}

@media (max-width: 480px) {
  .skin-panel {
    width: min(84vw, 296px);
  }
}
</style>
