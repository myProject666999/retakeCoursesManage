<template>
  <div class="retake-courses">
    <el-card shadow="never">
      <template #header>
        <div class="card-header">
          <span>重修课程列表</span>
          <el-select v-model="filterStatus" placeholder="状态筛选" clearable @change="fetchData" style="width: 150px; margin-left: 10px;">
            <el-option label="全部" value="" />
            <el-option label="进行中" value="active" />
            <el-option label="已结束" value="finished" />
          </el-select>
        </div>
      </template>
      <el-table :data="tableData" v-loading="loading" style="width: 100%">
        <el-table-column prop="id" label="ID" width="80" />
        <el-table-column prop="course" label="课程名称" min-width="200">
          <template #default="scope">
            {{ scope.row.course?.name || '-' }}
          </template>
        </el-table-column>
        <el-table-column prop="teacher" label="授课教师" width="120">
          <template #default="scope">
            {{ scope.row.teacher?.real_name || '-' }}
          </template>
        </el-table-column>
        <el-table-column prop="semester" label="学期" width="120" />
        <el-table-column prop="classroom" label="教室" width="100" />
        <el-table-column prop="max_students" label="最大人数" width="100" />
        <el-table-column prop="status" label="状态" width="100">
          <template #default="scope">
            <el-tag :type="scope.row.status === 'active' ? 'success' : 'info'">
              {{ scope.row.status === 'active' ? '进行中' : scope.row.status }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column label="操作" width="150" fixed="right">
          <template #default="scope">
            <el-button type="primary" size="small" @click="handleApply(scope.row)" :disabled="scope.row.status !== 'active'">
              申请重修
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

    <el-dialog v-model="applyDialogVisible" title="申请重修" width="500px">
      <el-form :model="applyForm" label-width="100px">
        <el-form-item label="课程名称">
          <el-input :value="currentCourse?.course?.name" disabled />
        </el-form-item>
        <el-form-item label="授课教师">
          <el-input :value="currentCourse?.teacher?.real_name" disabled />
        </el-form-item>
        <el-form-item label="申请理由">
          <el-input v-model="applyForm.reason" type="textarea" :rows="4" placeholder="请输入申请理由" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="applyDialogVisible = false">取消</el-button>
        <el-button type="primary" @click="submitApplication" :loading="submitting">
          提交申请
        </el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { getRetakeCourses, createApplication } from '@/api'

const loading = ref(false)
const tableData = ref([])
const filterStatus = ref('')

const pagination = reactive({
  page: 1,
  pageSize: 10,
  total: 0
})

const applyDialogVisible = ref(false)
const currentCourse = ref(null)
const submitting = ref(false)

const applyForm = reactive({
  reason: ''
})

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
    const res = await getRetakeCourses(params)
    tableData.value = res.data || []
    pagination.total = res.total || 0
  } catch (error) {
    console.error(error)
  } finally {
    loading.value = false
  }
}

const handleApply = (row) => {
  currentCourse.value = row
  applyForm.reason = ''
  applyDialogVisible.value = true
}

const submitApplication = async () => {
  submitting.value = true
  try {
    await createApplication({
      retake_course_id: currentCourse.value.id,
      reason: applyForm.reason
    })
    ElMessage.success('申请提交成功')
    applyDialogVisible.value = false
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
.retake-courses {
  padding: 0;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}
</style>
