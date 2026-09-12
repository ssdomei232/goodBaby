<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'
import { useRoute } from 'vue-router'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Plus, VideoPlay } from '@element-plus/icons-vue'
import { accountApi, gatewayApi, gatewayRuleApi } from '@/api'
import { ApiError } from '@/api/client'
import type { Account, GatewayRule, MessageGateway } from '@/api/types'
import { useMetaStore } from '@/stores/meta'
import { formatDateTime } from '@/utils/format'
import ConfigForm from '@/components/ConfigForm.vue'
import { useIsMobile } from '@/composables/useBreakpoint'

const isMobile = useIsMobile()
const route = useRoute()

const metaStore = useMetaStore()

const rules = ref<GatewayRule[]>([])
const gateways = ref<MessageGateway[]>([])
const accounts = ref<Account[]>([])
const loading = ref(false)
const testingId = ref<number | null>(null)
const filterGatewayId = ref<number | undefined>(undefined)

// ---- 创建 / 编辑对话框 ----
const dialogVisible = ref(false)
const saving = ref(false)
const editingId = ref<number | null>(null)

const form = ref<{
  name: string
  type: string
  gateway_id: number
  account_id: number | undefined
  enabled: boolean
  config_json: string
}>({
  name: '',
  type: '',
  gateway_id: 0,
  account_id: undefined,
  enabled: true,
  config_json: '',
})

const dialogTitle = computed(() => (editingId.value === null ? '创建网关规则' : '编辑网关规则'))
const currentMeta = computed(() => metaStore.ruleMeta(form.value.type))

/** 由投递请求提供的字段，页面上不让用户填 */
const gatewayMessageFields = computed(() =>
  (currentMeta.value?.fields ?? []).filter((f) => f.gateway_message),
)

/** 需要用户填写的字段 */
const configFields = computed(() => (currentMeta.value?.fields ?? []).filter((f) => !f.gateway_message))

/** 当前规则类型可用的账号(按类型过滤) */
const availableAccounts = computed(() => {
  const requiredType = currentMeta.value?.account_type
  if (!requiredType) return []
  return accounts.value.filter((a) => a.type === requiredType)
})

const filteredRules = computed(() => {
  if (!filterGatewayId.value) return rules.value
  return rules.value.filter((r) => r.gateway_id === filterGatewayId.value)
})

function gatewayName(id: number): string {
  return gateways.value.find((g) => g.id === id)?.name ?? `#${id}`
}

function accountName(id: number): string {
  if (!id) return '-'
  return accounts.value.find((a) => a.id === id)?.name ?? `#${id}`
}

async function refresh() {
  loading.value = true
  try {
    const [ruleList, gatewayList, accountList] = await Promise.all([
      gatewayRuleApi.list(),
      gatewayApi.list(),
      accountApi.list(),
    ])
    rules.value = ruleList
    gateways.value = gatewayList
    accounts.value = accountList
  } catch (error) {
    ElMessage.error(error instanceof ApiError ? error.message : '加载失败')
  } finally {
    loading.value = false
  }
}

function openCreate() {
  if (gateways.value.length === 0) {
    ElMessage.warning('请先在“消息网关”页面创建一个网关')
    return
  }
  editingId.value = null
  form.value = {
    name: '',
    type: metaStore.ruleMetas[0]?.type ?? '',
    gateway_id: filterGatewayId.value ?? gateways.value[0]?.id ?? 0,
    account_id: undefined,
    enabled: true,
    config_json: '',
  }
  dialogVisible.value = true
}

function openEdit(rule: GatewayRule) {
  editingId.value = rule.id
  form.value = {
    name: rule.name,
    type: rule.type,
    gateway_id: rule.gateway_id,
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
  if (!form.value.gateway_id) {
    ElMessage.warning('请选择消息网关')
    return
  }
  if (currentMeta.value?.account_type && !form.value.account_id) {
    ElMessage.warning(`该规则类型需要关联一个「${metaStore.accountLabel(currentMeta.value.account_type)}」账号`)
    return
  }

  const body = {
    name: form.value.name,
    type: form.value.type,
    gateway_id: form.value.gateway_id,
    account_id: form.value.account_id ?? 0,
    enabled: form.value.enabled,
    config_json: form.value.config_json,
  }

  saving.value = true
  try {
    if (editingId.value === null) {
      await gatewayRuleApi.create(body)
      ElMessage.success('网关规则已创建')
    } else {
      await gatewayRuleApi.update(editingId.value, body)
      ElMessage.success('网关规则已更新')
    }
    dialogVisible.value = false
    await refresh()
  } catch (error) {
    ElMessage.error(error instanceof ApiError ? error.message : '保存失败')
  } finally {
    saving.value = false
  }
}

async function toggleEnabled(rule: GatewayRule) {
  try {
    await gatewayRuleApi.update(rule.id, {
      name: rule.name,
      type: rule.type,
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

async function test(rule: GatewayRule) {
  try {
    await ElMessageBox.confirm(
      '测试会按规则里保存的内容真实执行一次（发送消息 / 修改仓库等），确定继续吗？',
      `测试网关规则「${rule.name}」`,
      { type: 'warning', confirmButtonText: '执行', cancelButtonText: '取消' },
    )
  } catch {
    return
  }

  testingId.value = rule.id
  try {
    await gatewayRuleApi.test(rule.id)
    ElMessage.success('执行成功，详情见执行日志')
  } catch (error) {
    ElMessage.error(error instanceof ApiError ? error.message : '执行失败')
  } finally {
    testingId.value = null
  }
}

async function remove(rule: GatewayRule) {
  try {
    await ElMessageBox.confirm(`确定删除网关规则「${rule.name}」吗？`, '删除网关规则', {
      type: 'warning',
      confirmButtonText: '删除',
      confirmButtonClass: 'el-button--danger',
      cancelButtonText: '取消',
    })
  } catch {
    return
  }

  try {
    await gatewayRuleApi.remove(rule.id)
    ElMessage.success('已删除')
    await refresh()
  } catch (error) {
    ElMessage.error(error instanceof ApiError ? error.message : '删除失败')
  }
}

onMounted(async () => {
  // 从「消息网关」页面跳转过来时可以带上要筛选的网关
  const fromQuery = Number(route.query.gateway_id)
  if (fromQuery) filterGatewayId.value = fromQuery

  await metaStore.ensureLoaded()
  await refresh()
})
</script>

<template>
  <div v-loading="loading">
    <div class="page-header">
      <div>
        <h2>网关规则</h2>
        <div class="muted">外部系统往消息网关投递消息时要执行的动作，只需配置「发给谁」，内容来自请求</div>
      </div>
      <div class="header-tools">
        <el-select
          v-model="filterGatewayId"
          placeholder="按消息网关筛选"
          clearable
          class="filter-select"
        >
          <el-option v-for="g in gateways" :key="g.id" :label="g.name" :value="g.id" />
        </el-select>
        <el-button type="primary" :icon="Plus" @click="openCreate">创建网关规则</el-button>
      </div>
    </div>

    <el-alert class="rule-tip" type="info" :closable="false" show-icon>
      <template #title>消息内容由投递请求提供</template>
      投递时把请求里的 <code>message</code> / <code>title</code> 填进规则的消息字段，因此规则里不需要预先写好内容。
    </el-alert>

    <el-empty v-if="!loading && filteredRules.length === 0" description="还没有网关规则">
      <el-button type="primary" @click="openCreate">创建网关规则</el-button>
    </el-empty>

    <el-card v-else class="table-card gb-rise">
      <el-table :data="filteredRules" :size="isMobile ? 'small' : 'default'">
        <el-table-column prop="name" label="名称" min-width="110" />
        <el-table-column label="类型" :width="isMobile ? 120 : 160">
          <template #default="{ row }">
            <el-tag>{{ metaStore.ruleLabel(row.type) }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column v-if="!isMobile" label="消息网关" min-width="120">
          <template #default="{ row }">
            <el-tag size="small" type="info" effect="light">{{ gatewayName(row.gateway_id) }}</el-tag>
          </template>
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
        <el-form-item label="消息网关" required>
          <el-select v-model="form.gateway_id" style="width: 100%" placeholder="选择消息网关">
            <el-option v-for="g in gateways" :key="g.id" :label="g.name" :value="g.id" />
          </el-select>
          <div class="muted">外部系统向这个网关投递消息时触发本规则</div>
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
          <div v-if="currentMeta && !gatewayMessageFields.length" class="muted">
            该规则类型没有消息字段，投递只作为触发，按下方的配置执行
          </div>
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

        <el-form-item v-if="gatewayMessageFields.length" label="消息内容">
          <div class="message-fields">
            <el-tag
              v-for="f in gatewayMessageFields"
              :key="f.key"
              size="small"
              type="info"
              effect="plain"
            >
              {{ f.label }}
            </el-tag>
            <div class="muted">
              由投递请求的 <code>title</code> / <code>message</code> 自动填充，无需在此填写
            </div>
          </div>
        </el-form-item>
        <ConfigForm v-if="currentMeta" v-model="form.config_json" :fields="configFields" />

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
.rule-tip {
  margin-bottom: 16px;
}

.rule-tip code {
  padding: 1px 5px;
  border-radius: 4px;
  background: var(--gb-bg);
  color: var(--el-text-color-primary);
  font-size: 12px;
}

.message-fields {
  display: flex;
  align-items: center;
  flex-wrap: wrap;
  gap: 8px;
  line-height: 1.6;
}

.message-fields .muted {
  width: 100%;
}

.header-tools {
  display: flex;
  gap: 12px;
}

.filter-select {
  width: 200px;
}

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
