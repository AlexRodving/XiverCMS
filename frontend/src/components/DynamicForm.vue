<template>
  <div class="space-y-6">
    <div
      v-for="fieldName in orderedFieldNames"
      :key="fieldName"
      class="space-y-2"
    >
      <template v-if="schema[fieldName]">
        {{ '' }}
      </template>
      <template v-else>
        {{ '' }}
      </template>
      <!-- Field Label -->
      <label
        :for="fieldName"
        class="block text-sm font-medium text-gray-700"
        v-if="schema[fieldName]"
      >
        {{ getFieldLabel(fieldName, schema[fieldName]) }}
        <span v-if="schema[fieldName].required" class="text-red-500">*</span>
        <span v-if="schema[fieldName].description" class="text-xs text-gray-500 font-normal ml-2">
          ({{ schema[fieldName].description }})
        </span>
      </label>

      <!-- String Field -->
      <input
        v-if="schema[fieldName] && (schema[fieldName].type === 'string' || schema[fieldName].type === 'email' || schema[fieldName].type === 'url')"
        :id="fieldName"
        v-model="formData[fieldName]"
        :type="getInputType(schema[fieldName].type)"
        :required="schema[fieldName].required"
        :placeholder="schema[fieldName].default || getPlaceholder(fieldName, schema[fieldName])"
        class="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
        :class="{ 'border-red-500': errors[fieldName] }"
      />

      <!-- Text (Long) Field -->
      <textarea
        v-if="schema[fieldName] && schema[fieldName].type === 'text'"
        :id="fieldName"
        v-model="formData[fieldName]"
        :required="schema[fieldName].required"
        :placeholder="schema[fieldName].default || getPlaceholder(fieldName, schema[fieldName])"
        rows="4"
        class="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
        :class="{ 'border-red-500': errors[fieldName] }"
      ></textarea>

      <!-- Rich Text Field -->
      <div v-if="schema[fieldName] && schema[fieldName].type === 'richtext'" class="space-y-2">
        <textarea
          :id="fieldName"
          v-model="formData[fieldName]"
          :required="schema[fieldName].required"
          rows="6"
          class="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500 font-mono text-sm"
          :class="{ 'border-red-500': errors[fieldName] }"
        ></textarea>
        <p class="text-xs text-gray-500">Rich text editor (Markdown/HTML supported)</p>
      </div>

      <!-- Number / Integer Field -->
      <input
        v-if="schema[fieldName] && (schema[fieldName].type === 'number' || schema[fieldName].type === 'integer')"
        :id="fieldName"
        v-model.number="formData[fieldName]"
        :type="schema[fieldName].type === 'integer' ? 'number' : 'number'"
        :step="schema[fieldName].type === 'integer' ? 1 : 0.01"
        :required="schema[fieldName].required"
        :min="schema[fieldName].min"
        :max="schema[fieldName].max"
        :placeholder="schema[fieldName].default || '0'"
        class="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
        :class="{ 'border-red-500': errors[fieldName] }"
      />

      <!-- Boolean Field -->
      <div v-if="schema[fieldName] && schema[fieldName].type === 'boolean'" class="flex items-center">
        <input
          :id="fieldName"
          v-model="formData[fieldName]"
          type="checkbox"
          class="w-4 h-4 text-blue-600 border-gray-300 rounded focus:ring-blue-500"
        />
        <label :for="fieldName" class="ml-2 text-sm text-gray-700">
          {{ getFieldLabel(fieldName, schema[fieldName]) }}
        </label>
      </div>

      <!-- Date Field -->
      <input
        v-if="schema[fieldName] && schema[fieldName].type === 'date'"
        :id="fieldName"
        v-model="formData[fieldName]"
        type="date"
        :required="schema[fieldName].required"
        class="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
        :class="{ 'border-red-500': errors[fieldName] }"
      />

      <!-- DateTime Field -->
      <input
        v-if="schema[fieldName] && schema[fieldName].type === 'datetime'"
        :id="fieldName"
        v-model="formData[fieldName]"
        type="datetime-local"
        :required="schema[fieldName].required"
        class="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
        :class="{ 'border-red-500': errors[fieldName] }"
      />

      <!-- Time Field -->
      <input
        v-if="schema[fieldName] && schema[fieldName].type === 'time'"
        :id="fieldName"
        v-model="formData[fieldName]"
        type="time"
        :required="schema[fieldName].required"
        class="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
        :class="{ 'border-red-500': errors[fieldName] }"
      />

      <!-- Enum Field -->
      <select
        v-if="schema[fieldName] && schema[fieldName].type === 'enum'"
        :id="fieldName"
        v-model="formData[fieldName]"
        :required="schema[fieldName].required"
        class="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
        :class="{ 'border-red-500': errors[fieldName] }"
      >
        <option value="">{{ $t('common.optional') }}</option>
        <option
          v-for="option in getEnumOptions(schema[fieldName])"
          :key="option"
          :value="option"
        >
          {{ option }}
        </option>
      </select>

      <!-- Relation Field -->
      <div v-if="schema[fieldName] && schema[fieldName].type === 'relation'" class="space-y-2">
        <label :for="fieldName" class="block text-sm font-medium text-gray-700 mb-1" v-if="!schema[fieldName].description">
          {{ getFieldLabel(fieldName, schema[fieldName]) }}
          <span v-if="schema[fieldName].required" class="text-red-500">*</span>
        </label>
        <div v-if="schema[fieldName].description" class="text-xs text-gray-500 mb-1">{{ schema[fieldName].description }}</div>
        <p v-if="!schema[fieldName].targetContentType" class="text-xs text-yellow-600 bg-yellow-50 p-2 rounded mb-2">
          {{ $t('fieldBuilder.relationNotConfigured') }}
        </p>
        <RelationSelector
          v-if="schema[fieldName].targetContentType"
          :target-content-type-uid="schema[fieldName].targetContentType"
          :relation-type="schema[fieldName].relationType || 'manyToOne'"
          :model-value="formData[fieldName]"
          @update:model-value="(val) => { formData[fieldName] = val }"
        />
      </div>

      <!-- Media Field -->
      <div v-if="schema[fieldName] && (schema[fieldName].type === 'media' || schema[fieldName].type === 'mediaMultiple')" class="space-y-2">
        <div class="flex items-center gap-2">
          <input
            :id="fieldName"
            type="file"
            :multiple="schema[fieldName].multiple || schema[fieldName].type === 'mediaMultiple'"
            @change="handleFileUpload($event, fieldName, schema[fieldName])"
            class="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
          />
        </div>
        <div v-if="formData[fieldName]" class="text-sm text-gray-600">
          {{ getMediaPreview(fieldName) }}
        </div>
      </div>

      <!-- JSON Field -->
      <div v-if="schema[fieldName] && schema[fieldName].type === 'json'" class="space-y-2">
        <textarea
          :id="fieldName"
          v-model="jsonFields[fieldName]"
          rows="6"
          class="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500 font-mono text-sm"
          :class="{ 'border-red-500': errors[fieldName] }"
          @input="updateJsonField(fieldName)"
        ></textarea>
      </div>

      <!-- Error message - show only once per field -->
      <p v-if="errors[fieldName]" class="text-sm text-red-600">{{ errors[fieldName] }}</p>
    </div>

    <!-- Empty State -->
    <div v-if="orderedFieldNames.length === 0" class="text-center py-8 text-gray-500">
      <p>{{ $t('contentTypes.noFields') }}</p>
      <p class="text-sm mt-2">{{ $t('contentTypes.addFieldsFirst') }}</p>
    </div>
  </div>
</template>

<script setup>
import { ref, reactive, watch, onMounted, computed } from 'vue'
import { useI18n } from 'vue-i18n'
import RelationSelector from './RelationSelector.vue'

const props = defineProps({
  schema: {
    type: Object,
    required: true,
    default: () => ({})
  },
  modelValue: {
    type: Object,
    default: () => ({})
  },
  availableContentTypes: {
    type: Array,
    default: () => []
  }
})

const emit = defineEmits(['update:modelValue', 'validation'])

const { t } = useI18n()
const formData = reactive({})
const errors = reactive({})
const jsonFields = reactive({})

// Compute ordered field names - preserve order from schema if it has _order, otherwise use Object.keys
const orderedFieldNames = computed(() => {
  // Try to get order from schema if it's stored
  if (props.schema._order && Array.isArray(props.schema._order)) {
    return props.schema._order.filter(name => props.schema[name] !== undefined)
  }
  // Otherwise use Object.keys which maintains insertion order in modern JS
  return Object.keys(props.schema).filter(key => key !== '_order')
})

// Initialize form data from modelValue and schema defaults
onMounted(() => {
  initializeFormData()
})

watch(() => props.modelValue, () => {
  initializeFormData()
}, { deep: true })

watch(formData, () => {
  emit('update:modelValue', { ...formData })
  // Debounce validation to avoid multiple calls
  clearTimeout(validationTimeout)
  validationTimeout = setTimeout(() => {
    validateAll()
  }, 300)
}, { deep: true })

let validationTimeout = null

const initializeFormData = () => {
  // Start with existing values
  Object.assign(formData, props.modelValue || {})
  
  // Apply defaults from schema
  orderedFieldNames.value.forEach(fieldName => {
    const fieldConfig = props.schema[fieldName]
    if (!fieldConfig) return
    
    if (formData[fieldName] === undefined || formData[fieldName] === null) {
      if (fieldConfig.default !== undefined) {
        formData[fieldName] = fieldConfig.default
      } else if (fieldConfig.type === 'boolean') {
        formData[fieldName] = false
      } else if (fieldConfig.type === 'number' || fieldConfig.type === 'integer') {
        formData[fieldName] = 0
      } else if (fieldConfig.type === 'array') {
        formData[fieldName] = []
      } else {
        formData[fieldName] = ''
      }
    }
    
    // Initialize JSON fields
    if (fieldConfig.type === 'json') {
      jsonFields[fieldName] = JSON.stringify(formData[fieldName] || {}, null, 2)
    }
  })
}

const getFieldLabel = (fieldName, fieldConfig) => {
  // Convert camelCase to Title Case
  return fieldName
    .replace(/([A-Z])/g, ' $1')
    .replace(/^./, str => str.toUpperCase())
    .trim()
}

const getInputType = (type) => {
  const typeMap = {
    string: 'text',
    email: 'email',
    url: 'url'
  }
  return typeMap[type] || 'text'
}

const getPlaceholder = (fieldName, fieldConfig) => {
  const placeholders = {
    title: 'Enter title',
    name: 'Enter name',
    email: 'example@email.com',
    url: 'https://example.com',
    slug: 'example-slug'
  }
  return placeholders[fieldName.toLowerCase()] || `Enter ${fieldName}`
}

const getEnumOptions = (fieldConfig) => {
  if (Array.isArray(fieldConfig.options)) {
    return fieldConfig.options
  }
  if (typeof fieldConfig.options === 'string') {
    return fieldConfig.options.split('\n').filter(o => o.trim())
  }
  return []
}

const validateField = (fieldName, fieldConfig) => {
  const value = formData[fieldName]
  const fieldErrors = []

  // Required validation - skip relation fields as they are handled separately
  if (fieldConfig.type !== 'relation' && fieldConfig.required && (value === undefined || value === null || value === '')) {
    fieldErrors.push(t('common.required'))
  }
  
  // Relation validation
  if (fieldConfig.type === 'relation' && fieldConfig.required) {
    const relationType = fieldConfig.relationType || 'manyToOne'
    if (relationType === 'oneToOne' || relationType === 'manyToOne') {
      if (value === null || value === undefined || value === '') {
        fieldErrors.push(t('common.required'))
      }
    } else {
      if (!value || (Array.isArray(value) && value.length === 0)) {
        fieldErrors.push(t('common.required'))
      }
    }
  }

  // Type-specific validations
  if (value !== undefined && value !== null && value !== '') {
    // String validations
    if (fieldConfig.type === 'string' || fieldConfig.type === 'text' || fieldConfig.type === 'email' || fieldConfig.type === 'url') {
      if (fieldConfig.minLength && String(value).length < fieldConfig.minLength) {
        fieldErrors.push(`Minimum length: ${fieldConfig.minLength}`)
      }
      if (fieldConfig.maxLength && String(value).length > fieldConfig.maxLength) {
        fieldErrors.push(`Maximum length: ${fieldConfig.maxLength}`)
      }
      if (fieldConfig.type === 'email' && !/^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(value)) {
        fieldErrors.push('Invalid email format')
      }
      if (fieldConfig.type === 'url' && !/^https?:\/\/.+/.test(value)) {
        fieldErrors.push('Invalid URL format')
      }
    }

    // Number validations
    if (fieldConfig.type === 'number' || fieldConfig.type === 'integer') {
      if (fieldConfig.min !== undefined && Number(value) < fieldConfig.min) {
        fieldErrors.push(`Minimum value: ${fieldConfig.min}`)
      }
      if (fieldConfig.max !== undefined && Number(value) > fieldConfig.max) {
        fieldErrors.push(`Maximum value: ${fieldConfig.max}`)
      }
    }
  }

  errors[fieldName] = fieldErrors.length > 0 ? fieldErrors.join(', ') : null
  return fieldErrors.length === 0
}

const validateAll = () => {
  // Clear all errors first
  Object.keys(errors).forEach(key => {
    errors[key] = null
  })
  
  let isValid = true
  orderedFieldNames.value.forEach(fieldName => {
    const fieldConfig = props.schema[fieldName]
    if (!fieldConfig) return
    const fieldValid = validateField(fieldName, fieldConfig)
    if (!fieldValid) {
      isValid = false
    }
  })
  emit('validation', { isValid, errors: { ...errors } })
  return isValid
}

const handleFileUpload = (event, fieldName, fieldConfig) => {
  const files = event.target.files
  if (files.length === 0) return

  // For now, just store file names
  // TODO: Implement actual file upload
  if (fieldConfig.multiple || fieldConfig.type === 'mediaMultiple') {
    formData[fieldName] = Array.from(files).map(f => f.name)
  } else {
    formData[fieldName] = files[0].name
  }
}

const getMediaPreview = (fieldName) => {
  const value = formData[fieldName]
  if (Array.isArray(value)) {
    return `${value.length} file(s) selected`
  }
  return value ? `File: ${value}` : 'No file selected'
}

const updateJsonField = (fieldName) => {
  try {
    formData[fieldName] = JSON.parse(jsonFields[fieldName])
    errors[fieldName] = null
  } catch (e) {
    errors[fieldName] = 'Invalid JSON format'
  }
}

const getData = () => {
  const data = { ...formData }
  
  // Ensure relation fields are sent as IDs, not objects
  orderedFieldNames.value.forEach(fieldName => {
    const fieldConfig = props.schema[fieldName]
    if (!fieldConfig) return

    if (fieldConfig.type === 'relation' && data[fieldName] !== undefined && data[fieldName] !== null) {
      const relationType = fieldConfig.relationType || 'manyToOne'
      if (relationType === 'oneToMany' || relationType === 'manyToMany') {
        // Ensure array of IDs
        if (Array.isArray(data[fieldName])) {
          data[fieldName] = data[fieldName].map(item => 
            typeof item === 'object' && item.id ? item.id : item
          ).filter(id => id !== null && id !== undefined)
        }
      } else {
        // Ensure single ID
        if (typeof data[fieldName] === 'object' && data[fieldName].id) {
          data[fieldName] = data[fieldName].id
        }
      }
    }
  })
  
  return data
}

// Expose validation method
defineExpose({
  validate: validateAll,
  getData
})
</script>

