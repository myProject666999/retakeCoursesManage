<template>
  <div class="teaching-courses">
    <el-card shadow="never">
      <template #header>
        <div class="card-header">
          <span>我的授课课程</span>
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
            <el-button type="primary" size="small" link @click="handleEdit(scope.row)">
              编辑
            </el-button>
            <el-button type="info" size="small" link @click="handleViewStudents(scope.row)">
              查看学生
            </el-button>
          </template>
        </el-table-column>
      </el-table>
    </el-card>

    <el-dialog v-model="editDialogVisible" title="编辑课程" width="500px">
      <el-form :model="editForm" label-width="100px">
        <el-form-item label="课程名称">
          <el-input :value="currentCourse?.course?.name" disabled />
        </el-form-item>
        <el-form-item label="学期">
          <el-input :value="currentCourse?.semester" disabled />
        </el-form-item>
        <el-form-item label="教室">
          <el-input v-model="editForm.classroom" placeholder="请输入教室" />
        </el-form-item>
        <el-form-item label="最大人数">
          <el-input-number v-model="editForm.max_students" :min="1" :max="200" />
        </el-form-item>
        <el-form-item label="状态">
          <el-select v-model="editForm.status" style="width: 100%">
            <el-option label="进行中" value="active" />
            <el-option label="已结束" value="finished" />
          </el-select>
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="editDialogVisible = false">取消</el-button>
        <el-button type="primary" @click="submitEdit" :loading="submitting">
          保存
        </el-button>
      </template>
    </el-dialog>

    <el-dialog v-model="studentsDialogVisible" title="选课学生列表" width="800px">
      <el-table :data="courseStudents" v-loading="studentsLoading" style="width: 100%">
        <el-table-column prop="student.username" label="学号" width="120">
          <template #default="scope">
            {{ scope.row.student?.username || '-' }}
          </template>
        </el-table-column>
        <el-table-column prop="student.real_name" label="姓名" width="100">
          <template #default="scope">
            {{ scope.row.student?.real_name || '-' }}
          </template>
        </el-table-column>
        <el-table-column prop="student.class" label="班级" width="150">
          <template #default="scope">
            {{ scope.row.student?.class?.name || '-' }}
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
        <el-table-column prop="status" label="状态" width="100">
          <template #default="scope">
            <el-tag :type="scope.row.status === 'submitted' ? 'primary' : 'warning'">
              {{ scope.row.status === 'submitted' ? '已提交' : '待录入' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column label="操作" width="120">
          <template #default="scope">
            <el-button type="primary" size="small" link @click="handleEnterGrade(scope.row)">
              录入成绩
            </el-button>
          </template>
        </el-table-column>
      </el-table>
    </el-dialog>

    <el-dialog v-model="gradeDialogVisible" title="录入成绩" width="400px">
      <el-form :model="gradeForm" label-width="80px">
        <el-form-item label="学生">
          <el-input :value="currentGrade?.student?.real_name" disabled />
        </el-form-item>
        <el-form-item label="成绩">
          <el-input-number v-model="gradeForm.score" :min="0" :max="100" :precision="1" style="width: 100%" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="gradeDialogVisible = false">取消</el-button>
        <el-button type="primary" @click="submitGrade" :loading="gradeSubmitting">
          确认
        </el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { ElMessage } from 'element-plus'
import { getMyTeachingCourses, getCourseStudents, updateTeacherRetakeCourse, enterGrade } from '@/api'

const loading = ref(false)
const tableData = ref([])
const filterStatus = ref('')

const editDialogVisible = ref(false)
const currentCourse = ref(null)
const submitting = ref(false)
const editForm = reactive({
  classroom: '',
  max_students: 0,
  status: ''
})

const studentsDialogVisible = ref(false)
const studentsLoading = ref(false)
const courseStudents = ref([])

const gradeDialogVisible = ref(false)
const currentGrade = ref(null)
const gradeSubmitting = ref(false)
const gradeForm = reactive({
  score: 0
})

const fetchData = async () => {
  loading.value = true
  try {
    const params = {}
    if (filterStatus.value) {
      params.status = filterStatus.value
    }
    const res = await getMyTeachingCourses(params)
    tableData.value = res || []
  } catch (error) {
    console.error(error)
  } finally {
    loading.value = false
  }
}

const handleEdit = (row) => {
  currentCourse.value = row
  editForm.classroom = row.classroom || ''
  editForm.max_students = row.max_students || 0
  editForm.status = row.status || ''
  editDialogVisible.value = true
}

const submitEdit = async () => {
  submitting.value = true
  try {
    await updateTeacherRetakeCourse(currentCourse.value.id, editForm)
    ElMessage.success('更新成功')
    editDialogVisible.value = false
    fetchData()
  } catch (error) {
    console.error(error)
  } finally {
    submitting.value = false
  }
}

const handleViewStudents = async (row) => {
  currentCourse.value = row
  studentsLoading.value = true
  try {
    const res = await getCourseStudents(row.id)
    courseStudents.value = res || []
    studentsDialogVisible.value = true
  } catch (error) {
    console.error(error)
  } finally {
    studentsLoading.value = false
  }
}

const handleEnterGrade = (row) => {
  currentGrade.value = row
  gradeForm.score = row.score || 0
  gradeDialogVisible.value = true
}

const submitGrade = async () => {
  gradeSubmitting.value = true
  try {
    await enterGrade(currentGrade.value.id, { score: gradeForm.score })
    ElMessage.success('成绩录入成功')
    gradeDialogVisible.value = false
    handleViewStudents(currentCourse.value)
  } catch (error) {
    console.error(error)
  } finally {
    gradeSubmitting.value = false
  }
}

onMounted(() => {
  fetchData()
})
</script>

<style scoped>
.teaching-courses {
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
