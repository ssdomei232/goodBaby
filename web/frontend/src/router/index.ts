import { createRouter, createWebHistory } from 'vue-router'
import { setUnauthorizedHandler } from '@/api/client'
import { useUserStore } from '@/stores/user'

const router = createRouter({
  history: createWebHistory(),
  routes: [
    {
      path: '/login',
      name: 'login',
      component: () => import('@/views/LoginView.vue'),
      meta: { public: true },
    },
    {
      path: '/',
      component: () => import('@/layouts/MainLayout.vue'),
      children: [
        { path: '', redirect: '/dashboard' },
        {
          path: 'dashboard',
          name: 'dashboard',
          component: () => import('@/views/DashboardView.vue'),
          meta: { title: '仪表盘' },
        },
        {
          path: 'timers',
          name: 'timers',
          component: () => import('@/views/TimersView.vue'),
          meta: { title: '定时器' },
        },
        {
          path: 'rules',
          name: 'rules',
          component: () => import('@/views/RulesView.vue'),
          meta: { title: '规则' },
        },
        {
          path: 'accounts',
          name: 'accounts',
          component: () => import('@/views/AccountsView.vue'),
          meta: { title: '账号' },
        },
        { path: 'gateways', name: 'gateways', component: () => import('@/views/GatewaysView.vue'), meta: { title: '消息网关' } },
        {
          path: 'logs',
          name: 'logs',
          component: () => import('@/views/LogsView.vue'),
          meta: { title: '执行日志' },
        },
        {
          path: 'settings',
          name: 'settings',
          component: () => import('@/views/SettingsView.vue'),
          meta: { title: '设置' },
        },
      ],
    },
    { path: '/:pathMatch(.*)*', redirect: '/dashboard' },
  ],
})

router.beforeEach(async (to) => {
  if (to.meta.public) return true

  const userStore = useUserStore()
  if (!userStore.loaded) {
    await userStore.fetchUser()
  }
  if (!userStore.user) {
    return { name: 'login', query: { redirect: to.fullPath } }
  }
  return true
})

// API 返回 401 时清空用户状态并跳到登录页
setUnauthorizedHandler(() => {
  const userStore = useUserStore()
  userStore.clear()
  if (router.currentRoute.value.name !== 'login') {
    router.push({ name: 'login' })
  }
})

export default router
