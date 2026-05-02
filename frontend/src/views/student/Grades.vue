<template>
  <div class="grades">
    <el-card shadow="never">
      <template #header>
        <div class="card-header">
          <span>我的成绩</span>
        </div>
      </template>
      <el-table :data="tableData" v-loading="loading" style="width: 100%">
        <el-table-column prop="id" label="ID" width="80" />
        <el-table-column prop="retake_course" label="课程名称" min-width="200">
          <template #default="scope">
            {{ scope.row.retake_course?.course?.name || '-' }}
          </template>
        </el-table-column>
        <el-table-column prop="semester" label="学期" width="120">
          <template #default="scope">
            {{ scope.row.retake_course?.semester || '-' }}
          </template>
        </el-table-column>
        <el-table-column prop="score" label="成绩" width="100">
          <template #default="scope">
            <span v-if="scope.row.score !== null && scope.row.score !== undefined">
              <el-tag :type="scope.row.score >= 60 ? 'success' : 'danger'">
                {{ scope.row.score }}
              </el-tag>
            </span>
            <span v-else class="text-gray">未录入</span>
          </template>
        </el-table-column>
        <el-table-column prop="status" label="状态" width="120">
          <template #default="scope">
            <el-tag :type="getStatusType(scope.row.status)">
              {{ getStatusText(scope.row.status) }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="created_at" label="创建时间" width="180">
          <template #default="scope">
            {{ formatDate(scope.row.created_at) }}
          </template>
        </el-table-column>
      </el-table>
      <el-pagination
        v-model:current-page="pagination.page"
        v-model:page-size="pagination.pageSize"
        :page-sizes="[10, 20, 50]"
        :total="pagination.total"
        layout="total, sizes, prev, pager, next, jumper"
        @size-change="fetchData"
        @current-change="fetchData"
        style="margin-top: 20px; justify-content: flex-end;"
      />
    </el-card>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { getMyGrades } from '@/api'
import dayjs from 'dayjs'

const loading = ref(false)
const tableData = ref([])

const pagination = reactive({
  page: 1,
  pageSize: 10,
  total: 0
})

const formatDate = (date) => {
  return dayjs(date).format('YYYY-MM-DD HH:mm:ss')
}

const getStatusType = (status) => {
  const types = {
    pending: 'warning',
    submitted: 'primary',
    published: 'success'
  }
  return types[status] || 'info'
}

const getStatusText = (status) => {
  const texts = {
    pending: '待录入',
    submitted: '已提交',
    published: '已发布'
  }
  return texts[status] || status
}

const fetchData = async () => {
  loading.value = true
  try {
    const res = await getMyGrades({
      page: pagination.page,
      page_size: pagination.pageSize
    })
    tableData.value = res.data || []
    pagination.total = res.total || 0
  } catch (error) {
    console.error(error)
  } finally {
    loading.value = false
  }
}

onMounted(() => {
  fetchData()
})
</script>

<style scoped>
.grades {
  padding: 0;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.text-gray {
  color: #909399;
}
</style>
