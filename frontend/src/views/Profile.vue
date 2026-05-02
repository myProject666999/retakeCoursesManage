<template>
  <div class="profile">
    <el-card shadow="never">
      <template #header>
        <span class="card-title">个人信息</span>
      </template>
      <el-form ref="profileForm" :model="profileForm" :rules="profileRules" label-width="100px" style="max-width: 500px;">
        <el-form-item label="用户名">
          <el-input v-model="profileForm.username" disabled />
        </el-form-item>
        <el-form-item label="姓名">
          <el-input v-model="profileForm.real_name" placeholder="请输入姓名" />
        </el-form-item>
        <el-form-item label="邮箱" prop="email">
          <el-input v-model="profileForm.email" placeholder="请输入邮箱" />
        </el-form-item>
        <el-form-item label="电话" prop="phone">
          <el-input v-model="profileForm.phone" placeholder="请输入电话" />
        </el-form-item>
        <el-form-item label="角色">
          <el-tag :type="roleTagType">{{ roleName }}</el-tag>
        </el-form-item>
        <el-form-item v-if="profileForm.class_name" label="班级">
          <el-input v-model="profileForm.class_name" disabled />
        </el-form-item>
        <el-form-item v-if="profileForm.major_name" label="专业">
          <el-input v-model="profileForm.major_name" disabled />
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="handleSubmit">保存修改</el-button>
        </el-form-item>
      </el-form>
    </el-card>
  </div>
</template>

<script setup>
import { ref, reactive, computed, onMounted } from 'vue'
import { ElMessage } from 'element-plus'
import { useUserStore } from '@/stores/user'
import { updateProfile } from '@/api'

const userStore = useUserStore()

const profileForm = reactive({
  username: '',
  real_name: '',
  email: '',
  phone: '',
  class_name: '',
  major_name: ''
})

const profileRules = {
  email: [
    { type: 'email', message: '请输入正确的邮箱地址', trigger: 'blur' }
  ],
  phone: [
    { pattern: /^1[3-9]\d{9}$/, message: '请输入正确的手机号', trigger: 'blur' }
  ]
}

const roleName = computed(() => {
  const roles = {
    student: '学生',
    teacher: '教师',
    admin: '管理员'
  }
  return roles[userStore.role] || '未知'
})

const roleTagType = computed(() => {
  const types = {
    student: 'warning',
    teacher: 'primary',
    admin: 'danger'
  }
  return types[userStore.role] || 'info'
})

const loadUserInfo = () => {
  if (userStore.user) {
    profileForm.username = userStore.user.username || ''
    profileForm.real_name = userStore.user.real_name || ''
    profileForm.email = userStore.user.email || ''
    profileForm.phone = userStore.user.phone || ''
    profileForm.class_name = userStore.user.class_name || ''
    profileForm.major_name = userStore.user.major_name || ''
  }
}

const handleSubmit = async () => {
  try {
    await updateProfile({
      real_name: profileForm.real_name,
      email: profileForm.email,
      phone: profileForm.phone
    })
    if (userStore.user) {
      userStore.user.real_name = profileForm.real_name
      userStore.user.email = profileForm.email
      userStore.user.phone = profileForm.phone
      localStorage.setItem('user', JSON.stringify(userStore.user))
    }
    ElMessage.success('修改成功')
  } catch (error) {
    console.error(error)
  }
}

onMounted(() => {
  loadUserInfo()
})
</script>

<style scoped>
.profile {
  padding: 0;
}

.card-title {
  font-size: 16px;
  font-weight: bold;
}
</style>
