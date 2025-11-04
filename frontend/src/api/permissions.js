import apiClient from './client'

export const permissionsAPI = {
  getAll: () => 
    apiClient.get('/permissions'),
  
  getOne: (id) => 
    apiClient.get(`/permissions/${id}`),
  
  create: (data) => 
    apiClient.post('/permissions', data),
  
  update: (id, data) => 
    apiClient.put(`/permissions/${id}`, data),
  
  delete: (id) => 
    apiClient.delete(`/permissions/${id}`),
}
