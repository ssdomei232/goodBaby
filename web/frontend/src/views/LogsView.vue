<script setup lang="ts">
import { onMounted, ref, watch } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Delete, Refresh } from '@element-plus/icons-vue'
import { logApi } from '@/api'
import { ApiError } from '@/api/client'
import type { ExecutionLog } from '@/api/types'
import { formatDateTime } from '@/utils/format'
import { useIsMobile } from '@/composables/useBreakpoint'

const isMobile = useIsMobile()

const logs = ref<ExecutionLog[]>([])
const total = ref(0)
const page = ref(1)
const pageSize = ref(20)
const successFilter = ref<'all' | 'true' | 'false'>('all')
const loading = ref(false)

const triggerLabels: Record<string, string> = {
  timer: '定时触发',
  manual: '手动测试',
  remind: '提醒',
}

async function refresh() {
  loading.value = true
  try {
    const result = await logApi.list({
      page: page.value,
      page_size: pageSize.value,
      success: successFilter.value === 'all' ? undefined : successFilter.value,
    })
    logs.value = result.items
    total.value = result.total
  } catch (error) {
    ElMessage.error(error instanceof ApiError ? error.message : '加载失败')
  } finally {
    loading.value = false
  }
}

async function clearAll() {
  try {
    await ElMessageBox.confirm('确定清空所有执行日志吗？该操作不可恢复。', '清空日志', {
      type: 'warning',
      confirmButtonText: '清空',
      confirmButtonClass: 'el-button--danger',
      cancelButtonText: '取消',
    })
  } catch {
    return
  }

  try {
    await logApi.clear()
    ElMessage.success('已清空')
    page.value = 1
    await refresh()
  } catch (error) {
    ElMessage.error(error instanceof ApiError ? error.message : '清空失败')
  }
}

watch([page, pageSize, successFilter], refresh)

onMounted(refresh)
</script>

<template>
  <div v-loading="loading">
    <div class="page-header">
      <div>
        <h2>执行日志</h2>
        <div class="muted">规则执行与提醒的历史记录</div>
      </div>
      <div class="header-tools">
        <el-select v-model="successFilter" class="filter-select">
          <el-option label="全部" value="all" />
          <el-option label="仅成功" value="true" />
          <el-option label="仅失败" value="false" />
        </el-select>
        <el-button :icon="Refresh" @click="refresh">刷新</el-button>
        <el-button type="danger" plain :icon="Delete" @click="clearAll">清空</el-button>
      </div>
    </div>

    <el-empty v-if="!loading && logs.length === 0" description="暂无执行记录" />

    <template v-else>
      <el-card class="gb-rise">
      <el-table :data="logs" :size="isMobile ? 'small' : 'default'">
        <el-table-column label="时间" :width="isMobile ? 130 : 175">
          <template #default="{ row }">
            {{ isMobile ? formatDateTime(row.create_at).slice(5) : formatDateTime(row.create_at) }}
          </template>
        </el-table-column>
        <el-table-column prop="rule_name" label="规则 / 定时器" min-width="130" />
        <!-- 窄屏隐藏次要列，避免横向滚动过长 -->
        <el-table-column v-if="!isMobile" prop="rule_type" label="类型" width="150" />
        <el-table-column v-if="!isMobile" label="触发方式" width="110">
          <template #default="{ row }">
            {{ triggerLabels[row.trigger] ?? row.trigger }}
          </template>
        </el-table-column>
        <el-table-column label="结果" :width="isMobile ? 70 : 90">
          <template #default="{ row }">
            <el-tag :type="row.success ? 'success' : 'danger'" size="small">
              {{ row.success ? '成功' : '失败' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column
          prop="message"
          label="详情"
          :min-width="isMobile ? 180 : 260"
          show-overflow-tooltip
        />
      </el-table>
      </el-card>

      <div class="pager">
        <el-pagination
          v-model:current-page="page"
          v-model:page-size="pageSize"
          :total="total"
          :page-sizes="[20, 50, 100]"
          :layout="isMobile ? 'prev, pager, next' : 'total, sizes, prev, pager, next'"
          :pager-count="isMobile ? 5 : 7"
          :small="isMobile"
          background
        />
      </div>
    </template>
  </div>
</template>

<style scoped>
.header-tools {
  display: flex;
  gap: 12px;
}

.filter-select {
  width: 140px;
}

.pager {
  margin-top: 16px;
  display: flex;
  justify-content: flex-end;
}

@media (max-width: 768px) {
  /* 筛选下拉占一行，两个按钮平分下一行 */
  .header-tools {
    flex-wrap: wrap;
    gap: 10px;
  }

  .filter-select {
    width: 100%;
  }

  .header-tools .el-button {
    flex: 1;
    margin-left: 0;
  }

  .pager {
    justify-content: center;
  }
}
</style>
