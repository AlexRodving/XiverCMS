<template>
  <div class="space-y-2">
    <label class="block text-sm font-medium text-gray-700">{{ label }}</label>
    
    <!-- Single selection -->
    <div v-if="relationType === 'oneToOne' || relationType === 'manyToOne'" class="space-y-2">
      <select
        v-model="selectedId"
        class="w-full px-3 py-2 border border-gray-300 rounded-md"
        @change="emitChange"
      >
        <option :value="null">-- None --</option>
        <option
          v-for="entry in entries"
          :key="entry.id"
          :value="entry.id"
        >
          {{ getEntryDisplayName(entry) }}
        </option>
      </select>
    </div>
    
    <!-- Multiple selection -->
    <div v-else class="space-y-2">
      <div class="border border-gray-300 rounded-md p-2 max-h-48 overflow-y-auto">
        <div v-if="entries.length === 0" class="text-sm text-gray-500 py-2">
          No entries available. Create some entries first.
        </div>
        <label
          v-for="entry in entries"
          :key="entry.id"
          class="flex items-center p-2 hover:bg-gray-50 rounded cursor-pointer"
        >
          <input
            type="checkbox"
            :value="entry.id"
            :checked="selectedIds.includes(entry.id)"
            @change="toggleSelection(entry.id)"
            class="mr-2"
          />
          <span class="text-sm">{{ getEntryDisplayName(entry) }}</span>
        </label>
      </div>
      <div v-if="selectedIds.length > 0" class="text-xs text-gray-500">
        {{ selectedIds.length }} selected
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, watch, onMounted } from 'vue'
import { contentAPI } from '../api/content'

const props = defineProps({
  modelValue: {
    type: [Number, Array],
    default: null
  },
  targetContentTypeUid: {
    type: String,
    required: true
  },
  relationType: {
    type: String,
    default: 'manyToOne'
  },
  label: {
    type: String,
    default: 'Select'
  }
})

const emit = defineEmits(['update:modelValue'])

const entries = ref([])
const selectedId = ref(null)
const selectedIds = ref([])

onMounted(async () => {
  await loadEntries()
  
  // Initialize selected values
  if (props.modelValue !== null && props.modelValue !== undefined) {
    if (props.relationType === 'oneToOne' || props.relationType === 'manyToOne') {
      selectedId.value = props.modelValue
    } else {
      selectedIds.value = Array.isArray(props.modelValue) ? props.modelValue : [props.modelValue]
    }
  }
})

watch(() => props.modelValue, (newValue) => {
  if (props.relationType === 'oneToOne' || props.relationType === 'manyToOne') {
    selectedId.value = newValue
  } else {
    selectedIds.value = Array.isArray(newValue) ? newValue : (newValue ? [newValue] : [])
  }
})

const loadEntries = async () => {
  try {
    const response = await contentAPI.getEntries(props.targetContentTypeUid, { status: 'published' })
    entries.value = response.data.data || []
  } catch (error) {
    console.error('Failed to load entries:', error)
  }
}

const emitChange = () => {
  emit('update:modelValue', selectedId.value)
}

const toggleSelection = (id) => {
  const index = selectedIds.value.indexOf(id)
  if (index > -1) {
    selectedIds.value.splice(index, 1)
  } else {
    selectedIds.value.push(id)
  }
  emit('update:modelValue', [...selectedIds.value])
}

const getEntryDisplayName = (entry) => {
  // Try to find a display field (title, name, etc.)
  if (entry.data) {
    if (entry.data.title) return entry.data.title
    if (entry.data.name) return entry.data.name
    if (entry.data.label) return entry.data.label
  }
  return `Entry #${entry.id}`
}
</script>

