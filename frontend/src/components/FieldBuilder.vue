<template>
  <div class="space-y-4">
    <div class="flex justify-between items-center">
      <h3 class="text-lg font-semibold text-gray-900">{{ $t('contentTypes.fields') }}</h3>
      <button
        @click="addField"
        type="button"
        class="px-3 py-1.5 text-sm bg-blue-600 text-white rounded-md hover:bg-blue-700"
      >
        + {{ $t('contentTypes.addField') }}
      </button>
    </div>

    <div v-if="fields.length === 0" class="text-center py-8 text-gray-500">
      {{ $t('contentTypes.noFields') }}. {{ $t('contentTypes.clickAddField') }}
    </div>

    <div v-else class="space-y-3">
      <div
        v-for="(field, index) in fields"
        :key="index"
        class="border border-gray-200 rounded-lg p-4 bg-white"
      >
        <div class="flex justify-between items-start mb-3">
          <div class="flex-1">
            <div class="flex items-center gap-2 mb-2">
              <input
                v-model="field.name"
                type="text"
                :placeholder="$t('contentTypes.fieldName')"
                class="px-2 py-1 border border-gray-300 rounded text-sm font-medium"
                :class="field.name ? 'bg-gray-50' : 'bg-yellow-50'"
              />
              <span
                v-if="field.type"
                class="px-2 py-1 text-xs bg-blue-100 text-blue-800 rounded"
              >
                {{ field.type }}
              </span>
            </div>
          </div>
          <div class="flex gap-2">
            <button
              @click="moveField(index, 'up')"
              type="button"
              :disabled="index === 0"
              class="p-1 text-gray-400 hover:text-gray-600 disabled:opacity-30"
              :title="$t('fieldBuilder.moveUp')"
            >
              ↑
            </button>
            <button
              @click="moveField(index, 'down')"
              type="button"
              :disabled="index === fields.length - 1"
              class="p-1 text-gray-400 hover:text-gray-600 disabled:opacity-30"
              :title="$t('fieldBuilder.moveDown')"
            >
              ↓
            </button>
            <button
              @click="removeField(index)"
              type="button"
              class="p-1 text-red-400 hover:text-red-600"
              :title="$t('fieldBuilder.removeField')"
            >
              ×
            </button>
          </div>
        </div>

        <div class="grid grid-cols-2 gap-3">
          <!-- Field Type -->
          <div>
            <label class="block text-xs font-medium text-gray-700 mb-1">{{ $t('fieldBuilder.type') }} *</label>
            <select
              v-model="field.type"
              class="w-full px-2 py-1.5 text-sm border border-gray-300 rounded-md"
              @change="onFieldTypeChange(field)"
            >
              <option value="">{{ $t('fieldBuilder.selectType') }}</option>
              <optgroup :label="$t('fieldBuilder.text')">
                <option value="string">{{ $t('fieldBuilder.string') }}</option>
                <option value="text">{{ $t('fieldBuilder.text') }}</option>
                <option value="richtext">{{ $t('fieldBuilder.richtext') }}</option>
                <option value="email">{{ $t('fieldBuilder.email') }}</option>
                <option value="url">{{ $t('fieldBuilder.url') }}</option>
              </optgroup>
              <optgroup :label="$t('fieldBuilder.number')">
                <option value="number">{{ $t('fieldBuilder.number') }}</option>
                <option value="integer">{{ $t('fieldBuilder.integer') }}</option>
              </optgroup>
              <optgroup :label="$t('fieldBuilder.dateTime')">
                <option value="date">{{ $t('fieldBuilder.date') }}</option>
                <option value="time">{{ $t('fieldBuilder.time') }}</option>
                <option value="datetime">{{ $t('fieldBuilder.datetime') }}</option>
              </optgroup>
              <optgroup :label="$t('fieldBuilder.booleanEnum')">
                <option value="boolean">{{ $t('fieldBuilder.boolean') }}</option>
                <option value="enum">{{ $t('fieldBuilder.enum') }}</option>
              </optgroup>
              <optgroup :label="$t('fieldBuilder.relations')">
                <option value="relation">{{ $t('fieldBuilder.relation') }}</option>
              </optgroup>
              <optgroup :label="$t('fieldBuilder.media')">
                <option value="media">{{ $t('fieldBuilder.mediaSingle') }}</option>
                <option value="mediaMultiple">{{ $t('fieldBuilder.mediaMultiple') }}</option>
              </optgroup>
              <optgroup :label="$t('fieldBuilder.structured')">
                <option value="json">{{ $t('fieldBuilder.json') }}</option>
                <option value="component">{{ $t('fieldBuilder.component') }}</option>
              </optgroup>
            </select>
          </div>

          <!-- Field Options -->
          <div class="flex flex-wrap gap-2 items-end">
            <label class="flex items-center text-xs">
              <input
                v-model="field.required"
                type="checkbox"
                class="mr-1"
              />
              {{ $t('contentTypes.required') }}
            </label>
            <label class="flex items-center text-xs">
              <input
                v-model="field.unique"
                type="checkbox"
                class="mr-1"
              />
              {{ $t('contentTypes.unique') }}
            </label>
          </div>
        </div>

        <!-- Description -->
        <div class="mt-2">
          <input
            v-model="field.description"
            type="text"
            :placeholder="$t('fieldBuilder.descriptionOptional')"
            class="w-full px-2 py-1 text-sm border border-gray-300 rounded-md"
          />
        </div>

        <!-- Default Value -->
        <div v-if="field.type && field.type !== 'relation' && field.type !== 'media' && field.type !== 'mediaMultiple'" class="mt-2">
          <input
            v-model="field.default"
            type="text"
            :placeholder="getDefaultPlaceholder(field.type)"
            class="w-full px-2 py-1 text-sm border border-gray-300 rounded-md"
          />
        </div>

        <!-- Relation Configuration -->
        <div v-if="field.type === 'relation'" class="mt-3 p-3 bg-blue-50 rounded-md border border-blue-200">
          <div class="mb-2">
            <p class="text-xs text-blue-700 font-medium mb-1">{{ $t('fieldBuilder.relationInfo') }}</p>
            <p class="text-xs text-blue-600">{{ $t('fieldBuilder.relationHint') }}</p>
          </div>
          <div class="grid grid-cols-2 gap-3 mt-3">
            <div>
              <label class="block text-xs font-medium text-gray-700 mb-1">{{ $t('fieldBuilder.relationType') }}</label>
              <select
                v-model="field.relationType"
                class="w-full px-2 py-1.5 text-sm border border-gray-300 rounded-md"
              >
                <option value="oneToOne">{{ $t('fieldBuilder.oneToOne') }}</option>
                <option value="oneToMany">{{ $t('fieldBuilder.oneToMany') }}</option>
                <option value="manyToOne">{{ $t('fieldBuilder.manyToOne') }}</option>
                <option value="manyToMany">{{ $t('fieldBuilder.manyToMany') }}</option>
              </select>
            </div>
            <div>
              <label class="block text-xs font-medium text-gray-700 mb-1">{{ $t('fieldBuilder.targetContentType') }}</label>
              <select
                v-model="field.targetContentType"
                class="w-full px-2 py-1.5 text-sm border border-gray-300 rounded-md"
              >
                <option value="">{{ $t('fieldBuilder.selectTargetType') }}</option>
                <option
                  v-for="ct in availableContentTypes"
                  :key="ct.uid"
                  :value="ct.uid"
                >
                  {{ ct.displayName }} ({{ ct.uid }})
                </option>
              </select>
            </div>
          </div>
        </div>

        <!-- Enum Options -->
        <div v-if="field.type === 'enum'" class="mt-3 p-3 bg-gray-50 rounded-md">
          <label class="block text-xs font-medium text-gray-700 mb-2">{{ $t('fieldBuilder.enumOptions') }}</label>
          <textarea
            v-model="field.enumOptions"
            :placeholder="$t('fieldBuilder.enumPlaceholder')"
            class="w-full px-2 py-1.5 text-sm border border-gray-300 rounded-md font-mono"
            rows="3"
            @input="updateEnumOptions(field)"
          ></textarea>
        </div>

        <!-- Validation -->
        <div v-if="field.type === 'string' || field.type === 'text' || field.type === 'email' || field.type === 'url'" class="mt-3 grid grid-cols-2 gap-2">
          <div>
            <label class="block text-xs font-medium text-gray-700 mb-1">{{ $t('fieldBuilder.minLength') }}</label>
            <input
              v-model.number="field.minLength"
              type="number"
              min="0"
              class="w-full px-2 py-1 text-sm border border-gray-300 rounded-md"
            />
          </div>
          <div>
            <label class="block text-xs font-medium text-gray-700 mb-1">{{ $t('fieldBuilder.maxLength') }}</label>
            <input
              v-model.number="field.maxLength"
              type="number"
              min="0"
              class="w-full px-2 py-1 text-sm border border-gray-300 rounded-md"
            />
          </div>
        </div>

        <div v-if="field.type === 'number' || field.type === 'integer'" class="mt-3 grid grid-cols-2 gap-2">
          <div>
            <label class="block text-xs font-medium text-gray-700 mb-1">{{ $t('fieldBuilder.min') }}</label>
            <input
              v-model.number="field.min"
              type="number"
              class="w-full px-2 py-1 text-sm border border-gray-300 rounded-md"
            />
          </div>
          <div>
            <label class="block text-xs font-medium text-gray-700 mb-1">{{ $t('fieldBuilder.max') }}</label>
            <input
              v-model.number="field.max"
              type="number"
              class="w-full px-2 py-1 text-sm border border-gray-300 rounded-md"
            />
          </div>
        </div>
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
    type: Object,
    default: () => ({})
  },
  availableContentTypes: {
    type: Array,
    default: () => []
  }
})

const emit = defineEmits(['update:modelValue'])

const fields = ref([])

// Load available content types if not provided
const availableContentTypes = ref(props.availableContentTypes || [])

onMounted(async () => {
  if (availableContentTypes.value.length === 0) {
    try {
      const response = await contentAPI.getContentTypes()
      availableContentTypes.value = response.data.data || []
    } catch (error) {
      console.error('Failed to load content types:', error)
    }
  }
  
  // Initialize fields from modelValue
  if (props.modelValue && Object.keys(props.modelValue).length > 0) {
    fields.value = Object.entries(props.modelValue).map(([name, config]) => ({
      name,
      type: config.type || '',
      required: config.required || false,
      unique: config.unique || false,
      description: config.description || '',
      default: config.default !== undefined ? String(config.default) : '',
      relationType: config.relationType || 'manyToOne',
      targetContentType: config.targetContentType || '',
      minLength: config.minLength,
      maxLength: config.maxLength,
      min: config.min,
      max: config.max,
      enumOptions: config.options ? config.options.join('\n') : ''
    }))
  }
})

// Watch fields and emit schema
watch(fields, (newFields) => {
  const schema = {}
  const order = []
  
  newFields.forEach(field => {
    if (field.name && field.type) {
      const fieldConfig = {
        type: field.type === 'mediaMultiple' ? 'media' : field.type
      }
      
      if (field.required) fieldConfig.required = true
      if (field.unique) fieldConfig.unique = true
      if (field.description) fieldConfig.description = field.description
      if (field.default !== '') {
        // Try to parse default value based on type
        if (field.type === 'number' || field.type === 'integer') {
          fieldConfig.default = parseFloat(field.default)
        } else if (field.type === 'boolean') {
          fieldConfig.default = field.default === 'true'
        } else {
          fieldConfig.default = field.default
        }
      }
      
      // Relation specific
      if (field.type === 'relation') {
        fieldConfig.relationType = field.relationType || 'manyToOne'
        fieldConfig.targetContentType = field.targetContentType
      }
      
      // Media specific
      if (field.type === 'mediaMultiple') {
        fieldConfig.multiple = true
      }
      
      // Validation
      if (field.minLength !== undefined && field.minLength !== '') fieldConfig.minLength = field.minLength
      if (field.maxLength !== undefined && field.maxLength !== '') fieldConfig.maxLength = field.maxLength
      if (field.min !== undefined && field.min !== '') fieldConfig.min = field.min
      if (field.max !== undefined && field.max !== '') fieldConfig.max = field.max
      
      // Enum options
      if (field.type === 'enum' && field.enumOptions) {
        fieldConfig.options = field.enumOptions.split('\n').filter(o => o.trim())
      }
      
      schema[field.name] = fieldConfig
      order.push(field.name)
    }
  })
  
  // Store order in schema for DynamicForm to use
  schema._order = order
  
  emit('update:modelValue', schema)
}, { deep: true })

const addField = () => {
  fields.value.push({
    name: '',
    type: '',
    required: false,
    unique: false,
    description: '',
    default: '',
    relationType: 'manyToOne',
    targetContentType: '',
    minLength: undefined,
    maxLength: undefined,
    min: undefined,
    max: undefined,
    enumOptions: ''
  })
}

const removeField = (index) => {
  fields.value.splice(index, 1)
}

const moveField = (index, direction) => {
  if (direction === 'up' && index > 0) {
    [fields.value[index], fields.value[index - 1]] = [fields.value[index - 1], fields.value[index]]
  } else if (direction === 'down' && index < fields.value.length - 1) {
    [fields.value[index], fields.value[index + 1]] = [fields.value[index + 1], fields.value[index]]
  }
}

const onFieldTypeChange = (field) => {
  // Reset relation-specific fields when type changes
  if (field.type !== 'relation') {
    field.relationType = 'manyToOne'
    field.targetContentType = ''
  }
}

const updateEnumOptions = (field) => {
  // Options are already stored in field.enumOptions
  // The watch handler will process them
}

const getDefaultPlaceholder = (type) => {
  const placeholders = {
    string: t('fieldBuilder.defaultString'),
    text: t('fieldBuilder.defaultText'),
    number: '0',
    integer: '0',
    boolean: t('fieldBuilder.defaultBoolean'),
    email: 'example@email.com',
    url: 'https://example.com',
    date: '2024-01-01'
  }
  return placeholders[type] || t('fieldBuilder.defaultValue')
}
</script>

