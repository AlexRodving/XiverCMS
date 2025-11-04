<template>
  <div>
    <div class="mb-8">
      <router-link to="/entries" class="text-blue-600 hover:text-blue-800 mb-4 inline-block">
        ← {{ $t('common.back') }} {{ $t('navigation.contentEntries') }}
      </router-link>
      <h1 class="text-3xl font-bold text-gray-900">{{ getEntryTitle() }}</h1>
    </div>

    <div v-if="loading" class="text-center py-12">
      <div class="text-gray-500">Loading...</div>
    </div>

    <div v-else-if="entry" class="bg-white rounded-lg shadow p-6">
      <div class="flex justify-between items-center mb-4">
        <h2 class="text-xl font-semibold">{{ $t('contentEntries.edit') }}</h2>
        <button
          @click="toggleEditorMode"
          class="px-3 py-1 text-sm border border-gray-300 rounded hover:bg-gray-50"
        >
          {{ useVisualEditor ? $t('contentTypes.jsonEditor') : $t('contentTypes.visualEditor') }}
        </button>
      </div>
      
      <form @submit.prevent="updateEntry">
        <div class="space-y-4">
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-1">{{ $t('contentEntries.status') }}</label>
            <select v-model="entry.status" class="w-full px-3 py-2 border border-gray-300 rounded-md">
              <option value="draft">{{ $t('contentEntries.draft') }}</option>
              <option value="published">{{ $t('contentEntries.published') }}</option>
            </select>
          </div>
          
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
              rows="15"
            ></textarea>
          </div>
        </div>
        <div class="mt-6 flex justify-end space-x-3">
          <button
            type="button"
            @click="router.push('/entries')"
            class="px-4 py-2 border border-gray-300 rounded-md text-gray-700 hover:bg-gray-50"
          >
            {{ $t('common.cancel') }}
          </button>
          <button
            type="submit"
            :disabled="!isFormValid"
            class="px-4 py-2 bg-blue-600 text-white rounded-md hover:bg-blue-700 disabled:opacity-50"
          >
            {{ $t('common.save') }}
          </button>
        </div>
      </form>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted, computed, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useI18n } from 'vue-i18n'
import { contentAPI } from '../api/content'
import DynamicForm from '../components/DynamicForm.vue'

const { t } = useI18n()
const route = useRoute()
const router = useRouter()
const contentTypeUID = computed(() => route.params.uid)
const entryId = computed(() => route.params.id)
const entry = ref(null)
const entryDataJson = ref('{}')
const entryData = ref({})
const contentTypeSchema = ref(null)
const availableContentTypes = ref([])
const useVisualEditor = ref(true)
const isFormValid = ref(true)
const loading = ref(true)
const dynamicFormRef = ref(null)

const loadContentType = async () => {
  try {
    const response = await contentAPI.getContentType(contentTypeUID.value)
    contentTypeSchema.value = response.data.schema || {}
  } catch (error) {
    console.error('Failed to load content type:', error)
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

onMounted(async () => {
  try {
    await Promise.all([
      loadContentType(),
      loadAvailableContentTypes()
    ])
    
    // Load entry with populated relations
    const response = await contentAPI.getEntry(contentTypeUID.value, entryId.value, { populate: true })
    entry.value = response.data
    
    // Process data: convert relation objects to IDs
    const processedData = { ...entry.value.data }
    if (contentTypeSchema.value) {
      Object.entries(contentTypeSchema.value).forEach(([fieldName, fieldConfig]) => {
        if (fieldConfig.type === 'relation' && processedData[fieldName]) {
          const relationValue = processedData[fieldName]
          const relationType = fieldConfig.relationType || 'manyToOne'
          if (relationType === 'oneToMany' || relationType === 'manyToMany') {
            // Convert array of entries to array of IDs
            if (Array.isArray(relationValue)) {
              processedData[fieldName] = relationValue.map(item => 
                typeof item === 'object' && item.id ? item.id : item
              )
            }
          } else {
            // Convert single entry to ID
            if (typeof relationValue === 'object' && relationValue.id) {
              processedData[fieldName] = relationValue.id
            }
          }
        }
      })
    }
    
    entryData.value = processedData
    entryDataJson.value = JSON.stringify(processedData, null, 2)
  } catch (error) {
    console.error('Failed to load entry:', error)
  } finally {
    loading.value = false
  }
})

const toggleEditorMode = () => {
  useVisualEditor.value = !useVisualEditor.value
  if (useVisualEditor.value) {
    try {
      entryData.value = JSON.parse(entryDataJson.value)
    } catch (e) {
      entryData.value = {}
    }
  } else {
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

const getEntryTitle = () => {
  if (!entry.value || !entry.value.data) {
    return `${t('contentEntries.edit')} #${entryId.value}`
  }
  
  // First, try to find required field from schema
  if (contentTypeSchema.value) {
    for (const [fieldName, fieldConfig] of Object.entries(contentTypeSchema.value)) {
      if (fieldConfig.required && entry.value.data[fieldName]) {
        return `${t('contentEntries.edit')}: ${entry.value.data[fieldName]}`
      }
    }
  }
  
  // Try to find a display field (title, name, label, etc.)
  const displayFields = ['title', 'name', 'label', 'heading', 'subject', 'surname', 'suname']
  for (const field of displayFields) {
    if (entry.value.data[field]) {
      return `${t('contentEntries.edit')}: ${entry.value.data[field]}`
    }
  }
  
  // Fallback to ID
  return `${t('contentEntries.edit')} #${entryId.value}`
}

const updateEntry = async () => {
  try {
    let data
    if (useVisualEditor.value && dynamicFormRef.value) {
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
      data = JSON.parse(entryDataJson.value)
    }
    
    await contentAPI.updateEntry(contentTypeUID.value, entryId.value, {
      data,
      status: entry.value.status,
    })
    
    if (window.showToast) {
      window.showToast.success(t('common.success'), t('contentEntries.updateSuccess'))
    }
    
    router.push('/entries')
  } catch (error) {
    const errorMsg = error.response?.data?.error || t('contentEntries.updateFailed')
    if (window.showToast) {
      window.showToast.error(t('common.error'), errorMsg)
    } else {
      alert(errorMsg)
    }
  }
}
</script>

