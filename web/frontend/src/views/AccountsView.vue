<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Connection, Plus } from '@element-plus/icons-vue'
import { accountApi } from '@/api'
import { ApiError } from '@/api/client'
import type { Account } from '@/api/types'
import { useMetaStore } from '@/stores/meta'
import { formatDateTime } from '@/utils/format'
import ConfigForm from '@/components/ConfigForm.vue'
import { useIsMobile } from '@/composables/useBreakpoint'

const isMobile = useIsMobile()

const metaStore = useMetaStore()

const accounts = ref<Account[]>([])
const loading = ref(false)
const testingId = ref<number | null>(null)

// ---- 创建 / 编辑对话框 ----
const dialogVisible = ref(false)
const saving = ref(false)
const editingId = ref<number | null>(null)

const form = ref({
  name: '',
  type: '',
  config: '',
})

const dialogTitle = computed(() => (editingId.value === null ? '添加账号' : '编辑账号'))
const currentMeta = computed(() => metaStore.accountMeta(form.value.type))

async function refresh() {
  loading.value = true
  try {
    accounts.value = await accountApi.list()
  } catch (error) {
    ElMessage.error(error instanceof ApiError ? error.message : '加载失败')
  } finally {
    loading.value = false
  }
}

function openCreate() {
  editingId.value = null
  form.value = { name: '', type: metaStore.accountMetas[0]?.type ?? '', config: '' }
  dialogVisible.value = true
}

function openEdit(account: Account) {
  editingId.value = account.id
  form.value = { name: account.name, type: account.type, config: account.config }
  dialogVisible.value = true
}

async function save() {
  if (!form.value.name || !form.value.type) {
    ElMessage.warning('请填写账号名称并选择类型')
    return
  }

  saving.value = true
  try {
    if (editingId.value === null) {
      await accountApi.create(form.value)
      ElMessage.success('账号已添加')
    } else {
      await accountApi.update(editingId.value, form.value)
      ElMessage.success('账号已更新')
    }
    dialogVisible.value = false
    await refresh()
  } catch (error) {
    ElMessage.error(error instanceof ApiError ? error.message : '保存失败')
  } finally {
    saving.value = false
  }
}

async function test(account: Account) {
  testingId.value = account.id
  try {
    await accountApi.test(account.id)
    ElMessage.success(`「${account.name}」测试通过`)
  } catch (error) {
    ElMessage.error(error instanceof ApiError ? error.message : '测试失败')
  } finally {
    testingId.value = null
  }
}

async function remove(account: Account) {
  try {
    const affected = await accountApi.checkDelete(account.id)
    const warning =
      affected.length > 0
        ? `删除后，使用该账号的 ${affected.length} 条规则也会被一并删除：${affected
            .map((r) => r.name)
            .join('、')}`
        : '没有规则使用该账号。'

    await ElMessageBox.confirm(warning, `删除账号「${account.name}」？`, {
      type: 'warning',
      confirmButtonText: '删除',
      confirmButtonClass: 'el-button--danger',
      cancelButtonText: '取消',
    })
  } catch {
    return
  }

  try {
    await accountApi.remove(account.id)
    ElMessage.success('已删除')
    await refresh()
  } catch (error) {
    ElMessage.error(error instanceof ApiError ? error.message : '删除失败')
  }
}

onMounted(async () => {
  await metaStore.ensureLoaded()
  await refresh()
})
</script>

<template>
  <div v-loading="loading">
    <div class="page-header">
      <div>
        <h2>账号</h2>
        <div class="muted">保存第三方服务的凭据，供规则执行时使用。敏感字段保存后显示为 ********</div>
      </div>
      <el-button type="primary" :icon="Plus" @click="openCreate">添加账号</el-button>
    </div>

    <el-empty v-if="!loading && accounts.length === 0" description="还没有账号">
      <el-button type="primary" @click="openCreate">添加账号</el-button>
    </el-empty>

    <el-card v-else class="gb-rise">
    <el-table :data="accounts" :size="isMobile ? 'small' : 'default'">
      <el-table-column prop="name" label="名称" min-width="110" />
      <el-table-column label="类型" :width="isMobile ? 110 : 160">
        <template #default="{ row }">
          <el-tag>{{ metaStore.accountLabel(row.type) }}</el-tag>
        </template>
      </el-table-column>
      <!-- 窄屏隐藏创建时间，优先保证操作列可见 -->
      <el-table-column v-if="!isMobile" label="创建时间" width="180">
        <template #default="{ row }">{{ formatDateTime(row.create_at) }}</template>
      </el-table-column>
      <el-table-column label="操作" :width="isMobile ? 210 : 260" fixed="right">
        <template #default="{ row }">
          <el-button
            v-if="metaStore.accountMeta(row.type)?.testable"
            size="small"
            :icon="Connection"
            :loading="testingId === row.id"
            @click="test(row)"
          >
            测试
          </el-button>
          <el-button size="small" @click="openEdit(row)">编辑</el-button>
          <el-button size="small" type="danger" plain @click="remove(row)">删除</el-button>
        </template>
      </el-table-column>
    </el-table>
    </el-card>

    <!-- 创建/编辑对话框 -->
    <el-dialog v-model="dialogVisible" :title="dialogTitle" width="520px">
      <el-form :label-width="isMobile ? 'auto' : '110px'" :label-position="isMobile ? 'top' : 'right'">
        <el-form-item label="账号名称" required>
          <el-input v-model="form.name" placeholder="给这个账号起个名字" maxlength="64" />
        </el-form-item>
        <el-form-item label="账号类型" required>
          <el-select
            v-model="form.type"
            :disabled="editingId !== null"
            style="width: 100%"
            @change="form.config = ''"
          >
            <el-option
              v-for="m in metaStore.accountMetas"
              :key="m.type"
              :label="m.label"
              :value="m.type"
            />
          </el-select>
          <div v-if="currentMeta?.description" class="muted">{{ currentMeta.description }}</div>
        </el-form-item>

        <ConfigForm
          v-if="currentMeta"
          v-model="form.config"
          :fields="currentMeta.fields"
        />
      </el-form>
      <template #footer>
        <el-button @click="dialogVisible = false">取消</el-button>
        <el-button type="primary" :loading="saving" @click="save">保存</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<style scoped>
</style>
