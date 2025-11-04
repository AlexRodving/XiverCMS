import apiClient from './client'

export const usersAPI = {
  getAll: (params = {}) => 
    apiClient.get('/users', { params }),
  
  getOne: (id) => 
    apiClient.get(`/users/${id}`),
  
  update: (id, data) => 
    apiClient.put(`/users/${id}`, data),
  
  delete: (id) => 
    apiClient.delete(`/users/${id}`),
  
  changePassword: (id, password) => 
    apiClient.put(`/users/${id}/password`, { password }),
}

