import apiClient from './client'

export const authAPI = {
  login: (email, password) => 
    apiClient.post('/auth/login', { email, password }),
  
  register: (data) => 
    apiClient.post('/auth/register', data),
  
  me: () => 
    apiClient.get('/auth/me'),
}

