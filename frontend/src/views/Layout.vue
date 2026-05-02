<template>
  <el-container class="layout-container">
    <el-aside :width="isCollapse ? '64px' : '220px'" class="layout-aside">
      <div class="logo">
        <el-icon v-if="isCollapse" size="30" color="#fff"><School /></el-icon>
        <span v-else class="logo-text">重修课程管理系统</span>
      </div>
      <el-menu
        :default-active="activeMenu"
        :collapse="isCollapse"
        :collapse-transition="false"
        router
        class="layout-menu"
        background-color="#304156"
        text-color="#bfcbd9"
        active-text-color="#409eff"
      >
        <el-menu-item index="/dashboard">
          <el-icon><HomeFilled /></el-icon>
          <template #title>首页</template>
        </el-menu-item>

        <el-sub-menu index="student" v-if="userStore.isStudent">
          <template #title>
            <el-icon><Reading /></el-icon>
            <span>学生功能</span>
          </template>
          <el-menu-item index="/announcements">
            <el-icon><Bell /></el-icon>
            <template #title>公告查询</template>
          </el-menu-item>
          <el-menu-item index="/student/retake-courses">
            <el-icon><Notebook /></el-icon>
            <template #title>重修课程查询</template>
          </el-menu-item>
          <el-menu-item index="/student/grades">
            <el-icon><DataAnalysis /></el-icon>
            <template #title>课程成绩查询</template>
          </el-menu-item>
          <el-menu-item index="/student/reports">
            <el-icon><Document /></el-icon>
            <template #title>建议报告</template>
          </el-menu-item>
          <el-menu-item index="/student/applications">
            <el-icon><EditPen /></el-icon>
            <template #title>重修申请管理</template>
          </el-menu-item>
        </el-sub-menu>

        <el-sub-menu index="teacher" v-if="userStore.isTeacher">
          <template #title>
            <el-icon><Notebook /></el-icon>
            <span>教师功能</span>
          </template>
          <el-menu-item index="/teacher/students">
            <el-icon><User /></el-icon>
            <template #title>学生查询</template>
          </el-menu-item>
          <el-menu-item index="/teacher/grades">
            <el-icon><Edit /></el-icon>
            <template #title>成绩管理</template>
          </el-menu-item>
          <el-menu-item index="/teacher/reports">
            <el-icon><Document /></el-icon>
            <template #title>报告管理</template>
          </el-menu-item>
          <el-menu-item index="/teacher/teaching-courses">
            <el-icon><Notebook /></el-icon>
            <template #title>重修课程管理</template>
          </el-menu-item>
          <el-menu-item index="/teacher/applications">
            <el-icon><EditPen /></el-icon>
            <template #title>重修申请管理</template>
          </el-menu-item>
        </el-sub-menu>

        <el-sub-menu index="admin" v-if="userStore.isAdmin">
          <template #title>
            <el-icon><Setting /></el-icon>
            <span>系统管理</span>
          </template>
          <el-menu-item index="/admin/majors">
            <el-icon><OfficeBuilding /></el-icon>
            <template #title>专业管理</template>
          </el-menu-item>
          <el-menu-item index="/admin/classes">
            <el-icon><Collection /></el-icon>
            <template #title>班级管理</template>
          </el-menu-item>
          <el-menu-item index="/admin/users">
            <el-icon><UserFilled /></el-icon>
            <template #title>用户管理</template>
          </el-menu-item>
          <el-menu-item index="/admin/courses">
            <el-icon><Reading /></el-icon>
            <template #title>课程管理</template>
          </el-menu-item>
          <el-menu-item index="/admin/retake-courses">
            <el-icon><Notebook /></el-icon>
            <template #title>重修课程管理</template>
          </el-menu-item>
          <el-menu-item index="/admin/grades">
            <el-icon><DataAnalysis /></el-icon>
            <template #title>成绩管理</template>
          </el-menu-item>
          <el-menu-item index="/admin/announcements">
            <el-icon><Bell /></el-icon>
            <template #title>公告管理</template>
          </el-menu-item>
          <el-menu-item index="/admin/reports">
            <el-icon><Document /></el-icon>
            <template #title>建议报告管理</template>
          </el-menu-item>
          <el-menu-item index="/admin/applications">
            <el-icon><EditPen /></el-icon>
            <template #title>重修申请管理</template>
          </el-menu-item>
        </el-sub-menu>

        <el-sub-menu index="profile">
          <template #title>
            <el-icon><User /></el-icon>
            <span>个人中心</span>
          </template>
          <el-menu-item index="/profile">
            <el-icon><User /></el-icon>
            <template #title>个人信息</template>
          </el-menu-item>
          <el-menu-item index="/password">
            <el-icon><Lock /></el-icon>
            <template #title>密码修改</template>
          </el-menu-item>
        </el-sub-menu>
      </el-menu>
    </el-aside>
    <el-container>
      <el-header class="layout-header">
        <el-icon @click="isCollapse = !isCollapse" class="collapse-btn">
          <Fold v-if="!isCollapse" />
          <Expand v-else />
        </el-icon>
        <div class="header-right">
          <el-dropdown @command="handleCommand">
            <span class="user-info">
              <el-icon><UserFilled /></el-icon>
              {{ userStore.user?.real_name || userStore.user?.username }}
            </span>
            <template #dropdown>
              <el-dropdown-menu>
                <el-dropdown-item command="profile">个人信息</el-dropdown-item>
                <el-dropdown-item command="password">修改密码</el-dropdown-item>
                <el-dropdown-item divided command="logout">退出登录</el-dropdown-item>
              </el-dropdown-menu>
            </template>
          </el-dropdown>
        </div>
      </el-header>
      <el-main class="main-content">
        <router-view />
      </el-main>
    </el-container>
  </el-container>
</template>

<script setup>
import { ref, computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useUserStore } from '@/stores/user'
import { ElMessageBox } from 'element-plus'

const route = useRoute()
const router = useRouter()
const userStore = useUserStore()

const isCollapse = ref(false)

const activeMenu = computed(() => route.path)

const handleCommand = (command) => {
  if (command === 'logout') {
    ElMessageBox.confirm('确定要退出登录吗？', '提示', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    }).then(() => {
      userStore.logout()
      router.push('/login')
    }).catch(() => {})
  } else if (command === 'profile') {
    router.push('/profile')
  } else if (command === 'password') {
    router.push('/password')
  }
}
</script>

<style scoped>
.layout-container {
  height: 100vh;
}

.layout-aside {
  background-color: #304156;
  transition: width 0.3s;
}

.logo {
  height: 60px;
  display: flex;
  align-items: center;
  justify-content: center;
  background-color: #263445;
}

.logo-text {
  color: #fff;
  font-size: 18px;
  font-weight: bold;
}

.layout-menu {
  border-right: none;
}

.layout-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  background-color: #fff;
  box-shadow: 0 1px 4px rgba(0, 21, 41, 0.08);
  padding: 0 20px;
}

.collapse-btn {
  font-size: 20px;
  cursor: pointer;
  color: #606266;
}

.header-right {
  display: flex;
  align-items: center;
}

.user-info {
  display: flex;
  align-items: center;
  cursor: pointer;
  color: #606266;
}

.user-info > * {
  margin-right: 5px;
}
</style>
