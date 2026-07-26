import { onUnmounted, ref } from 'vue'

/** 移动端断点，与 main.css 中的媒体查询保持一致 */
export const MOBILE_BREAKPOINT = 768

/**
 * 响应式的移动端判断。
 *
 * 用于那些光靠 CSS 改不了的地方，比如表格列的显隐、分页器布局、表单标签位置。
 *
 * 除了监听 matchMedia 的 change，额外挂一个 resize 兜底：
 * 某些环境(内嵌 WebView、开发者工具的设备模拟)改视口时不会派发 change，
 * 只靠 change 会导致旋转屏幕后布局不跟着变。
 */
export function useIsMobile() {
  const query = window.matchMedia(`(max-width: ${MOBILE_BREAKPOINT}px)`)
  const isMobile = ref(query.matches)

  const sync = () => {
    if (isMobile.value !== query.matches) {
      isMobile.value = query.matches
    }
  }

  query.addEventListener('change', sync)
  window.addEventListener('resize', sync, { passive: true })
  window.addEventListener('orientationchange', sync)

  onUnmounted(() => {
    query.removeEventListener('change', sync)
    window.removeEventListener('resize', sync)
    window.removeEventListener('orientationchange', sync)
  })

  return isMobile
}
