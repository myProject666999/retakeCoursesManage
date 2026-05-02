<template>
  <div class="teacher-grades">
    <el-card shadow="never">
      <template #header>
        <div class="card-header">
          <span>成绩管理</span>
          <div class="search-box">
            <el-select v-model="filterCourse" placeholder="选择课程" clearable @change="fetchData" style="width: 200px;">
              <el-option
                v-for="course in myCourses"
                :key="course.id"
                :label="course.course?.name"
                :value="course.id"
              />
            </el-select>
          </div>
        </div>
      </template>
      <el-table :data="tableData" v-loading="loading" style="width: 100%">
        <el-table-column prop="id" label="ID" width="80" />
        <el-table-column prop="student" label="学生信息" min-width="180">
          <template #default="scope">
            <div>
              <div>{{ scope.row.student?.real_name || '-' }}</div>
              <div class="text-sm text-gray">{{ scope.row.student?.username || '-' }}</div>
            </div>
          </template>
        </el-table-column>
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
        <el-table-column prop="score" label="成绩" width="120">
          <template #default="scope">
            <span v-if="scope.row.score !== null && scope.row.score !== undefined">
              <el-tag :type="scope.row.score >= 60 ? 'success' : 'danger'">
                {{ scope.row.score }}
              </el-tag>
            </span>
            <span v-else class="text-gray">未录入</span>
          </template>
        </el-table-column>
        <el-table-column prop="status" label="状态" width="100">
          <template #default="scope">
            <el-tag :type="getStatusType(scope.row.status)">
              {{ getStatusText(scope.row.status) }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column label="操作" width="120" fixed="right">
          <template #default="scope">
            <el-button type="primary" size="small" link @click="handleEditGrade(scope.row)">
              录入/修改
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

    <el-dialog v-model="gradeDialogVisible" title="录入成绩" width="400px">
      <el-form :model="gradeForm" label-width="80px">
        <el-form-item label="学生">
          <el-input :value="currentGrade?.student?.real_name" disabled />
        </el-form-item>
        <el-form-item label="课程">
          <el-input :value="currentGrade?.retake_course?.course?.name" disabled />
        </el-form-item>
        <el-form-item label="成绩">
          <el-input-number v-model="gradeForm.score" :min="0" :max="100" :precision="1" style="width: 100%" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="gradeDialogVisible = false">取消</el-button>
        <el-button type="primary" @click="submitGrade" :loading="submitting">
          确认
        </el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { ElMessage } from 'element-plus'
import { getTeacherGrades, getMyTeachingCourses, enterGrade } from '@/api'

const loading = ref(false)
const tableData = ref([])
const myCourses = ref([])
const filterCourse = ref(null)

const pagination = reactive({
  page: 1,
  pageSize: 10,
  total: 0
})

const gradeDialogVisible = ref(false)
const currentGrade = ref(null)
const submitting = ref(false)
const gradeForm = reactive({
  score: 0
})

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

const fetchCourses = async () => {
  try {
    const res = await getMyTeachingCourses({})
    myCourses.value = res || []
  } catch (error) {
    console.error(error)
  }
}

const fetchData = async () => {
  loading.value = true
  try {
    const params = {
      page: pagination.page,
      page_size: pagination.pageSize
    }
    if (filterCourse.value) {
      params.course_id = filterCourse.value
    }
    const res = await getTeacherGrades(params)
    tableData.value = res.data || []
    pagination.total = res.total || 0
  } catch (error) {
    console.error(error)
  } finally {
    loading.value = false
  }
}

const handleEditGrade = (row) => {
  currentGrade.value = row
  gradeForm.score = row.score || 0
  gradeDialogVisible.value = true
}

const submitGrade = async () => {
  submitting.value = true
  try {
    await enterGrade(currentGrade.value.id, { score: gradeForm.score })
    ElMessage.success('成绩录入成功')
    gradeDialogVisible.value = false
    fetchData()
  } catch (error) {
    console.error(error)
  } finally {
    submitting.value = false
  }
}

onMounted(() => {
  fetchCourses()
  fetchData()
})
</script>

<style scoped>
.teacher-grades {
  padding: 0;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.search-box {
  display: flex;
  align-items: center;
}

.text-sm {
  font-size: 12px;
}

.text-gray {
  color: #909399;
}
</style>
