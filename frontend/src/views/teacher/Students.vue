<template>
  <div class="students">
    <el-card shadow="never">
      <template #header>
        <div class="card-header">
          <span>学生列表</span>
          <div class="search-box">
            <el-input v-model="keyword" placeholder="搜索学号/姓名" style="width: 200px" clearable @keyup.enter="fetchData">
              <template #prefix>
                <el-icon><Search /></el-icon>
              </template>
            </el-input>
            <el-button type="primary" @click="fetchData" style="margin-left: 10px;">
              <el-icon><Search /></el-icon>
              搜索
            </el-button>
          </div>
        </div>
      </template>
      <el-table :data="tableData" v-loading="loading" style="width: 100%">
        <el-table-column prop="id" label="ID" width="80" />
        <el-table-column prop="username" label="学号" width="120" />
        <el-table-column prop="real_name" label="姓名" width="100" />
        <el-table-column prop="class" label="班级" width="150">
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
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { getStudents } from '@/api'

const loading = ref(false)
const tableData = ref([])
const keyword = ref('')

const pagination = reactive({
  page: 1,
  pageSize: 10,
  total: 0
})

const fetchData = async () => {
  loading.value = true
  try {
    const params = {
      page: pagination.page,
      page_size: pagination.pageSize
    }
    if (keyword.value) {
      params.keyword = keyword.value
    }
    const res = await getStudents(params)
    tableData.value = res.data || []
    pagination.total = res.total || 0
  } catch (error) {
    console.error(error)
  } finally {
    loading.value = false
  }
}

onMounted(() => {
  fetchData()
})
</script>

<style scoped>
.students {
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
