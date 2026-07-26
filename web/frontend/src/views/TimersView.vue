<script setup lang="ts">
import { computed, onMounted, onUnmounted, ref } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Check, Plus, VideoPlay } from '@element-plus/icons-vue'
import { timerApi } from '@/api'
import { ApiError } from '@/api/client'
import type { Timer, TimerRequest } from '@/api/types'
import { durationPresets, formatDateTime, formatDuration } from '@/utils/format'

const timers = ref<Timer[]>([])
const loading = ref(false)

const nowSeconds = ref(Math.floor(Date.now() / 1000))
let ticker: number | undefined

// ---- 创建 / 编辑对话框 ----
const dialogVisible = ref(false)
const saving = ref(false)
const editingId = ref<number | null>(null)

interface TimerForm {
  name: string
  description: string
  enabled: boolean
  durationDays: number
  remindHours: number
}

const emptyForm = (): TimerForm => ({
  name: '',
  description: '',
  enabled: true,
  durationDays: 7,
  remindHours: 24,
})

const form = ref<TimerForm>(emptyForm())

const dialogTitle = computed(() => (editingId.value === null ? '创建定时器' : '编辑定时器'))

function secondsLeft(timer: Timer): number {
  return timer.last_sign + timer.sign_deration_seconds - nowSeconds.value
}

/** 剩余时间占比(0~1)，驱动卡片进度条 */
function fractionLeft(timer: Timer): number {
  return Math.min(1, Math.max(0, secondsLeft(timer) / timer.sign_deration_seconds))
}

function barColor(timer: Timer): string {
  if (!timer.enabled) return 'var(--el-fill-color-dark)'
  const f = fractionLeft(timer)
  if (timer.triggered || f <= 0.1) return '#ef4444'
  if (f <= 0.3) return '#f59e0b'
  return 'var(--gb-primary)'
}

function statusOf(timer: Timer): { label: string; type: 'success' | 'warning' | 'danger' | 'info' } {
  if (!timer.enabled) return { label: '已暂停', type: 'info' }
  if (timer.triggered) return { label: '已触发', type: 'danger' }
  const left = secondsLeft(timer)
  if (left <= 0) return { label: '即将触发', type: 'danger' }
  if (left <= timer.remind_time_seconds) return { label: '临近到期', type: 'warning' }
  return { label: '正常', type: 'success' }
}

async function refresh() {
  loading.value = true
  try {
    timers.value = await timerApi.list()
  } catch (error) {
    ElMessage.error(error instanceof ApiError ? error.message : '加载失败')
  } finally {
    loading.value = false
  }
}

function openCreate() {
  editingId.value = null
  form.value = emptyForm()
  dialogVisible.value = true
}

function openEdit(timer: Timer) {
  editingId.value = timer.id
  form.value = {
    name: timer.name,
    description: timer.description,
    enabled: timer.enabled,
    durationDays: timer.sign_deration_seconds / 86400,
    remindHours: timer.remind_time_seconds / 3600,
  }
  dialogVisible.value = true
}

async function save() {
  const body: TimerRequest = {
    name: form.value.name,
    description: form.value.description,
    enabled: form.value.enabled,
    sign_deration_seconds: Math.round(form.value.durationDays * 86400),
    remind_time_seconds: Math.round(form.value.remindHours * 3600),
  }

  saving.value = true
  try {
    if (editingId.value === null) {
      await timerApi.create(body)
      ElMessage.success('定时器已创建，从现在开始计时')
    } else {
      await timerApi.update(editingId.value, body)
      ElMessage.success('定时器已更新')
    }
    dialogVisible.value = false
    await refresh()
  } catch (error) {
    ElMessage.error(error instanceof ApiError ? error.message : '保存失败')
  } finally {
    saving.value = false
  }
}

async function sign(timer: Timer) {
  try {
    await timerApi.sign(timer.id)
    ElMessage.success(`「${timer.name}」签到成功，重新计时`)
    await refresh()
  } catch (error) {
    ElMessage.error(error instanceof ApiError ? error.message : '签到失败')
  }
}

const triggeringId = ref<number | null>(null)

/** 手动触发：真实执行该定时器下所有启用的规则，用于调试 */
async function trigger(timer: Timer) {
  try {
    await ElMessageBox.confirm(
      '将立即执行该定时器下所有启用的规则（真实发送消息 / 修改仓库），不影响签到状态。确定继续吗？',
      `手动触发「${timer.name}」`,
      { type: 'warning', confirmButtonText: '执行', cancelButtonText: '取消' },
    )
  } catch {
    return
  }

  triggeringId.value = timer.id
  try {
    const result = await timerApi.trigger(timer.id)
    const failed = result.failed ?? []
    if (failed.length === 0) {
      ElMessage.success(`已执行 ${result.total} 条规则，全部成功`)
    } else {
      ElMessage.warning(`已执行 ${result.total} 条规则，${failed.length} 条失败，详情见执行日志`)
    }
  } catch (error) {
    ElMessage.error(error instanceof ApiError ? error.message : '触发失败')
  } finally {
    triggeringId.value = null
  }
}

async function toggleEnabled(timer: Timer) {
  try {
    await timerApi.update(timer.id, {
      name: timer.name,
      description: timer.description,
      enabled: timer.enabled,
      sign_deration_seconds: timer.sign_deration_seconds,
      remind_time_seconds: timer.remind_time_seconds,
    })
    ElMessage.success(timer.enabled ? '已启用' : '已暂停')
  } catch (error) {
    timer.enabled = !timer.enabled
    ElMessage.error(error instanceof ApiError ? error.message : '操作失败')
  }
}

async function remove(timer: Timer) {
  try {
    const affected = await timerApi.checkDelete(timer.id)
    const warning =
      affected.length > 0
        ? `删除后，关联的 ${affected.length} 条规则也会被一并删除：${affected
            .map((r) => r.name)
            .join('、')}`
        : '该定时器下没有关联规则。'

    await ElMessageBox.confirm(warning, `删除定时器「${timer.name}」？`, {
      type: 'warning',
      confirmButtonText: '删除',
      confirmButtonClass: 'el-button--danger',
      cancelButtonText: '取消',
    })
  } catch {
    return // 用户取消
  }

  try {
    await timerApi.remove(timer.id)
    ElMessage.success('已删除')
    await refresh()
  } catch (error) {
    ElMessage.error(error instanceof ApiError ? error.message : '删除失败')
  }
}

function applyPreset(seconds: number) {
  form.value.durationDays = seconds / 86400
}

onMounted(() => {
  refresh()
  ticker = window.setInterval(() => {
    nowSeconds.value = Math.floor(Date.now() / 1000)
  }, 1000)
})

onUnmounted(() => {
  if (ticker) window.clearInterval(ticker)
})
</script>

<template>
  <div v-loading="loading">
    <div class="page-header">
      <div>
        <h2>定时器</h2>
        <div class="muted">在签到周期内完成签到；超时未签到时将执行关联的规则</div>
      </div>
      <el-button type="primary" :icon="Plus" @click="openCreate">创建定时器</el-button>
    </div>

    <el-empty v-if="!loading && timers.length === 0" description="还没有定时器">
      <el-button type="primary" @click="openCreate">创建定时器</el-button>
    </el-empty>

    <div class="timer-grid">
      <el-card
        v-for="(timer, i) in timers"
        :key="timer.id"
        class="timer-card gb-rise"
        :style="{ animationDelay: `${i * 0.05}s` }"
      >
        <div class="timer-head">
          <div class="timer-name">{{ timer.name }}</div>
          <el-tag :type="statusOf(timer).type" size="small" effect="light">
            {{ statusOf(timer).label }}
          </el-tag>
        </div>

        <div class="muted timer-desc">{{ timer.description || ' ' }}</div>

        <div
          class="timer-countdown"
          :class="{ danger: timer.triggered || secondsLeft(timer) <= timer.remind_time_seconds }"
        >
          <template v-if="timer.triggered">已触发，规则已执行；请重新签到</template>
          <template v-else-if="!timer.enabled">已暂停</template>
          <template v-else-if="secondsLeft(timer) > 0">
            {{ formatDuration(secondsLeft(timer)) }} 后到期
          </template>
          <template v-else><span class="gb-pulse">已到期，即将执行规则</span></template>
        </div>

        <!-- 剩余时间进度条 -->
        <div class="timer-bar">
          <div
            class="timer-bar-fill"
            :style="{
              width: `${fractionLeft(timer) * 100}%`,
              background: barColor(timer),
            }"
          />
        </div>

        <div class="muted timer-meta">
          <span>周期 {{ formatDuration(timer.sign_deration_seconds) }}</span>
          <span>提前 {{ formatDuration(timer.remind_time_seconds) }} 提醒</span>
          <span>上次 {{ formatDateTime(timer.last_sign) }}</span>
        </div>

        <div class="timer-actions">
          <el-button type="primary" :icon="Check" size="small" @click="sign(timer)">
            签到
          </el-button>
          <el-button
            size="small"
            :icon="VideoPlay"
            :loading="triggeringId === timer.id"
            title="手动执行该定时器下的所有规则，用于调试"
            @click="trigger(timer)"
          >
            触发
          </el-button>
          <el-button size="small" @click="openEdit(timer)">编辑</el-button>
          <el-button size="small" type="danger" plain @click="remove(timer)">删除</el-button>
          <el-switch
            v-model="timer.enabled"
            size="small"
            class="timer-switch"
            @change="toggleEnabled(timer)"
          />
        </div>
      </el-card>
    </div>

    <!-- 创建/编辑对话框 -->
    <el-dialog v-model="dialogVisible" :title="dialogTitle" width="480px">
      <el-form label-width="90px">
        <el-form-item label="名称" required>
          <el-input v-model="form.name" placeholder="例如：每周报平安" maxlength="64" />
        </el-form-item>
        <el-form-item label="描述">
          <el-input v-model="form.description" type="textarea" :rows="2" placeholder="备注(可选)" />
        </el-form-item>
        <el-form-item label="签到周期" required>
          <el-input-number v-model="form.durationDays" :min="0.001" :step="1" style="width: 160px" />
          <span class="unit">天</span>
          <div class="presets">
            <el-tag
              v-for="preset in durationPresets"
              :key="preset.seconds"
              class="preset-tag"
              :effect="Math.round(form.durationDays * 86400) === preset.seconds ? 'dark' : 'plain'"
              :type="Math.round(form.durationDays * 86400) === preset.seconds ? 'primary' : 'info'"
              @click="applyPreset(preset.seconds)"
            >
              {{ preset.label }}
            </el-tag>
          </div>
        </el-form-item>
        <el-form-item label="提前提醒" required>
          <el-input-number v-model="form.remindHours" :min="0" :step="1" style="width: 160px" />
          <span class="unit">小时</span>
          <div class="muted">到期前会通过提醒渠道(设置页)提醒你签到</div>
        </el-form-item>
        <el-form-item label="启用">
          <el-switch v-model="form.enabled" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="dialogVisible = false">取消</el-button>
        <el-button type="primary" :loading="saving" @click="save">保存</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<style scoped>
.timer-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(320px, 1fr));
  gap: 16px;
}

.timer-card:hover {
  transform: translateY(-3px);
  box-shadow: var(--gb-shadow-hover);
}

.timer-head {
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.timer-name {
  font-size: 16px;
  font-weight: 700;
  letter-spacing: -0.01em;
}

.timer-desc {
  margin: 6px 0 12px;
  min-height: 18px;
}

.timer-countdown {
  font-size: 18px;
  font-weight: 700;
  letter-spacing: -0.01em;
  color: var(--gb-primary-deep);
  margin-bottom: 10px;
}

.timer-countdown.danger {
  color: var(--el-color-danger);
}

.timer-bar {
  height: 6px;
  border-radius: 999px;
  background: var(--el-fill-color);
  overflow: hidden;
  margin-bottom: 12px;
}

.timer-bar-fill {
  height: 100%;
  border-radius: 999px;
  transition:
    width 1s linear,
    background 0.5s var(--gb-ease);
}

.timer-meta {
  display: flex;
  flex-wrap: wrap;
  gap: 4px 14px;
  margin-bottom: 14px;
}

.timer-actions {
  display: flex;
  align-items: center;
}

.timer-switch {
  margin-left: auto;
}

.unit {
  margin-left: 8px;
}

.presets {
  margin-top: 8px;
  width: 100%;
}

.preset-tag {
  margin-right: 6px;
  cursor: pointer;
  transition: all 0.2s var(--gb-ease);
}
</style>
