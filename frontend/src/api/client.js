import axios from 'axios'

// Use VITE_API_URL from environment if set, otherwise use relative path (proxy)
const apiUrl = import.meta.env.VITE_API_URL || ''
const baseURL = apiUrl 
  ? `${apiUrl.replace(/\/$/, '')}/api`  // Remove trailing slash if present
  : '/api'

console.log('API Client initialized:')
console.log('  VITE_API_URL:', import.meta.env.VITE_API_URL || '(not set, using proxy)')
console.log('  baseURL:', baseURL)

const apiClient = axios.create({
  baseURL: baseURL,
  headers: {
    'Content-Type': 'application/json',
  },
})

// Request interceptor to add auth token
apiClient.interceptors.request.use(
  (config) => {
    const token = localStorage.getItem('token')
    if (token) {
      config.headers.Authorization = `Bearer ${token}`
    }
    return config
  },
  (error) => {
    return Promise.reject(error)
  }
)

// Response interceptor to handle errors
apiClient.interceptors.response.use(
  (response) => response,
  (error) => {
    if (error.response?.status === 401) {
      localStorage.removeItem('token')
      localStorage.removeItem('user')
      window.location.href = '/login'
    }
    return Promise.reject(error)
  }
)

export default apiClient

