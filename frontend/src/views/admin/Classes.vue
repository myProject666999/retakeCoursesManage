<template>
  <div class="admin-classes">
    <el-card shadow="never">
      <template #header>
        <div class="card-header">
          <span>班级管理</span>
          <div class="search-box">
            <el-select v-model="filterMajor" placeholder="选择专业" clearable @change="fetchData" style="width: 180px;">
              <el-option
                v-for="major in majors"
                :key="major.id"
                :label="major.name"
                :value="major.id"
              />
            </el-select>
            <el-button type="primary" @click="handleCreate" style="margin-left: 10px;">
              <el-icon><Plus /></el-icon>
              新增班级
            </el-button>
          </div>
        </div>
      </template>
      <el-table :data="tableData" v-loading="loading" style="width: 100%">
        <el-table-column prop="id" label="ID" width="80" />
        <el-table-column prop="name" label="班级名称" min-width="200" />
        <el-table-column prop="major" label="所属专业" min-width="180">
          <template #default="scope">
            {{ scope.row.major?.name || '-' }}
          </template>
        </el-table-column>
        <el-table-column prop="created_at" label="创建时间" width="180">
          <template #default="scope">
            {{ formatDate(scope.row.created_at) }}
          </template>
        </el-table-column>
        <el-table-column label="操作" width="150" fixed="right">
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

    <el-dialog v-model="dialogVisible" :title="isEdit ? '编辑班级' : '新增班级'" width="500px">
      <el-form :model="form" :rules="rules" ref="formRef" label-width="80px">
        <el-form-item label="班级名称" prop="name">
          <el-input v-model="form.name" placeholder="请输入班级名称" />
        </el-form-item>
        <el-form-item label="所属专业" prop="major_id">
          <el-select v-model="form.major_id" placeholder="请选择专业" style="width: 100%">
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
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { getClasses, createClass, updateClass, deleteClass, getMajors } from '@/api'
import dayjs from 'dayjs'

const loading = ref(false)
const tableData = ref([])
const majors = ref([])
const filterMajor = ref(null)
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
  name: '',
  major_id: null
})

const rules = {
  name: [{ required: true, message: '请输入班级名称', trigger: 'blur' }],
  major_id: [{ required: true, message: '请选择专业', trigger: 'change' }]
}

const formatDate = (date) => {
  return dayjs(date).format('YYYY-MM-DD HH:mm:ss')
}

const fetchMajors = async () => {
  try {
    const res = await getMajors({ page_size: 1000 })
    majors.value = res.data || []
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
    if (filterMajor.value) {
      params.major_id = filterMajor.value
    }
    const res = await getClasses(params)
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
  form.name = ''
  form.major_id = null
  dialogVisible.value = true
}

const handleEdit = (row) => {
  isEdit.value = true
  form.id = row.id
  form.name = row.name
  form.major_id = row.major_id
  dialogVisible.value = true
}

const handleDelete = async (row) => {
  try {
    await ElMessageBox.confirm('确定要删除这个班级吗？', '提示', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    })
    await deleteClass(row.id)
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
      await updateClass(form.id, { name: form.name, major_id: form.major_id })
      ElMessage.success('更新成功')
    } else {
      await createClass({ name: form.name, major_id: form.major_id })
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
  fetchMajors()
  fetchData()
})
</script>

<style scoped>
.admin-classes {
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
