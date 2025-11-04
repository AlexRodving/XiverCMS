<template>
  <div>
    <div class="flex justify-between items-center mb-8">
      <h1 class="text-3xl font-bold text-gray-900">{{ $t('users.title') }}</h1>
    </div>

    <!-- Loading State -->
    <div v-if="loading" class="bg-white rounded-lg shadow p-12 text-center">
      <div class="text-gray-500">{{ $t('common.loading') }}...</div>
    </div>

    <!-- Users List -->
    <div v-else-if="users.length > 0" class="bg-white rounded-lg shadow overflow-hidden">
      <table class="min-w-full divide-y divide-gray-200">
        <thead class="bg-gray-50">
          <tr>
            <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
              {{ $t('users.email') }}
            </th>
            <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
              {{ $t('users.username') }}
            </th>
            <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
              {{ $t('users.name') }}
            </th>
            <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
              {{ $t('users.status') }}
            </th>
            <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
              {{ $t('users.roles') }}
            </th>
            <th class="px-6 py-3 text-right text-xs font-medium text-gray-500 uppercase tracking-wider">
              {{ $t('common.actions') }}
            </th>
          </tr>
        </thead>
        <tbody class="bg-white divide-y divide-gray-200">
          <tr v-for="user in users" :key="user.id" class="hover:bg-gray-50">
            <td class="px-6 py-4 whitespace-nowrap text-sm font-medium text-gray-900">
              {{ user.email }}
            </td>
            <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
              {{ user.username }}
            </td>
            <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
              {{ user.firstName }} {{ user.lastName }}
            </td>
            <td class="px-6 py-4 whitespace-nowrap">
              <span
                class="px-2 inline-flex text-xs leading-5 font-semibold rounded-full"
                :class="user.isActive ? 'bg-green-100 text-green-800' : 'bg-red-100 text-red-800'"
              >
                {{ user.isActive ? $t('users.active') : $t('users.inactive') }}
              </span>
              <span v-if="user.isSuperAdmin" class="ml-2 px-2 py-1 text-xs font-semibold rounded-full bg-purple-100 text-purple-800">
                {{ $t('users.superAdmin') }}
              </span>
            </td>
            <td class="px-6 py-4 text-sm text-gray-500">
              <span v-if="user.roles && user.roles.length > 0">
                <span v-for="(role, index) in user.roles" :key="role.id" class="mr-2">
                  {{ role.name }}<span v-if="index < user.roles.length - 1">,</span>
                </span>
              </span>
              <span v-else class="text-gray-400">-</span>
            </td>
            <td class="px-6 py-4 whitespace-nowrap text-right text-sm font-medium space-x-2">
              <button
                @click="editUser(user)"
                class="text-blue-600 hover:text-blue-900"
              >
                {{ $t('common.edit') }}
              </button>
              <button
                @click="deleteUser(user.id)"
                class="text-red-600 hover:text-red-900"
                :disabled="user.isSuperAdmin"
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
      <p class="text-gray-500">{{ $t('users.noUsers') }}</p>
    </div>

    <!-- Edit Modal -->
    <div v-if="showEditModal && editingUser" class="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50">
      <div class="bg-white rounded-lg p-6 max-w-2xl w-full mx-4 max-h-[90vh] overflow-y-auto">
        <h2 class="text-2xl font-bold mb-4">{{ $t('users.edit') }}</h2>
        
        <form @submit.prevent="saveUser">
          <div class="space-y-4">
            <div>
              <label class="block text-sm font-medium text-gray-700 mb-1">
                {{ $t('users.email') }} *
              </label>
              <input
                v-model="userForm.email"
                type="email"
                required
                class="w-full px-3 py-2 border border-gray-300 rounded-md"
              />
            </div>

            <div>
              <label class="block text-sm font-medium text-gray-700 mb-1">
                {{ $t('users.username') }} *
              </label>
              <input
                v-model="userForm.username"
                type="text"
                required
                class="w-full px-3 py-2 border border-gray-300 rounded-md"
              />
            </div>

            <div class="grid grid-cols-2 gap-4">
              <div>
                <label class="block text-sm font-medium text-gray-700 mb-1">
                  {{ $t('users.firstName') }}
                </label>
                <input
                  v-model="userForm.firstName"
                  type="text"
                  class="w-full px-3 py-2 border border-gray-300 rounded-md"
                />
              </div>
              <div>
                <label class="block text-sm font-medium text-gray-700 mb-1">
                  {{ $t('users.lastName') }}
                </label>
                <input
                  v-model="userForm.lastName"
                  type="text"
                  class="w-full px-3 py-2 border border-gray-300 rounded-md"
                />
              </div>
            </div>

            <div>
              <label class="flex items-center">
                <input
                  v-model="userForm.isActive"
                  type="checkbox"
                  class="mr-2"
                />
                <span class="text-sm font-medium text-gray-700">{{ $t('users.isActive') }}</span>
              </label>
            </div>

            <div>
              <label class="block text-sm font-medium text-gray-700 mb-2">
                {{ $t('users.roles') }}
              </label>
              <div class="border border-gray-300 rounded-md p-3 max-h-60 overflow-y-auto">
                <div v-if="availableRoles.length === 0" class="text-sm text-gray-500">
                  {{ $t('roles.noRoles') }}
                </div>
                <label
                  v-for="role in availableRoles"
                  :key="role.id"
                  class="flex items-center p-2 hover:bg-gray-50 rounded cursor-pointer"
                >
                  <input
                    type="checkbox"
                    :value="role.id"
                    v-model="userForm.roleIds"
                    class="mr-2"
                  />
                  <div class="flex-1">
                    <div class="text-sm font-medium text-gray-900">
                      {{ role.name }}
                    </div>
                    <div v-if="role.description" class="text-xs text-gray-500">
                      {{ role.description }}
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
import { usersAPI } from '../api/users'
import { rolesAPI } from '../api/roles'

const { t } = useI18n()
const users = ref([])
const availableRoles = ref([])
const loading = ref(true)
const showEditModal = ref(false)
const editingUser = ref(null)

const userForm = ref({
  email: '',
  username: '',
  firstName: '',
  lastName: '',
  isActive: true,
  roleIds: []
})

const loadUsers = async () => {
  loading.value = true
  try {
    const response = await usersAPI.getAll()
    users.value = response.data.data || []
  } catch (error) {
    console.error('Failed to load users:', error)
  } finally {
    loading.value = false
  }
}

const loadRoles = async () => {
  try {
    const response = await rolesAPI.getAll()
    availableRoles.value = response.data.data || []
  } catch (error) {
    console.error('Failed to load roles:', error)
  }
}

const editUser = (user) => {
  editingUser.value = user
  userForm.value = {
    email: user.email,
    username: user.username,
    firstName: user.firstName || '',
    lastName: user.lastName || '',
    isActive: user.isActive,
    roleIds: user.roles ? user.roles.map(r => r.id) : []
  }
  showEditModal.value = true
}

const closeModal = () => {
  showEditModal.value = false
  editingUser.value = null
  userForm.value = {
    email: '',
    username: '',
    firstName: '',
    lastName: '',
    isActive: true,
    roleIds: []
  }
}

const saveUser = async () => {
  try {
    await usersAPI.update(editingUser.value.id, userForm.value)
    if (window.showToast) {
      window.showToast.success(t('common.success'), t('users.updateSuccess'))
    }
    closeModal()
    await loadUsers()
  } catch (error) {
    const errorMsg = error.response?.data?.error || t('users.updateFailed')
    if (window.showToast) {
      window.showToast.error(t('common.error'), errorMsg)
    } else {
      alert(errorMsg)
    }
  }
}

const deleteUser = async (id) => {
  if (!confirm(t('users.confirmDelete'))) {
    return
  }
  try {
    await usersAPI.delete(id)
    await loadUsers()
    if (window.showToast) {
      window.showToast.success(t('common.success'), t('users.deleteSuccess'))
    }
  } catch (error) {
    const errorMsg = error.response?.data?.error || t('users.deleteFailed')
    if (window.showToast) {
      window.showToast.error(t('common.error'), errorMsg)
    } else {
      alert(errorMsg)
    }
  }
}

onMounted(async () => {
  await Promise.all([
    loadUsers(),
    loadRoles()
  ])
})
</script>
