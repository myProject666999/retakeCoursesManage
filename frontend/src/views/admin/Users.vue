<template>
  <div class="admin-users">
    <el-card shadow="never">
      <template #header>
        <div class="card-header">
          <span>用户管理</span>
          <div class="action-box">
            <el-upload
              ref="uploadRef"
              :auto-upload="false"
              :on-change="handleImportFileChange"
              :limit="1"
              accept=".csv"
              :show-file-list="false"
              style="display: inline-block;"
            >
              <el-button type="success" @click="handleImport">
                <el-icon><Upload /></el-icon>
                导入用户
              </el-button>
            </el-upload>
            <el-button type="info" @click="handleExport" style="margin-left: 10px;">
              <el-icon><Download /></el-icon>
              导出用户
            </el-button>
            <el-button type="primary" @click="handleCreate" style="margin-left: 10px;">
              <el-icon><Plus /></el-icon>
              新增用户
            </el-button>
          </div>
        </div>
      </template>
      <el-form :inline="true" style="margin-bottom: 20px;">
        <el-form-item label="角色">
          <el-select v-model="filterRole" placeholder="全部角色" clearable @change="fetchData" style="width: 120px;">
            <el-option label="学生" value="student" />
            <el-option label="教师" value="teacher" />
            <el-option label="管理员" value="admin" />
          </el-select>
        </el-form-item>
        <el-form-item label="搜索">
          <el-input v-model="filterKeyword" placeholder="用户名/姓名" clearable @keyup.enter="fetchData" style="width: 200px;">
            <template #prefix>
              <el-icon><Search /></el-icon>
            </template>
          </el-input>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="fetchData">搜索</el-button>
        </el-form-item>
      </el-form>
      <el-table :data="tableData" v-loading="loading" style="width: 100%">
        <el-table-column prop="id" label="ID" width="80" />
        <el-table-column prop="username" label="用户名" width="120" />
        <el-table-column prop="real_name" label="姓名" width="100" />
        <el-table-column prop="role" label="角色" width="100">
          <template #default="scope">
            <el-tag :type="getRoleType(scope.row.role)">
              {{ getRoleText(scope.row.role) }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="class" label="班级" width="120">
          <template #default="scope">
            {{ scope.row.class?.name || '-' }}
          </template>
        </el-table-column>
        <el-table-column prop="major" label="专业" width="150">
          <template #default="scope">
            {{ scope.row.major?.name || '-' }}
          </template>
        </el-table-column>
        <el-table-column prop="email" label="邮箱" width="180" show-overflow-tooltip />
        <el-table-column prop="phone" label="电话" width="130" />
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

    <el-dialog v-model="dialogVisible" :title="isEdit ? '编辑用户' : '新增用户'" width="600px">
      <el-form :model="form" :rules="rules" ref="formRef" label-width="100px">
        <el-form-item label="用户名" prop="username">
          <el-input v-model="form.username" placeholder="请输入用户名" :disabled="isEdit" />
        </el-form-item>
        <el-form-item v-if="!isEdit" label="密码" prop="password">
          <el-input v-model="form.password" type="password" placeholder="请输入密码" show-password />
        </el-form-item>
        <el-form-item label="姓名" prop="real_name">
          <el-input v-model="form.real_name" placeholder="请输入姓名" />
        </el-form-item>
        <el-form-item label="角色" prop="role">
          <el-select v-model="form.role" placeholder="请选择角色" style="width: 100%">
            <el-option label="学生" value="student" />
            <el-option label="教师" value="teacher" />
            <el-option label="管理员" value="admin" />
          </el-select>
        </el-form-item>
        <el-form-item label="邮箱" prop="email">
          <el-input v-model="form.email" placeholder="请输入邮箱" />
        </el-form-item>
        <el-form-item label="电话" prop="phone">
          <el-input v-model="form.phone" placeholder="请输入电话" />
        </el-form-item>
        <el-form-item v-if="form.role === 'student'" label="班级">
          <el-select v-model="form.class_id" placeholder="请选择班级" clearable style="width: 100%">
            <el-option
              v-for="cls in classes"
              :key="cls.id"
              :label="cls.name"
              :value="cls.id"
            />
          </el-select>
        </el-form-item>
        <el-form-item v-if="form.role === 'student'" label="专业">
          <el-select v-model="form.major_id" placeholder="请选择专业" clearable style="width: 100%">
            <el-option
              v-for="major in majors"
              :key="major.id"
              :label="major.name"
              :value="major.id"
            />
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

    <el-dialog v-model="importDialogVisible" title="导入用户" width="500px">
      <el-upload
        ref="importUploadRef"
        :auto-upload="false"
        :on-change="handleImportFileChange"
        :limit="1"
        accept=".csv"
        drag
      >
        <el-icon class="el-icon--upload"><UploadFilled /></el-icon>
        <div class="el-upload__text">
          将文件拖到此处，或<em>点击上传</em>
        </div>
        <template #tip>
          <div class="el-upload__tip">
            只能上传CSV文件，格式：用户名,姓名,角色,邮箱,电话
          </div>
        </template>
      </el-upload>
      <div v-if="importFile" style="margin-top: 20px;">
        <el-tag type="info">{{ importFile.name }}</el-tag>
      </div>
      <template #footer>
        <el-button @click="importDialogVisible = false">取消</el-button>
        <el-button type="primary" @click="submitImport" :loading="importing">
          确认导入
        </el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted, watch } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import {
  getAllUsers, createUser, updateUser, deleteUser, exportUsers, importUsers,
  getMajors, getClasses
} from '@/api'
import dayjs from 'dayjs'

const loading = ref(false)
const tableData = ref([])
const majors = ref([])
const classes = ref([])
const filterRole = ref('')
const filterKeyword = ref('')
const dialogVisible = ref(false)
const importDialogVisible = ref(false)
const isEdit = ref(false)
const submitting = ref(false)
const importing = ref(false)
const formRef = ref(null)
const uploadRef = ref(null)
const importUploadRef = ref(null)
const importFile = ref(null)

const pagination = reactive({
  page: 1,
  pageSize: 10,
  total: 0
})

const form = reactive({
  id: null,
  username: '',
  password: '',
  real_name: '',
  role: 'student',
  email: '',
  phone: '',
  class_id: null,
  major_id: null
})

const rules = {
  username: [{ required: true, message: '请输入用户名', trigger: 'blur' }],
  password: [{ required: true, message: '请输入密码', trigger: 'blur' }],
  real_name: [{ required: true, message: '请输入姓名', trigger: 'blur' }],
  role: [{ required: true, message: '请选择角色', trigger: 'change' }]
}

const getRoleType = (role) => {
  const types = {
    student: 'warning',
    teacher: 'primary',
    admin: 'danger'
  }
  return types[role] || 'info'
}

const getRoleText = (role) => {
  const texts = {
    student: '学生',
    teacher: '教师',
    admin: '管理员'
  }
  return texts[role] || role
}

watch(() => form.role, (newRole) => {
  if (newRole !== 'student') {
    form.class_id = null
    form.major_id = null
  }
})

const fetchMajors = async () => {
  try {
    const res = await getMajors({ page_size: 1000 })
    majors.value = res.data || []
  } catch (error) {
    console.error(error)
  }
}

const fetchClasses = async () => {
  try {
    const res = await getClasses({ page_size: 1000 })
    classes.value = res.data || []
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
    if (filterRole.value) {
      params.role = filterRole.value
    }
    if (filterKeyword.value) {
      params.keyword = filterKeyword.value
    }
    const res = await getAllUsers(params)
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
  form.username = ''
  form.password = ''
  form.real_name = ''
  form.role = 'student'
  form.email = ''
  form.phone = ''
  form.class_id = null
  form.major_id = null
  dialogVisible.value = true
}

const handleEdit = (row) => {
  isEdit.value = true
  form.id = row.id
  form.username = row.username
  form.password = ''
  form.real_name = row.real_name
  form.role = row.role
  form.email = row.email
  form.phone = row.phone
  form.class_id = row.class_id
  form.major_id = row.major_id
  dialogVisible.value = true
}

const handleDelete = async (row) => {
  try {
    await ElMessageBox.confirm('确定要删除这个用户吗？', '提示', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    })
    await deleteUser(row.id)
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
        real_name: form.real_name,
        role: form.role,
        email: form.email,
        phone: form.phone
      }
      if (form.class_id) data.class_id = form.class_id
      if (form.major_id) data.major_id = form.major_id
      await updateUser(form.id, data)
      ElMessage.success('更新成功')
    } else {
      const data = {
        username: form.username,
        password: form.password,
        real_name: form.real_name,
        role: form.role,
        email: form.email,
        phone: form.phone
      }
      if (form.class_id) data.class_id = form.class_id
      if (form.major_id) data.major_id = form.major_id
      await createUser(data)
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

const handleImport = () => {
  importFile.value = null
  importDialogVisible.value = true
}

const handleImportFileChange = (file) => {
  importFile.value = file.raw
}

const submitImport = async () => {
  if (!importFile.value) {
    ElMessage.warning('请选择要导入的文件')
    return
  }
  importing.value = true
  try {
    const formData = new FormData()
    formData.append('file', importFile.value)
    const res = await importUsers(formData)
    ElMessage.success(res.message || '导入成功')
    importDialogVisible.value = false
    fetchData()
  } catch (error) {
    console.error(error)
  } finally {
    importing.value = false
  }
}

const handleExport = async () => {
  try {
    const res = await exportUsers()
    const blob = new Blob([res], { type: 'text/csv;charset=utf-8;' })
    const link = document.createElement('a')
    const url = URL.createObjectURL(blob)
    link.setAttribute('href', url)
    link.setAttribute('download', `users_${dayjs().format('YYYYMMDDHHmmss')}.csv`)
    link.style.visibility = 'hidden'
    document.body.appendChild(link)
    link.click()
    document.body.removeChild(link)
    ElMessage.success('导出成功')
  } catch (error) {
    console.error(error)
  }
}

onMounted(() => {
  fetchMajors()
  fetchClasses()
  fetchData()
})
</script>

<style scoped>
.admin-users {
  padding: 0;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.action-box {
  display: flex;
  align-items: center;
}
</style>
