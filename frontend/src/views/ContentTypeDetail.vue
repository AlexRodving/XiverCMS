<template>
  <div>
    <div class="mb-8">
      <router-link to="/content-types" class="text-blue-600 hover:text-blue-800 mb-4 inline-block">
        ← Back to Content Types
      </router-link>
      <h1 class="text-3xl font-bold text-gray-900">{{ contentType?.displayName || 'Content Type' }}</h1>
    </div>

    <div v-if="loading" class="text-center py-12">
      <div class="text-gray-500">Loading...</div>
    </div>

    <div v-else-if="contentType" class="space-y-6">
      <div class="bg-white rounded-lg shadow p-6">
        <h2 class="text-xl font-semibold mb-4">Schema</h2>
        <pre class="bg-gray-50 p-4 rounded-md overflow-auto text-sm font-mono">{{ JSON.stringify(contentType.schema, null, 2) }}</pre>
      </div>

      <div class="bg-white rounded-lg shadow p-6">
        <h2 class="text-xl font-semibold mb-4">Public API Endpoints</h2>
        <div class="space-y-3">
          <div>
            <dt class="text-sm font-medium text-gray-500 mb-1">Get Entries List</dt>
            <dd class="flex items-center gap-2">
              <code class="flex-1 bg-gray-50 px-3 py-2 rounded-md text-sm font-mono">{{ getAPIBaseURL() }}/api/{{ contentType.uid }}</code>
              <button
                @click="copyToClipboard(`${getAPIBaseURL()}/api/${contentType.uid}`)"
                class="px-3 py-2 bg-blue-600 text-white rounded-md hover:bg-blue-700 text-sm"
                title="Copy URL"
              >
                📋
              </button>
            </dd>
          </div>
          <div>
            <dt class="text-sm font-medium text-gray-500 mb-1">Get Single Entry (example)</dt>
            <dd class="flex items-center gap-2">
              <code class="flex-1 bg-gray-50 px-3 py-2 rounded-md text-sm font-mono">{{ getAPIBaseURL() }}/api/{{ contentType.uid }}/:id</code>
              <button
                @click="copyToClipboard(`${getAPIBaseURL()}/api/${contentType.uid}/:id`)"
                class="px-3 py-2 bg-blue-600 text-white rounded-md hover:bg-blue-700 text-sm"
                title="Copy URL"
              >
                📋
              </button>
            </dd>
          </div>
          <div class="mt-4 p-3 bg-blue-50 rounded-md text-sm text-blue-800">
            <p class="font-medium mb-1">ℹ️ Access Control:</p>
            <p v-if="contentType.accessType === 'public'" class="text-xs">This endpoint is <strong>public</strong> - no authentication required.</p>
            <p v-else-if="contentType.accessType === 'authenticated'" class="text-xs">This endpoint requires <strong>authentication</strong> - add <code>Authorization: Bearer TOKEN</code> header.</p>
            <p v-else-if="contentType.accessType === 'moderator'" class="text-xs">This endpoint requires <strong>Moderator or Admin</strong> role.</p>
            <p v-else-if="contentType.accessType === 'admin'" class="text-xs">This endpoint requires <strong>Admin</strong> role.</p>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRoute } from 'vue-router'
import { contentAPI } from '../api/content'

const route = useRoute()
const contentType = ref(null)
const loading = ref(true)

const getAPIBaseURL = () => {
  // Get base URL from environment or use current origin
  return import.meta.env.VITE_API_URL || window.location.origin.replace(':5173', ':8080')
}

const copyToClipboard = async (text) => {
  try {
    await navigator.clipboard.writeText(text)
    alert('URL copied to clipboard!')
  } catch (err) {
    // Fallback for older browsers
    const textarea = document.createElement('textarea')
    textarea.value = text
    document.body.appendChild(textarea)
    textarea.select()
    document.execCommand('copy')
    document.body.removeChild(textarea)
    alert('URL copied to clipboard!')
  }
}

onMounted(async () => {
  try {
    const response = await contentAPI.getContentType(route.params.uid)
    contentType.value = response.data
  } catch (error) {
    console.error('Failed to load content type:', error)
  } finally {
    loading.value = false
  }
})
</script>

