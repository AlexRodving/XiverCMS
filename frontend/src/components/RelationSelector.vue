<template>
  <div class="space-y-2">
    <!-- Single selection -->
    <div v-if="relationType === 'oneToOne' || relationType === 'manyToOne'" class="space-y-2">
      <select
        v-model="selectedId"
        class="w-full px-3 py-2 border border-gray-300 rounded-md"
        @change="emitChange"
      >
        <option :value="null">{{ $t('relationSelector.none') }}</option>
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
          {{ $t('relationSelector.noEntries') }}
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
        {{ $t('relationSelector.selectedCount', { count: selectedIds.length }) }}
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, watch, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { contentAPI } from '../api/content'

const { t } = useI18n()

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
  }
})

const emit = defineEmits(['update:modelValue'])

const entries = ref([])
const selectedId = ref(null)
const selectedIds = ref([])
const targetContentType = ref(null)

onMounted(async () => {
  await Promise.all([
    loadTargetContentType(),
    loadEntries()
  ])
  
  // Initialize selected values
  if (props.modelValue !== null && props.modelValue !== undefined) {
    if (props.relationType === 'oneToOne' || props.relationType === 'manyToOne') {
      selectedId.value = props.modelValue
    } else {
      selectedIds.value = Array.isArray(props.modelValue) ? props.modelValue : [props.modelValue]
    }
  }
})

const loadTargetContentType = async () => {
  try {
    const response = await contentAPI.getContentType(props.targetContentTypeUid)
    targetContentType.value = response.data
  } catch (error) {
    console.error('Failed to load target content type:', error)
  }
}

watch(() => props.modelValue, (newValue) => {
  if (props.relationType === 'oneToOne' || props.relationType === 'manyToOne') {
    selectedId.value = newValue
  } else {
    selectedIds.value = Array.isArray(newValue) ? newValue : (newValue ? [newValue] : [])
  }
})

const loadEntries = async () => {
  try {
    // Load all entries (both draft and published)
    const response = await contentAPI.getEntries(props.targetContentTypeUid)
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
  if (!entry.data) {
    return `Entry #${entry.id}`
  }
  
  // If we have target content type schema, use it to find display field
  if (targetContentType.value && targetContentType.value.schema) {
    const schema = targetContentType.value.schema
    
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
  if (entry.data.heading) return String(entry.data.heading)
  
  return `Entry #${entry.id}`
}
</script>

