<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { ElMessage } from 'element-plus'
import { userApi } from '@/api'
import { ApiError } from '@/api/client'
import { useUserStore } from '@/stores/user'
import { formatDateTime } from '@/utils/format'

const userStore = useUserStore()

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
      <el-descriptions :column="2">
        <el-descriptions-item label="用户名">
          {{ userStore.user?.username }}
        </el-descriptions-item>
        <el-descriptions-item label="注册时间">
          {{ formatDateTime(userStore.user?.create_at ?? 0) }}
        </el-descriptions-item>
      </el-descriptions>
    </el-card>

    <el-card class="section">
      <template #header>提醒渠道（钉钉机器人）</template>
      <p class="muted">
        定时器临近到期时，通过钉钉自定义机器人提醒你签到。在钉钉群里添加“自定义机器人”，把
        Webhook 中的 access_token 填到下面；如果机器人开启了“加签”，还需要填写 Secret。
      </p>
      <el-form label-width="130px" style="max-width: 560px">
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
      <el-form label-width="130px" style="max-width: 560px">
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
</style>
