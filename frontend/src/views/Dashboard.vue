<template>
  <div>
    <h1 class="text-3xl font-bold text-gray-900 mb-8">{{ $t('dashboard.title') }}</h1>
    
    <div class="grid grid-cols-1 md:grid-cols-3 gap-6 mb-8">
      <div class="bg-white rounded-lg shadow p-6">
        <h3 class="text-sm font-medium text-gray-500 mb-2">{{ $t('dashboard.totalContentTypes') }}</h3>
        <p class="text-3xl font-bold text-gray-900">{{ stats.contentTypes }}</p>
      </div>
      <div class="bg-white rounded-lg shadow p-6">
        <h3 class="text-sm font-medium text-gray-500 mb-2">{{ $t('dashboard.totalEntries') }}</h3>
        <p class="text-3xl font-bold text-gray-900">{{ stats.totalEntries }}</p>
      </div>
      <div class="bg-white rounded-lg shadow p-6">
        <h3 class="text-sm font-medium text-gray-500 mb-2">{{ $t('dashboard.totalUsers') }}</h3>
        <p class="text-3xl font-bold text-gray-900">{{ stats.users }}</p>
      </div>
    </div>

    <div class="bg-white rounded-lg shadow">
      <div class="p-6 border-b border-gray-200">
        <h2 class="text-xl font-semibold text-gray-900">{{ $t('dashboard.welcome') }}</h2>
      </div>
      <div class="p-6">
        <p class="text-gray-600 mb-4">
          {{ $t('dashboard.description') }}
        </p>
        <div class="space-y-2">
          <h3 class="font-semibold text-gray-900">{{ $t('dashboard.features') }}:</h3>
          <ul class="list-disc list-inside text-gray-600 space-y-1">
            <li>{{ $t('dashboard.feature1') }}</li>
            <li>{{ $t('dashboard.feature2') }}</li>
            <li>{{ $t('dashboard.feature3') }}</li>
            <li>{{ $t('dashboard.feature4') }}</li>
            <li>{{ $t('dashboard.feature5') }}</li>
          </ul>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { contentAPI } from '../api/content'
import { usersAPI } from '../api/users'

const { t } = useI18n()

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

