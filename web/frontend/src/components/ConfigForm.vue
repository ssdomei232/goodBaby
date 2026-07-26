<script setup lang="ts">
// 根据驱动元数据(MetaField[])动态渲染配置表单。
// 通过 v-model 双向绑定一个 JSON 字符串，内部维护解析后的对象。
import { ref, watch } from 'vue'
import type { MetaField } from '@/api/types'

const props = defineProps<{
  fields: MetaField[]
  modelValue: string
}>()

const emit = defineEmits<{
  (e: 'update:modelValue', value: string): void
}>()

type ConfigValue = Record<string, unknown>

const form = ref<ConfigValue>({})

/** 字符串数组/数字数组字段在 UI 中用多行文本编辑，每行一项 */
const listDrafts = ref<Record<string, string>>({})

function defaultFor(field: MetaField): unknown {
  if (field.default !== undefined && field.default !== null) return field.default
  switch (field.type) {
    case 'number':
      return undefined
    case 'bool':
      return false
    case 'string-list':
    case 'number-list':
      return []
    default:
      return ''
  }
}

function parseIncoming(raw: string) {
  let parsed: ConfigValue = {}
  if (raw) {
    try {
      parsed = JSON.parse(raw) as ConfigValue
    } catch {
      parsed = {}
    }
  }

  const next: ConfigValue = {}
  const drafts: Record<string, string> = {}
  for (const field of props.fields) {
    const value = parsed[field.key] ?? defaultFor(field)
    next[field.key] = value
    if (field.type === 'string-list' || field.type === 'number-list') {
      drafts[field.key] = Array.isArray(value) ? value.join('\n') : ''
    }
  }
  form.value = next
  listDrafts.value = drafts
}

function serialize(): string {
  const out: ConfigValue = {}
  for (const field of props.fields) {
    let value = form.value[field.key]
    if (field.type === 'string-list') {
      value = (listDrafts.value[field.key] ?? '')
        .split('\n')
        .map((s) => s.trim())
        .filter((s) => s !== '')
    } else if (field.type === 'number-list') {
      value = (listDrafts.value[field.key] ?? '')
        .split('\n')
        .map((s) => s.trim())
        .filter((s) => s !== '')
        .map((s) => Number(s))
        .filter((n) => !Number.isNaN(n))
    } else if (field.type === 'number') {
      if (value === '' || value === undefined || value === null) continue
      value = Number(value)
    }
    if (value === '' || value === undefined) continue
    out[field.key] = value
  }
  return JSON.stringify(out)
}

let syncing = false

watch(
  () => [props.modelValue, props.fields] as const,
  () => {
    if (syncing) return
    parseIncoming(props.modelValue)
  },
  { immediate: true, deep: false },
)

watch(
  [form, listDrafts],
  () => {
    syncing = true
    emit('update:modelValue', serialize())
    // 下一轮微任务再解除，避免回环触发 parseIncoming
    queueMicrotask(() => {
      syncing = false
    })
  },
  { deep: true },
)
</script>

<template>
  <div class="config-form">
    <el-form-item
      v-for="field in fields"
      :key="field.key"
      :label="field.label"
      :required="field.required"
    >
      <!-- 单行文本 -->
      <el-input
        v-if="field.type === 'string'"
        v-model="form[field.key] as string"
        :placeholder="field.placeholder"
        clearable
      />

      <!-- 密码 -->
      <el-input
        v-else-if="field.type === 'password'"
        v-model="form[field.key] as string"
        type="password"
        show-password
        :placeholder="field.placeholder"
      />

      <!-- 多行文本 -->
      <el-input
        v-else-if="field.type === 'textarea'"
        v-model="form[field.key] as string"
        type="textarea"
        :rows="4"
        :placeholder="field.placeholder"
      />

      <!-- 数字 -->
      <el-input-number
        v-else-if="field.type === 'number'"
        v-model="form[field.key] as number"
        :min="0"
        :controls="false"
        :placeholder="field.placeholder"
        style="width: 100%"
      />

      <!-- 开关 -->
      <el-switch v-else-if="field.type === 'bool'" v-model="form[field.key] as boolean" />

      <!-- 列表：每行一项 -->
      <el-input
        v-else
        v-model="listDrafts[field.key]"
        type="textarea"
        :rows="3"
        :placeholder="field.placeholder ? `${field.placeholder}\n(每行一项)` : '每行一项'"
      />

      <div v-if="field.help" class="muted field-help">{{ field.help }}</div>
    </el-form-item>
  </div>
</template>

<style scoped>
.field-help {
  line-height: 1.5;
  margin-top: 4px;
  width: 100%;
}
</style>
