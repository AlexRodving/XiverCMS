import apiClient from './client'

export const contentAPI = {
  // Content Types
  getContentTypes: (params = {}) => 
    apiClient.get('/content-types', { params }),
  
  getContentType: (uid) => 
    apiClient.get(`/content-types/${uid}`),
  
  createContentType: (data) => 
    apiClient.post('/content-types', data),
  
  updateContentType: (uid, data) => 
    apiClient.put(`/content-types/${uid}`, data),
  
  deleteContentType: (uid) => 
    apiClient.delete(`/content-types/${uid}`),
  
  // Content Entries (admin endpoints - require authentication)
  getEntries: (uid, params = {}) => 
    apiClient.get(`/admin/content-types/${uid}/entries`, { params }),
  
  getEntry: (uid, id, params = {}) => 
    apiClient.get(`/admin/content-types/${uid}/entries/${id}`, { params }),
  
  createEntry: (uid, data) => 
    apiClient.post(`/admin/content-types/${uid}/entries`, data),
  
  updateEntry: (uid, id, data) => 
    apiClient.put(`/admin/content-types/${uid}/entries/${id}`, data),
  
  deleteEntry: (uid, id) => 
    apiClient.delete(`/admin/content-types/${uid}/entries/${id}`),
}

