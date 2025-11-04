<template>
  <div>
    <div class="flex justify-between items-center mb-8">
      <h1 class="text-3xl font-bold text-gray-900">{{ $t('permissions.title') }}</h1>
      <button
        @click="showCreateModal = true"
        class="px-4 py-2 bg-blue-600 text-white rounded-lg hover:bg-blue-700 transition-colors"
      >
        + {{ $t('permissions.create') }}
      </button>
    </div>

    <!-- Loading State -->
    <div v-if="loading" class="bg-white rounded-lg shadow p-12 text-center">
      <div class="text-gray-500">{{ $t('common.loading') }}...</div>
    </div>

    <!-- Permissions List -->
    <div v-else-if="permissions.length > 0" class="bg-white rounded-lg shadow overflow-hidden">
      <table class="min-w-full divide-y divide-gray-200">
        <thead class="bg-gray-50">
          <tr>
            <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
              {{ $t('permissions.action') }}
            </th>
            <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
              {{ $t('permissions.subject') }}
            </th>
            <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
              {{ $t('permissions.roles') }}
            </th>
            <th class="px-6 py-3 text-right text-xs font-medium text-gray-500 uppercase tracking-wider">
              {{ $t('common.actions') }}
            </th>
          </tr>
        </thead>
        <tbody class="bg-white divide-y divide-gray-200">
          <tr v-for="permission in permissions" :key="permission.id" class="hover:bg-gray-50">
            <td class="px-6 py-4 whitespace-nowrap">
              <span class="px-2 py-1 text-xs font-semibold rounded-full bg-blue-100 text-blue-800">
                {{ permission.action }}
              </span>
            </td>
            <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-900">
              {{ permission.subject }}
            </td>
            <td class="px-6 py-4 text-sm text-gray-500">
              <span v-if="permission.roles && permission.roles.length > 0">
                {{ permission.roles.map(r => r.name).join(', ') }}
              </span>
              <span v-else class="text-gray-400">-</span>
            </td>
            <td class="px-6 py-4 whitespace-nowrap text-right text-sm font-medium space-x-2">
              <button
                @click="editPermission(permission)"
                class="text-blue-600 hover:text-blue-900"
              >
                {{ $t('common.edit') }}
              </button>
              <button
                @click="deletePermission(permission.id)"
                class="text-red-600 hover:text-red-900"
              >
                {{ $t('common.delete') }}
              </button>
            </td>
          </tr>
        </tbody>
      </table>
    </div>

    <!-- Empty State -->
    <div v-else class="bg-white rounded-lg shadow p-12 text-center">
      <p class="text-gray-500">{{ $t('permissions.noPermissions') }}</p>
    </div>

    <!-- Create/Edit Modal -->
    <div v-if="showCreateModal || editingPermission" class="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50">
      <div class="bg-white rounded-lg p-6 max-w-xl w-full mx-4">
        <h2 class="text-2xl font-bold mb-4">
          {{ editingPermission ? $t('permissions.edit') : $t('permissions.create') }}
        </h2>
        
        <form @submit.prevent="savePermission">
          <div class="space-y-4">
            <div>
              <label class="block text-sm font-medium text-gray-700 mb-1">
                {{ $t('permissions.action') }} *
              </label>
              <select
                v-model="permissionForm.action"
                required
                class="w-full px-3 py-2 border border-gray-300 rounded-md"
              >
                <option value="">{{ $t('common.select') }}</option>
                <option value="create">create</option>
                <option value="read">read</option>
                <option value="update">update</option>
                <option value="delete">delete</option>
                <option value="publish">publish</option>
                <option value="all">all</option>
              </select>
            </div>

            <div>
              <label class="block text-sm font-medium text-gray-700 mb-1">
                {{ $t('permissions.subject') }} *
              </label>
              <input
                v-model="permissionForm.subject"
                type="text"
                required
                placeholder="e.g. content-type, user, all, content-type:article"
                class="w-full px-3 py-2 border border-gray-300 rounded-md"
              />
              <p class="mt-1 text-xs text-gray-500">
                {{ $t('permissions.subjectHint') }}
              </p>
            </div>

            <div>
              <label class="block text-sm font-medium text-gray-700 mb-1">
                {{ $t('permissions.properties') }} (JSON)
              </label>
              <textarea
                v-model="permissionForm.propertiesJson"
                rows="3"
                class="w-full px-3 py-2 border border-gray-300 rounded-md font-mono text-sm"
                placeholder='{"field": "value"}'
              ></textarea>
            </div>

            <div>
              <label class="block text-sm font-medium text-gray-700 mb-1">
                {{ $t('permissions.conditions') }} (JSON)
              </label>
              <textarea
                v-model="permissionForm.conditionsJson"
                rows="3"
                class="w-full px-3 py-2 border border-gray-300 rounded-md font-mono text-sm"
                placeholder='{"condition": "value"}'
              ></textarea>
            </div>
          </div>

          <div class="mt-6 flex justify-end space-x-3">
            <button
              type="button"
              @click="closeModal"
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
import { ref, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { permissionsAPI } from '../api/permissions'

const { t } = useI18n()

const permissions = ref([])
const loading = ref(true)
const showCreateModal = ref(false)
const editingPermission = ref(null)

const permissionForm = ref({
  action: '',
  subject: '',
  propertiesJson: '',
  conditionsJson: ''
})

const loadPermissions = async () => {
  loading.value = true
  try {
    const response = await permissionsAPI.getAll()
    permissions.value = response.data || []
  } catch (error) {
    console.error('Failed to load permissions:', error)
    if (window.showToast) {
      window.showToast.error(t('common.error'), t('permissions.loadFailed'))
    }
  } finally {
    loading.value = false
  }
}

const editPermission = (permission) => {
  editingPermission.value = permission
  permissionForm.value = {
    action: permission.action,
    subject: permission.subject,
    propertiesJson: permission.properties ? (typeof permission.properties === 'string' ? JSON.stringify(JSON.parse(permission.properties), null, 2) : JSON.stringify(permission.properties, null, 2)) : '',
    conditionsJson: permission.conditions ? (typeof permission.conditions === 'string' ? JSON.stringify(JSON.parse(permission.conditions), null, 2) : JSON.stringify(permission.conditions, null, 2)) : ''
  }
  showCreateModal.value = true
}

const closeModal = () => {
  showCreateModal.value = false
  editingPermission.value = null
  permissionForm.value = {
    action: '',
    subject: '',
    propertiesJson: '',
    conditionsJson: ''
  }
}

const savePermission = async () => {
  try {
    const data = {
      action: permissionForm.value.action,
      subject: permissionForm.value.subject,
    }

    if (permissionForm.value.propertiesJson) {
      try {
        data.properties = JSON.parse(permissionForm.value.propertiesJson)
      } catch (e) {
        if (window.showToast) {
          window.showToast.error(t('common.error'), t('permissions.invalidPropertiesJson'))
        }
        return
      }
    }

    if (permissionForm.value.conditionsJson) {
      try {
        data.conditions = JSON.parse(permissionForm.value.conditionsJson)
      } catch (e) {
        if (window.showToast) {
          window.showToast.error(t('common.error'), t('permissions.invalidConditionsJson'))
        }
        return
      }
    }

    if (editingPermission.value) {
      await permissionsAPI.update(editingPermission.value.id, data)
      if (window.showToast) {
        window.showToast.success(t('common.success'), t('permissions.updateSuccess'))
      }
    } else {
      await permissionsAPI.create(data)
      if (window.showToast) {
        window.showToast.success(t('common.success'), t('permissions.createSuccess'))
      }
    }
    closeModal()
    await loadPermissions()
  } catch (error) {
    const errorMsg = error.response?.data?.error || t('permissions.saveFailed')
    if (window.showToast) {
      window.showToast.error(t('common.error'), errorMsg)
    } else {
      alert(errorMsg)
    }
  }
}

const deletePermission = async (id) => {
  if (!confirm(t('permissions.confirmDelete'))) {
    return
  }
  try {
    await permissionsAPI.delete(id)
    if (window.showToast) {
      window.showToast.success(t('common.success'), t('permissions.deleteSuccess'))
    }
    await loadPermissions()
  } catch (error) {
    const errorMsg = error.response?.data?.error || t('permissions.deleteFailed')
    if (window.showToast) {
      window.showToast.error(t('common.error'), errorMsg)
    } else {
      alert(errorMsg)
    }
  }
}

onMounted(loadPermissions)
</script>
