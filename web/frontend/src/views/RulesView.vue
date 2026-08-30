<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Plus, VideoPlay } from '@element-plus/icons-vue'
import { accountApi, gatewayApi, ruleApi, timerApi } from '@/api'
import { ApiError } from '@/api/client'
import type { Account, MessageGateway, Rule, Timer } from '@/api/types'
import { useMetaStore } from '@/stores/meta'
import { formatDateTime } from '@/utils/format'
import ConfigForm from '@/components/ConfigForm.vue'
import { useIsMobile } from '@/composables/useBreakpoint'

const isMobile = useIsMobile()

const metaStore = useMetaStore()

const rules = ref<Rule[]>([])
const timers = ref<Timer[]>([])
const gateways = ref<MessageGateway[]>([])
const accounts = ref<Account[]>([])
const loading = ref(false)
const testingId = ref<number | null>(null)
const filterTimerId = ref<number | undefined>(undefined)

// ---- 创建 / 编辑对话框 ----
const dialogVisible = ref(false)
const saving = ref(false)
const editingId = ref<number | null>(null)

const form = ref<{
  name: string
  type: string
  timer_id: number
  gateway_id: number
  account_id: number | undefined
  enabled: boolean
  config_json: string
}>({
  name: '',
  type: '',
  timer_id: 0,
  gateway_id: 0,
  account_id: undefined,
  enabled: true,
  config_json: '',
})

const dialogTitle = computed(() => (editingId.value === null ? '创建规则' : '编辑规则'))
const currentMeta = computed(() => metaStore.ruleMeta(form.value.type))

/** 当前规则类型可用的账号(按类型过滤) */
const availableAccounts = computed(() => {
  const requiredType = currentMeta.value?.account_type
  if (!requiredType) return []
  return accounts.value.filter((a) => a.type === requiredType)
})

const filteredRules = computed(() => {
  if (!filterTimerId.value) return rules.value
  return rules.value.filter((r) => r.timer_id === filterTimerId.value)
})

function timerName(id: number): string {
  return timers.value.find((t) => t.id === id)?.name ?? `#${id}`
}

function accountName(id: number): string {
  if (!id) return '-'
  return accounts.value.find((a) => a.id === id)?.name ?? `#${id}`
}

async function refresh() {
  loading.value = true
  try {
    const [ruleList, timerList, accountList, gatewayList] = await Promise.all([
      ruleApi.list(),
      timerApi.list(),
      accountApi.list(),
      gatewayApi.list(),
    ])
    rules.value = ruleList
    timers.value = timerList
    accounts.value = accountList
    gateways.value = gatewayList
  } catch (error) {
    ElMessage.error(error instanceof ApiError ? error.message : '加载失败')
  } finally {
    loading.value = false
  }
}

function openCreate() {
  if (timers.value.length === 0 && gateways.value.length === 0) {
    ElMessage.warning('请先在“定时器”页面创建一个定时器')
    return
  }
  editingId.value = null
  form.value = {
    name: '',
    type: metaStore.ruleMetas[0]?.type ?? '',
    timer_id: timers.value[0]?.id ?? 0,
    gateway_id: 0,
    account_id: undefined,
    enabled: true,
    config_json: '',
  }
  dialogVisible.value = true
}

function openEdit(rule: Rule) {
  editingId.value = rule.id
  form.value = {
    name: rule.name,
    type: rule.type,
    timer_id: rule.timer_id,
    gateway_id: rule.gateway_id || 0,
    account_id: rule.account_id || undefined,
    enabled: rule.enabled,
    config_json: rule.config_json,
  }
  dialogVisible.value = true
}

function onTypeChange() {
  form.value.config_json = ''
  form.value.account_id = undefined
}

async function save() {
  if (!form.value.name) {
    ElMessage.warning('请填写规则名称')
    return
  }
  if (currentMeta.value?.account_type && !form.value.account_id) {
    ElMessage.warning(`该规则类型需要关联一个「${metaStore.accountLabel(currentMeta.value.account_type)}」账号`)
    return
  }

  const body = {
    name: form.value.name,
    type: form.value.type,
    timer_id: form.value.timer_id,
    gateway_id: form.value.gateway_id,
    account_id: form.value.account_id ?? 0,
    enabled: form.value.enabled,
    config_json: form.value.config_json,
  }

  saving.value = true
  try {
    if (editingId.value === null) {
      await ruleApi.create(body)
      ElMessage.success('规则已创建')
    } else {
      await ruleApi.update(editingId.value, body)
      ElMessage.success('规则已更新')
    }
    dialogVisible.value = false
    await refresh()
  } catch (error) {
    ElMessage.error(error instanceof ApiError ? error.message : '保存失败')
  } finally {
    saving.value = false
  }
}

async function toggleEnabled(rule: Rule) {
  try {
    await ruleApi.update(rule.id, {
      name: rule.name,
      type: rule.type,
      timer_id: rule.timer_id,
      gateway_id: rule.gateway_id,
      account_id: rule.account_id,
      enabled: rule.enabled,
      config_json: rule.config_json,
    })
    ElMessage.success(rule.enabled ? '已启用' : '已停用')
  } catch (error) {
    rule.enabled = !rule.enabled
    ElMessage.error(error instanceof ApiError ? error.message : '操作失败')
  }
}

async function test(rule: Rule) {
  try {
    await ElMessageBox.confirm(
      '测试会真实执行该规则（发送消息 / 修改仓库等），确定继续吗？',
      `测试规则「${rule.name}」`,
      { type: 'warning', confirmButtonText: '执行', cancelButtonText: '取消' },
    )
  } catch {
    return
  }

  testingId.value = rule.id
  try {
    await ruleApi.test(rule.id)
    ElMessage.success('执行成功，详情见执行日志')
  } catch (error) {
    ElMessage.error(error instanceof ApiError ? error.message : '执行失败')
  } finally {
    testingId.value = null
  }
}

async function remove(rule: Rule) {
  try {
    await ElMessageBox.confirm(`确定删除规则「${rule.name}」吗？`, '删除规则', {
      type: 'warning',
      confirmButtonText: '删除',
      confirmButtonClass: 'el-button--danger',
      cancelButtonText: '取消',
    })
  } catch {
    return
  }

  try {
    await ruleApi.remove(rule.id)
    ElMessage.success('已删除')
    await refresh()
  } catch (error) {
    ElMessage.error(error instanceof ApiError ? error.message : '删除失败')
  }
}

onMounted(async () => {
  await metaStore.ensureLoaded()
  await refresh()
})
</script>

<template>
  <div v-loading="loading">
    <div class="page-header">
      <div>
        <h2>规则</h2>
        <div class="muted">定时器到期时要执行的动作：发送邮件、QQ / 钉钉消息、B 站动态，或公开 GitHub 仓库</div>
      </div>
      <div class="header-tools">
        <el-select
          v-model="filterTimerId"
          placeholder="按定时器筛选"
          clearable
          class="filter-select"
        >
          <el-option v-for="t in timers" :key="t.id" :label="t.name" :value="t.id" />
        </el-select>
        <el-button type="primary" :icon="Plus" @click="openCreate">创建规则</el-button>
      </div>
    </div>

    <el-empty v-if="!loading && filteredRules.length === 0" description="还没有规则">
      <el-button type="primary" @click="openCreate">创建规则</el-button>
    </el-empty>

    <el-card v-else class="table-card gb-rise">
    <el-table :data="filteredRules" :size="isMobile ? 'small' : 'default'">
      <el-table-column prop="name" label="名称" min-width="110" />
      <el-table-column label="类型" :width="isMobile ? 120 : 160">
        <template #default="{ row }">
          <el-tag>{{ metaStore.ruleLabel(row.type) }}</el-tag>
        </template>
      </el-table-column>
      <!-- 窄屏隐藏定时器/账号/创建时间，保留名称、类型与操作 -->
      <el-table-column v-if="!isMobile" label="定时器" min-width="120">
        <template #default="{ row }">{{ timerName(row.timer_id) }}</template>
      </el-table-column>
      <el-table-column v-if="!isMobile" label="账号" min-width="120">
        <template #default="{ row }">{{ accountName(row.account_id) }}</template>
      </el-table-column>
      <el-table-column v-if="!isMobile" label="创建时间" width="170">
        <template #default="{ row }">{{ formatDateTime(row.create_at) }}</template>
      </el-table-column>
      <el-table-column label="启用" :width="isMobile ? 60 : 80">
        <template #default="{ row }">
          <el-switch v-model="row.enabled" size="small" @change="toggleEnabled(row)" />
        </template>
      </el-table-column>
      <el-table-column label="操作" :width="isMobile ? 200 : 240" fixed="right">
        <template #default="{ row }">
          <el-button
            size="small"
            :icon="VideoPlay"
            :loading="testingId === row.id"
            @click="test(row)"
          >
            测试
          </el-button>
          <el-button size="small" @click="openEdit(row)">编辑</el-button>
          <el-button size="small" type="danger" plain @click="remove(row)">删除</el-button>
        </template>
      </el-table-column>
    </el-table>
    </el-card>

    <!-- 创建/编辑对话框 -->
    <el-dialog v-model="dialogVisible" :title="dialogTitle" width="560px">
      <el-form :label-width="isMobile ? 'auto' : '110px'" :label-position="isMobile ? 'top' : 'right'">
        <el-form-item label="规则名称" required>
          <el-input v-model="form.name" placeholder="给这条规则起个名字" maxlength="64" />
        </el-form-item>
        <el-form-item label="规则类型" required>
          <el-select
            v-model="form.type"
            :disabled="editingId !== null"
            style="width: 100%"
            @change="onTypeChange"
          >
            <el-option
              v-for="m in metaStore.ruleMetas"
              :key="m.type"
              :label="m.label"
              :value="m.type"
            />
          </el-select>
          <div v-if="currentMeta?.description" class="muted">{{ currentMeta.description }}</div>
        </el-form-item>
        <el-form-item label="关联定时器/消息网关" required>
          <div class="source-hint muted">规则只能选择一种触发方式</div>
          <el-select v-model="form.timer_id" class="source-select" clearable placeholder="定时器触发">
            <el-option v-for="t in timers" :key="t.id" :label="t.name" :value="t.id" />
          </el-select>
          <el-select v-model="form.gateway_id" clearable placeholder="消息网关（可选）" style="width: 100%; margin-top: 8px" @change="form.gateway_id && (form.timer_id = 0)">
            <el-option v-for="g in gateways" :key="g.id" :label="g.name" :value="g.id" />
          </el-select>
        </el-form-item>
        <el-form-item v-if="currentMeta?.account_type" label="关联账号" required>
          <el-select
            v-model="form.account_id"
            style="width: 100%"
            :placeholder="availableAccounts.length ? '选择账号' : '没有可用账号，请先到“账号”页面添加'"
          >
            <el-option
              v-for="a in availableAccounts"
              :key="a.id"
              :label="a.name"
              :value="a.id"
            />
          </el-select>
          <div class="muted">
            需要「{{ metaStore.accountLabel(currentMeta.account_type) }}」类型的账号
          </div>
        </el-form-item>

        <ConfigForm
          v-if="currentMeta"
          v-model="form.config_json"
          :fields="currentMeta.fields"
        />

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
.tip {
  margin-bottom: 16px;
}

.header-tools {
  display: flex;
  gap: 12px;
}

.filter-select {
  width: 200px;
}

.source-select { width: 100%; margin-top: 8px; }
.source-hint { margin-bottom: 2px; }

@media (max-width: 768px) {
  .header-tools {
    flex-wrap: wrap;
    gap: 10px;
  }

  .filter-select {
    width: 100%;
  }
}
</style>
