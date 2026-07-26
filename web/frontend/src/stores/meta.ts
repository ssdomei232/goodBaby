import { defineStore } from 'pinia'
import { computed, ref } from 'vue'
import { providerApi } from '@/api'
import type { AccountMeta, Providers, RuleMeta } from '@/api/types'

/** 驱动元数据：账号/规则类型的表单描述，登录后拉取一次全局共享 */
export const useMetaStore = defineStore('meta', () => {
  const providers = ref<Providers | null>(null)

  async function ensureLoaded() {
    if (!providers.value) {
      providers.value = await providerApi.all()
    }
    return providers.value
  }

  const accountMetas = computed<AccountMeta[]>(() => providers.value?.accounts ?? [])
  const ruleMetas = computed<RuleMeta[]>(() => providers.value?.rules ?? [])

  function accountMeta(type: string): AccountMeta | undefined {
    return accountMetas.value.find((m) => m.type === type)
  }

  function ruleMeta(type: string): RuleMeta | undefined {
    return ruleMetas.value.find((m) => m.type === type)
  }

  function accountLabel(type: string): string {
    return accountMeta(type)?.label ?? type
  }

  function ruleLabel(type: string): string {
    return ruleMeta(type)?.label ?? type
  }

  return {
    providers,
    accountMetas,
    ruleMetas,
    ensureLoaded,
    accountMeta,
    ruleMeta,
    accountLabel,
    ruleLabel,
  }
})
