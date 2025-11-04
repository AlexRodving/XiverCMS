<template>
  <div>
    <div class="flex justify-between items-center mb-8">
      <h1 class="text-3xl font-bold text-gray-900">{{ $t('roles.title') }}</h1>
      <button
        @click="showCreateModal = true"
        class="px-4 py-2 bg-blue-600 text-white rounded-lg hover:bg-blue-700 transition-colors"
      >
        + {{ $t('roles.create') }}
      </button>
    </div>

    <!-- Loading State -->
    <div v-if="loading" class="bg-white rounded-lg shadow p-12 text-center">
      <div class="text-gray-500">{{ $t('common.loading') }}...</div>
    </div>

    <!-- Roles List -->
    <div v-else-if="roles.length > 0" class="bg-white rounded-lg shadow overflow-hidden">
      <table class="min-w-full divide-y divide-gray-200">
        <thead class="bg-gray-50">
          <tr>
            <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
              {{ $t('roles.name') }}
            </th>
            <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
              {{ $t('roles.type') }}
            </th>
            <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
              {{ $t('roles.description') }}
            </th>
            <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
              {{ $t('roles.permissions') }}
            </th>
            <th class="px-6 py-3 text-right text-xs font-medium text-gray-500 uppercase tracking-wider">
              {{ $t('common.actions') }}
            </th>
          </tr>
        </thead>
        <tbody class="bg-white divide-y divide-gray-200">
          <tr v-for="role in roles" :key="role.id" class="hover:bg-gray-50">
            <td class="px-6 py-4 whitespace-nowrap">
              <div class="text-sm font-medium text-gray-900">{{ role.name }}</div>
            </td>
            <td class="px-6 py-4 whitespace-nowrap">
              <span
                class="px-2 inline-flex text-xs leading-5 font-semibold rounded-full"
                :class="role.type === 'public' ? 'bg-green-100 text-green-800' : 'bg-blue-100 text-blue-800'"
              >
                {{ role.type === 'public' ? $t('roles.public') : $t('roles.custom') }}
              </span>
            </td>
            <td class="px-6 py-4 text-sm text-gray-500">
              {{ role.description || '-' }}
            </td>
            <td class="px-6 py-4 text-sm text-gray-500">
              <span v-if="role.permissions && role.permissions.length > 0">
                {{ role.permissions.length }} {{ $t('roles.permissions') }}
              </span>
              <span v-else class="text-gray-400">-</span>
            </td>
            <td class="px-6 py-4 whitespace-nowrap text-right text-sm font-medium space-x-2">
              <button
                @click="editRole(role)"
                class="text-blue-600 hover:text-blue-900"
              >
                {{ $t('common.edit') }}
              </button>
              <button
                @click="deleteRole(role.id)"
                class="text-red-600 hover:text-red-900"
                :disabled="role.type === 'public' && ['Public', 'Authenticated'].includes(role.name)"
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
      <p class="text-gray-500">{{ $t('roles.noRoles') }}</p>
    </div>

    <!-- Create/Edit Modal -->
    <div v-if="showCreateModal || editingRole" class="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50">
      <div class="bg-white rounded-lg p-6 max-w-2xl w-full mx-4 max-h-[90vh] overflow-y-auto">
        <h2 class="text-2xl font-bold mb-4">
          {{ editingRole ? $t('roles.edit') : $t('roles.create') }}
        </h2>
        
        <form @submit.prevent="saveRole">
          <div class="space-y-4">
            <div>
              <label class="block text-sm font-medium text-gray-700 mb-1">
                {{ $t('roles.name') }} *
              </label>
              <input
                v-model="roleForm.name"
                type="text"
                required
                class="w-full px-3 py-2 border border-gray-300 rounded-md"
                :disabled="editingRole && editingRole.type === 'public' && ['Public', 'Authenticated'].includes(editingRole.name)"
              />
            </div>

            <div>
              <label class="block text-sm font-medium text-gray-700 mb-1">
                {{ $t('roles.description') }}
              </label>
              <textarea
                v-model="roleForm.description"
                rows="3"
                class="w-full px-3 py-2 border border-gray-300 rounded-md"
              ></textarea>
            </div>

            <div>
              <label class="block text-sm font-medium text-gray-700 mb-1">
                {{ $t('roles.type') }} *
              </label>
              <select
                v-model="roleForm.type"
                required
                class="w-full px-3 py-2 border border-gray-300 rounded-md"
                :disabled="editingRole && editingRole.type === 'public'"
              >
                <option value="custom">{{ $t('roles.custom') }}</option>
                <option value="public">{{ $t('roles.public') }}</option>
              </select>
            </div>

            <div>
              <label class="block text-sm font-medium text-gray-700 mb-2">
                {{ $t('roles.permissions') }}
              </label>
              <div class="border border-gray-300 rounded-md p-3 max-h-60 overflow-y-auto">
                <div v-if="permissions.length === 0" class="text-sm text-gray-500">
                  {{ $t('permissions.noPermissions') }}
                </div>
                <label
                  v-for="permission in permissions"
                  :key="permission.id"
                  class="flex items-center p-2 hover:bg-gray-50 rounded cursor-pointer"
                >
                  <input
                    type="checkbox"
                    :value="permission.id"
                    v-model="roleForm.permissionIds"
                    class="mr-2"
                  />
                  <div class="flex-1">
                    <div class="text-sm font-medium text-gray-900">
                      {{ permission.action }} - {{ permission.subject }}
                    </div>
                  </div>
                </label>
              </div>
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
import { rolesAPI } from '../api/roles'
import { permissionsAPI } from '../api/permissions'

const { t } = useI18n()

const roles = ref([])
const permissions = ref([])
const loading = ref(true)
const showCreateModal = ref(false)
const editingRole = ref(null)

const roleForm = ref({
  name: '',
  description: '',
  type: 'custom',
  permissionIds: []
})

const loadRoles = async () => {
  loading.value = true
  try {
    const response = await rolesAPI.getAll()
    roles.value = response.data.data || []
  } catch (error) {
    console.error('Failed to load roles:', error)
    if (window.showToast) {
      window.showToast.error(t('common.error'), t('roles.loadFailed'))
    }
  } finally {
    loading.value = false
  }
}

const loadPermissions = async () => {
  try {
    const response = await permissionsAPI.getAll()
    permissions.value = response.data || []
  } catch (error) {
    console.error('Failed to load permissions:', error)
  }
}

const editRole = (role) => {
  editingRole.value = role
  roleForm.value = {
    name: role.name,
    description: role.description || '',
    type: role.type || 'custom',
    permissionIds: role.permissions ? role.permissions.map(p => p.id) : []
  }
  showCreateModal.value = true
}

const closeModal = () => {
  showCreateModal.value = false
  editingRole.value = null
  roleForm.value = {
    name: '',
    description: '',
    type: 'custom',
    permissionIds: []
  }
}

const saveRole = async () => {
  try {
    if (editingRole.value) {
      await rolesAPI.update(editingRole.value.id, roleForm.value)
      if (window.showToast) {
        window.showToast.success(t('common.success'), t('roles.updateSuccess'))
      }
    } else {
      await rolesAPI.create(roleForm.value)
      if (window.showToast) {
        window.showToast.success(t('common.success'), t('roles.createSuccess'))
      }
    }
    closeModal()
    await loadRoles()
  } catch (error) {
    const errorMsg = error.response?.data?.error || t('roles.saveFailed')
    if (window.showToast) {
      window.showToast.error(t('common.error'), errorMsg)
    } else {
      alert(errorMsg)
    }
  }
}

const deleteRole = async (id) => {
  if (!confirm(t('roles.confirmDelete'))) {
    return
  }
  try {
    await rolesAPI.delete(id)
    if (window.showToast) {
      window.showToast.success(t('common.success'), t('roles.deleteSuccess'))
    }
    await loadRoles()
  } catch (error) {
    const errorMsg = error.response?.data?.error || t('roles.deleteFailed')
    if (window.showToast) {
      window.showToast.error(t('common.error'), errorMsg)
    } else {
      alert(errorMsg)
    }
  }
}

onMounted(async () => {
  await Promise.all([
    loadRoles(),
    loadPermissions()
  ])
})
</script>
