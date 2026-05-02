import { defineStore } from 'pinia'
import { login, getCurrentUser } from '@/api'

export const useUserStore = defineStore('user', {
  state: () => ({
    token: localStorage.getItem('token') || '',
    user: JSON.parse(localStorage.getItem('user') || 'null'),
    role: localStorage.getItem('role') || ''
  }),

  getters: {
    isLoggedIn: (state) => !!state.token,
    isStudent: (state) => state.role === 'student',
    isTeacher: (state) => state.role === 'teacher',
    isAdmin: (state) => state.role === 'admin'
  },

  actions: {
    async login(loginForm) {
      const res = await login(loginForm)
      this.token = res.token
      this.user = res.user
      this.role = res.user.role

      localStorage.setItem('token', res.token)
      localStorage.setItem('user', JSON.stringify(res.user))
      localStorage.setItem('role', res.user.role)

      return res
    },

    async fetchUserInfo() {
      const res = await getCurrentUser()
      this.user = res
      localStorage.setItem('user', JSON.stringify(res))
      return res
    },

    logout() {
      this.token = ''
      this.user = null
      this.role = ''

      localStorage.removeItem('token')
      localStorage.removeItem('user')
      localStorage.removeItem('role')
    }
  }
})
