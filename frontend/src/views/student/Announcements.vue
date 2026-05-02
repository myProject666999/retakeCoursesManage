<template>
  <div class="announcements">
    <el-card shadow="never">
      <template #header>
        <div class="card-header">
          <span>公告列表</span>
        </div>
      </template>
      <el-table :data="tableData" v-loading="loading" style="width: 100%">
        <el-table-column prop="id" label="ID" width="80" />
        <el-table-column prop="title" label="标题" min-width="200">
          <template #default="scope">
            <el-link type="primary" @click="viewDetail(scope.row)">{{ scope.row.title }}</el-link>
          </template>
        </el-table-column>
        <el-table-column prop="author" label="发布人" width="120">
          <template #default="scope">
            {{ scope.row.author?.real_name || '-' }}
          </template>
        </el-table-column>
        <el-table-column prop="is_top" label="是否置顶" width="100">
          <template #default="scope">
            <el-tag v-if="scope.row.is_top" type="danger">置顶</el-tag>
            <span v-else>-</span>
          </template>
        </el-table-column>
        <el-table-column prop="created_at" label="发布时间" width="180">
          <template #default="scope">
            {{ formatDate(scope.row.created_at) }}
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
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { getAnnouncements } from '@/api'
import dayjs from 'dayjs'

const router = useRouter()

const loading = ref(false)
const tableData = ref([])

const pagination = reactive({
  page: 1,
  pageSize: 10,
  total: 0
})

const formatDate = (date) => {
  return dayjs(date).format('YYYY-MM-DD HH:mm:ss')
}

const viewDetail = (row) => {
  router.push(`/announcements/${row.id}`)
}

const fetchData = async () => {
  loading.value = true
  try {
    const res = await getAnnouncements({
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

onMounted(() => {
  fetchData()
})
</script>

<style scoped>
.announcements {
  padding: 0;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  font-size: 16px;
  font-weight: bold;
}
</style>
