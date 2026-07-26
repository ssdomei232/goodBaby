import { api } from './client'
import type {
  Account,
  AccountRequest,
  AdminConfig,
  AdminConfigResponse,
  DashboardOverview,
  ExecutionLog,
  LogPage,
  Providers,
  Rule,
  RuleRequest,
  SiteInfo,
  Timer,
  TimerRequest,
  UserInfo,
} from './types'

export const siteApi = {
  info: () => api.get<SiteInfo>('/site'),
}

export const userApi = {
  register: (username: string, password: string) =>
    api.post<string>('/user/registry', { username, password }),
  login: (username: string, password: string) =>
    api.post<string>('/user/login', { username, password }),
  logout: () => api.post<string>('/user/logout'),
  info: () => api.get<UserInfo>('/user/info'),
  changePassword: (oldPassword: string, newPassword: string) =>
    api.post<string>('/user/password', {
      old_password: oldPassword,
      new_password: newPassword,
    }),
  updateNotify: (dingtalkConfig: string | null) =>
    api.put<string>('/user/notify', { dingtalk_config: dingtalkConfig }),
}

export const dashboardApi = {
  overview: () => api.get<DashboardOverview>('/dashboard'),
}

export const providerApi = {
  all: () => api.get<Providers>('/providers'),
}

export const timerApi = {
  list: () => api.get<Timer[]>('/timers/'),
  get: (id: number) => api.get<{ timer: Timer; rule_count: number }>(`/timers/${id}`),
  create: (body: TimerRequest) => api.post<Timer>('/timers/', body),
  update: (id: number, body: TimerRequest) => api.put<Timer>(`/timers/${id}`, body),
  sign: (id: number) => api.post<Timer>(`/timers/${id}/sign`),
  signAll: () => api.post<{ signed: number; last_sign: number }>('/timers/sign'),
  trigger: (id: number) => api.post<{ total: number; failed: string[] }>(`/timers/${id}/trigger`),
  checkDelete: (id: number) => api.get<Rule[]>(`/timers/${id}/check`),
  remove: (id: number) => api.delete<string>(`/timers/${id}`),
}

export const ruleApi = {
  list: (timerId?: number) =>
    api.get<Rule[]>('/rules/', timerId ? { timer_id: timerId } : undefined),
  create: (body: RuleRequest) => api.post<Rule>('/rules/', body),
  update: (id: number, body: RuleRequest) => api.put<Rule>(`/rules/${id}`, body),
  test: (id: number) => api.post<string>(`/rules/${id}/test`),
  remove: (id: number) => api.delete<string>(`/rules/${id}`),
}

export const accountApi = {
  list: (type?: string) =>
    api.get<Account[]>('/accounts/', type ? { type } : undefined),
  create: (body: AccountRequest) => api.post<Account>('/accounts/', body),
  update: (id: number, body: AccountRequest) => api.put<Account>(`/accounts/${id}`, body),
  test: (id: number) => api.post<string>(`/accounts/${id}/test`),
  checkDelete: (id: number) => api.get<Rule[]>(`/accounts/${id}/check`),
  remove: (id: number) => api.delete<string>(`/accounts/${id}`),
}

export const logApi = {
  list: (params: { page?: number; page_size?: number; rule_id?: number; success?: string }) =>
    api.get<LogPage>('/logs/', params),
  clear: () => api.delete<string>('/logs/'),
}

export const adminApi = {
  getConfig: () => api.get<AdminConfigResponse>('/admin/config'),
  updateConfig: (body: AdminConfig) => api.put<AdminConfig>('/admin/config', body),
}

export type { ExecutionLog }
