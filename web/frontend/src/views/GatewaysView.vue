<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Connection, CopyDocument, Delete, Link, Plus, Setting } from '@element-plus/icons-vue'
import { gatewayApi, gatewayRuleApi } from '@/api'
import type { GatewayRule, MessageGateway } from '@/api/types'
import { useMetaStore } from '@/stores/meta'
import { formatDateTime } from '@/utils/format'

const router = useRouter()
const metaStore = useMetaStore()

const gateways = ref<MessageGateway[]>([])
const rules = ref<GatewayRule[]>([])
const loading = ref(false)
const saving = ref(false)
const dialogVisible = ref(false)
const name = ref('')
const type = ref('')
const origin = window.location.origin

const currentMeta = computed(() => metaStore.gatewayMeta(type.value))
const enabledRuleCount = computed(() => rules.value.filter((r) => r.enabled).length)

async function refresh() {
  loading.value = true
  try {
    ;[gateways.value, rules.value] = await Promise.all([gatewayApi.list(), gatewayRuleApi.list()])
  } catch {
    ElMessage.error('加载消息网关失败')
  } finally {
    loading.value = false
  }
}

function openCreate() {
  name.value = ''
  type.value = metaStore.gatewayMetas[0]?.type ?? 'webhook'
  dialogVisible.value = true
}

async function create() {
  if (!name.value.trim()) {
    ElMessage.warning('请输入网关名称')
    return
  }

  saving.value = true
  try {
    await gatewayApi.create(name.value.trim(), type.value)
    dialogVisible.value = false
    await refresh()
    ElMessage.success('消息网关已创建')
  } catch (e) {
    ElMessage.error(e instanceof Error ? e.message : '创建失败')
  } finally {
    saving.value = false
  }
}

async function remove(item: MessageGateway) {
  const count = ruleCount(item)
  const detail = count > 0 ? `关联的 ${count} 条网关规则也会被一并删除。` : ''
  try {
    await ElMessageBox.confirm(`删除“${item.name}”？${detail}`, '删除消息网关', { type: 'warning' })
  } catch {
    return
  }

  try {
    await gatewayApi.remove(item.id)
    await refresh()
    ElMessage.success('消息网关已删除')
  } catch {
    ElMessage.error('删除失败')
  }
}

async function copy(url: string) {
  await navigator.clipboard?.writeText(url)
  ElMessage.success('Webhook 地址已复制')
}

function endpoint(item: MessageGateway) {
  return `${origin}/api/v1/gateways/${item.token}/webhook`
}

function ruleCount(item: MessageGateway) {
  return rules.value.filter((r) => r.gateway_id === item.id).length
}

function payloadHint(item: MessageGateway) {
  return metaStore.gatewayMeta(item.type)?.payload_hint ?? '{ "message": "..." }'
}

function manageRules(item: MessageGateway) {
  router.push({ name: 'gateway-rules', query: { gateway_id: String(item.id) } })
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
        <h2>消息网关</h2>
        <div class="muted">为外部系统生成 Webhook 地址，投递消息即可触发绑定在网关上的规则</div>
      </div>
      <el-button type="primary" :icon="Plus" @click="openCreate">新建网关</el-button>
    </div>

    <div class="gw-stats gb-rise">
      <div class="gw-stat">
        <span class="gw-stat-label">网关</span>
        <span class="gw-stat-value">{{ gateways.length }}</span>
      </div>
      <div class="gw-stat">
        <span class="gw-stat-label">网关规则</span>
        <span class="gw-stat-value">{{ rules.length }}</span>
      </div>
      <div class="gw-stat">
        <span class="gw-stat-label">已启用规则</span>
        <span class="gw-stat-value">{{ enabledRuleCount }}</span>
      </div>
    </div>

    <el-empty v-if="!loading && gateways.length === 0" description="还没有消息网关">
      <el-button type="primary" :icon="Plus" @click="openCreate">新建网关</el-button>
    </el-empty>

    <div v-else class="gateway-grid">
      <el-card
        v-for="(item, i) in gateways"
        :key="item.id"
        class="gateway-card gb-rise"
        :style="{ animationDelay: `${i * 0.05}s` }"
      >
        <div class="gateway-head">
          <div class="gateway-ident">
            <div class="gateway-icon">
              <el-icon :size="18"><Connection /></el-icon>
            </div>
            <div class="gateway-titles">
              <div class="gateway-name">{{ item.name }}</div>
              <div class="gateway-sub">
                <el-tag size="small" type="info" effect="light">
                  {{ metaStore.gatewayLabel(item.type) }}
                </el-tag>
                <span>{{ ruleCount(item) }} 条网关规则</span>
                <span>{{ formatDateTime(item.create_at) }}</span>
              </div>
            </div>
          </div>
          <el-button
            circle
            text
            type="danger"
            :icon="Delete"
            title="删除网关"
            @click="remove(item)"
          />
        </div>

        <div class="endpoint-label">Webhook 地址</div>
        <div class="endpoint">
          <code>{{ endpoint(item) }}</code>
          <el-button text :icon="CopyDocument" title="复制地址" @click="copy(endpoint(item))" />
        </div>

        <div class="endpoint-label payload-label">请求示例</div>
        <pre class="payload"><code>{{ payloadHint(item) }}</code></pre>

        <div class="gateway-actions">
          <el-button size="small" :icon="Setting" @click="manageRules(item)">管理网关规则</el-button>
        </div>
      </el-card>
    </div>

    <!-- 新建网关 -->
    <el-dialog v-model="dialogVisible" title="新建消息网关" width="520px">
      <el-form label-width="90px">
        <el-form-item label="网关类型" required>
          <el-select v-model="type" style="width: 100%">
            <el-option
              v-for="m in metaStore.gatewayMetas"
              :key="m.type"
              :label="m.label"
              :value="m.type"
            />
          </el-select>
          <div v-if="currentMeta?.description" class="muted field-hint">
            {{ currentMeta.description }}
          </div>
        </el-form-item>
        <el-form-item label="网关名称" required>
          <el-input v-model="name" placeholder="例如：生产环境告警" maxlength="64" clearable />
        </el-form-item>
      </el-form>

      <div class="dialog-hint">
        <el-icon><Link /></el-icon>
        <span>创建后到「网关规则」页面，把通知规则关联到这个网关</span>
      </div>

      <template #footer>
        <el-button @click="dialogVisible = false">取消</el-button>
        <el-button type="primary" :loading="saving" @click="create">创建</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<style scoped>
.gw-stats {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(150px, 1fr));
  gap: 12px;
  margin-bottom: 18px;
}

.gw-stat {
  display: flex;
  align-items: baseline;
  gap: 10px;
  padding: 14px 18px;
  border-radius: var(--gb-radius-card);
  background: var(--gb-card);
  box-shadow: var(--gb-shadow-card);
}

.gw-stat-label {
  color: var(--el-text-color-secondary);
  font-size: 13px;
}

.gw-stat-value {
  margin-left: auto;
  font-size: 22px;
  font-weight: 700;
  font-variant-numeric: tabular-nums;
  color: var(--gb-primary-deep);
}

.gateway-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(340px, 1fr));
  gap: 16px;
}

.gateway-card {
  min-width: 0;
}

.gateway-card:hover {
  transform: translateY(-3px);
  box-shadow: var(--gb-shadow-hover);
}

.gateway-head {
  display: flex;
  align-items: flex-start;
  justify-content: space-between;
  gap: 12px;
  margin-bottom: 18px;
}

.gateway-ident {
  display: flex;
  align-items: center;
  gap: 12px;
  min-width: 0;
}

.gateway-icon {
  width: 38px;
  height: 38px;
  flex: none;
  display: grid;
  place-items: center;
  border-radius: 11px;
  color: var(--gb-primary-deep);
  background: var(--el-color-primary-light-9);
}

.gateway-titles {
  min-width: 0;
}

.gateway-name {
  font-size: 16px;
  font-weight: 700;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.gateway-sub {
  display: flex;
  align-items: center;
  flex-wrap: wrap;
  gap: 8px;
  margin-top: 6px;
  color: var(--el-text-color-secondary);
  font-size: 12px;
}

.endpoint-label {
  color: var(--el-text-color-secondary);
  font-size: 12px;
}

.endpoint {
  display: flex;
  align-items: center;
  gap: 6px;
  min-width: 0;
  margin-top: 6px;
  padding: 8px 10px;
  border-radius: 10px;
  background: var(--gb-bg);
}

.endpoint code {
  flex: 1;
  min-width: 0;
  overflow-wrap: anywhere;
  color: var(--el-text-color-primary);
  font-size: 12px;
}

.payload-label {
  margin-top: 14px;
}

.payload {
  margin: 6px 0 0;
  padding: 10px 12px;
  border-radius: 10px;
  background: var(--gb-bg);
  overflow-x: auto;
}

.payload code {
  color: var(--el-text-color-regular);
  font-size: 12px;
  white-space: pre-wrap;
  overflow-wrap: anywhere;
}

.gateway-actions {
  margin-top: 16px;
}

.field-hint {
  line-height: 1.5;
  margin-top: 4px;
}

.dialog-hint {
  display: flex;
  align-items: center;
  gap: 8px;
  margin-top: 4px;
  padding: 10px 12px;
  border-radius: 10px;
  background: var(--el-color-primary-light-9);
  color: var(--el-text-color-regular);
  font-size: 12px;
}

@media (max-width: 600px) {
  .gateway-grid {
    grid-template-columns: 1fr;
  }
}
</style>
