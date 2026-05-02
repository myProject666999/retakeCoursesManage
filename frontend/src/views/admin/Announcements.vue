<template>
  <div class="admin-announcements">
    <el-card shadow="never">
      <template #header>
        <div class="card-header">
          <span>公告管理</span>
          <el-button type="primary" @click="handleCreate">
            <el-icon><Plus /></el-icon>
            新增公告
          </el-button>
        </div>
      </template>
      <el-table :data="tableData" v-loading="loading" style="width: 100%">
        <el-table-column prop="id" label="ID" width="80" />
        <el-table-column prop="title" label="标题" min-width="250" />
        <el-table-column prop="author" label="发布人" width="120">
          <template #default="scope">
            {{ scope.row.author?.real_name || '-' }}
          </template>
        </el-table-column>
        <el-table-column prop="is_top" label="是否置顶" width="100">
          <template #default="scope">
            <el-tag v-if="scope.row.is_top" type="danger">置顶</el-tag>
            <span v-else>否</span>
          </template>
        </el-table-column>
        <el-table-column prop="created_at" label="创建时间" width="180">
          <template #default="scope">
            {{ formatDate(scope.row.created_at) }}
          </template>
        </el-table-column>
        <el-table-column label="操作" width="200" fixed="right">
          <template #default="scope">
            <el-button type="primary" size="small" link @click="handleView(scope.row)">
              查看
            </el-button>
            <el-button type="success" size="small" link @click="handleEdit(scope.row)">
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

    <el-dialog v-model="dialogVisible" :title="isView ? '查看公告' : (isEdit ? '编辑公告' : '新增公告')" width="600px">
      <el-form :model="form" :rules="rules" ref="formRef" label-width="80px">
        <el-form-item label="标题" prop="title">
          <el-input v-model="form.title" placeholder="请输入标题" :disabled="isView" />
        </el-form-item>
        <el-form-item label="是否置顶">
          <el-switch v-model="form.is_top" :disabled="isView" />
        </el-form-item>
        <el-form-item label="内容" prop="content">
          <el-input v-model="form.content" type="textarea" :rows="8" placeholder="请输入内容" :disabled="isView" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="dialogVisible = false">关闭</el-button>
        <el-button v-if="!isView" type="primary" @click="handleSubmit" :loading="submitting">
          确定
        </el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { getAllAnnouncements, createAnnouncement, updateAnnouncement, deleteAnnouncement } from '@/api'
import dayjs from 'dayjs'

const loading = ref(false)
const tableData = ref([])
const dialogVisible = ref(false)
const isEdit = ref(false)
const isView = ref(false)
const submitting = ref(false)
const formRef = ref(null)

const pagination = reactive({
  page: 1,
  pageSize: 10,
  total: 0
})

const form = reactive({
  id: null,
  title: '',
  content: '',
  is_top: false
})

const rules = {
  title: [{ required: true, message: '请输入标题', trigger: 'blur' }],
  content: [{ required: true, message: '请输入内容', trigger: 'blur' }]
}

const formatDate = (date) => {
  return dayjs(date).format('YYYY-MM-DD HH:mm:ss')
}

const fetchData = async () => {
  loading.value = true
  try {
    const res = await getAllAnnouncements({
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

const handleCreate = () => {
  isEdit.value = false
  isView.value = false
  form.id = null
  form.title = ''
  form.content = ''
  form.is_top = false
  dialogVisible.value = true
}

const handleView = (row) => {
  isView.value = true
  isEdit.value = false
  form.id = row.id
  form.title = row.title
  form.content = row.content
  form.is_top = row.is_top
  dialogVisible.value = true
}

const handleEdit = (row) => {
  isView.value = false
  isEdit.value = true
  form.id = row.id
  form.title = row.title
  form.content = row.content
  form.is_top = row.is_top
  dialogVisible.value = true
}

const handleDelete = async (row) => {
  try {
    await ElMessageBox.confirm('确定要删除这个公告吗？', '提示', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    })
    await deleteAnnouncement(row.id)
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
      const data = { title: form.title, content: form.content }
      if (form.is_top !== undefined) data.is_top = form.is_top
      await updateAnnouncement(form.id, data)
      ElMessage.success('更新成功')
    } else {
      await createAnnouncement({ title: form.title, content: form.content, is_top: form.is_top })
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
  fetchData()
})
</script>

<style scoped>
.admin-announcements {
  padding: 0;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}
</style>
