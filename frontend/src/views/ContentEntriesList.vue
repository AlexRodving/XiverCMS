<template>
  <div>
    <div class="flex justify-between items-center mb-8">
      <h1 class="text-3xl font-bold text-gray-900">{{ $t('contentEntries.title') }}</h1>
    </div>

    <!-- Loading State -->
    <div v-if="loading" class="bg-white rounded-lg shadow p-12 text-center">
      <div class="text-gray-500">{{ $t('common.loading') }}...</div>
    </div>

    <!-- Empty State -->
    <div v-else-if="contentTypes.length === 0" class="bg-white rounded-lg shadow p-12 text-center">
      <p class="text-gray-500">{{ $t('contentTypes.noContentTypes') }}</p>
      <p class="text-sm text-gray-400 mt-2">{{ $t('contentTypes.createFirst') }}</p>
      <router-link
        to="/content-types"
        class="mt-4 inline-block px-4 py-2 bg-blue-600 text-white rounded-lg hover:bg-blue-700"
      >
        {{ $t('contentTypes.create') }}
      </router-link>
    </div>

    <!-- Content Types Grid -->
    <div v-else class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
      <div
        v-for="contentType in contentTypes"
        :key="contentType.id"
        class="bg-white rounded-lg shadow hover:shadow-lg transition-shadow p-6 flex flex-col"
      >
        <div class="flex items-start justify-between mb-4">
          <div class="flex-1">
            <h3 class="text-lg font-semibold text-gray-900 mb-1">
              {{ contentType.displayName }}
            </h3>
            <p class="text-sm text-gray-500 font-mono">{{ contentType.uid }}</p>
            <p v-if="contentType.description" class="text-sm text-gray-600 mt-2">
              {{ contentType.description }}
            </p>
          </div>
          <span
            class="px-2 py-1 text-xs font-semibold rounded-full"
            :class="contentType.kind === 'collectionType' ? 'bg-blue-100 text-blue-800' : 'bg-purple-100 text-purple-800'"
          >
            {{ contentType.kind === 'collectionType' ? $t('contentTypes.collectionType') : $t('contentTypes.singleType') }}
          </span>
        </div>

        <div class="flex items-center justify-between text-sm text-gray-500 mb-4 flex-shrink-0">
          <span>{{ $t('contentEntries.entriesCount') }}: {{ contentType.entriesCount || 0 }}</span>
          <span v-if="contentType.isVisible" class="text-green-600">{{ $t('contentTypes.visible') }}</span>
          <span v-else class="text-gray-400">{{ $t('common.no') }}</span>
        </div>

        <div class="flex gap-2 items-center mt-auto">
          <router-link
            :to="`/content-types/${contentType.uid}/entries`"
            class="flex-1 px-4 py-2.5 bg-blue-600 text-white rounded-lg hover:bg-blue-700 text-center text-sm font-medium transition-colors"
          >
            {{ $t('contentEntries.viewEntries') }}
          </router-link>
          <button
            @click="createEntry(contentType.uid)"
            class="px-4 py-2.5 bg-green-600 text-white rounded-lg hover:bg-green-700 text-sm font-medium transition-colors min-w-[44px]"
            :title="$t('contentEntries.create')"
          >
            +
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useI18n } from 'vue-i18n'
import { contentAPI } from '../api/content'

const { t } = useI18n()
const router = useRouter()
const contentTypes = ref([])
const loading = ref(true)

const loadContentTypes = async () => {
  loading.value = true
  try {
    const response = await contentAPI.getContentTypes()
    const types = response.data.data || []
    
    // Load entries count for each content type
    for (const contentType of types) {
      try {
        const entriesResponse = await contentAPI.getEntries(contentType.uid, { pageSize: 1 })
        contentType.entriesCount = entriesResponse.data.meta?.pagination?.total || 0
      } catch (error) {
        contentType.entriesCount = 0
      }
    }
    
    contentTypes.value = types
  } catch (error) {
    console.error('Failed to load content types:', error)
  } finally {
    loading.value = false
  }
}

const createEntry = (uid) => {
  router.push(`/content-types/${uid}/entries`)
}

onMounted(loadContentTypes)
</script>

