import apiClient from './client'

export const rolesAPI = {
  getAll: (params = {}) => 
    apiClient.get('/roles', { params }),
  
  getOne: (id) => 
    apiClient.get(`/roles/${id}`),
  
  create: (data) => 
    apiClient.post('/roles', data),
  
  update: (id, data) => 
    apiClient.put(`/roles/${id}`, data),
  
  delete: (id) => 
    apiClient.delete(`/roles/${id}`),
  
  getPublic: () => 
    apiClient.get('/roles/public'),
}
