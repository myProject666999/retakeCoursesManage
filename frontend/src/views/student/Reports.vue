<template>
  <div class="reports">
    <el-card shadow="never">
      <template #header>
        <div class="card-header">
          <span>我的建议/报告</span>
          <el-button type="primary" @click="handleCreate">
            <el-icon><Plus /></el-icon>
            新建
          </el-button>
        </div>
      </template>
      <el-table :data="tableData" v-loading="loading" style="width: 100%">
        <el-table-column prop="id" label="ID" width="80" />
        <el-table-column prop="title" label="标题" min-width="200" />
        <el-table-column prop="type" label="类型" width="100">
          <template #default="scope">
            <el-tag :type="scope.row.type === 'suggestion' ? 'warning' : 'primary'">
              {{ scope.row.type === 'suggestion' ? '建议' : '报告' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="status" label="状态" width="100">
          <template #default="scope">
            <el-tag :type="scope.row.status === 'replied' ? 'success' : 'warning'">
              {{ scope.row.status === 'replied' ? '已回复' : '待回复' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="created_at" label="创建时间" width="180">
          <template #default="scope">
            {{ formatDate(scope.row.created_at) }}
          </template>
        </el-table-column>
        <el-table-column label="操作" width="150" fixed="right">
          <template #default="scope">
            <el-button type="primary" size="small" link @click="handleView(scope.row)">
              查看详情
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

    <el-dialog v-model="dialogVisible" :title="isView ? '查看详情' : '新建'" width="600px">
      <el-form :model="form" :rules="rules" ref="formRef" label-width="80px">
        <el-form-item label="标题" prop="title">
          <el-input v-model="form.title" placeholder="请输入标题" :disabled="isView" />
        </el-form-item>
        <el-form-item label="类型" prop="type">
          <el-select v-model="form.type" placeholder="请选择类型" style="width: 100%" :disabled="isView">
            <el-option label="建议" value="suggestion" />
            <el-option label="报告" value="report" />
          </el-select>
        </el-form-item>
        <el-form-item label="内容" prop="content">
          <el-input v-model="form.content" type="textarea" :rows="6" placeholder="请输入内容" :disabled="isView" />
        </el-form-item>
        <el-form-item label="回复" v-if="isView && currentReport.reply">
          <el-input v-model="currentReport.reply" type="textarea" :rows="4" disabled />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="dialogVisible = false">取消</el-button>
        <el-button type="primary" @click="handleSubmit" :loading="submitting" v-if="!isView">
          提交
        </el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { ElMessage } from 'element-plus'
import { getMyReports, createReport } from '@/api'
import dayjs from 'dayjs'

const loading = ref(false)
const tableData = ref([])
const dialogVisible = ref(false)
const isView = ref(false)
const submitting = ref(false)
const currentReport = ref(null)
const formRef = ref(null)

const pagination = reactive({
  page: 1,
  pageSize: 10,
  total: 0
})

const form = reactive({
  title: '',
  type: 'suggestion',
  content: ''
})

const rules = {
  title: [{ required: true, message: '请输入标题', trigger: 'blur' }],
  type: [{ required: true, message: '请选择类型', trigger: 'change' }],
  content: [{ required: true, message: '请输入内容', trigger: 'blur' }]
}

const formatDate = (date) => {
  return dayjs(date).format('YYYY-MM-DD HH:mm:ss')
}

const fetchData = async () => {
  loading.value = true
  try {
    const res = await getMyReports({
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
  isView.value = false
  form.title = ''
  form.type = 'suggestion'
  form.content = ''
  dialogVisible.value = true
}

const handleView = (row) => {
  isView.value = true
  currentReport.value = row
  form.title = row.title
  form.type = row.type
  form.content = row.content
  dialogVisible.value = true
}

const handleSubmit = async () => {
  const valid = await formRef.value.validate().catch(() => false)
  if (!valid) return

  submitting.value = true
  try {
    await createReport({
      title: form.title,
      type: form.type,
      content: form.content
    })
    ElMessage.success('提交成功')
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
.reports {
  padding: 0;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}
</style>
