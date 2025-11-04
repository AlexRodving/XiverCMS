<template>
  <div>
    <div class="flex justify-between items-center mb-8">
      <div>
        <router-link to="/content-types" class="text-blue-600 hover:text-blue-800 mb-2 inline-block">
          ← {{ $t('common.back') }} {{ $t('navigation.contentTypes') }}
        </router-link>
        <h1 class="text-3xl font-bold text-gray-900">{{ contentTypeUID }} {{ $t('contentEntries.title') }}</h1>
      </div>
      <button
        @click="showCreateModal = true"
        class="px-4 py-2 bg-blue-600 text-white rounded-lg hover:bg-blue-700 transition-colors"
      >
        + {{ $t('contentEntries.create') }}
      </button>
    </div>

    <!-- Search and Filters -->
    <div class="bg-white rounded-lg shadow p-4 mb-6">
      <div class="grid grid-cols-1 md:grid-cols-3 gap-4">
        <!-- Search -->
        <div class="md:col-span-2">
          <div class="relative">
            <div class="absolute inset-y-0 left-0 pl-3 flex items-center pointer-events-none">
              <MagnifyingGlassIcon class="h-5 w-5 text-gray-400" />
            </div>
            <input
              v-model="searchQuery"
              type="text"
              :placeholder="$t('common.search')"
              class="block w-full pl-10 pr-3 py-2 border border-gray-300 rounded-md leading-5 bg-white placeholder-gray-500 focus:outline-none focus:placeholder-gray-400 focus:ring-1 focus:ring-blue-500 focus:border-blue-500 sm:text-sm"
              @input="debouncedSearch"
            />
          </div>
        </div>
        
        <!-- Status Filter -->
        <div>
          <select
            v-model="statusFilter"
            @change="applyFilters"
            class="block w-full px-3 py-2 border border-gray-300 rounded-md bg-white focus:outline-none focus:ring-blue-500 focus:border-blue-500 sm:text-sm"
          >
            <option value="">{{ $t('filters.allStatuses') }}</option>
            <option value="draft">{{ $t('contentEntries.draft') }}</option>
            <option value="published">{{ $t('contentEntries.published') }}</option>
          </select>
        </div>
      </div>
    </div>

    <!-- Empty State -->
    <div v-if="!loading && entries.length === 0" class="bg-white rounded-lg shadow p-12 text-center">
      <p class="text-gray-500">{{ $t('contentEntries.noEntries') }}</p>
      <p class="text-sm text-gray-400 mt-2">{{ $t('contentEntries.createFirst') }}</p>
    </div>

    <!-- Loading State -->
    <div v-if="loading" class="bg-white rounded-lg shadow p-12 text-center">
      <div class="text-gray-500">{{ $t('common.loading') }}...</div>
    </div>

    <!-- Table -->
    <div v-else-if="entries.length > 0" class="bg-white rounded-lg shadow overflow-hidden">
      <table class="min-w-full divide-y divide-gray-200">
        <thead class="bg-gray-50">
          <tr>
            <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
              ID
            </th>
            <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
              {{ $t('contentEntries.data') }}
            </th>
            <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
              {{ $t('contentEntries.status') }}
            </th>
            <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
              {{ $t('contentEntries.createdAt') }}
            </th>
            <th class="px-6 py-3 text-right text-xs font-medium text-gray-500 uppercase tracking-wider">
              {{ $t('common.actions') }}
            </th>
          </tr>
        </thead>
        <tbody class="bg-white divide-y divide-gray-200">
          <tr v-for="entry in entries" :key="entry.id" class="hover:bg-gray-50">
            <td class="px-6 py-4 whitespace-nowrap text-sm font-medium text-gray-900">
              {{ entry.id }}
            </td>
            <td class="px-6 py-4 text-sm text-gray-500">
              <div class="max-w-md">
                <div v-if="getEntryDisplayName(entry)" class="font-semibold text-gray-900 mb-1">
                  {{ getEntryDisplayName(entry) }}
                </div>
                <div v-for="(value, key) in getDisplayFields(entry.data)" :key="key" class="truncate text-xs">
                  <span class="font-medium text-gray-600">{{ key }}:</span>
                  <span class="ml-2">{{ formatFieldValue(value) }}</span>
                </div>
                <div v-if="Object.keys(entry.data || {}).length === 0" class="text-gray-400 italic text-xs">
                  {{ $t('contentEntries.noData') }}
                </div>
              </div>
            </td>
            <td class="px-6 py-4 whitespace-nowrap">
              <span
                class="px-2 inline-flex text-xs leading-5 font-semibold rounded-full"
                :class="entry.status === 'published' ? 'bg-green-100 text-green-800' : 'bg-yellow-100 text-yellow-800'"
              >
                {{ entry.status === 'published' ? $t('contentEntries.published') : $t('contentEntries.draft') }}
              </span>
            </td>
            <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
              {{ formatDate(entry.createdAt) }}
            </td>
            <td class="px-6 py-4 whitespace-nowrap text-right text-sm font-medium space-x-2">
              <router-link
                :to="`/content-types/${contentTypeUID}/entries/${entry.id}`"
                class="text-blue-600 hover:text-blue-900"
              >
                {{ $t('common.edit') }}
              </router-link>
              <button
                @click="deleteEntry(entry.id)"
                class="text-red-600 hover:text-red-900"
              >
                {{ $t('common.delete') }}
              </button>
            </td>
          </tr>
        </tbody>
      </table>
      
      <!-- Pagination -->
      <Pagination
        v-if="pagination.total > 0"
        :current-page="pagination.page"
        :page-size="pagination.pageSize"
        :total="pagination.total"
        @change-page="changePage"
      />
    </div>

    <!-- Create Modal -->
    <div v-if="showCreateModal" class="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50">
      <div class="bg-white rounded-lg p-6 max-w-4xl w-full mx-4 max-h-[95vh] overflow-y-auto">
        <div class="flex justify-between items-center mb-4">
          <h2 class="text-2xl font-bold">{{ $t('contentEntries.create') }}</h2>
          <button
            @click="toggleEditorMode"
            class="px-3 py-1 text-sm border border-gray-300 rounded hover:bg-gray-50"
          >
            {{ useVisualEditor ? $t('contentTypes.jsonEditor') : $t('contentTypes.visualEditor') }}
          </button>
        </div>
        
        <form @submit.prevent="createEntry">
          <div class="space-y-4">
            <!-- Visual Editor -->
            <div v-if="useVisualEditor && contentTypeSchema">
              <DynamicForm
                ref="dynamicFormRef"
                :schema="contentTypeSchema"
                v-model="entryData"
                :available-content-types="availableContentTypes"
                @validation="onValidation"
              />
            </div>
            
            <!-- JSON Editor (Fallback) -->
            <div v-else>
              <label class="block text-sm font-medium text-gray-700 mb-1">{{ $t('contentEntries.data') }} *</label>
              <textarea
                v-model="entryDataJson"
                required
                class="w-full px-3 py-2 border border-gray-300 rounded-md font-mono text-sm"
                rows="10"
                placeholder='{"title": "Example", "content": "Example content"}'
              ></textarea>
            </div>
            
            <div>
              <label class="block text-sm font-medium text-gray-700 mb-1">{{ $t('contentEntries.status') }}</label>
              <select v-model="entryStatus" class="w-full px-3 py-2 border border-gray-300 rounded-md">
                <option value="draft">{{ $t('contentEntries.draft') }}</option>
                <option value="published">{{ $t('contentEntries.published') }}</option>
              </select>
            </div>
          </div>
          <div class="mt-6 flex justify-end space-x-3">
            <button
              type="button"
              @click="showCreateModal = false; resetForm()"
              class="px-4 py-2 border border-gray-300 rounded-md text-gray-700 hover:bg-gray-50"
            >
              {{ $t('common.cancel') }}
            </button>
            <button
              type="submit"
              :disabled="!isFormValid"
              class="px-4 py-2 bg-blue-600 text-white rounded-md hover:bg-blue-700 disabled:opacity-50"
            >
              {{ $t('common.create') }}
            </button>
          </div>
        </form>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted, computed, watch } from 'vue'
import { useRoute } from 'vue-router'
import { useI18n } from 'vue-i18n'
import { contentAPI } from '../api/content'
import DynamicForm from '../components/DynamicForm.vue'
import Pagination from '../components/Pagination.vue'
import { MagnifyingGlassIcon } from '@heroicons/vue/24/outline'

const { t } = useI18n()
const route = useRoute()
const contentTypeUID = computed(() => route.params.uid)
const entries = ref([])
const showCreateModal = ref(false)
const entryDataJson = ref('{}')
const entryData = ref({})
const entryStatus = ref('draft')
const contentTypeSchema = ref(null)
const contentType = ref(null)
const availableContentTypes = ref([])
const useVisualEditor = ref(true)
const isFormValid = ref(true)
const dynamicFormRef = ref(null)
const loading = ref(false)

// Search and filters
const searchQuery = ref('')
const statusFilter = ref('')
const searchTimeout = ref(null)

// Pagination
const pagination = ref({
  page: 1,
  pageSize: 10,
  total: 0
})

const loadContentType = async () => {
  try {
    console.log('Loading content type:', contentTypeUID.value)
    const response = await contentAPI.getContentType(contentTypeUID.value)
    contentType.value = response.data
    contentTypeSchema.value = response.data.schema || {}
    console.log('Content type loaded:', contentType.value)
    console.log('Schema:', contentTypeSchema.value)
    
    if (!contentTypeSchema.value || Object.keys(contentTypeSchema.value).length === 0) {
      console.warn('Content type schema is empty!')
      if (window.showToast) {
        window.showToast.warning(t('common.warning'), t('contentEntries.emptySchema'))
      }
    }
  } catch (error) {
    console.error('Failed to load content type:', error)
    console.error('Error response:', error.response?.data)
    if (window.showToast) {
      window.showToast.error(t('common.error'), t('contentEntries.loadContentTypeFailed'))
    }
  }
}

const loadAvailableContentTypes = async () => {
  try {
    const response = await contentAPI.getContentTypes()
    availableContentTypes.value = response.data.data || []
  } catch (error) {
    console.error('Failed to load content types:', error)
  }
}

const loadEntries = async () => {
  loading.value = true
  try {
    const params = {
      page: pagination.value.page,
      pageSize: pagination.value.pageSize
    }
    
    if (statusFilter.value) {
      params.status = statusFilter.value
    }
    
    const response = await contentAPI.getEntries(contentTypeUID.value, params)
    entries.value = response.data.data || []
    
    // Update pagination from response
    if (response.data.meta?.pagination) {
      pagination.value = {
        page: response.data.meta.pagination.page,
        pageSize: response.data.meta.pagination.pageSize,
        total: response.data.meta.pagination.total
      }
    }
  } catch (error) {
    console.error('Failed to load entries:', error)
  } finally {
    loading.value = false
  }
}

const changePage = (page) => {
  pagination.value.page = page
  loadEntries()
  // Scroll to top
  window.scrollTo({ top: 0, behavior: 'smooth' })
}

const applyFilters = () => {
  pagination.value.page = 1 // Reset to first page
  loadEntries()
}

const debouncedSearch = () => {
  if (searchTimeout.value) {
    clearTimeout(searchTimeout.value)
  }
  
  searchTimeout.value = setTimeout(() => {
    // For now, search is client-side
    // TODO: Implement server-side search
    filterEntries()
  }, 300)
}

const filterEntries = async () => {
  // Client-side filtering for now
  // TODO: Implement server-side search API endpoint
  if (!searchQuery.value.trim()) {
    // If search is cleared, reload from server
    await loadEntries()
    return
  }
  
  // For now, we'll reload and filter client-side
  // In future, this should be server-side search
  await loadEntries()
  const query = searchQuery.value.toLowerCase().trim()
  if (query) {
    entries.value = entries.value.filter(entry => {
      const dataStr = JSON.stringify(entry.data).toLowerCase()
      return dataStr.includes(query) || 
             String(entry.id).includes(query) ||
             entry.status.toLowerCase().includes(query)
    })
  }
}

const toggleEditorMode = () => {
  useVisualEditor.value = !useVisualEditor.value
  if (useVisualEditor.value) {
    // Sync JSON to form data
    try {
      entryData.value = JSON.parse(entryDataJson.value)
    } catch (e) {
      entryData.value = {}
    }
  } else {
    // Sync form data to JSON
    entryDataJson.value = JSON.stringify(entryData.value, null, 2)
  }
}

const onValidation = ({ isValid }) => {
  isFormValid.value = isValid
}

watch(entryData, () => {
  if (useVisualEditor.value) {
    entryDataJson.value = JSON.stringify(entryData.value, null, 2)
  }
}, { deep: true })

watch(entryDataJson, () => {
  if (!useVisualEditor.value) {
    try {
      entryData.value = JSON.parse(entryDataJson.value)
    } catch (e) {
      // Invalid JSON, ignore
    }
  }
})

const resetForm = () => {
  entryData.value = {}
  entryDataJson.value = '{}'
  entryStatus.value = 'draft'
  isFormValid.value = true
}

const createEntry = async () => {
  try {
    // Check if content type is loaded
    if (!contentTypeSchema.value || Object.keys(contentTypeSchema.value).length === 0) {
      const errorMsg = t('contentEntries.contentTypeNotLoaded')
      console.error('Content type schema is not loaded:', contentTypeUID.value)
      if (window.showToast) {
        window.showToast.error(t('common.error'), errorMsg)
      } else {
        alert(errorMsg)
      }
      return
    }
    
    let data
    if (useVisualEditor.value && dynamicFormRef.value) {
      // Validate form
      const validation = dynamicFormRef.value.validate()
      if (!validation) {
        if (window.showToast) {
          window.showToast.error(t('common.error'), t('contentEntries.validationFailed'))
        } else {
          alert(t('common.error') + ': ' + t('contentEntries.validationFailed'))
        }
        return
      }
      data = dynamicFormRef.value.getData()
    } else {
      try {
        data = JSON.parse(entryDataJson.value)
      } catch (e) {
        const errorMsg = t('contentEntries.invalidJson')
        console.error('Invalid JSON:', e)
        if (window.showToast) {
          window.showToast.error(t('common.error'), errorMsg)
        } else {
          alert(errorMsg)
        }
        return
      }
    }
    
    console.log('Creating entry for content type:', contentTypeUID.value, 'with data:', data)
    
    await contentAPI.createEntry(contentTypeUID.value, {
      data,
      status: entryStatus.value,
    })
    showCreateModal.value = false
    resetForm()
    await loadEntries()
    
    if (window.showToast) {
      window.showToast.success(t('common.success'), t('contentEntries.createSuccess'))
    }
  } catch (error) {
    console.error('Error creating entry:', error)
    console.error('Error response:', error.response?.data)
    const errorMsg = error.response?.data?.error || error.message || t('contentEntries.createFailed')
    if (window.showToast) {
      window.showToast.error(t('common.error'), errorMsg)
    } else {
      alert(errorMsg + '\n\nCheck console for details.')
    }
  }
}

const deleteEntry = async (id) => {
  if (!confirm(t('contentEntries.confirmDelete'))) {
    return
  }
  try {
    await contentAPI.deleteEntry(contentTypeUID.value, id)
    await loadEntries()
    
    if (window.showToast) {
      window.showToast.success(t('common.success'), t('contentEntries.deleteSuccess'))
    }
  } catch (error) {
    const errorMsg = error.response?.data?.error || t('contentEntries.deleteFailed')
    if (window.showToast) {
      window.showToast.error(t('common.error'), errorMsg)
    } else {
      alert(errorMsg)
    }
  }
}

const getEntryDisplayName = (entry) => {
  if (!entry || !entry.data) return null
  
  // Use contentType schema if available
  if (contentType.value && contentType.value.schema) {
    const schema = contentType.value.schema
    
    // First priority: required field
    for (const [fieldName, fieldConfig] of Object.entries(schema)) {
      if (fieldConfig.required && entry.data[fieldName]) {
        return String(entry.data[fieldName])
      }
    }
    
    // Second priority: common display fields
    const displayFields = ['title', 'name', 'label', 'heading', 'subject', 'slug', 'surname', 'suname']
    for (const field of displayFields) {
      if (schema[field] && entry.data[field]) {
        return String(entry.data[field])
      }
    }
    
    // Third priority: first string field in schema
    for (const [fieldName, fieldConfig] of Object.entries(schema)) {
      if ((fieldConfig.type === 'string' || fieldConfig.type === 'text') && entry.data[fieldName]) {
        return String(entry.data[fieldName])
      }
    }
  }
  
  // Fallback: try common fields
  if (entry.data.title) return String(entry.data.title)
  if (entry.data.name) return String(entry.data.name)
  if (entry.data.label) return String(entry.data.label)
  if (entry.data.surname) return String(entry.data.surname)
  if (entry.data.suname) return String(entry.data.suname)
  
  return null
}

const getDisplayFields = (data) => {
  if (!data || typeof data !== 'object') return {}
  
  // Get first few fields for display (title, name, etc. are prioritized)
  const priorityFields = ['title', 'name', 'label', 'slug', 'description']
  const display = {}
  
  // Add priority fields first (skip main display field)
  const mainDisplayField = getEntryDisplayName({ data })
  priorityFields.forEach(field => {
    if (data[field] !== undefined && String(data[field]) !== mainDisplayField) {
      display[field] = data[field]
    }
  })
  
  // Add other fields up to 3 total
  Object.keys(data).forEach(key => {
    if (!priorityFields.includes(key) && Object.keys(display).length < 3 && String(data[key]) !== mainDisplayField) {
      display[key] = data[key]
    }
  })
  
  return display
}

const formatFieldValue = (value) => {
  if (value === null || value === undefined) return '-'
  if (typeof value === 'object') return JSON.stringify(value).substring(0, 50) + '...'
  if (typeof value === 'boolean') return value ? '✓' : '✗'
  return String(value).substring(0, 100)
}

const formatDate = (dateString) => {
  if (!dateString) return '-'
  const date = new Date(dateString)
  return date.toLocaleDateString() + ' ' + date.toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' })
}

onMounted(async () => {
  await Promise.all([
    loadContentType(),
    loadAvailableContentTypes(),
    loadEntries()
  ])
})
</script>

