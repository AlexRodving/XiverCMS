<template>
  <div>
    <div class="flex justify-between items-center mb-8">
      <h1 class="text-3xl font-bold text-gray-900">Content Types</h1>
      <button
        @click="showCreateModal = true"
        class="px-4 py-2 bg-blue-600 text-white rounded-lg hover:bg-blue-700 transition-colors"
      >
        + Create Content Type
      </button>
    </div>

    <div class="bg-white rounded-lg shadow overflow-hidden">
      <table class="min-w-full divide-y divide-gray-200">
        <thead class="bg-gray-50">
          <tr>
            <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
              UID
            </th>
            <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
              Display Name
            </th>
            <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
              Kind
            </th>
            <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
              Visible
            </th>
            <th class="px-6 py-3 text-right text-xs font-medium text-gray-500 uppercase tracking-wider">
              Actions
            </th>
          </tr>
        </thead>
        <tbody class="bg-white divide-y divide-gray-200">
          <tr v-for="contentType in contentTypes" :key="contentType.id">
            <td class="px-6 py-4 whitespace-nowrap text-sm font-medium text-gray-900">
              {{ contentType.uid }}
            </td>
            <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
              {{ contentType.displayName }}
            </td>
            <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
              {{ contentType.kind }}
            </td>
            <td class="px-6 py-4 whitespace-nowrap">
              <span
                class="px-2 inline-flex text-xs leading-5 font-semibold rounded-full"
                :class="contentType.isVisible ? 'bg-green-100 text-green-800' : 'bg-gray-100 text-gray-800'"
              >
                {{ contentType.isVisible ? 'Yes' : 'No' }}
              </span>
            </td>
            <td class="px-6 py-4 whitespace-nowrap text-right text-sm font-medium">
              <router-link
                :to="`/content-types/${contentType.uid}`"
                class="text-blue-600 hover:text-blue-900 mr-4"
              >
                View
              </router-link>
              <button
                @click="deleteContentType(contentType.uid)"
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
      <div class="bg-white rounded-lg p-6 max-w-4xl w-full mx-4 max-h-[95vh] overflow-y-auto">
        <h2 class="text-2xl font-bold mb-4">Create Content Type</h2>
        <form @submit.prevent="createContentType">
          <div class="space-y-4">
            <div>
              <label class="block text-sm font-medium text-gray-700 mb-1">UID *</label>
              <input
                v-model="newContentType.uid"
                type="text"
                required
                class="w-full px-3 py-2 border border-gray-300 rounded-md"
                placeholder="e.g., article, product"
              />
            </div>
            <div>
              <label class="block text-sm font-medium text-gray-700 mb-1">Display Name *</label>
              <input
                v-model="newContentType.displayName"
                type="text"
                required
                class="w-full px-3 py-2 border border-gray-300 rounded-md"
              />
            </div>
            <div>
              <label class="block text-sm font-medium text-gray-700 mb-1">Description</label>
              <textarea
                v-model="newContentType.description"
                class="w-full px-3 py-2 border border-gray-300 rounded-md"
                rows="3"
              ></textarea>
            </div>
            <div>
              <label class="block text-sm font-medium text-gray-700 mb-1">Kind</label>
              <select v-model="newContentType.kind" class="w-full px-3 py-2 border border-gray-300 rounded-md">
                <option value="collectionType">Collection Type</option>
                <option value="singleType">Single Type</option>
              </select>
            </div>
            <div>
              <label class="flex items-center">
                <input
                  v-model="newContentType.isVisible"
                  type="checkbox"
                  class="mr-2"
                />
                <span class="text-sm font-medium text-gray-700">Is Visible</span>
              </label>
            </div>
            <div class="border-t pt-4">
              <div class="flex items-center justify-between mb-2">
                <label class="block text-sm font-medium text-gray-700">Schema (Fields) *</label>
                <div class="flex gap-2">
                  <button
                    type="button"
                    @click="showJsonEditor = !showJsonEditor"
                    class="px-2 py-1 text-xs border border-gray-300 rounded hover:bg-gray-50"
                  >
                    {{ showJsonEditor ? 'Visual Editor' : 'JSON Editor' }}
                  </button>
                </div>
              </div>
              
              <!-- Visual Field Builder -->
              <div v-if="!showJsonEditor">
                <FieldBuilder
                  v-model="schema"
                  :available-content-types="contentTypes"
                />
              </div>
              
              <!-- JSON Editor (fallback) -->
              <div v-else>
                <textarea
                  v-model="schemaJson"
                  class="w-full px-3 py-2 border border-gray-300 rounded-md font-mono text-sm"
                  rows="10"
                  placeholder='{"title": {"type": "string"}, "content": {"type": "text"}}'
                  @input="updateSchemaFromJson"
                ></textarea>
                <p class="mt-1 text-xs text-gray-500">Edit schema as JSON</p>
              </div>
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
import { ref, onMounted, watch } from 'vue'
import { contentAPI } from '../api/content'
import FieldBuilder from '../components/FieldBuilder.vue'

const contentTypes = ref([])
const showCreateModal = ref(false)
const showJsonEditor = ref(false)
const newContentType = ref({
  uid: '',
  displayName: '',
  description: '',
  kind: 'collectionType',
  isVisible: true,
  schema: {},
})
const schema = ref({})
const schemaJson = ref('{}')

const loadContentTypes = async () => {
  try {
    const response = await contentAPI.getContentTypes()
    contentTypes.value = response.data.data || []
  } catch (error) {
    console.error('Failed to load content types:', error)
  }
}

// Watch schema changes from visual editor
watch(schema, (newSchema) => {
  if (Object.keys(newSchema).length > 0) {
    schemaJson.value = JSON.stringify(newSchema, null, 2)
  }
}, { deep: true })

const updateSchemaFromJson = () => {
  try {
    const parsed = JSON.parse(schemaJson.value)
    schema.value = parsed
  } catch (e) {
    // Invalid JSON, ignore
  }
}

const createContentType = async () => {
  try {
    // Use schema from visual editor or JSON
    const finalSchema = Object.keys(schema.value).length > 0 ? schema.value : JSON.parse(schemaJson.value)
    
    if (Object.keys(finalSchema).length === 0) {
      alert('Please add at least one field to the schema')
      return
    }
    
    await contentAPI.createContentType({
      ...newContentType.value,
      schema: finalSchema,
    })
    showCreateModal.value = false
    schemaJson.value = '{}'
    schema.value = {}
    Object.assign(newContentType.value, {
      uid: '',
      displayName: '',
      description: '',
      kind: 'collectionType',
      isVisible: true,
      schema: {},
    })
    await loadContentTypes()
  } catch (error) {
    alert(error.response?.data?.error || 'Failed to create content type')
  }
}

const deleteContentType = async (uid) => {
  if (!confirm(`Are you sure you want to delete content type "${uid}"?`)) {
    return
  }
  try {
    await contentAPI.deleteContentType(uid)
    await loadContentTypes()
  } catch (error) {
    alert(error.response?.data?.error || 'Failed to delete content type')
  }
}

onMounted(loadContentTypes)
</script>

