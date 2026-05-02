<template>
  <div class="applications">
    <el-card shadow="never">
      <template #header>
        <div class="card-header">
          <span>我的重修申请</span>
          <el-select v-model="filterStatus" placeholder="状态筛选" clearable @change="fetchData" style="width: 150px; margin-left: 10px;">
            <el-option label="全部" value="" />
            <el-option label="待审核" value="pending" />
            <el-option label="已通过" value="approved" />
            <el-option label="已拒绝" value="rejected" />
            <el-option label="已取消" value="cancelled" />
          </el-select>
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
        <el-table-column prop="reason" label="申请理由" min-width="150" show-overflow-tooltip />
        <el-table-column prop="status" label="状态" width="120">
          <template #default="scope">
            <el-tag :type="getStatusType(scope.row.status)">
              {{ getStatusText(scope.row.status) }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="review_comment" label="审核意见" min-width="150" show-overflow-tooltip>
          <template #default="scope">
            {{ scope.row.review_comment || '-' }}
          </template>
        </el-table-column>
        <el-table-column prop="created_at" label="申请时间" width="180">
          <template #default="scope">
            {{ formatDate(scope.row.created_at) }}
          </template>
        </el-table-column>
        <el-table-column label="操作" width="100" fixed="right">
          <template #default="scope">
            <el-button
              v-if="scope.row.status === 'pending'"
              type="danger"
              size="small"
              link
              @click="handleCancel(scope.row)"
            >
              取消申请
            </el-button>
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
import { ElMessage, ElMessageBox } from 'element-plus'
import { getMyApplications, cancelApplication } from '@/api'
import dayjs from 'dayjs'

const loading = ref(false)
const tableData = ref([])
const filterStatus = ref('')

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
    approved: 'success',
    rejected: 'danger',
    cancelled: 'info'
  }
  return types[status] || 'info'
}

const getStatusText = (status) => {
  const texts = {
    pending: '待审核',
    approved: '已通过',
    rejected: '已拒绝',
    cancelled: '已取消'
  }
  return texts[status] || status
}

const fetchData = async () => {
  loading.value = true
  try {
    const params = {
      page: pagination.page,
      page_size: pagination.pageSize
    }
    if (filterStatus.value) {
      params.status = filterStatus.value
    }
    const res = await getMyApplications(params)
    tableData.value = res.data || []
    pagination.total = res.total || 0
  } catch (error) {
    console.error(error)
  } finally {
    loading.value = false
  }
}

const handleCancel = async (row) => {
  try {
    await ElMessageBox.confirm('确定要取消这个申请吗？', '提示', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    })
    await cancelApplication(row.id)
    ElMessage.success('申请已取消')
    fetchData()
  } catch (error) {
    if (error !== 'cancel') {
      console.error(error)
    }
  }
}

onMounted(() => {
  fetchData()
})
</script>

<style scoped>
.applications {
  padding: 0;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}
</style>
