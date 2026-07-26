<script setup lang="ts">
import { computed, onMounted, onUnmounted, ref } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import {
  Check,
  Refresh,
  Timer as TimerIcon,
  Operation,
  User,
  Warning,
} from '@element-plus/icons-vue'
import { dashboardApi, timerApi } from '@/api'
import { ApiError } from '@/api/client'
import type { DashboardOverview } from '@/api/types'
import { formatDateTime, formatDuration } from '@/utils/format'
import { useIsMobile } from '@/composables/useBreakpoint'

const isMobile = useIsMobile()

const router = useRouter()

const overview = ref<DashboardOverview | null>(null)
const loading = ref(false)
const signing = ref(false)

/** 本地时钟：用服务器时间校准，倒计时每秒刷新 */
const clockOffset = ref(0)
const nowSeconds = ref(Math.floor(Date.now() / 1000))
let clockTimer: number | undefined

const urgentSecondsLeft = computed(() => {
  if (!overview.value?.urgent_timer) return null
  const timer = overview.value.urgent_timer
  return timer.last_sign + timer.sign_deration_seconds - (nowSeconds.value + clockOffset.value)
})

/** 剩余时间占比(0~1)，驱动进度环 */
const urgentFraction = computed(() => {
  if (!overview.value?.urgent_timer || urgentSecondsLeft.value === null) return 0
  const total = overview.value.urgent_timer.sign_deration_seconds
  return Math.min(1, Math.max(0, urgentSecondsLeft.value / total))
})

const RING_R = 52
const ringCircumference = 2 * Math.PI * RING_R

const ringColor = computed(() => {
  const f = urgentFraction.value
  if (f <= 0.1) return '#ef4444'
  if (f <= 0.3) return '#f59e0b'
  return 'var(--gb-primary)'
})

/** 环中心的百分比文字用深一档的颜色，浅色背景上更易读 */
const ringTextColor = computed(() => {
  const f = urgentFraction.value
  if (f <= 0.1) return '#ef4444'
  if (f <= 0.3) return '#f59e0b'
  return 'var(--gb-primary-deep)'
})

const stats = computed(() => [
  {
    label: '定时器',
    sub: `已启用 ${overview.value?.enabled_timers ?? 0} 个`,
    value: overview.value?.timer_count,
    icon: TimerIcon,
    color: '#66ccff',
    bg: 'rgb(102 204 255 / 0.14)',
    to: '/timers',
  },
  {
    label: '规则',
    sub: '触发时执行的动作',
    value: overview.value?.rule_count,
    icon: Operation,
    color: '#6366f1',
    bg: 'rgb(99 102 241 / 0.12)',
    to: '/rules',
  },
  {
    label: '账号',
    sub: '第三方凭据',
    value: overview.value?.account_count,
    icon: User,
    color: '#f59e0b',
    bg: 'rgb(245 158 11 / 0.12)',
    to: '/accounts',
  },
  {
    label: '已触发',
    sub: '规则已执行',
    value: overview.value?.triggered_count,
    icon: Warning,
    color: (overview.value?.triggered_count ?? 0) > 0 ? '#ef4444' : '#94a3b8',
    bg:
      (overview.value?.triggered_count ?? 0) > 0
        ? 'rgb(239 68 68 / 0.12)'
        : 'rgb(148 163 184 / 0.12)',
    to: '/logs',
  },
])

async function refresh() {
  loading.value = true
  try {
    overview.value = await dashboardApi.overview()
    clockOffset.value = overview.value.server_time - Math.floor(Date.now() / 1000)
  } catch (error) {
    ElMessage.error(error instanceof ApiError ? error.message : '加载失败')
  } finally {
    loading.value = false
  }
}

async function signAll() {
  signing.value = true
  try {
    const result = await timerApi.signAll()
    ElMessage.success(`签到成功，${result.signed} 个定时器重新计时`)
    await refresh()
  } catch (error) {
    ElMessage.error(error instanceof ApiError ? error.message : '签到失败')
  } finally {
    signing.value = false
  }
}

onMounted(() => {
  refresh()
  clockTimer = window.setInterval(() => {
    nowSeconds.value = Math.floor(Date.now() / 1000)
  }, 1000)
})

onUnmounted(() => {
  if (clockTimer) window.clearInterval(clockTimer)
})
</script>

<template>
  <div v-loading="loading">
    <div class="page-header">
      <div>
        <h2>仪表盘</h2>
        <div class="muted">总览与快捷签到</div>
      </div>
      <div class="header-btns">
        <el-button :icon="Refresh" circle @click="refresh" />
        <el-button type="primary" size="large" :icon="Check" :loading="signing" @click="signAll">
          一键签到
        </el-button>
      </div>
    </div>

    <!-- 已触发警示 -->
    <el-alert
      v-if="overview && overview.triggered_count > 0"
      type="error"
      show-icon
      :closable="false"
      class="urgent-alert gb-rise"
      :title="`有 ${overview.triggered_count} 个定时器已触发，关联规则已执行。请前往定时器页面重新签到。`"
    />

    <!-- 最近到期的定时器 -->
    <el-card v-if="overview?.urgent_timer && urgentSecondsLeft !== null" class="hero-card gb-rise">
      <div class="hero-body">
        <div class="hero-ring">
          <svg width="128" height="128" viewBox="0 0 128 128">
            <circle
              cx="64"
              cy="64"
              :r="RING_R"
              fill="none"
              stroke="var(--el-fill-color)"
              stroke-width="9"
            />
            <circle
              cx="64"
              cy="64"
              :r="RING_R"
              fill="none"
              :stroke="ringColor"
              stroke-width="9"
              stroke-linecap="round"
              :stroke-dasharray="ringCircumference"
              :stroke-dashoffset="ringCircumference * (1 - urgentFraction)"
              transform="rotate(-90 64 64)"
              class="ring-progress"
            />
          </svg>
          <div class="ring-label">
            <span class="ring-pct" :style="{ color: ringTextColor }">
              {{ Math.round(urgentFraction * 100) }}%
            </span>
            <span class="muted">剩余</span>
          </div>
        </div>

        <div class="hero-info">
          <div class="muted">下一个到期的定时器</div>
          <div class="hero-name">{{ overview.urgent_timer.name }}</div>
          <div class="hero-countdown" :class="{ danger: urgentFraction <= 0.1 }">
            <template v-if="urgentSecondsLeft > 0">
              <span class="countdown-num">{{ formatDuration(urgentSecondsLeft) }}</span>
              <span class="muted">后到期</span>
            </template>
            <template v-else><span class="countdown-num gb-pulse">已到期</span></template>
          </div>
          <div class="muted">上次签到：{{ formatDateTime(overview.urgent_timer.last_sign) }}</div>
        </div>
      </div>
    </el-card>

    <!-- 统计卡片 -->
    <div class="stats-grid">
      <el-card
        v-for="(stat, i) in stats"
        :key="stat.label"
        class="stat-card gb-rise"
        :style="{ animationDelay: `${i * 0.05}s` }"
        @click="router.push(stat.to)"
      >
        <div class="stat-body">
          <div class="stat-icon" :style="{ background: stat.bg, color: stat.color }">
            <el-icon :size="20"><component :is="stat.icon" /></el-icon>
          </div>
          <div>
            <div class="stat-value">{{ stat.value ?? '-' }}</div>
            <div class="stat-label">{{ stat.label }}</div>
            <div class="muted">{{ stat.sub }}</div>
          </div>
        </div>
      </el-card>
    </div>

    <!-- 最近日志 -->
    <el-card class="gb-rise" style="animation-delay: 0.15s">
      <template #header>最近执行</template>
      <el-empty
        v-if="!overview?.recent_logs?.length"
        description="暂无执行记录"
        :image-size="80"
      />
      <el-table v-else :data="overview.recent_logs" :size="isMobile ? 'small' : 'default'">
        <el-table-column label="时间" :width="isMobile ? 130 : 175">
          <template #default="{ row }">
            {{ isMobile ? formatDateTime(row.create_at).slice(5) : formatDateTime(row.create_at) }}
          </template>
        </el-table-column>
        <el-table-column prop="rule_name" label="规则" min-width="100" />
        <el-table-column label="结果" :width="isMobile ? 70 : 90">
          <template #default="{ row }">
            <el-tag :type="row.success ? 'success' : 'danger'" size="small" effect="light">
              {{ row.success ? '成功' : '失败' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="message" label="详情" min-width="240" show-overflow-tooltip />
      </el-table>
    </el-card>
  </div>
</template>

<style scoped>
.header-btns {
  display: flex;
  align-items: center;
  gap: 10px;
}

.urgent-alert {
  margin-bottom: 16px;
}

/* ---------- 主倒计时卡 ---------- */

.hero-card {
  margin-bottom: 20px;
}

.hero-body {
  display: flex;
  align-items: center;
  gap: 32px;
}

.hero-ring {
  position: relative;
  width: 128px;
  height: 128px;
  flex-shrink: 0;
}

.ring-progress {
  transition:
    stroke-dashoffset 0.8s var(--gb-ease),
    stroke 0.5s var(--gb-ease);
}

.ring-label {
  position: absolute;
  inset: 0;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  gap: 2px;
}

.ring-pct {
  font-size: 24px;
  font-weight: 800;
  letter-spacing: -0.02em;
  transition: color 0.5s var(--gb-ease);
}

.hero-info {
  display: flex;
  flex-direction: column;
  gap: 6px;
}

.hero-name {
  font-size: 22px;
  font-weight: 700;
  letter-spacing: -0.02em;
}

.hero-countdown {
  color: var(--gb-primary-deep);
}

.hero-countdown.danger {
  color: var(--el-color-danger);
}

.countdown-num {
  font-size: 26px;
  font-weight: 800;
  letter-spacing: -0.02em;
  margin-right: 8px;
}

/* ---------- 统计卡片 ---------- */

.stats-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
  gap: 16px;
  margin-bottom: 20px;
}

.stat-card {
  cursor: pointer;
}

.stat-card:hover {
  transform: translateY(-3px);
  box-shadow: var(--gb-shadow-hover);
}

.stat-body {
  display: flex;
  align-items: center;
  gap: 16px;
}

.stat-icon {
  width: 46px;
  height: 46px;
  border-radius: 13px;
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
}

.stat-value {
  font-size: 26px;
  font-weight: 800;
  line-height: 1.1;
  letter-spacing: -0.02em;
}

.stat-label {
  font-size: 13px;
  font-weight: 600;
  color: var(--el-text-color-regular);
  margin: 2px 0 1px;
}

/* ---------- 移动端 ---------- */

@media (max-width: 768px) {
  .header-btns {
    width: 100%;
  }

  /* 倒计时环与文字改为上下排列并居中 */
  .hero-body {
    flex-direction: column;
    gap: 18px;
    text-align: center;
  }

  .hero-info {
    align-items: center;
  }

  .hero-name {
    font-size: 20px;
  }

  .countdown-num {
    font-size: 22px;
  }

  /* 统计卡两列，避免一列太长要滑很久 */
  .stats-grid {
    grid-template-columns: repeat(2, 1fr);
    gap: 12px;
  }

  .stat-body {
    gap: 12px;
  }

  .stat-icon {
    width: 40px;
    height: 40px;
    border-radius: 11px;
  }

  .stat-value {
    font-size: 22px;
  }

  /* 触屏没有 hover，去掉位移避免点击时抖动 */
  .stat-card:hover {
    transform: none;
    box-shadow: var(--gb-shadow-card);
  }
}

/* 只有极窄屏才退回单列，375px 的机型两列仍然放得下 */
@media (max-width: 340px) {
  .stats-grid {
    grid-template-columns: 1fr;
  }
}
</style>
