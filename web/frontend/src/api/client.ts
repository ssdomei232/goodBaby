import axios, { AxiosError } from 'axios'
import type { ApiResponse } from './types'

/** 业务错误：后端返回了非 200 的 code，message 可直接展示给用户 */
export class ApiError extends Error {
  code: number

  constructor(code: number, message: string) {
    super(message)
    this.code = code
  }
}

const http = axios.create({
  baseURL: '/api/v1',
  timeout: 60_000,
  withCredentials: true,
})

/** 未登录时的回调，由 router 注入，跳转到登录页 */
let onUnauthorized: (() => void) | null = null

export function setUnauthorizedHandler(handler: () => void) {
  onUnauthorized = handler
}

async function request<T>(promise: Promise<{ data: ApiResponse<T> }>): Promise<T> {
  try {
    const response = await promise
    return response.data.data
  } catch (error) {
    const axiosError = error as AxiosError<ApiResponse<string>>
    const status = axiosError.response?.status ?? 0
    const data = axiosError.response?.data

    if (status === 401) {
      onUnauthorized?.()
    }

    const message =
      typeof data?.data === 'string' && data.data !== ''
        ? data.data
        : axiosError.message || '网络错误'
    throw new ApiError(status, message)
  }
}

export const api = {
  get: <T>(url: string, params?: Record<string, unknown>) =>
    request<T>(http.get<ApiResponse<T>>(url, { params })),
  post: <T>(url: string, body?: unknown) =>
    request<T>(http.post<ApiResponse<T>>(url, body)),
  put: <T>(url: string, body?: unknown) =>
    request<T>(http.put<ApiResponse<T>>(url, body)),
  delete: <T>(url: string) => request<T>(http.delete<ApiResponse<T>>(url)),
}
