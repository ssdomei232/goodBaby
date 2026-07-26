<script setup lang="ts">
import { computed, ref, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import {
  Odometer,
  Timer as TimerIcon,
  Operation,
  User,
  Document,
  Setting,
  SwitchButton,
  Moon,
  Sunny,
  Fold,
} from '@element-plus/icons-vue'
import { useUserStore } from '@/stores/user'
import { useTheme } from '@/composables/useTheme'
import LogoMark from '@/components/LogoMark.vue'

const route = useRoute()
const router = useRouter()
const userStore = useUserStore()

/** 移动端侧边栏抽屉的展开状态 */
const drawerOpen = ref(false)

// 切换路由后自动收起抽屉
watch(
  () => route.path,
  () => {
    drawerOpen.value = false
  },
)

const menus = [
  { path: '/dashboard', title: '仪表盘', icon: Odometer },
  { path: '/timers', title: '定时器', icon: TimerIcon },
  { path: '/rules', title: '规则', icon: Operation },
  { path: '/accounts', title: '账号', icon: User },
  { path: '/logs', title: '执行日志', icon: Document },
  { path: '/settings', title: '设置', icon: Setting },
]

const activePath = computed(() => `/${route.path.split('/')[1] ?? ''}`)

// 主题状态是全局单例，初始化在 main.ts 里完成
const { isDark, toggleTheme } = useTheme()

async function handleLogout() {
  await userStore.logout()
  ElMessage.success('已退出登录')
  router.push({ name: 'login' })
}
</script>

<template>
  <div class="layout">
    <!-- 移动端抽屉遮罩 -->
    <Transition name="fade">
      <div v-if="drawerOpen" class="drawer-mask" @click="drawerOpen = false" />
    </Transition>

    <!-- 侧边栏：桌面端常驻，移动端为抽屉 -->
    <aside class="aside" :class="{ open: drawerOpen }">
      <div class="logo">
        <div class="logo-mark"><LogoMark :size="22" /></div>
        <div class="logo-name">
          <span class="logo-text">goodBaby</span>
          <span class="logo-sub">摇篮系统</span>
        </div>
      </div>

      <nav class="nav">
        <RouterLink
          v-for="item in menus"
          :key="item.path"
          :to="item.path"
          class="nav-item"
          :class="{ active: activePath === item.path }"
        >
          <el-icon :size="17"><component :is="item.icon" /></el-icon>
          <span>{{ item.title }}</span>
        </RouterLink>
      </nav>

      <div class="aside-footer">
        <div class="footer-quote">goodBaby v2</div>
      </div>
    </aside>

    <!-- 主区域 -->
    <div class="content">
      <header class="header">
        <div class="header-left">
          <button class="icon-btn menu-btn" title="菜单" @click="drawerOpen = true">
            <el-icon :size="18"><Fold /></el-icon>
          </button>
          <div class="header-title">{{ route.meta.title ?? '' }}</div>
        </div>
        <div class="header-actions">
          <button
            class="icon-btn theme-btn"
            :title="isDark ? '切换到亮色' : '切换到暗色'"
            @click="toggleTheme"
          >
            <Transition name="theme-icon" mode="out-in">
              <el-icon v-if="isDark" :size="17" key="sun"><Sunny /></el-icon>
              <el-icon v-else :size="17" key="moon"><Moon /></el-icon>
            </Transition>
          </button>
          <el-dropdown>
            <span class="user-chip">
              <span class="user-avatar">{{ userStore.user?.username?.[0]?.toUpperCase() }}</span>
              <span class="user-name">{{ userStore.user?.username }}</span>
            </span>
            <template #dropdown>
              <el-dropdown-menu>
                <el-dropdown-item @click="router.push('/settings')">
                  <el-icon><Setting /></el-icon>设置
                </el-dropdown-item>
                <el-dropdown-item divided @click="handleLogout">
                  <el-icon><SwitchButton /></el-icon>退出登录
                </el-dropdown-item>
              </el-dropdown-menu>
            </template>
          </el-dropdown>
        </div>
      </header>

      <main class="main">
        <RouterView v-slot="{ Component }">
          <Transition name="page" mode="out-in">
            <component :is="Component" />
          </Transition>
        </RouterView>
      </main>
    </div>
  </div>
</template>

<style scoped>
.layout {
  height: 100%;
  display: flex;
}

/* 移动端抽屉遮罩 */
.drawer-mask {
  position: fixed;
  inset: 0;
  z-index: 1999;
  background: rgb(0 0 0 / 0.45);
  backdrop-filter: blur(2px);
}

.fade-enter-active,
.fade-leave-active {
  transition: opacity 0.25s var(--gb-ease);
}

.fade-enter-from,
.fade-leave-to {
  opacity: 0;
}

/* 汉堡按钮只在移动端出现。
   用复合选择器压过后面的 .icon-btn { display: flex }，避免受声明顺序影响。 */
.icon-btn.menu-btn {
  display: none;
}

.header-left {
  display: flex;
  align-items: center;
  gap: 8px;
  min-width: 0;
}

/* ---------- 侧边栏：墨色 + 药丸导航 ---------- */

.aside {
  width: 224px;
  flex-shrink: 0;
  background: linear-gradient(180deg, var(--gb-ink) 0%, var(--gb-ink-soft) 100%);
  display: flex;
  flex-direction: column;
  padding: 20px 14px;
  box-sizing: border-box;
}

.logo {
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 4px 8px 22px;
}

.logo-mark {
  width: 38px;
  height: 38px;
  border-radius: 12px;
  background: linear-gradient(135deg, var(--gb-primary) 0%, #3fa9e0 100%);
  display: flex;
  align-items: center;
  justify-content: center;
  color: #0b2b3d;
  box-shadow: 0 4px 12px rgb(102 204 255 / 0.45);
}

.logo-name {
  display: flex;
  flex-direction: column;
  line-height: 1.25;
}

.logo-text {
  color: #fff;
  font-size: 17px;
  font-weight: 700;
  letter-spacing: -0.01em;
}

.logo-sub {
  color: rgb(255 255 255 / 0.45);
  font-size: 11px;
  letter-spacing: 0.2em;
}

.nav {
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.nav-item {
  display: flex;
  align-items: center;
  gap: 11px;
  padding: 10px 13px;
  border-radius: 10px;
  color: rgb(255 255 255 / 0.62);
  text-decoration: none;
  font-size: 14px;
  transition:
    background 0.2s var(--gb-ease),
    color 0.2s var(--gb-ease),
    transform 0.2s var(--gb-ease);
}

.nav-item:hover {
  background: rgb(255 255 255 / 0.07);
  color: rgb(255 255 255 / 0.9);
  transform: translateX(2px);
}

.nav-item.active {
  background: linear-gradient(135deg, var(--gb-primary) 0%, #3fa9e0 100%);
  color: #0b2b3d;
  font-weight: 600;
  box-shadow: 0 4px 14px rgb(102 204 255 / 0.4);
}

.aside-footer {
  margin-top: auto;
  padding: 12px 8px 4px;
}

.footer-quote {
  color: rgb(255 255 255 / 0.28);
  font-size: 12px;
  line-height: 1.9;
  letter-spacing: 0.04em;
}

/* ---------- 主区域 ---------- */

.content {
  flex: 1;
  min-width: 0;
  display: flex;
  flex-direction: column;
}

.header {
  height: 60px;
  flex-shrink: 0;
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 0 28px;
  background: color-mix(in srgb, var(--gb-card) 82%, transparent);
  backdrop-filter: blur(10px);
  border-bottom: 1px solid var(--el-border-color-extra-light);
  position: sticky;
  top: 0;
  z-index: 10;
}

.header-title {
  font-size: 15px;
  font-weight: 600;
  color: var(--el-text-color-primary);
}

.header-actions {
  display: flex;
  align-items: center;
  gap: 14px;
}

.icon-btn {
  width: 34px;
  height: 34px;
  border: none;
  border-radius: 10px;
  background: transparent;
  color: var(--el-text-color-regular);
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  transition: background 0.2s var(--gb-ease);
}

.icon-btn:hover {
  background: var(--el-fill-color);
}

/* 主题按钮：hover 时图标轻微旋转，点击有回弹 */
.theme-btn .el-icon {
  transition: transform 0.4s var(--gb-ease);
}

.theme-btn:hover .el-icon {
  transform: rotate(25deg) scale(1.1);
}

.theme-btn:active {
  transform: scale(0.9);
}

/* 图标切换：旧图标转出，新图标转入 */
.theme-icon-enter-active,
.theme-icon-leave-active {
  transition:
    opacity 0.2s var(--gb-ease),
    transform 0.28s var(--gb-ease);
}

.theme-icon-enter-from {
  opacity: 0;
  transform: rotate(-90deg) scale(0.5);
}

.theme-icon-leave-to {
  opacity: 0;
  transform: rotate(90deg) scale(0.5);
}

.user-chip {
  display: inline-flex;
  align-items: center;
  gap: 8px;
  cursor: pointer;
  color: var(--el-text-color-primary);
  font-size: 14px;
  outline: none;
  padding: 4px 10px 4px 4px;
  border-radius: 999px;
  transition: background 0.2s var(--gb-ease);
}

.user-chip:hover {
  background: var(--el-fill-color);
}

.user-avatar {
  width: 28px;
  height: 28px;
  border-radius: 50%;
  background: linear-gradient(135deg, var(--gb-primary) 0%, #3fa9e0 100%);
  color: #0b2b3d;
  font-size: 13px;
  font-weight: 700;
  display: flex;
  align-items: center;
  justify-content: center;
}

.main {
  flex: 1;
  overflow-y: auto;
  padding: 24px 28px 40px;
  -webkit-overflow-scrolling: touch;
}

/* ---------- 移动端 ---------- */

@media (max-width: 768px) {
  /* 侧边栏脱离文档流，变成从左侧滑出的抽屉 */
  .aside {
    position: fixed;
    top: 0;
    left: 0;
    bottom: 0;
    z-index: 2000;
    width: 250px;
    transform: translateX(-100%);
    transition: transform 0.28s var(--gb-ease);
    box-shadow: 4px 0 24px rgb(0 0 0 / 0.25);
  }

  .aside.open {
    transform: translateX(0);
  }

  .icon-btn.menu-btn {
    display: flex;
  }

  .header {
    height: 54px;
    padding: 0 14px;
  }

  .header-title {
    font-size: 14px;
    white-space: nowrap;
    overflow: hidden;
    text-overflow: ellipsis;
  }

  .header-actions {
    gap: 6px;
  }

  /* 窄屏只留头像，省出空间给标题 */
  .user-name {
    display: none;
  }

  .user-chip {
    padding: 4px;
  }

  .main {
    padding: 16px 14px 32px;
  }
}
</style>
