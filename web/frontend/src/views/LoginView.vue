<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import { siteApi, userApi } from '@/api'
import { ApiError } from '@/api/client'
import type { SiteInfo } from '@/api/types'
import { useUserStore } from '@/stores/user'
import LogoMark from '@/components/LogoMark.vue'

const route = useRoute()
const router = useRouter()
const userStore = useUserStore()

const site = ref<SiteInfo | null>(null)
const mode = ref<'login' | 'register'>('login')
const loading = ref(false)

const form = ref({
  username: '',
  password: '',
  confirm: '',
})

onMounted(async () => {
  try {
    site.value = await siteApi.info()
    // 全新部署引导创建第一个账号
    if (site.value.need_initial_user) {
      mode.value = 'register'
    }
  } catch {
    // 站点信息拉取失败不阻塞登录
  }
})

async function submit() {
  if (!form.value.username || !form.value.password) {
    ElMessage.warning('请输入用户名和密码')
    return
  }
  if (mode.value === 'register' && form.value.password !== form.value.confirm) {
    ElMessage.warning('两次输入的密码不一致')
    return
  }

  loading.value = true
  try {
    if (mode.value === 'register') {
      await userApi.register(form.value.username, form.value.password)
      ElMessage.success('注册成功')
    } else {
      await userApi.login(form.value.username, form.value.password)
    }
    await userStore.fetchUser()
    const redirect = typeof route.query.redirect === 'string' ? route.query.redirect : '/dashboard'
    router.push(redirect)
  } catch (error) {
    ElMessage.error(error instanceof ApiError ? error.message : '操作失败')
  } finally {
    loading.value = false
  }
}
</script>

<template>
  <div class="login-page">
    <!-- 左侧：项目介绍 -->
    <div class="intro">
      <div class="intro-inner gb-rise">
        <div class="intro-mark"><LogoMark :size="28" /></div>
        <h1 class="intro-title">goodBaby</h1>
        <p class="intro-slogan">摇篮系统 · Dead Man's Switch</p>
        <p class="intro-desc">
          设定签到周期并定期签到；<br />
          超时未签到时，系统会自动执行你预设的规则。
        </p>
        <ul class="intro-points">
          <li>到期前通过钉钉机器人提醒签到</li>
          <li>支持发送邮件、QQ、钉钉消息、B 站动态</li>
          <li>支持自动公开 GitHub 仓库</li>
        </ul>
      </div>
    </div>

    <!-- 右侧：登录表单 -->
    <div class="panel">
      <div class="panel-card gb-rise">
        <h2 class="panel-title">
          {{ site?.need_initial_user ? '创建初始账号' : mode === 'login' ? '登录' : '注册' }}
        </h2>
        <p class="panel-sub">
          {{
            site?.need_initial_user
              ? '首次使用，请先创建管理员账号'
              : mode === 'login'
                ? '登录你的 goodBaby 账号'
                : '创建一个新账号'
          }}
        </p>

        <el-form label-position="top" size="large" @keyup.enter="submit">
          <el-form-item label="用户名">
            <el-input v-model="form.username" placeholder="2-32 个字符" autofocus />
          </el-form-item>
          <el-form-item label="密码">
            <el-input
              v-model="form.password"
              type="password"
              show-password
              placeholder="6-64 个字符"
            />
          </el-form-item>
          <el-form-item v-if="mode === 'register'" label="确认密码">
            <el-input v-model="form.confirm" type="password" show-password placeholder="再输入一次" />
          </el-form-item>

          <el-button type="primary" size="large" class="submit-btn" :loading="loading" @click="submit">
            {{ mode === 'login' ? '登 录' : '注 册' }}
          </el-button>
        </el-form>

        <div v-if="site?.enable_registry && !site?.need_initial_user" class="switch-mode">
          <el-link v-if="mode === 'login'" type="primary" :underline="false" @click="mode = 'register'">
            没有账号？注册
          </el-link>
          <el-link v-else type="primary" :underline="false" @click="mode = 'login'">
            已有账号？登录
          </el-link>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
.login-page {
  height: 100%;
  display: flex;
}

/* ---------- 左侧理念区 ---------- */

.intro {
  flex: 1.15;
  background:
    radial-gradient(ellipse at 20% 20%, rgb(102 204 255 / 0.16) 0%, transparent 55%),
    radial-gradient(ellipse at 85% 80%, rgb(102 204 255 / 0.12) 0%, transparent 50%),
    linear-gradient(160deg, #14161c 0%, #1b1f27 100%);
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 48px;
}

.intro-inner {
  max-width: 420px;
}

.intro-mark {
  width: 56px;
  height: 56px;
  border-radius: 16px;
  background: linear-gradient(135deg, var(--gb-primary) 0%, #3fa9e0 100%);
  display: flex;
  align-items: center;
  justify-content: center;
  color: #0b2b3d;
  box-shadow: 0 8px 24px rgb(102 204 255 / 0.5);
  margin-bottom: 24px;
}

.intro-title {
  margin: 0;
  color: #fff;
  font-size: 34px;
  font-weight: 800;
  letter-spacing: -0.02em;
}

.intro-slogan {
  margin: 10px 0 20px;
  color: #66ccff;
  font-size: 17px;
  font-weight: 600;
}

.intro-desc {
  margin: 0 0 26px;
  color: rgb(255 255 255 / 0.66);
  font-size: 15px;
  line-height: 2;
}

.intro-points {
  margin: 0;
  padding: 0;
  list-style: none;
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.intro-points li {
  color: rgb(255 255 255 / 0.5);
  font-size: 13px;
  padding-left: 20px;
  position: relative;
}

.intro-points li::before {
  content: '';
  position: absolute;
  left: 0;
  top: 6px;
  width: 7px;
  height: 7px;
  border-radius: 50%;
  background: var(--gb-primary);
}

/* ---------- 右侧表单区 ---------- */

.panel {
  flex: 1;
  display: flex;
  align-items: center;
  justify-content: center;
  background: var(--gb-bg);
  padding: 32px;
}

.panel-card {
  width: 100%;
  max-width: 360px;
}

.panel-title {
  margin: 0;
  font-size: 26px;
  font-weight: 700;
  letter-spacing: -0.02em;
  color: var(--el-text-color-primary);
}

.panel-sub {
  margin: 8px 0 28px;
  color: var(--el-text-color-secondary);
  font-size: 14px;
}

.submit-btn {
  width: 100%;
  margin-top: 8px;
  font-weight: 600;
  letter-spacing: 0.1em;
}

.switch-mode {
  text-align: center;
  margin-top: 20px;
}

/* 窄屏时隐藏左侧介绍，登录框居中占满 */
@media (max-width: 860px) {
  .intro {
    display: none;
  }
}

@media (max-width: 768px) {
  .login-page {
    /* 移动浏览器地址栏会挤压 100%，用 dvh 保证真实可视高度 */
    min-height: 100dvh;
  }

  .panel {
    padding: 24px 20px;
    align-items: flex-start;
    padding-top: 12vh;
  }

  .panel-card {
    max-width: none;
  }

  .panel-title {
    font-size: 23px;
  }

  /* 输入框放大到 16px，避免 iOS Safari 聚焦时自动缩放页面 */
  .panel :deep(.el-input__inner) {
    font-size: 16px;
  }
}
</style>
