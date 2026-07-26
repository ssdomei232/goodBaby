/** 时间与时长的展示工具 */

export function formatDateTime(unixSeconds: number): string {
  if (!unixSeconds) return '-'
  const date = new Date(unixSeconds * 1000)
  const pad = (n: number) => String(n).padStart(2, '0')
  return `${date.getFullYear()}-${pad(date.getMonth() + 1)}-${pad(date.getDate())} ${pad(date.getHours())}:${pad(date.getMinutes())}:${pad(date.getSeconds())}`
}

/** 把秒数格式化为 "3天 4小时" 这样的可读时长 */
export function formatDuration(totalSeconds: number): string {
  if (totalSeconds <= 0) return '0秒'

  const units: Array<[number, string]> = [
    [86400, '天'],
    [3600, '小时'],
    [60, '分钟'],
    [1, '秒'],
  ]

  const parts: string[] = []
  let remain = Math.floor(totalSeconds)
  for (const [size, label] of units) {
    if (remain >= size) {
      parts.push(`${Math.floor(remain / size)}${label}`)
      remain %= size
    }
    if (parts.length >= 2) break
  }
  return parts.join(' ') || '0秒'
}

/** 常用签到周期预设(秒) */
export const durationPresets: Array<{ label: string; seconds: number }> = [
  { label: '1 天', seconds: 86400 },
  { label: '3 天', seconds: 3 * 86400 },
  { label: '7 天', seconds: 7 * 86400 },
  { label: '14 天', seconds: 14 * 86400 },
  { label: '30 天', seconds: 30 * 86400 },
  { label: '90 天', seconds: 90 * 86400 },
]
