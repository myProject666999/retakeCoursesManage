<template>
  <div class="admin-retake-courses">
    <el-card shadow="never">
      <template #header>
        <div class="card-header">
          <span>重修课程管理</span>
          <div class="search-box">
            <el-select v-model="filterStatus" placeholder="状态筛选" clearable @change="fetchData" style="width: 120px;">
              <el-option label="进行中" value="active" />
              <el-option label="已结束" value="finished" />
            </el-select>
            <el-button type="primary" @click="handleCreate" style="margin-left: 10px;">
              <el-icon><Plus /></el-icon>
              新增课程
            </el-button>
          </div>
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
        <el-table-column prop="created_at" label="创建时间" width="180">
          <template #default="scope">
            {{ formatDate(scope.row.created_at) }}
          </template>
        </el-table-column>
        <el-table-column label="操作" width="180" fixed="right">
          <template #default="scope">
            <el-button type="primary" size="small" link @click="handleEdit(scope.row)">
              编辑
            </el-button>
            <el-button type="danger" size="small" link @click="handleDelete(scope.row)">
              删除
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

    <el-dialog v-model="dialogVisible" :title="isEdit ? '编辑重修课程' : '新增重修课程'" width="600px">
      <el-form :model="form" :rules="rules" ref="formRef" label-width="100px">
        <el-form-item label="基础课程" prop="course_id">
          <el-select v-model="form.course_id" placeholder="请选择基础课程" style="width: 100%">
            <el-option
              v-for="course in courses"
              :key="course.id"
              :label="course.name"
              :value="course.id"
            />
          </el-select>
        </el-form-item>
        <el-form-item label="授课教师" prop="teacher_id">
          <el-select v-model="form.teacher_id" placeholder="请选择授课教师" style="width: 100%">
            <el-option
              v-for="teacher in teachers"
              :key="teacher.id"
              :label="teacher.real_name"
              :value="teacher.id"
            />
          </el-select>
        </el-form-item>
        <el-form-item label="学期" prop="semester">
          <el-input v-model="form.semester" placeholder="请输入学期，如：2024-2025-1" />
        </el-form-item>
        <el-form-item label="教室">
          <el-input v-model="form.classroom" placeholder="请输入教室" />
        </el-form-item>
        <el-form-item label="最大人数">
          <el-input-number v-model="form.max_students" :min="1" :max="200" style="width: 100%" />
        </el-form-item>
        <el-form-item label="状态">
          <el-select v-model="form.status" style="width: 100%">
            <el-option label="进行中" value="active" />
            <el-option label="已结束" value="finished" />
          </el-select>
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="dialogVisible = false">取消</el-button>
        <el-button type="primary" @click="handleSubmit" :loading="submitting">
          确定
        </el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import {
  getAllRetakeCourses, createRetakeCourse, updateRetakeCourse, deleteRetakeCourse,
  getAllCourses, getAllTeachers
} from '@/api'
import dayjs from 'dayjs'

const loading = ref(false)
const tableData = ref([])
const courses = ref([])
const teachers = ref([])
const filterStatus = ref('')
const dialogVisible = ref(false)
const isEdit = ref(false)
const submitting = ref(false)
const formRef = ref(null)

const pagination = reactive({
  page: 1,
  pageSize: 10,
  total: 0
})

const form = reactive({
  id: null,
  course_id: null,
  teacher_id: null,
  semester: '',
  classroom: '',
  max_students: 50,
  status: 'active'
})

const rules = {
  course_id: [{ required: true, message: '请选择基础课程', trigger: 'change' }],
  teacher_id: [{ required: true, message: '请选择授课教师', trigger: 'change' }],
  semester: [{ required: true, message: '请输入学期', trigger: 'blur' }]
}

const formatDate = (date) => {
  return dayjs(date).format('YYYY-MM-DD HH:mm:ss')
}

const fetchCourses = async () => {
  try {
    const res = await getAllCourses({ page_size: 1000 })
    courses.value = res.data || []
  } catch (error) {
    console.error(error)
  }
}

const fetchTeachers = async () => {
  try {
    const res = await getAllTeachers()
    teachers.value = res || []
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
    if (filterStatus.value) {
      params.status = filterStatus.value
    }
    const res = await getAllRetakeCourses(params)
    tableData.value = res.data || []
    pagination.total = res.total || 0
  } catch (error) {
    console.error(error)
  } finally {
    loading.value = false
  }
}

const handleCreate = () => {
  isEdit.value = false
  form.id = null
  form.course_id = null
  form.teacher_id = null
  form.semester = ''
  form.classroom = ''
  form.max_students = 50
  form.status = 'active'
  dialogVisible.value = true
}

const handleEdit = (row) => {
  isEdit.value = true
  form.id = row.id
  form.course_id = row.course_id
  form.teacher_id = row.teacher_id
  form.semester = row.semester
  form.classroom = row.classroom || ''
  form.max_students = row.max_students || 50
  form.status = row.status || 'active'
  dialogVisible.value = true
}

const handleDelete = async (row) => {
  try {
    await ElMessageBox.confirm('确定要删除这个重修课程吗？', '提示', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    })
    await deleteRetakeCourse(row.id)
    ElMessage.success('删除成功')
    fetchData()
  } catch (error) {
    if (error !== 'cancel') {
      console.error(error)
    }
  }
}

const handleSubmit = async () => {
  const valid = await formRef.value.validate().catch(() => false)
  if (!valid) return

  submitting.value = true
  try {
    if (isEdit.value) {
      const data = {
        course_id: form.course_id,
        teacher_id: form.teacher_id,
        semester: form.semester,
        classroom: form.classroom,
        max_students: form.max_students,
        status: form.status
      }
      await updateRetakeCourse(form.id, data)
      ElMessage.success('更新成功')
    } else {
      await createRetakeCourse({
        course_id: form.course_id,
        teacher_id: form.teacher_id,
        semester: form.semester,
        classroom: form.classroom,
        max_students: form.max_students,
        status: form.status
      })
      ElMessage.success('创建成功')
    }
    dialogVisible.value = false
    fetchData()
  } catch (error) {
    console.error(error)
  } finally {
    submitting.value = false
  }
}

onMounted(() => {
  fetchCourses()
  fetchTeachers()
  fetchData()
})
</script>

<style scoped>
.admin-retake-courses {
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
</style>
