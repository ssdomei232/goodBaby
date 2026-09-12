import { computed, ref } from 'vue'

/** 明暗模式的存储键 */
const DARK_STORAGE_KEY = 'gb-theme'
/** 外观皮肤（主题）的存储键 */
const SKIN_STORAGE_KEY = 'gb-skin'

/**
 * 可选皮肤：
 * - default：原有的晴空配色
 * - paranoia：致敬《妄想症 Paranoia》系列曲的外观主题
 */
export type ThemeSkin = 'default' | 'paranoia'

export interface SkinOption {
  id: ThemeSkin
  /** 皮肤名 */
  name: string
  /** 一句话说明 */
  sub: string
  /** 预览用的色板 */
  swatches: string[]
}

export const skinOptions: SkinOption[] = [
  {
    id: 'default',
    name: '晴空',
    sub: 'goodBaby 默认外观',
    swatches: ['#66ccff', '#3fa9e0', '#232630'],
  },
  {
    id: 'paranoia',
    name: '妄想症',
    sub: 'Paranoia · 九重档案',
    swatches: ['#e0245e', '#7ee8e0', '#f2d16b'],
  },
]

/** 全局共享的主题状态，模块级单例，任何组件读到的都是同一份 */
const isDark = ref(false)
const skin = ref<ThemeSkin>('default')
const isParanoia = computed(() => skin.value === 'paranoia')

function readStorage(key: string): string | null {
  try {
    return localStorage.getItem(key)
  } catch {
    // 隐私模式 / 禁用存储时静默降级为「本次会话内有效」
    return null
  }
}

function writeStorage(key: string, value: string) {
  try {
    localStorage.setItem(key, value)
  } catch {
    // 同上
  }
}

/** 把明暗模式应用到 <html> 上 */
function applyDark(dark: boolean) {
  document.documentElement.classList.toggle('dark', dark)
  isDark.value = dark
}

/** 把皮肤应用到 <html> 上，配色覆盖写在 theme-paranoia.css 里 */
function applySkin(next: ThemeSkin) {
  document.documentElement.classList.toggle('paranoia', next === 'paranoia')
  skin.value = next
}

/**
 * 在应用启动时调用一次：恢复上次选择，没有记录时跟随系统。
 */
export function initTheme() {
  const savedDark = readStorage(DARK_STORAGE_KEY)
  const prefersDark = window.matchMedia('(prefers-color-scheme: dark)').matches
  applyDark(savedDark ? savedDark === 'dark' : prefersDark)

  applySkin(readStorage(SKIN_STORAGE_KEY) === 'paranoia' ? 'paranoia' : 'default')

  // 标记浏览器是否支持圆形揭示，供 CSS 决定要不要退化成淡变
  if (typeof document.startViewTransition !== 'function') {
    document.documentElement.classList.add('no-vt')
  }

  // 用户没有手动选择过时，跟随系统变化
  window.matchMedia('(prefers-color-scheme: dark)').addEventListener('change', (e) => {
    if (!readStorage(DARK_STORAGE_KEY)) {
      applyDark(e.matches)
    }
  })
}

/**
 * 用一次圆形揭示动画包住外观切换。
 *
 * 支持 View Transitions 时以点击处为圆心做圆形揭示；不支持（Firefox / Safari 旧版）
 * 或用户要求减少动效时直接切换。
 *
 * @param mutate 真正改变外观状态的函数
 * @param event  触发切换的点击事件，用来决定圆心
 * @param mode   dark/light 沿用原有的进/退方向；skin 需要额外把新画面提到最上层
 */
async function reveal(
  mutate: () => void,
  event: MouseEvent | undefined,
  mode: 'dark' | 'light' | 'skin',
) {
  const startViewTransition = document.startViewTransition?.bind(document)
  const reduceMotion = window.matchMedia('(prefers-reduced-motion: reduce)').matches

  if (!startViewTransition || reduceMotion) {
    mutate()
    return
  }

  // 换皮肤时新旧画面可能同为暗色，靠这个类强制让新画面浮到上面
  if (mode === 'skin') {
    document.documentElement.classList.add('vt-skin')
  }

  // 从点击点向外扩散；没有事件对象时从右上角（切换按钮所在位置）扩散
  const x = event?.clientX ?? window.innerWidth - 60
  const y = event?.clientY ?? 30
  // 半径取到最远的那个角，保证能覆盖整个视口
  const endRadius = Math.hypot(
    Math.max(x, window.innerWidth - x),
    Math.max(y, window.innerHeight - y),
  )

  const transition = startViewTransition(mutate)

  try {
    await transition.ready
  } catch {
    // 浏览器跳过了这次过渡（页面不可见、连续快速点击等）。
    // 此时外观已由 mutate 应用好，只是没有揭示动画。
    document.documentElement.classList.remove('vt-skin')
    return
  }

  const clipPath = [`circle(0px at ${x}px ${y}px)`, `circle(${endRadius}px at ${x}px ${y}px)`]

  // 切到深色时让新画面扩散进来；切回浅色时让旧画面收缩出去
  const pseudoElement =
    mode === 'light' ? '::view-transition-old(root)' : '::view-transition-new(root)'

  document.documentElement.animate(
    { clipPath: mode === 'light' ? [...clipPath].reverse() : clipPath },
    {
      duration: 480,
      easing: 'cubic-bezier(0.4, 0, 0.2, 1)',
      pseudoElement,
    },
  )

  try {
    await transition.finished
  } catch {
    // finished 在过渡被跳过时也会抛错，忽略即可
  }
  document.documentElement.classList.remove('vt-skin')
}

export function useTheme() {
  /**
   * 切换明暗主题。
   */
  async function toggleTheme(event?: MouseEvent) {
    const next = !isDark.value
    await reveal(
      () => {
        applyDark(next)
        writeStorage(DARK_STORAGE_KEY, next ? 'dark' : 'light')
      },
      event,
      next ? 'dark' : 'light',
    )
  }

  /**
   * 切换外观皮肤。传入点击事件时以点击处为圆心做揭示动画。
   */
  async function setSkin(next: ThemeSkin, event?: MouseEvent) {
    if (next === skin.value) return
    await reveal(
      () => {
        applySkin(next)
        writeStorage(SKIN_STORAGE_KEY, next)
      },
      event,
      'skin',
    )
  }

  return { isDark, isParanoia, skin, toggleTheme, setSkin }
}
