<template>
  <div class="announcement-detail">
    <el-card shadow="never">
      <template #header>
        <div class="card-header">
          <el-button type="text" @click="goBack">
            <el-icon><ArrowLeft /></el-icon>
            返回列表
          </el-button>
        </div>
      </template>
      <div v-loading="loading" class="detail-content">
        <h2 class="detail-title">{{ announcement.title }}</h2>
        <div class="detail-meta">
          <span>发布人：{{ announcement.author?.real_name || '-' }}</span>
          <span>发布时间：{{ formatDate(announcement.created_at) }}</span>
          <el-tag v-if="announcement.is_top" type="danger">置顶</el-tag>
        </div>
        <el-divider />
        <div class="detail-content-text">{{ announcement.content }}</div>
      </div>
    </el-card>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { getAnnouncementDetail } from '@/api'
import dayjs from 'dayjs'

const route = useRoute()
const router = useRouter()

const loading = ref(false)
const announcement = ref({
  title: '',
  content: '',
  author: null,
  is_top: false,
  created_at: ''
})

const formatDate = (date) => {
  return dayjs(date).format('YYYY-MM-DD HH:mm:ss')
}

const goBack = () => {
  router.push('/announcements')
}

const fetchData = async () => {
  loading.value = true
  try {
    const res = await getAnnouncementDetail(route.params.id)
    announcement.value = res
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
.announcement-detail {
  padding: 0;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.detail-content {
  max-width: 800px;
  margin: 0 auto;
}

.detail-title {
  text-align: center;
  font-size: 24px;
  color: #303133;
  margin-bottom: 20px;
}

.detail-meta {
  display: flex;
  align-items: center;
  gap: 20px;
  justify-content: center;
  color: #909399;
  font-size: 14px;
  flex-wrap: wrap;
}

.detail-content-text {
  line-height: 2;
  color: #606266;
  font-size: 16px;
  white-space: pre-wrap;
}
</style>
