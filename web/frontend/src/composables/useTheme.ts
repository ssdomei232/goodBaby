import { ref } from 'vue'

const STORAGE_KEY = 'gb-theme'

/** 全局共享的主题状态，模块级单例，任何组件读到的都是同一份 */
const isDark = ref(false)

/** 把主题应用到 <html> 上 */
function applyTheme(dark: boolean) {
  document.documentElement.classList.toggle('dark', dark)
  isDark.value = dark
}

/**
 * 在应用启动时调用一次：恢复上次选择，没有记录时跟随系统。
 */
export function initTheme() {
  const saved = localStorage.getItem(STORAGE_KEY)
  const prefersDark = window.matchMedia('(prefers-color-scheme: dark)').matches
  applyTheme(saved ? saved === 'dark' : prefersDark)

  // 标记浏览器是否支持圆形揭示，供 CSS 决定要不要退化成淡变
  if (typeof document.startViewTransition !== 'function') {
    document.documentElement.classList.add('no-vt')
  }

  // 用户没有手动选择过时，跟随系统变化
  window.matchMedia('(prefers-color-scheme: dark)').addEventListener('change', (e) => {
    if (!localStorage.getItem(STORAGE_KEY)) {
      applyTheme(e.matches)
    }
  })
}

export function useTheme() {
  /**
   * 切换明暗主题。
   *
   * 支持 View Transitions 时，以点击处为圆心做一次圆形揭示动画；
   * 不支持(Firefox / Safari 旧版)或用户要求减少动效时直接切换。
   */
  async function toggleTheme(event?: MouseEvent) {
    const next = !isDark.value
    const persist = () => {
      applyTheme(next)
      localStorage.setItem(STORAGE_KEY, next ? 'dark' : 'light')
    }

    // 绑定后再判空，TS 才能正确收窄这个可选 API
    const startViewTransition = document.startViewTransition?.bind(document)
    const reduceMotion = window.matchMedia('(prefers-reduced-motion: reduce)').matches

    if (!startViewTransition || reduceMotion) {
      persist()
      return
    }

    // 以点击点为圆心；没有事件对象时从右上角(切换按钮所在位置)扩散
    const x = event?.clientX ?? window.innerWidth - 60
    const y = event?.clientY ?? 30
    // 半径取到最远的那个角，保证能覆盖整个视口
    const endRadius = Math.hypot(
      Math.max(x, window.innerWidth - x),
      Math.max(y, window.innerHeight - y),
    )

    const transition = startViewTransition(persist)

    try {
      await transition.ready
    } catch {
      // 浏览器跳过了这次过渡（页面不可见、连续快速点击等）。
      // 此时主题已经通过 persist() 应用好了，只是没有揭示动画，直接返回。
      return
    }

    const clipPath = [`circle(0px at ${x}px ${y}px)`, `circle(${endRadius}px at ${x}px ${y}px)`]

    document.documentElement.animate(
      // 切到深色时让新画面扩散进来；切回浅色时让旧画面收缩出去
      { clipPath: next ? clipPath : [...clipPath].reverse() },
      {
        duration: 480,
        easing: 'cubic-bezier(0.4, 0, 0.2, 1)',
        pseudoElement: next ? '::view-transition-new(root)' : '::view-transition-old(root)',
      },
    )
  }

  return { isDark, toggleTheme }
}
