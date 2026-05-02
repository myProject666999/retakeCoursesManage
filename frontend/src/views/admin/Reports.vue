<template>
  <div class="admin-reports">
    <el-card shadow="never">
      <template #header>
        <div class="card-header">
          <span>建议报告管理</span>
          <div class="search-box">
            <el-select v-model="filterType" placeholder="类型筛选" clearable @change="fetchData" style="width: 120px;">
              <el-option label="建议" value="suggestion" />
              <el-option label="报告" value="report" />
            </el-select>
            <el-select v-model="filterStatus" placeholder="状态筛选" clearable @change="fetchData" style="width: 120px; margin-left: 10px;">
              <el-option label="待回复" value="pending" />
              <el-option label="已回复" value="replied" />
            </el-select>
          </div>
        </div>
      </template>
      <el-table :data="tableData" v-loading="loading" style="width: 100%">
        <el-table-column prop="id" label="ID" width="80" />
        <el-table-column prop="title" label="标题" min-width="200" />
        <el-table-column prop="user" label="提交人" width="120">
          <template #default="scope">
            {{ scope.row.user?.real_name || '-' }}
          </template>
        </el-table-column>
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
        <el-table-column prop="created_at" label="提交时间" width="180">
          <template #default="scope">
            {{ formatDate(scope.row.created_at) }}
          </template>
        </el-table-column>
        <el-table-column label="操作" width="200" fixed="right">
          <template #default="scope">
            <el-button type="primary" size="small" link @click="handleView(scope.row)">
              查看详情
            </el-button>
            <el-button
              v-if="scope.row.status !== 'replied'"
              type="success"
              size="small"
              link
              @click="handleReply(scope.row)"
            >
              回复
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

    <el-dialog v-model="dialogVisible" :title="isReply ? '回复' : '查看详情'" width="600px">
      <el-descriptions :column="1" border>
        <el-descriptions-item label="标题">{{ currentReport.title }}</el-descriptions-item>
        <el-descriptions-item label="提交人">{{ currentReport.user?.real_name || '-' }}</el-descriptions-item>
        <el-descriptions-item label="类型">
          <el-tag :type="currentReport.type === 'suggestion' ? 'warning' : 'primary'">
            {{ currentReport.type === 'suggestion' ? '建议' : '报告' }}
          </el-tag>
        </el-descriptions-item>
        <el-descriptions-item label="状态">
          <el-tag :type="currentReport.status === 'replied' ? 'success' : 'warning'">
            {{ currentReport.status === 'replied' ? '已回复' : '待回复' }}
          </el-tag>
        </el-descriptions-item>
        <el-descriptions-item label="内容">
          <div style="white-space: pre-wrap;">{{ currentReport.content }}</div>
        </el-descriptions-item>
        <el-descriptions-item v-if="currentReport.reply" label="回复">
          <div style="white-space: pre-wrap;">{{ currentReport.reply }}</div>
        </el-descriptions-item>
      </el-descriptions>
      <el-form v-if="isReply" :model="replyForm" label-width="80px" style="margin-top: 20px;">
        <el-form-item label="回复">
          <el-input v-model="replyForm.reply" type="textarea" :rows="4" placeholder="请输入回复内容" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="dialogVisible = false">关闭</el-button>
        <el-button v-if="isReply" type="primary" @click="submitReply" :loading="submitting">
          提交回复
        </el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { getAllReports, adminReplyReport, deleteReport } from '@/api'
import dayjs from 'dayjs'

const loading = ref(false)
const tableData = ref([])
const filterType = ref('')
const filterStatus = ref('')
const dialogVisible = ref(false)
const isReply = ref(false)
const currentReport = ref({})
const submitting = ref(false)

const pagination = reactive({
  page: 1,
  pageSize: 10,
  total: 0
})

const replyForm = reactive({
  reply: ''
})

const formatDate = (date) => {
  return dayjs(date).format('YYYY-MM-DD HH:mm:ss')
}

const fetchData = async () => {
  loading.value = true
  try {
    const params = {
      page: pagination.page,
      page_size: pagination.pageSize
    }
    if (filterType.value) {
      params.type = filterType.value
    }
    if (filterStatus.value) {
      params.status = filterStatus.value
    }
    const res = await getAllReports(params)
    tableData.value = res.data || []
    pagination.total = res.total || 0
  } catch (error) {
    console.error(error)
  } finally {
    loading.value = false
  }
}

const handleView = (row) => {
  isReply.value = false
  currentReport.value = row
  dialogVisible.value = true
}

const handleReply = (row) => {
  isReply.value = true
  currentReport.value = row
  replyForm.reply = ''
  dialogVisible.value = true
}

const handleDelete = async (row) => {
  try {
    await ElMessageBox.confirm('确定要删除这个建议/报告吗？', '提示', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    })
    await deleteReport(row.id)
    ElMessage.success('删除成功')
    fetchData()
  } catch (error) {
    if (error !== 'cancel') {
      console.error(error)
    }
  }
}

const submitReply = async () => {
  if (!replyForm.reply.trim()) {
    ElMessage.warning('请输入回复内容')
    return
  }
  submitting.value = true
  try {
    await adminReplyReport(currentReport.value.id, { reply: replyForm.reply })
    ElMessage.success('回复成功')
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
.admin-reports {
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
