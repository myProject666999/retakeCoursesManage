<template>
  <div class="dashboard">
    <el-row :gutter="20">
      <el-col :span="6">
        <el-card shadow="hover" class="stat-card">
          <div class="stat-content">
            <div class="stat-icon" style="background: #409eff">
              <el-icon size="30"><Bell /></el-icon>
            </div>
            <div class="stat-info">
              <div class="stat-value">{{ statData.announcements }}</div>
              <div class="stat-label">公告数量</div>
            </div>
          </div>
        </el-card>
      </el-col>
      <el-col :span="6">
        <el-card shadow="hover" class="stat-card">
          <div class="stat-content">
            <div class="stat-icon" style="background: #67c23a">
              <el-icon size="30"><Notebook /></el-icon>
            </div>
            <div class="stat-info">
              <div class="stat-value">{{ statData.courses }}</div>
              <div class="stat-label">课程数量</div>
            </div>
          </div>
        </el-card>
      </el-col>
      <el-col :span="6">
        <el-card shadow="hover" class="stat-card">
          <div class="stat-content">
            <div class="stat-icon" style="background: #e6a23c">
              <el-icon size="30"><UserFilled /></el-icon>
            </div>
            <div class="stat-info">
              <div class="stat-value">{{ statData.students }}</div>
              <div class="stat-label">学生数量</div>
            </div>
          </div>
        </el-card>
      </el-col>
      <el-col :span="6">
        <el-card shadow="hover" class="stat-card">
          <div class="stat-content">
            <div class="stat-icon" style="background: #f56c6c">
              <el-icon size="30"><EditPen /></el-icon>
            </div>
            <div class="stat-info">
              <div class="stat-value">{{ statData.applications }}</div>
              <div class="stat-label">待审核申请</div>
            </div>
          </div>
        </el-card>
      </el-col>
    </el-row>

    <el-row :gutter="20" style="margin-top: 20px">
      <el-col :span="12">
        <el-card shadow="hover">
          <template #header>
            <span class="card-header">欢迎使用</span>
          </template>
          <div class="welcome-content">
            <h2>重修课程信息管理系统</h2>
            <p>当前用户：{{ userStore.user?.real_name || userStore.user?.username }}</p>
            <p>用户角色：{{ roleName }}</p>
            <el-divider />
            <div class="quick-links">
              <template v-if="userStore.isStudent">
                <el-button type="primary" @click="goTo('/student/retake-courses')">
                  查看重修课程
                </el-button>
                <el-button type="success" @click="goTo('/student/applications')">
                  我的重修申请
                </el-button>
                <el-button type="warning" @click="goTo('/student/grades')">
                  查看成绩
                </el-button>
              </template>
              <template v-if="userStore.isTeacher">
                <el-button type="primary" @click="goTo('/teacher/teaching-courses')">
                  我的授课课程
                </el-button>
                <el-button type="success" @click="goTo('/teacher/grades')">
                  成绩录入
                </el-button>
                <el-button type="warning" @click="goTo('/teacher/applications')">
                  申请审核
                </el-button>
              </template>
              <template v-if="userStore.isAdmin">
                <el-button type="primary" @click="goTo('/admin/users')">
                  用户管理
                </el-button>
                <el-button type="success" @click="goTo('/admin/retake-courses')">
                  重修课程管理
                </el-button>
                <el-button type="warning" @click="goTo('/admin/announcements')">
                  公告管理
                </el-button>
              </template>
            </div>
          </div>
        </el-card>
      </el-col>
      <el-col :span="12">
        <el-card shadow="hover">
          <template #header>
            <span class="card-header">最新公告</span>
          </template>
          <el-table :data="latestAnnouncements" v-loading="loading" style="width: 100%">
            <el-table-column prop="title" label="标题" min-width="200">
              <template #default="scope">
                <span class="announcement-title" @click="viewAnnouncement(scope.row)">{{ scope.row.title }}</span>
              </template>
            </el-table-column>
            <el-table-column prop="created_at" label="发布时间" width="180">
              <template #default="scope">
                {{ formatDate(scope.row.created_at) }}
              </template>
            </el-table-column>
          </el-table>
        </el-card>
      </el-col>
    </el-row>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useUserStore } from '@/stores/user'
import { getAnnouncements } from '@/api'
import dayjs from 'dayjs'

const router = useRouter()
const userStore = useUserStore()

const loading = ref(false)
const latestAnnouncements = ref([])

const statData = ref({
  announcements: 0,
  courses: 0,
  students: 0,
  applications: 0
})

const roleName = computed(() => {
  const roles = {
    student: '学生',
    teacher: '教师',
    admin: '管理员'
  }
  return roles[userStore.role] || '未知'
})

const formatDate = (date) => {
  return dayjs(date).format('YYYY-MM-DD HH:mm:ss')
}

const goTo = (path) => {
  router.push(path)
}

const viewAnnouncement = (row) => {
  router.push(`/announcements/${row.id}`)
}

const fetchAnnouncements = async () => {
  loading.value = true
  try {
    const res = await getAnnouncements({ page_size: 5 })
    latestAnnouncements.value = res.data || []
    statData.value.announcements = res.total || 0
  } catch (error) {
    console.error(error)
  } finally {
    loading.value = false
  }
}

onMounted(() => {
  fetchAnnouncements()
})
</script>

<style scoped>
.dashboard {
  padding: 0;
}

.stat-card {
  margin-bottom: 20px;
}

.stat-content {
  display: flex;
  align-items: center;
}

.stat-icon {
  width: 60px;
  height: 60px;
  border-radius: 10px;
  display: flex;
  align-items: center;
  justify-content: center;
  color: #fff;
}

.stat-info {
  margin-left: 20px;
}

.stat-value {
  font-size: 28px;
  font-weight: bold;
  color: #303133;
}

.stat-label {
  font-size: 14px;
  color: #909399;
  margin-top: 5px;
}

.card-header {
  font-size: 16px;
  font-weight: bold;
  color: #303133;
}

.welcome-content {
  text-align: center;
}

.welcome-content h2 {
  color: #409eff;
  margin-bottom: 20px;
}

.welcome-content p {
  color: #606266;
  margin: 10px 0;
}

.quick-links {
  display: flex;
  gap: 10px;
  justify-content: center;
  flex-wrap: wrap;
}

.announcement-title {
  color: #409eff;
  cursor: pointer;
}

.announcement-title:hover {
  text-decoration: underline;
}
</style>
