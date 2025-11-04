<template>
  <div>
    <div class="flex justify-between items-center mb-8">
      <h1 class="text-3xl font-bold text-gray-900">{{ $t('contentTypes.title') }}</h1>
      <button
        @click="showCreateModal = true"
        class="px-4 py-2 bg-blue-600 text-white rounded-lg hover:bg-blue-700 transition-colors"
      >
        + {{ $t('contentTypes.create') }}
      </button>
    </div>

    <div class="bg-white rounded-lg shadow overflow-hidden">
      <table class="min-w-full divide-y divide-gray-200">
        <thead class="bg-gray-50">
          <tr>
            <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
              {{ $t('contentTypes.uid') }}
            </th>
            <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
              {{ $t('contentTypes.displayName') }}
            </th>
            <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
              {{ $t('contentTypes.kind') }}
            </th>
            <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
              {{ $t('contentTypes.visible') }}
            </th>
            <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
              {{ $t('contentTypes.accessType') }}
            </th>
            <th class="px-6 py-3 text-right text-xs font-medium text-gray-500 uppercase tracking-wider">
              {{ $t('common.actions') }}
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
                {{ contentType.isVisible ? $t('common.yes') : $t('common.no') }}
              </span>
            </td>
            <td class="px-6 py-4 whitespace-nowrap">
              <div class="flex flex-col gap-1">
                <span
                  class="px-2 inline-flex text-xs leading-5 font-semibold rounded-full w-fit"
                  :class="{
                    'bg-green-100 text-green-800': contentType.accessType === 'public',
                    'bg-blue-100 text-blue-800': contentType.accessType === 'authenticated',
                    'bg-yellow-100 text-yellow-800': contentType.accessType === 'moderator',
                    'bg-red-100 text-red-800': contentType.accessType === 'admin'
                  }"
                >
                  {{ contentType.accessType || 'public' }}
                </span>
                <code class="text-xs text-gray-500 font-mono" :title="$t('contentTypes.publicAPIUrl')">
                  /api/content-types/{{ contentType.uid }}/entries
                </code>
              </div>
            </td>
            <td class="px-6 py-4 whitespace-nowrap text-right text-sm font-medium">
              <router-link
                :to="`/content-types/${contentType.uid}`"
                class="text-blue-600 hover:text-blue-900 mr-4"
              >
                {{ $t('common.view') }}
              </router-link>
              <button
                @click="editContentType(contentType)"
                class="text-green-600 hover:text-green-900 mr-4"
              >
                {{ $t('common.edit') }}
              </button>
              <button
                @click="deleteContentType(contentType.uid)"
                class="text-red-600 hover:text-red-900"
              >
                {{ $t('common.delete') }}
              </button>
            </td>
          </tr>
        </tbody>
      </table>
    </div>

    <!-- Create Modal -->
    <div v-if="showCreateModal" class="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50">
      <div class="bg-white rounded-lg p-6 max-w-4xl w-full mx-4 max-h-[95vh] overflow-y-auto">
        <h2 class="text-2xl font-bold mb-4">{{ $t('contentTypes.create') }}</h2>
        <form @submit.prevent="createContentType">
          <div class="space-y-4">
            <div>
              <label class="block text-sm font-medium text-gray-700 mb-1">
                {{ $t('contentTypes.uid') }} *
                <span class="text-xs text-gray-500 font-normal ml-1">({{ $t('contentTypes.uidDescription') }})</span>
              </label>
              <input
                v-model="newContentType.uid"
                type="text"
                required
                pattern="[a-z0-9_-]+"
                class="w-full px-3 py-2 border border-gray-300 rounded-md font-mono text-sm"
                :placeholder="$t('contentTypes.uidPlaceholder')"
                @input="formatUID"
              />
              <p class="mt-1 text-xs text-gray-500">
                {{ $t('contentTypes.uidHint') }}
              </p>
            </div>
            <div>
              <label class="block text-sm font-medium text-gray-700 mb-1">{{ $t('contentTypes.displayName') }} *</label>
              <input
                v-model="newContentType.displayName"
                type="text"
                required
                class="w-full px-3 py-2 border border-gray-300 rounded-md"
                :placeholder="$t('contentTypes.displayNamePlaceholder')"
              />
            </div>
            <div>
              <label class="block text-sm font-medium text-gray-700 mb-1">{{ $t('contentTypes.description') }}</label>
              <textarea
                v-model="newContentType.description"
                class="w-full px-3 py-2 border border-gray-300 rounded-md"
                rows="3"
                :placeholder="$t('contentTypes.descriptionPlaceholder')"
              ></textarea>
            </div>
            <div>
              <label class="block text-sm font-medium text-gray-700 mb-1">{{ $t('contentTypes.kind') }}</label>
              <select v-model="newContentType.kind" class="w-full px-3 py-2 border border-gray-300 rounded-md">
                <option value="collectionType">{{ $t('contentTypes.collectionType') }}</option>
                <option value="singleType">{{ $t('contentTypes.singleType') }}</option>
              </select>
            </div>
            <div>
              <label class="flex items-center">
                <input
                  v-model="newContentType.isVisible"
                  type="checkbox"
                  class="mr-2"
                />
                <span class="text-sm font-medium text-gray-700">{{ $t('contentTypes.isVisible') }}</span>
              </label>
            </div>
            <div>
              <label class="block text-sm font-medium text-gray-700 mb-1">
                {{ $t('contentTypes.accessType') }} *
                <span class="text-xs text-gray-500 font-normal ml-1">({{ $t('contentTypes.accessTypeDescription') }})</span>
              </label>
              <select v-model="newContentType.accessType" class="w-full px-3 py-2 border border-gray-300 rounded-md" required>
                <option value="public">{{ $t('contentTypes.accessTypePublic') }}</option>
                <option value="authenticated">{{ $t('contentTypes.accessTypeAuthenticated') }}</option>
                <option value="moderator">{{ $t('contentTypes.accessTypeModerator') }}</option>
                <option value="admin">{{ $t('contentTypes.accessTypeAdmin') }}</option>
              </select>
              <p class="mt-1 text-xs text-gray-500">
                {{ $t('contentTypes.accessTypeHint') }}
              </p>
            </div>
            <div class="border-t pt-4">
              <div class="flex items-center justify-between mb-2">
                <label class="block text-sm font-medium text-gray-700">{{ $t('contentTypes.schema') }} ({{ $t('contentTypes.fields') }}) *</label>
                <div class="flex gap-2">
                  <button
                    type="button"
                    @click="showJsonEditor = !showJsonEditor"
                    class="px-2 py-1 text-xs border border-gray-300 rounded hover:bg-gray-50"
                  >
                    {{ showJsonEditor ? $t('contentTypes.visualEditor') : $t('contentTypes.jsonEditor') }}
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
                  :placeholder="$t('contentTypes.schemaPlaceholder')"
                  @input="updateSchemaFromJson"
                ></textarea>
                <p class="mt-1 text-xs text-gray-500">{{ $t('contentTypes.editSchemaJson') }}</p>
              </div>
            </div>
          </div>
          <div class="mt-6 flex justify-end space-x-3">
            <button
              type="button"
              @click="showCreateModal = false"
              class="px-4 py-2 border border-gray-300 rounded-md text-gray-700 hover:bg-gray-50"
            >
              {{ $t('common.cancel') }}
            </button>
            <button
              type="submit"
              class="px-4 py-2 bg-blue-600 text-white rounded-md hover:bg-blue-700"
            >
              {{ $t('common.create') }}
            </button>
          </div>
        </form>
      </div>
    </div>

    <!-- Edit Modal -->
    <div v-if="showEditModal && editingContentType" class="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50">
      <div class="bg-white rounded-lg p-6 max-w-4xl w-full mx-4 max-h-[95vh] overflow-y-auto">
        <h2 class="text-2xl font-bold mb-4">{{ $t('contentTypes.edit') }}</h2>
        <form @submit.prevent="updateContentType">
          <div class="space-y-4">
            <div>
              <label class="block text-sm font-medium text-gray-700 mb-1">
                {{ $t('contentTypes.uid') }}
                <span class="text-xs text-gray-500 font-normal ml-1">({{ $t('contentTypes.uidReadonly') }})</span>
              </label>
              <input
                :value="editingContentType.uid"
                type="text"
                disabled
                class="w-full px-3 py-2 border border-gray-300 rounded-md font-mono text-sm bg-gray-100 text-gray-500"
              />
            </div>
            <div>
              <label class="block text-sm font-medium text-gray-700 mb-1">{{ $t('contentTypes.displayName') }} *</label>
              <input
                v-model="editingContentType.displayName"
                type="text"
                required
                class="w-full px-3 py-2 border border-gray-300 rounded-md"
                :placeholder="$t('contentTypes.displayNamePlaceholder')"
              />
            </div>
            <div>
              <label class="block text-sm font-medium text-gray-700 mb-1">{{ $t('contentTypes.description') }}</label>
              <textarea
                v-model="editingContentType.description"
                class="w-full px-3 py-2 border border-gray-300 rounded-md"
                rows="3"
                :placeholder="$t('contentTypes.descriptionPlaceholder')"
              ></textarea>
            </div>
            <div>
              <label class="block text-sm font-medium text-gray-700 mb-1">{{ $t('contentTypes.kind') }}</label>
              <select v-model="editingContentType.kind" class="w-full px-3 py-2 border border-gray-300 rounded-md">
                <option value="collectionType">{{ $t('contentTypes.collectionType') }}</option>
                <option value="singleType">{{ $t('contentTypes.singleType') }}</option>
              </select>
            </div>
            <div>
              <label class="flex items-center">
                <input
                  v-model="editingContentType.isVisible"
                  type="checkbox"
                  class="mr-2"
                />
                <span class="text-sm font-medium text-gray-700">{{ $t('contentTypes.isVisible') }}</span>
              </label>
            </div>
            <div>
              <label class="block text-sm font-medium text-gray-700 mb-1">
                {{ $t('contentTypes.accessType') }} *
                <span class="text-xs text-gray-500 font-normal ml-1">({{ $t('contentTypes.accessTypeDescription') }})</span>
              </label>
              <select v-model="editingContentType.accessType" class="w-full px-3 py-2 border border-gray-300 rounded-md" required>
                <option value="public">{{ $t('contentTypes.accessTypePublic') }}</option>
                <option value="authenticated">{{ $t('contentTypes.accessTypeAuthenticated') }}</option>
                <option value="moderator">{{ $t('contentTypes.accessTypeModerator') }}</option>
                <option value="admin">{{ $t('contentTypes.accessTypeAdmin') }}</option>
              </select>
              <p class="mt-1 text-xs text-gray-500">
                {{ $t('contentTypes.accessTypeHint') }}
              </p>
            </div>
            <div class="border-t pt-4">
              <div class="flex items-center justify-between mb-2">
                <label class="block text-sm font-medium text-gray-700">{{ $t('contentTypes.schema') }} ({{ $t('contentTypes.fields') }}) *</label>
                <div class="flex gap-2">
                  <button
                    type="button"
                    @click="showJsonEditor = !showJsonEditor"
                    class="px-2 py-1 text-xs border border-gray-300 rounded hover:bg-gray-50"
                  >
                    {{ showJsonEditor ? $t('contentTypes.visualEditor') : $t('contentTypes.jsonEditor') }}
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
                  :placeholder="$t('contentTypes.schemaPlaceholder')"
                  @input="updateSchemaFromJson"
                ></textarea>
                <p class="mt-1 text-xs text-gray-500">{{ $t('contentTypes.editSchemaJson') }}</p>
              </div>
            </div>
          </div>
          <div class="mt-6 flex justify-end space-x-3">
            <button
              type="button"
              @click="showEditModal = false; editingContentType = null"
              class="px-4 py-2 border border-gray-300 rounded-md text-gray-700 hover:bg-gray-50"
            >
              {{ $t('common.cancel') }}
            </button>
            <button
              type="submit"
              class="px-4 py-2 bg-blue-600 text-white rounded-md hover:bg-blue-700"
            >
              {{ $t('common.save') }}
            </button>
          </div>
        </form>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted, watch } from 'vue'
import { useI18n } from 'vue-i18n'
import { contentAPI } from '../api/content'
import FieldBuilder from '../components/FieldBuilder.vue'

const { t } = useI18n()

const contentTypes = ref([])
const showCreateModal = ref(false)
const showEditModal = ref(false)
const editingContentType = ref(null)
const showJsonEditor = ref(false)
const newContentType = ref({
  uid: '',
  displayName: '',
  description: '',
  kind: 'collectionType',
  isVisible: true,
  accessType: 'public',
  schema: {},
})
const schema = ref({})
const schemaJson = ref('{}')

const formatUID = (event) => {
  // Convert to lowercase and replace spaces with hyphens
  let value = event.target.value.toLowerCase().trim()
  // Remove invalid characters (keep only lowercase letters, numbers, hyphens, underscores)
  value = value.replace(/[^a-z0-9_-]/g, '-')
  // Replace multiple hyphens with single hyphen
  value = value.replace(/-+/g, '-')
  // Remove leading/trailing hyphens
  value = value.replace(/^-+|-+$/g, '')
  newContentType.value.uid = value
}

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
      if (window.showToast) {
        window.showToast.warning(t('common.error'), t('contentTypes.addFieldFirst'))
      } else {
        alert(t('contentTypes.addFieldFirst'))
      }
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
      accessType: 'public',
      schema: {},
    })
    await loadContentTypes()
    if (window.showToast) {
      window.showToast.success(t('common.success'), t('contentTypes.createSuccess'))
    }
  } catch (error) {
    const errorMsg = error.response?.data?.error || t('contentTypes.createFailed')
    if (window.showToast) {
      window.showToast.error(t('common.error'), errorMsg)
    } else {
      alert(errorMsg)
    }
  }
}

const editContentType = async (contentType) => {
  editingContentType.value = { ...contentType }
  schema.value = contentType.schema || {}
  schemaJson.value = JSON.stringify(contentType.schema || {}, null, 2)
  showJsonEditor.value = false
  showEditModal.value = true
}

const updateContentType = async () => {
  try {
    // Use schema from visual editor or JSON
    const finalSchema = Object.keys(schema.value).length > 0 ? schema.value : JSON.parse(schemaJson.value)
    
    if (Object.keys(finalSchema).length === 0) {
      if (window.showToast) {
        window.showToast.warning(t('common.error'), t('contentTypes.addFieldFirst'))
      } else {
        alert(t('contentTypes.addFieldFirst'))
      }
      return
    }
    
    await contentAPI.updateContentType(editingContentType.value.uid, {
      displayName: editingContentType.value.displayName,
      description: editingContentType.value.description,
      kind: editingContentType.value.kind,
      isVisible: editingContentType.value.isVisible,
      accessType: editingContentType.value.accessType || 'public',
      schema: finalSchema,
    })
    
    showEditModal.value = false
    editingContentType.value = null
    await loadContentTypes()
    
    if (window.showToast) {
      window.showToast.success(t('common.success'), t('contentTypes.updateSuccess'))
    }
  } catch (error) {
    const errorMsg = error.response?.data?.error || t('contentTypes.updateFailed')
    if (window.showToast) {
      window.showToast.error(t('common.error'), errorMsg)
    } else {
      alert(errorMsg)
    }
  }
}

const deleteContentType = async (uid) => {
  if (!confirm(t('contentTypes.confirmDelete').replace('{uid}', uid))) {
    return
  }
  try {
    await contentAPI.deleteContentType(uid)
    await loadContentTypes()
    if (window.showToast) {
      window.showToast.success(t('common.success'), t('contentTypes.deleteSuccess'))
    }
  } catch (error) {
    const errorMsg = error.response?.data?.error || t('contentTypes.deleteFailed')
    if (window.showToast) {
      window.showToast.error(t('common.error'), errorMsg)
    } else {
      alert(errorMsg)
    }
  }
}

onMounted(loadContentTypes)
</script>

