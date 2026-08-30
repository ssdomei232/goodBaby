// 与后端 model / internal/meta 对应的类型定义

export interface ApiResponse<T = unknown> {
  code: number
  data: T
}

export interface UserInfo {
  id: number
  create_at: number
  username: string
  is_admin: boolean
	dingtalk_config: string | null
	api_key: string
}

export interface Timer {
  id: number
  uid: number
  name: string
  description: string
  enabled: boolean
  sign_deration_seconds: number
  last_sign: number
  remind_time_seconds: number
  last_remind: number
  triggered: boolean
  last_trigger: number
  create_at: number
}

export interface TimerRequest {
  name: string
  description: string
  enabled?: boolean
  sign_deration_seconds: number
  remind_time_seconds: number
}

export interface Rule {
  id: number
  uid: number
  timer_id: number
  gateway_id: number
  name: string
  account_id: number
  type: string
  config_json: string
  enabled: boolean
  create_at: number
}

export interface MessageGateway { id: number; uid: number; name: string; token: string; create_at: number }

export interface RuleRequest {
  name: string
  timer_id: number
  gateway_id?: number
  account_id: number
  type: string
  config_json: string
  enabled?: boolean
}

export interface Account {
  id: number
  uid: number
  name: string
  type: string
  config: string
  create_at: number
}

export interface AccountRequest {
  name: string
  type: string
  config: string
}

export type FieldType =
  | 'string'
  | 'password'
  | 'textarea'
  | 'number'
  | 'string-list'
  | 'number-list'
  | 'bool'

export interface MetaField {
  key: string
  label: string
  type: FieldType
  required: boolean
  placeholder?: string
  help?: string
  secret?: boolean
  default?: unknown
}

export interface AccountMeta {
  type: string
  label: string
  description?: string
  docs?: string
  testable: boolean
  fields: MetaField[]
}

export interface RuleMeta {
  type: string
  label: string
  description?: string
  docs?: string
  account_type: string
  fields: MetaField[]
}

export interface Providers {
  accounts: AccountMeta[]
  rules: RuleMeta[]
}

export interface ExecutionLog {
  id: number
  uid: number
  rule_id: number
  rule_name: string
  rule_type: string
  timer_id: number
  trigger: 'timer' | 'manual' | 'remind'
  success: boolean
  message: string
  create_at: number
}

export interface LogPage {
  total: number
  page: number
  page_size: number
  items: ExecutionLog[]
}

export interface DashboardOverview {
  timer_count: number
  enabled_timers: number
  triggered_count: number
  rule_count: number
  account_count: number
  server_time: number
  recent_logs: ExecutionLog[]
  urgent_timer?: Timer
  urgent_seconds_left?: number
}

export interface SiteInfo {
  name: string
  enable_registry: boolean
  need_initial_user: boolean
  check_interval_minutes: number
}

/** 管理员可在 WebUI 修改的系统配置 */
export interface AdminConfig {
  enable_registry: boolean
  timeout_duration_hours: number
  check_interval_minutes: number
  log_retain_count: number
}

export interface AdminConfigResponse {
  config: AdminConfig
  /** 只读的运行时信息，只能改配置文件 */
  readonly: {
    listen_addr: string
    database_driver: string
  }
}

/** 后端对 secret 字段的掩码占位符，原样提交时表示保留旧值 */
export const MASK_PLACEHOLDER = '********'
