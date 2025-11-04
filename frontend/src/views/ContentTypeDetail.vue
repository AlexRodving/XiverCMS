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
        <h2 class="text-xl font-semibold mb-4">Details</h2>
        <dl class="grid grid-cols-2 gap-4">
          <div>
            <dt class="text-sm font-medium text-gray-500">UID</dt>
            <dd class="mt-1 text-sm text-gray-900">{{ contentType.uid }}</dd>
          </div>
          <div>
            <dt class="text-sm font-medium text-gray-500">Kind</dt>
            <dd class="mt-1 text-sm text-gray-900">{{ contentType.kind }}</dd>
          </div>
          <div>
            <dt class="text-sm font-medium text-gray-500">Display Name</dt>
            <dd class="mt-1 text-sm text-gray-900">{{ contentType.displayName }}</dd>
          </div>
          <div>
            <dt class="text-sm font-medium text-gray-500">Visible</dt>
            <dd class="mt-1 text-sm text-gray-900">{{ contentType.isVisible ? 'Yes' : 'No' }}</dd>
          </div>
          <div class="col-span-2">
            <dt class="text-sm font-medium text-gray-500">Description</dt>
            <dd class="mt-1 text-sm text-gray-900">{{ contentType.description || 'No description' }}</dd>
          </div>
        </dl>
      </div>

      <div class="bg-white rounded-lg shadow p-6">
        <h2 class="text-xl font-semibold mb-4">Schema</h2>
        <pre class="bg-gray-50 p-4 rounded-md overflow-auto text-sm">{{ JSON.stringify(contentType.schema, null, 2) }}</pre>
      </div>

      <div class="bg-white rounded-lg shadow p-6">
        <h2 class="text-xl font-semibold mb-4">Public API Endpoints</h2>
        <div>
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
        </div>
      </div>

      <div class="flex justify-end space-x-4">
        <router-link
          :to="`/content-types/${contentType.uid}/entries`"
          class="px-4 py-2 bg-blue-600 text-white rounded-lg hover:bg-blue-700"
        >
          View Entries
        </router-link>
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

