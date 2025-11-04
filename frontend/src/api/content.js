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
  
  // Content Entries
  getEntries: (uid, params = {}) => 
    apiClient.get(`/content-types/${uid}/entries`, { params }),
  
  getEntry: (uid, id) => 
    apiClient.get(`/content-types/${uid}/entries/${id}`),
  
  createEntry: (uid, data) => 
    apiClient.post(`/content-types/${uid}/entries`, data),
  
  updateEntry: (uid, id, data) => 
    apiClient.put(`/content-types/${uid}/entries/${id}`, data),
  
  deleteEntry: (uid, id) => 
    apiClient.delete(`/content-types/${uid}/entries/${id}`),
}

