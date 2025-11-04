<template>
  <div>
    <div class="mb-8">
      <router-link :to="`/content-types/${contentTypeUID}/entries`" class="text-blue-600 hover:text-blue-800 mb-4 inline-block">
        ← Back to Entries
      </router-link>
      <h1 class="text-3xl font-bold text-gray-900">Edit Entry #{{ entryId }}</h1>
    </div>

    <div v-if="loading" class="text-center py-12">
      <div class="text-gray-500">Loading...</div>
    </div>

    <div v-else-if="entry" class="bg-white rounded-lg shadow p-6">
      <form @submit.prevent="updateEntry">
        <div class="space-y-4">
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-1">Status</label>
            <select v-model="entry.status" class="w-full px-3 py-2 border border-gray-300 rounded-md">
              <option value="draft">Draft</option>
              <option value="published">Published</option>
            </select>
          </div>
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-1">Data (JSON) *</label>
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
            @click="$router.back()"
            class="px-4 py-2 border border-gray-300 rounded-md text-gray-700 hover:bg-gray-50"
          >
            Cancel
          </button>
          <button
            type="submit"
            class="px-4 py-2 bg-blue-600 text-white rounded-md hover:bg-blue-700"
          >
            Save
          </button>
        </div>
      </form>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted, computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { contentAPI } from '../api/content'

const route = useRoute()
const router = useRouter()
const contentTypeUID = computed(() => route.params.uid)
const entryId = computed(() => route.params.id)
const entry = ref(null)
const entryDataJson = ref('{}')
const loading = ref(true)

onMounted(async () => {
  try {
    const response = await contentAPI.getEntry(contentTypeUID.value, entryId.value)
    entry.value = response.data
    entryDataJson.value = JSON.stringify(entry.value.data, null, 2)
  } catch (error) {
    console.error('Failed to load entry:', error)
  } finally {
    loading.value = false
  }
})

const updateEntry = async () => {
  try {
    const data = JSON.parse(entryDataJson.value)
    await contentAPI.updateEntry(contentTypeUID.value, entryId.value, {
      data,
      status: entry.value.status,
    })
    router.push(`/content-types/${contentTypeUID.value}/entries`)
  } catch (error) {
    alert(error.response?.data?.error || 'Failed to update entry')
  }
}
</script>

