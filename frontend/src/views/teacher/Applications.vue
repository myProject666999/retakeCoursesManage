<template>
  <div class="teacher-applications">
    <el-card shadow="never">
      <template #header>
        <div class="card-header">
          <span>重修申请审核</span>
          <el-select v-model="filterStatus" placeholder="状态筛选" clearable @change="fetchData" style="width: 150px;">
            <el-option label="待审核" value="pending" />
            <el-option label="已通过" value="approved" />
            <el-option label="已拒绝" value="rejected" />
          </el-select>
        </div>
      </template>
      <el-table :data="tableData" v-loading="loading" style="width: 100%">
        <el-table-column prop="id" label="ID" width="80" />
        <el-table-column prop="student" label="学生信息" min-width="150">
          <template #default="scope">
            <div>
              <div>{{ scope.row.student?.real_name || '-' }}</div>
              <div class="text-sm text-gray">{{ scope.row.student?.username || '-' }}</div>
            </div>
          </template>
        </el-table-column>
        <el-table-column prop="retake_course" label="课程名称" min-width="180">
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
        <el-table-column prop="status" label="状态" width="100">
          <template #default="scope">
            <el-tag :type="getStatusType(scope.row.status)">
              {{ getStatusText(scope.row.status) }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="created_at" label="申请时间" width="180">
          <template #default="scope">
            {{ formatDate(scope.row.created_at) }}
          </template>
        </el-table-column>
        <el-table-column label="操作" width="200" fixed="right">
          <template #default="scope">
            <template v-if="scope.row.status === 'pending'">
              <el-button type="success" size="small" link @click="handleReview(scope.row, 'approved')">
                通过
              </el-button>
              <el-button type="danger" size="small" link @click="handleReview(scope.row, 'rejected')">
                拒绝
              </el-button>
            </template>
            <el-button type="primary" size="small" link @click="handleView(scope.row)">
              详情
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

    <el-dialog v-model="reviewDialogVisible" :title="isView ? '查看详情' : '审核'" width="500px">
      <el-descriptions :column="1" border>
        <el-descriptions-item label="学生">{{ currentApp.student?.real_name || '-' }}</el-descriptions-item>
        <el-descriptions-item label="课程">{{ currentApp.retake_course?.course?.name || '-' }}</el-descriptions-item>
        <el-descriptions-item label="学期">{{ currentApp.retake_course?.semester || '-' }}</el-descriptions-item>
        <el-descriptions-item label="申请理由">{{ currentApp.reason || '-' }}</el-descriptions-item>
        <el-descriptions-item v-if="currentApp.review_comment" label="审核意见">
          {{ currentApp.review_comment }}
        </el-descriptions-item>
      </el-descriptions>
      <el-form v-if="!isView" :model="reviewForm" label-width="80px" style="margin-top: 20px;">
        <el-form-item label="审核意见">
          <el-input v-model="reviewForm.review_comment" type="textarea" :rows="3" placeholder="请输入审核意见（可选）" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="reviewDialogVisible = false">取消</el-button>
        <template v-if="!isView">
          <el-button type="success" @click="submitReview('approved')" :loading="submitting">
            通过
          </el-button>
          <el-button type="danger" @click="submitReview('rejected')" :loading="submitting">
            拒绝
          </el-button>
        </template>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { getTeacherApplications, reviewApplication } from '@/api'
import dayjs from 'dayjs'

const loading = ref(false)
const tableData = ref([])
const filterStatus = ref('')

const pagination = reactive({
  page: 1,
  pageSize: 10,
  total: 0
})

const reviewDialogVisible = ref(false)
const isView = ref(false)
const currentApp = ref({})
const submitting = ref(false)
const reviewForm = reactive({
  review_comment: ''
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
    const res = await getTeacherApplications(params)
    tableData.value = res.data || []
    pagination.total = res.total || 0
  } catch (error) {
    console.error(error)
  } finally {
    loading.value = false
  }
}

const handleView = (row) => {
  isView.value = true
  currentApp.value = row
  reviewDialogVisible.value = true
}

const handleReview = (row, status) => {
  isView.value = false
  currentApp.value = row
  reviewForm.review_comment = ''
  reviewDialogVisible.value = true
}

const submitReview = async (status) => {
  submitting.value = true
  try {
    await reviewApplication(currentApp.value.id, {
      status: status,
      review_comment: reviewForm.review_comment
    })
    ElMessage.success(status === 'approved' ? '已通过' : '已拒绝')
    reviewDialogVisible.value = false
    fetchData()
  } catch (error) {
    console.error(error)
  } finally {
    submitting.value = false
  }
}

onMounted(() => {
  fetchData()
})
</script>

<style scoped>
.teacher-applications {
  padding: 0;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.text-sm {
  font-size: 12px;
}

.text-gray {
  color: #909399;
}
</style>
