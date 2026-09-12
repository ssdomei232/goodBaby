<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import { siteApi, userApi } from '@/api'
import { ApiError } from '@/api/client'
import type { SiteInfo } from '@/api/types'
import { useUserStore } from '@/stores/user'
import LogoMark from '@/components/LogoMark.vue'
import ThemePicker from '@/components/ThemePicker.vue'
import ParanoiaBackdrop from '@/components/ParanoiaBackdrop.vue'
import { useTheme } from '@/composables/useTheme'
import { paranoiaCharacters } from '@/data/paranoia'
import loginArt from '@/assets/paranoia/ch09.jpg'

const route = useRoute()
const router = useRouter()
const userStore = useUserStore()

// 登录页也能挑皮肤，这样不登录就能先看新主题
const { isParanoia } = useTheme()

/** 左侧理念区的文案：跟着皮肤换一副口吻 */
const copy = computed(() =>
  isParanoia.value
    ? {
        title: '妄想症',
        slogan: 'Paranoia · 九重档案',
        desc: '被害者泠珞、加害者与守护者颜语、背叛者零羽，\n九重正篇与四百七十五天的一场层层叠叠的戏剧。',
        points: paranoiaCharacters
          .slice(0, 3)
          .map((item) => `${item.name}【${item.vocal}】· ${item.role}`),
        quote: '谎言重复一千次就变成真理。所以，现实也是这样诞生的。',
      }
    : {
        title: 'goodBaby',
        slogan: "摇篮系统 · Dead Man's Switch",
        desc: '设定签到周期并定期签到；\n超时未签到时，系统会自动执行你预设的规则。',
        points: [
          '到期前通过钉钉机器人提醒签到',
          '支持发送邮件、QQ、钉钉消息、B 站动态',
          '支持自动公开 GitHub 仓库',
        ],
        quote: '',
      },
)

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
  <div
    class="login-page"
    :class="{ 'skin-paranoia': isParanoia }"
    :style="isParanoia ? { '--pa-login-art': `url(${loginArt})` } : undefined"
  >
    <!-- 外观主题：登录前就能先看一眼新皮肤 -->
    <div class="skin-corner"><ThemePicker /></div>

    <!-- 左侧：项目介绍 -->
    <div class="intro">
      <!-- 妄想症皮肤下，花瓣落在曲绘之上 -->
      <ParanoiaBackdrop v-if="isParanoia" />

      <div class="intro-inner gb-rise">
        <div class="intro-mark"><LogoMark :size="28" /></div>
        <h1 class="intro-title">{{ copy.title }}</h1>
        <p class="intro-slogan">{{ copy.slogan }}</p>
        <p class="intro-desc">{{ copy.desc }}</p>
        <p v-if="copy.quote" class="intro-quote">{{ copy.quote }}</p>
        <ul class="intro-points">
          <li v-for="point in copy.points" :key="point">{{ point }}</li>
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
  position: relative;
}

/* 皮肤开关固定在右上角，两种皮肤下都在同一位置 */
.skin-corner {
  position: absolute;
  top: 16px;
  right: 16px;
  z-index: 20;
}

/* ---------- 左侧理念区 ---------- */

.intro {
  flex: 1.15;
  position: relative;
  /* 隔离出图层，让飘落的花瓣落在曲绘之上、文案之下 */
  isolation: isolate;
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
  position: relative;
  z-index: 1;
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
  /* 文案里用 \n 换行，避免模板里散落 <br /> */
  white-space: pre-line;
}

.intro-quote {
  margin: 0 0 24px;
  padding-left: 14px;
  border-left: 2px solid var(--gb-primary);
  color: rgb(255 255 255 / 0.78);
  font-size: 13px;
  line-height: 1.95;
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

/* ---------- 妄想症皮肤：登录页换成一张曲绘 ---------- */

.login-page.skin-paranoia .intro {
  /* 曲绘垫在最底下，压一层暗色的膜，让文字依然读得清 */
  background:
    linear-gradient(
      180deg,
      rgb(9 6 12 / 0.82) 0%,
      rgb(9 6 12 / 0.9) 55%,
      rgb(22 7 14 / 0.97) 100%
    ),
    radial-gradient(ellipse at 22% 18%, rgb(224 36 94 / 0.4), transparent 58%),
    var(--pa-login-art) center / cover no-repeat,
    linear-gradient(160deg, #14060c 0%, #08070c 100%);
}

.login-page.skin-paranoia .intro-mark {
  border-radius: 6px;
  background: linear-gradient(135deg, #e0245e 0%, #7a0f2c 100%);
  color: #fff;
  box-shadow: 0 8px 28px rgb(224 36 94 / 0.5);
}

.login-page.skin-paranoia .intro-title {
  font-family: var(--pa-serif);
  font-size: 38px;
  letter-spacing: 0.12em;
}

.login-page.skin-paranoia .intro-slogan {
  font-family: var(--pa-mono);
  font-size: 12px;
  letter-spacing: 0.28em;
  color: #f2568a;
}

.login-page.skin-paranoia .intro-desc {
  color: rgb(255 255 255 / 0.7);
  font-size: 14px;
}

.login-page.skin-paranoia .intro-quote {
  font-family: var(--pa-serif);
  color: rgb(255 255 255 / 0.86);
}

.login-page.skin-paranoia .intro-points li {
  color: rgb(255 255 255 / 0.56);
}

/* 列表点改成小小的菱形，像彼岸花的花瓣 */
.login-page.skin-paranoia .intro-points li::before {
  border-radius: 1px;
  background: #e0245e;
  transform: rotate(45deg);
}

.login-page.skin-paranoia .panel {
  position: relative;
  z-index: 1;
}

.login-page.skin-paranoia .panel-title {
  font-family: var(--pa-serif);
  letter-spacing: 0.1em;
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
