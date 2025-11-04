<template>
  <div>
    <h1 class="text-3xl font-bold text-gray-900 mb-8">Dashboard</h1>
    
    <div class="grid grid-cols-1 md:grid-cols-3 gap-6 mb-8">
      <div class="bg-white rounded-lg shadow p-6">
        <h3 class="text-sm font-medium text-gray-500 mb-2">Content Types</h3>
        <p class="text-3xl font-bold text-gray-900">{{ stats.contentTypes }}</p>
      </div>
      <div class="bg-white rounded-lg shadow p-6">
        <h3 class="text-sm font-medium text-gray-500 mb-2">Total Entries</h3>
        <p class="text-3xl font-bold text-gray-900">{{ stats.totalEntries }}</p>
      </div>
      <div class="bg-white rounded-lg shadow p-6">
        <h3 class="text-sm font-medium text-gray-500 mb-2">Users</h3>
        <p class="text-3xl font-bold text-gray-900">{{ stats.users }}</p>
      </div>
    </div>

    <div class="bg-white rounded-lg shadow">
      <div class="p-6 border-b border-gray-200">
        <h2 class="text-xl font-semibold text-gray-900">Welcome to XiverCRM</h2>
      </div>
      <div class="p-6">
        <p class="text-gray-600 mb-4">
          XiverCRM is a headless CMS built with Go, Gin, GORM, and Vue3. 
          It provides a powerful and flexible content management system similar to Strapi.
        </p>
        <div class="space-y-2">
          <h3 class="font-semibold text-gray-900">Features:</h3>
          <ul class="list-disc list-inside text-gray-600 space-y-1">
            <li>Dynamic Content Types</li>
            <li>User & Role Management</li>
            <li>RESTful API</li>
            <li>JWT Authentication</li>
            <li>Modern Admin Interface</li>
          </ul>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { contentAPI } from '../api/content'
import { usersAPI } from '../api/users'

const stats = ref({
  contentTypes: 0,
  totalEntries: 0,
  users: 0,
})

onMounted(async () => {
  try {
    const [contentTypesRes, usersRes] = await Promise.all([
      contentAPI.getContentTypes(),
      usersAPI.getAll(),
    ])
    
    stats.value.contentTypes = contentTypesRes.data.meta?.pagination?.total || 0
    stats.value.users = usersRes.data.meta?.pagination?.total || 0
    
    // Count total entries
    let totalEntries = 0
    for (const contentType of contentTypesRes.data.data || []) {
      try {
        const entriesRes = await contentAPI.getEntries(contentType.uid)
        totalEntries += entriesRes.data.meta?.pagination?.total || 0
      } catch (e) {
        // Ignore errors
      }
    }
    stats.value.totalEntries = totalEntries
  } catch (error) {
    console.error('Failed to load stats:', error)
  }
})
</script>

