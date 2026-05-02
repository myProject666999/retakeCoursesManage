import { createRouter, createWebHistory } from 'vue-router'
import { useUserStore } from '@/stores/user'
import { ElMessage } from 'element-plus'

const routes = [
  {
    path: '/login',
    name: 'Login',
    component: () => import('@/views/Login.vue'),
    meta: { title: '登录' }
  },
  {
    path: '/',
    component: () => import('@/views/Layout.vue'),
    redirect: '/dashboard',
    children: [
      {
        path: 'dashboard',
        name: 'Dashboard',
        component: () => import('@/views/Dashboard.vue'),
        meta: { title: '首页', requiresAuth: true }
      },
      {
        path: 'profile',
        name: 'Profile',
        component: () => import('@/views/Profile.vue'),
        meta: { title: '个人信息', requiresAuth: true }
      },
      {
        path: 'password',
        name: 'Password',
        component: () => import('@/views/Password.vue'),
        meta: { title: '密码修改', requiresAuth: true }
      },
      {
        path: 'announcements',
        name: 'Announcements',
        component: () => import('@/views/student/Announcements.vue'),
        meta: { title: '公告查询', requiresAuth: true, roles: ['student', 'teacher', 'admin'] }
      },
      {
        path: 'announcements/:id',
        name: 'AnnouncementDetail',
        component: () => import('@/views/student/AnnouncementDetail.vue'),
        meta: { title: '公告详情', requiresAuth: true, roles: ['student', 'teacher', 'admin'] }
      },
      {
        path: 'student/retake-courses',
        name: 'StudentRetakeCourses',
        component: () => import('@/views/student/RetakeCourses.vue'),
        meta: { title: '重修课程查询', requiresAuth: true, roles: ['student'] }
      },
      {
        path: 'student/grades',
        name: 'StudentGrades',
        component: () => import('@/views/student/Grades.vue'),
        meta: { title: '课程成绩查询', requiresAuth: true, roles: ['student'] }
      },
      {
        path: 'student/reports',
        name: 'StudentReports',
        component: () => import('@/views/student/Reports.vue'),
        meta: { title: '建议报告', requiresAuth: true, roles: ['student'] }
      },
      {
        path: 'student/applications',
        name: 'StudentApplications',
        component: () => import('@/views/student/Applications.vue'),
        meta: { title: '重修申请管理', requiresAuth: true, roles: ['student'] }
      },
      {
        path: 'teacher/students',
        name: 'TeacherStudents',
        component: () => import('@/views/teacher/Students.vue'),
        meta: { title: '学生查询', requiresAuth: true, roles: ['teacher'] }
      },
      {
        path: 'teacher/teaching-courses',
        name: 'TeacherTeachingCourses',
        component: () => import('@/views/teacher/TeachingCourses.vue'),
        meta: { title: '重修课程管理', requiresAuth: true, roles: ['teacher'] }
      },
      {
        path: 'teacher/grades',
        name: 'TeacherGrades',
        component: () => import('@/views/teacher/Grades.vue'),
        meta: { title: '成绩管理', requiresAuth: true, roles: ['teacher'] }
      },
      {
        path: 'teacher/reports',
        name: 'TeacherReports',
        component: () => import('@/views/teacher/Reports.vue'),
        meta: { title: '报告管理', requiresAuth: true, roles: ['teacher'] }
      },
      {
        path: 'teacher/applications',
        name: 'TeacherApplications',
        component: () => import('@/views/teacher/Applications.vue'),
        meta: { title: '重修申请管理', requiresAuth: true, roles: ['teacher'] }
      },
      {
        path: 'admin/majors',
        name: 'AdminMajors',
        component: () => import('@/views/admin/Majors.vue'),
        meta: { title: '专业管理', requiresAuth: true, roles: ['admin'] }
      },
      {
        path: 'admin/classes',
        name: 'AdminClasses',
        component: () => import('@/views/admin/Classes.vue'),
        meta: { title: '班级管理', requiresAuth: true, roles: ['admin'] }
      },
      {
        path: 'admin/users',
        name: 'AdminUsers',
        component: () => import('@/views/admin/Users.vue'),
        meta: { title: '用户管理', requiresAuth: true, roles: ['admin'] }
      },
      {
        path: 'admin/courses',
        name: 'AdminCourses',
        component: () => import('@/views/admin/Courses.vue'),
        meta: { title: '课程管理', requiresAuth: true, roles: ['admin'] }
      },
      {
        path: 'admin/retake-courses',
        name: 'AdminRetakeCourses',
        component: () => import('@/views/admin/RetakeCourses.vue'),
        meta: { title: '重修课程管理', requiresAuth: true, roles: ['admin'] }
      },
      {
        path: 'admin/grades',
        name: 'AdminGrades',
        component: () => import('@/views/admin/Grades.vue'),
        meta: { title: '成绩管理', requiresAuth: true, roles: ['admin'] }
      },
      {
        path: 'admin/announcements',
        name: 'AdminAnnouncements',
        component: () => import('@/views/admin/Announcements.vue'),
        meta: { title: '公告管理', requiresAuth: true, roles: ['admin'] }
      },
      {
        path: 'admin/reports',
        name: 'AdminReports',
        component: () => import('@/views/admin/Reports.vue'),
        meta: { title: '建议报告管理', requiresAuth: true, roles: ['admin'] }
      },
      {
        path: 'admin/applications',
        name: 'AdminApplications',
        component: () => import('@/views/admin/Applications.vue'),
        meta: { title: '重修申请管理', requiresAuth: true, roles: ['admin'] }
      }
    ]
  }
]

const router = createRouter({
  history: createWebHistory(),
  routes
})

router.beforeEach(async (to, from, next) => {
  document.title = to.meta.title ? `${to.meta.title} - 重修课程信息管理系统` : '重修课程信息管理系统'

  const userStore = useUserStore()

  if (to.meta.requiresAuth) {
    if (!userStore.isLoggedIn) {
      ElMessage.warning('请先登录')
      next('/login')
      return
    }

    if (!userStore.user) {
      try {
        await userStore.fetchUserInfo()
      } catch (error) {
        ElMessage.warning('登录状态已失效，请重新登录')
        userStore.logout()
        next('/login')
        return
      }
    }

    if (to.meta.roles && to.meta.roles.length > 0) {
      const hasRole = to.meta.roles.includes(userStore.role)
      if (!hasRole) {
        ElMessage.error('权限不足')
        next(from.path || '/dashboard')
        return
      }
    }
  }

  if (to.path === '/login' && userStore.isLoggedIn) {
    next('/dashboard')
    return
  }

  next()
})

export default router
