<template>
  <div>
    <div class="flex justify-between items-center mb-8">
      <div>
        <router-link to="/content-types" class="text-blue-600 hover:text-blue-800 mb-2 inline-block">
          ← Back to Content Types
        </router-link>
        <h1 class="text-3xl font-bold text-gray-900">{{ contentTypeUID }} Entries</h1>
      </div>
      <button
        @click="showCreateModal = true"
        class="px-4 py-2 bg-blue-600 text-white rounded-lg hover:bg-blue-700 transition-colors"
      >
        + Create Entry
      </button>
    </div>

    <div class="bg-white rounded-lg shadow overflow-hidden">
      <table class="min-w-full divide-y divide-gray-200">
        <thead class="bg-gray-50">
          <tr>
            <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
              ID
            </th>
            <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
              Data
            </th>
            <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
              Status
            </th>
            <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
              Created
            </th>
            <th class="px-6 py-3 text-right text-xs font-medium text-gray-500 uppercase tracking-wider">
              Actions
            </th>
          </tr>
        </thead>
        <tbody class="bg-white divide-y divide-gray-200">
          <tr v-for="entry in entries" :key="entry.id">
            <td class="px-6 py-4 whitespace-nowrap text-sm font-medium text-gray-900">
              {{ entry.id }}
            </td>
            <td class="px-6 py-4 text-sm text-gray-500">
              <pre class="text-xs">{{ JSON.stringify(entry.data, null, 2).substring(0, 100) }}...</pre>
            </td>
            <td class="px-6 py-4 whitespace-nowrap">
              <span
                class="px-2 inline-flex text-xs leading-5 font-semibold rounded-full"
                :class="entry.status === 'published' ? 'bg-green-100 text-green-800' : 'bg-yellow-100 text-yellow-800'"
              >
                {{ entry.status }}
              </span>
            </td>
            <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
              {{ new Date(entry.createdAt).toLocaleDateString() }}
            </td>
            <td class="px-6 py-4 whitespace-nowrap text-right text-sm font-medium">
              <router-link
                :to="`/content-types/${contentTypeUID}/entries/${entry.id}`"
                class="text-blue-600 hover:text-blue-900 mr-4"
              >
                Edit
              </router-link>
              <button
                @click="deleteEntry(entry.id)"
                class="text-red-600 hover:text-red-900"
              >
                Delete
              </button>
            </td>
          </tr>
        </tbody>
      </table>
    </div>

    <!-- Create Modal -->
    <div v-if="showCreateModal" class="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50">
      <div class="bg-white rounded-lg p-6 max-w-2xl w-full mx-4 max-h-[90vh] overflow-y-auto">
        <h2 class="text-2xl font-bold mb-4">Create Entry</h2>
        <form @submit.prevent="createEntry">
          <div class="space-y-4">
            <div>
              <label class="block text-sm font-medium text-gray-700 mb-1">Data (JSON) *</label>
              <textarea
                v-model="entryDataJson"
                required
                class="w-full px-3 py-2 border border-gray-300 rounded-md font-mono text-sm"
                rows="10"
                placeholder='{"title": "Example", "content": "Example content"}'
              ></textarea>
            </div>
            <div>
              <label class="block text-sm font-medium text-gray-700 mb-1">Status</label>
              <select v-model="entryStatus" class="w-full px-3 py-2 border border-gray-300 rounded-md">
                <option value="draft">Draft</option>
                <option value="published">Published</option>
              </select>
            </div>
          </div>
          <div class="mt-6 flex justify-end space-x-3">
            <button
              type="button"
              @click="showCreateModal = false"
              class="px-4 py-2 border border-gray-300 rounded-md text-gray-700 hover:bg-gray-50"
            >
              Cancel
            </button>
            <button
              type="submit"
              class="px-4 py-2 bg-blue-600 text-white rounded-md hover:bg-blue-700"
            >
              Create
            </button>
          </div>
        </form>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted, computed } from 'vue'
import { useRoute } from 'vue-router'
import { contentAPI } from '../api/content'

const route = useRoute()
const contentTypeUID = computed(() => route.params.uid)
const entries = ref([])
const showCreateModal = ref(false)
const entryDataJson = ref('{}')
const entryStatus = ref('draft')

const loadEntries = async () => {
  try {
    const response = await contentAPI.getEntries(contentTypeUID.value)
    entries.value = response.data.data || []
  } catch (error) {
    console.error('Failed to load entries:', error)
  }
}

const createEntry = async () => {
  try {
    const data = JSON.parse(entryDataJson.value)
    await contentAPI.createEntry(contentTypeUID.value, {
      data,
      status: entryStatus.value,
    })
    showCreateModal.value = false
    entryDataJson.value = '{}'
    entryStatus.value = 'draft'
    await loadEntries()
  } catch (error) {
    alert(error.response?.data?.error || 'Failed to create entry')
  }
}

const deleteEntry = async (id) => {
  if (!confirm('Are you sure you want to delete this entry?')) {
    return
  }
  try {
    await contentAPI.deleteEntry(contentTypeUID.value, id)
    await loadEntries()
  } catch (error) {
    alert(error.response?.data?.error || 'Failed to delete entry')
  }
}

onMounted(loadEntries)
</script>

