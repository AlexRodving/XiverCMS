import { defineStore } from 'pinia'
import { authAPI } from '../api/auth'
import router from '../router'

export const useAuthStore = defineStore('auth', {
  state: () => ({
    token: localStorage.getItem('token') || null,
    user: JSON.parse(localStorage.getItem('user') || 'null'),
    isAuthenticated: !!localStorage.getItem('token'),
  }),

  actions: {
    async login(email, password) {
      try {
        const response = await authAPI.login(email, password)
        this.token = response.data.jwt
        this.user = response.data.user
        this.isAuthenticated = true
        
        localStorage.setItem('token', this.token)
        localStorage.setItem('user', JSON.stringify(this.user))
        
        return { success: true }
      } catch (error) {
        return { 
          success: false, 
          error: error.response?.data?.error || 'Login failed' 
        }
      }
    },

    async register(data) {
      try {
        const response = await authAPI.register(data)
        this.token = response.data.jwt
        this.user = response.data.user
        this.isAuthenticated = true
        
        localStorage.setItem('token', this.token)
        localStorage.setItem('user', JSON.stringify(this.user))
        
        return { success: true }
      } catch (error) {
        return { 
          success: false, 
          error: error.response?.data?.error || 'Registration failed' 
        }
      }
    },

    async fetchMe() {
      try {
        const response = await authAPI.me()
        this.user = response.data
        localStorage.setItem('user', JSON.stringify(this.user))
        return { success: true }
      } catch (error) {
        this.logout()
        return { success: false }
      }
    },

    logout() {
      this.token = null
      this.user = null
      this.isAuthenticated = false
      localStorage.removeItem('token')
      localStorage.removeItem('user')
      router.push('/login')
    },

    isSuperAdmin() {
      return this.user?.isSuperAdmin || false
    },
  },
})

