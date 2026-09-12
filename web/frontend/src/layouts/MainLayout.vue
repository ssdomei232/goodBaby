<script setup lang="ts">
import { computed, ref, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import {
  Odometer,
  Timer as TimerIcon,
  Operation,
  User,
  Connection,
  Promotion,
  Document,
  Setting,
  SwitchButton,
  Moon,
  Sunny,
  Fold,
  Reading,
} from '@element-plus/icons-vue'
import { useUserStore } from '@/stores/user'
import { useTheme } from '@/composables/useTheme'
import LogoMark from '@/components/LogoMark.vue'
import ThemePicker from '@/components/ThemePicker.vue'
import ParanoiaBackdrop from '@/components/ParanoiaBackdrop.vue'
import { paranoiaNavMarks } from '@/data/paranoia'

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

const baseMenus = [
  { path: '/dashboard', title: '仪表盘', icon: Odometer },
  { path: '/timers', title: '定时器', icon: TimerIcon },
  { path: '/rules', title: '规则', icon: Operation },
  { path: '/accounts', title: '账号', icon: User },
  { path: '/gateways', title: '消息网关', icon: Connection },
  { path: '/gateway-rules', title: '网关规则', icon: Promotion },
  { path: '/logs', title: '执行日志', icon: Document },
  { path: '/settings', title: '设置', icon: Setting },
]

const activePath = computed(() => `/${route.path.split('/')[1] ?? ''}`)

// 主题状态是全局单例，初始化在 main.ts 里完成
const { isDark, isParanoia, toggleTheme } = useTheme()

/** 妄想症皮肤下多出一个「妄想症」展柜，其余菜单不变 */
const menus = computed(() =>
  isParanoia.value
    ? [...baseMenus, { path: '/paranoia', title: '妄想症', icon: Reading }]
    : baseMenus,
)

/** 每个功能对应的章节记号，只在妄想症皮肤下显示 */
function markOf(path: string) {
  return paranoiaNavMarks[path.slice(1)]
}

/** 当前所处章节，用于顶栏与侧边栏落款 */
const currentMark = computed(() => markOf(activePath.value) ?? paranoiaNavMarks.dashboard)

// 展柜属于皮肤的一部分：皮肤关掉（含直接输入地址进入）就退回仪表盘
watch(
  isParanoia,
  (on) => {
    if (!on && route.path.startsWith('/paranoia')) {
      router.push('/dashboard')
    }
  },
  { immediate: true },
)

async function handleLogout() {
  await userStore.logout()
  ElMessage.success('已退出登录')
  router.push({ name: 'login' })
}
</script>

<template>
  <div class="layout" :class="{ 'skin-paranoia': isParanoia }">
    <!-- 妄想症皮肤的背景层：噪点、暗角与飘落的石蒜花瓣 -->
    <ParanoiaBackdrop v-if="isParanoia" />

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
          <span class="logo-sub">{{ isParanoia ? '妄想症 · PARANOIA' : '摇篮系统' }}</span>
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
          <span class="nav-text">{{ item.title }}</span>
          <span v-if="isParanoia" class="nav-mark mono">
            {{ markOf(item.path)?.sigil }}{{ markOf(item.path)?.label }}
          </span>
        </RouterLink>
      </nav>

      <div class="aside-footer">
        <template v-if="isParanoia">
          <div class="footer-chapter">
            <span class="footer-chapter__sigil mono">{{ currentMark?.sigil }} {{ currentMark?.label }}</span>
            <span class="footer-chapter__title">{{ currentMark?.title }}</span>
          </div>
          <div class="footer-quote">九重妄想 · 幸存者 goodBaby v2</div>
        </template>
        <div v-else class="footer-quote">goodBaby v2</div>
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
          <span v-if="isParanoia" class="pa-mark header-mark">
            <em>{{ currentMark?.sigil }}</em>{{ currentMark?.title }}
          </span>
        </div>
        <div class="header-actions">
          <ThemePicker />
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
                  <el-icon><SwitchButton /></el-icon>
                  {{ isParanoia ? '八重回归 · 退出' : '退出登录' }}
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

/* ---------- 妄想症皮肤：侧边栏变成卷宗封皮 ---------- */

/* 背景层待在内容之下 */
.layout.skin-paranoia {
  position: relative;
  isolation: isolate;
}

.layout.skin-paranoia > .content {
  position: relative;
  z-index: 1;
}

.layout.skin-paranoia .aside {
  background:
    linear-gradient(180deg, rgb(224 36 94 / 0.18) 0%, transparent 36%),
    linear-gradient(180deg, var(--gb-ink) 0%, var(--gb-ink-soft) 100%);
  border-right: 1px solid rgb(224 36 94 / 0.18);
}

.layout.skin-paranoia .logo-mark {
  border-radius: 6px;
  background: linear-gradient(135deg, #e0245e 0%, #7a0f2c 100%);
  color: #fff;
  box-shadow: 0 4px 16px rgb(224 36 94 / 0.45);
}

.layout.skin-paranoia .user-avatar {
  background: linear-gradient(135deg, #e0245e 0%, #7a2ea0 100%);
  color: #fff;
}

.layout.skin-paranoia .logo-text {
  font-family: var(--pa-serif);
  letter-spacing: 0.04em;
}

.layout.skin-paranoia .logo-sub {
  color: rgb(242 86 138 / 0.75);
  font-family: var(--pa-mono);
  font-size: 10px;
  letter-spacing: 0.14em;
}

.layout.skin-paranoia .nav-item {
  border: 1px solid transparent;
  border-radius: 4px;
}

.layout.skin-paranoia .nav-item:hover {
  background: rgb(224 36 94 / 0.1);
  border-color: rgb(224 36 94 / 0.2);
}

.layout.skin-paranoia .nav-item.active {
  background: linear-gradient(135deg, #e0245e 0%, #8f1230 100%);
  color: #fff;
  box-shadow: 0 4px 16px rgb(224 36 94 / 0.4);
}

/* 每一项右侧标出它对应的「重数」 */
.layout.skin-paranoia .nav-mark {
  margin-left: auto;
  font-size: 10px;
  letter-spacing: 0.06em;
  opacity: 0.45;
}

.layout.skin-paranoia .nav-item.active .nav-mark {
  opacity: 0.85;
}

.layout.skin-paranoia .footer-chapter {
  display: flex;
  flex-direction: column;
  gap: 3px;
  padding: 0 0 10px;
  margin-bottom: 10px;
  border-bottom: 1px solid rgb(224 36 94 / 0.2);
}

.layout.skin-paranoia .footer-chapter__sigil {
  font-size: 10px;
  letter-spacing: 0.16em;
  color: rgb(242 86 138 / 0.85);
}

.layout.skin-paranoia .footer-chapter__title {
  font-family: var(--pa-serif);
  font-size: 15px;
  font-weight: 700;
  letter-spacing: 0.1em;
  color: rgb(255 255 255 / 0.88);
}

.layout.skin-paranoia .header-mark {
  margin-left: 4px;
  font-size: 10px;
}

/* ---------- 移动端 ---------- */

@media (min-width: 769px) {
  /* 桌面端侧边栏常驻，压到背景层上面 */
  .layout.skin-paranoia > .aside {
    position: relative;
    z-index: 1;
  }
}

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

  /* 窄屏标题优先，章节记号召回抽屉里看 */
  .layout.skin-paranoia .header-mark {
    display: none;
  }
}
</style>
