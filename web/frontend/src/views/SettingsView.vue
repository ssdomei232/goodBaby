<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'
import { ElMessage } from 'element-plus'
import { adminApi, userApi } from '@/api'
import { ApiError } from '@/api/client'
import type { AdminConfig, AdminConfigResponse } from '@/api/types'
import { useUserStore } from '@/stores/user'
import { formatDateTime } from '@/utils/format'
import { useIsMobile } from '@/composables/useBreakpoint'

const isMobile = useIsMobile()

const userStore = useUserStore()

const isAdmin = computed(() => userStore.user?.is_admin === true)

// ---- 系统配置（仅管理员） ----
const adminLoading = ref(false)
const adminSaving = ref(false)
const adminReadonly = ref<AdminConfigResponse['readonly'] | null>(null)
const adminConfig = ref<AdminConfig>({
  enable_registry: true,
  timeout_duration_hours: 6,
  check_interval_minutes: 10,
  log_retain_count: 500,
})

async function loadAdminConfig() {
  if (!isAdmin.value) return
  adminLoading.value = true
  try {
    const result = await adminApi.getConfig()
    adminConfig.value = result.config
    adminReadonly.value = result.readonly
  } catch (error) {
    ElMessage.error(error instanceof ApiError ? error.message : '读取系统配置失败')
  } finally {
    adminLoading.value = false
  }
}

async function saveAdminConfig() {
  adminSaving.value = true
  try {
    adminConfig.value = await adminApi.updateConfig(adminConfig.value)
    ElMessage.success('系统配置已保存并生效')
  } catch (error) {
    ElMessage.error(error instanceof ApiError ? error.message : '保存失败')
  } finally {
    adminSaving.value = false
  }
}

const driverLabels: Record<string, string> = {
  sqlite: 'SQLite',
  postgres: 'PostgreSQL',
}

// ---- 钉钉提醒配置 ----
const notifySaving = ref(false)
const dingtalk = ref({
  enabled: false,
  access_token: '',
  secret: '',
})

function loadNotifyFromUser() {
  const raw = userStore.user?.dingtalk_config
  if (!raw) {
    dingtalk.value = { enabled: false, access_token: '', secret: '' }
    return
  }
  try {
    const parsed = JSON.parse(raw) as { access_token?: string; secret?: string }
    dingtalk.value = {
      enabled: true,
      access_token: parsed.access_token ?? '',
      secret: parsed.secret ?? '',
    }
  } catch {
    dingtalk.value = { enabled: false, access_token: '', secret: '' }
  }
}

async function saveNotify() {
  if (dingtalk.value.enabled && !dingtalk.value.access_token) {
    ElMessage.warning('请填写钉钉机器人的 Access Token')
    return
  }

  notifySaving.value = true
  try {
    const payload = dingtalk.value.enabled
      ? JSON.stringify({
          access_token: dingtalk.value.access_token,
          secret: dingtalk.value.secret,
        })
      : null
    await userApi.updateNotify(payload)
    ElMessage.success('提醒配置已保存')
    await userStore.fetchUser()
  } catch (error) {
    ElMessage.error(error instanceof ApiError ? error.message : '保存失败')
  } finally {
    notifySaving.value = false
  }
}

// ---- 修改密码 ----
const passwordSaving = ref(false)
const passwordForm = ref({
  old: '',
  new1: '',
  new2: '',
})

async function changePassword() {
  if (!passwordForm.value.old || !passwordForm.value.new1) {
    ElMessage.warning('请填写完整')
    return
  }
  if (passwordForm.value.new1 !== passwordForm.value.new2) {
    ElMessage.warning('两次输入的新密码不一致')
    return
  }

  passwordSaving.value = true
  try {
    await userApi.changePassword(passwordForm.value.old, passwordForm.value.new1)
    ElMessage.success('密码修改成功')
    passwordForm.value = { old: '', new1: '', new2: '' }
  } catch (error) {
    ElMessage.error(error instanceof ApiError ? error.message : '修改失败')
  } finally {
    passwordSaving.value = false
  }
}

onMounted(async () => {
  if (!userStore.user) {
    await userStore.fetchUser()
  }
  loadNotifyFromUser()
  await loadAdminConfig()
})
</script>

<template>
  <div class="settings">
    <div class="page-header">
      <div>
        <h2>设置</h2>
        <div class="muted">提醒渠道与账号安全</div>
      </div>
    </div>

    <el-card class="section">
      <template #header>账号信息</template>
      <el-descriptions :column="isMobile ? 1 : 2">
        <el-descriptions-item label="用户名">
          {{ userStore.user?.username }}
          <el-tag v-if="isAdmin" type="primary" size="small" effect="light" class="admin-tag">
            管理员
          </el-tag>
        </el-descriptions-item>
        <el-descriptions-item label="注册时间">
          {{ formatDateTime(userStore.user?.create_at ?? 0) }}
        </el-descriptions-item>
      </el-descriptions>
    </el-card>

    <!-- 系统配置：仅管理员可见 -->
    <el-card v-if="isAdmin" v-loading="adminLoading" class="section">
      <template #header>系统配置</template>
      <p class="muted">
        这些配置对整个站点生效，保存后立即生效，无需重启。监听地址、数据库等启动项只能改配置文件。
      </p>

      <el-form
        :label-width="isMobile ? 'auto' : '150px'"
        :label-position="isMobile ? 'top' : 'right'"
        style="max-width: 620px"
      >
        <el-form-item label="开放注册">
          <el-switch v-model="adminConfig.enable_registry" />
          <div class="muted">
            关闭后新用户无法注册。系统内还没有任何用户时始终允许注册，避免全新部署无法创建账号。
          </div>
        </el-form-item>

        <el-form-item label="检查间隔">
          <el-input-number
            v-model="adminConfig.check_interval_minutes"
            :min="1"
            :max="1440"
            class="cfg-number"
          />
          <span class="unit">分钟</span>
          <div class="muted">多久检查一次定时器是否到期。间隔越短，触发越及时。</div>
        </el-form-item>

        <el-form-item label="规则重试时长">
          <el-input-number
            v-model="adminConfig.timeout_duration_hours"
            :min="1"
            :max="72"
            class="cfg-number"
          />
          <span class="unit">小时</span>
          <div class="muted">规则执行失败后按指数退避重试的最长时间。</div>
        </el-form-item>

        <el-form-item label="日志保留条数">
          <el-input-number
            v-model="adminConfig.log_retain_count"
            :min="0"
            :max="100000"
            :step="100"
            class="cfg-number"
          />
          <span class="unit">条 / 用户</span>
          <div class="muted">超出后自动清理最旧的记录，0 表示不限制。</div>
        </el-form-item>

        <el-form-item v-if="adminReadonly" label="运行环境">
          <div class="muted readonly-info">
            <span>监听地址 <code class="mono">{{ adminReadonly.listen_addr }}</code></span>
            <span>
              数据库
              <code class="mono">
                {{ driverLabels[adminReadonly.database_driver] ?? adminReadonly.database_driver }}
              </code>
            </span>
          </div>
        </el-form-item>

        <el-form-item>
          <el-button type="primary" :loading="adminSaving" @click="saveAdminConfig">
            保存系统配置
          </el-button>
        </el-form-item>
      </el-form>
    </el-card>

    <el-card class="section">
      <template #header>提醒渠道（钉钉机器人）</template>
      <p class="muted">
        定时器临近到期时，通过钉钉自定义机器人提醒你签到。在钉钉群里添加“自定义机器人”，把
        Webhook 中的 access_token 填到下面；如果机器人开启了“加签”，还需要填写 Secret。
      </p>
      <el-form :label-width="isMobile ? 'auto' : '130px'" :label-position="isMobile ? 'top' : 'right'" style="max-width: 560px">
        <el-form-item label="启用钉钉提醒">
          <el-switch v-model="dingtalk.enabled" />
        </el-form-item>
        <template v-if="dingtalk.enabled">
          <el-form-item label="Access Token" required>
            <el-input v-model="dingtalk.access_token" show-password type="password" />
          </el-form-item>
          <el-form-item label="加签 Secret">
            <el-input v-model="dingtalk.secret" show-password type="password" />
          </el-form-item>
        </template>
        <el-form-item>
          <el-button type="primary" :loading="notifySaving" @click="saveNotify">保存</el-button>
        </el-form-item>
      </el-form>
    </el-card>

    <el-card class="section">
      <template #header>修改密码</template>
      <el-form :label-width="isMobile ? 'auto' : '130px'" :label-position="isMobile ? 'top' : 'right'" style="max-width: 560px">
        <el-form-item label="原密码" required>
          <el-input v-model="passwordForm.old" type="password" show-password />
        </el-form-item>
        <el-form-item label="新密码" required>
          <el-input
            v-model="passwordForm.new1"
            type="password"
            show-password
            placeholder="6-64 个字符"
          />
        </el-form-item>
        <el-form-item label="确认新密码" required>
          <el-input v-model="passwordForm.new2" type="password" show-password />
        </el-form-item>
        <el-form-item>
          <el-button type="primary" :loading="passwordSaving" @click="changePassword">
            修改密码
          </el-button>
        </el-form-item>
      </el-form>
    </el-card>
  </div>
</template>

<style scoped>
.section {
  margin-bottom: 16px;
}

.admin-tag {
  margin-left: 6px;
}

.cfg-number {
  width: 150px;
}

.unit {
  margin-left: 8px;
  color: var(--el-text-color-secondary);
  font-size: 13px;
}

.readonly-info {
  display: flex;
  flex-wrap: wrap;
  gap: 4px 18px;
}

.readonly-info code {
  color: var(--el-text-color-primary);
}

@media (max-width: 768px) {
  .cfg-number {
    width: 100%;
  }

  .unit {
    display: inline-block;
    margin: 6px 0 0;
  }
}
</style>
