<template>
  <TransitionGroup
    name="toast"
    tag="div"
    class="fixed top-4 right-4 z-50 space-y-2"
  >
    <div
      v-for="toast in toasts"
      :key="toast.id"
      :class="[
        'min-w-[320px] max-w-md w-full bg-white shadow-lg rounded-lg pointer-events-auto ring-1 ring-black ring-opacity-5',
        toast.type === 'success' ? 'border-l-4 border-green-500' : '',
        toast.type === 'error' ? 'border-l-4 border-red-500' : '',
        toast.type === 'warning' ? 'border-l-4 border-yellow-500' : '',
        toast.type === 'info' ? 'border-l-4 border-blue-500' : ''
      ]"
    >
      <div class="p-5">
        <div class="flex items-start">
          <div class="flex-shrink-0">
            <CheckCircleIcon
              v-if="toast.type === 'success'"
              class="h-6 w-6 text-green-400"
            />
            <XCircleIcon
              v-else-if="toast.type === 'error'"
              class="h-6 w-6 text-red-400"
            />
            <ExclamationTriangleIcon
              v-else-if="toast.type === 'warning'"
              class="h-6 w-6 text-yellow-400"
            />
            <InformationCircleIcon
              v-else
              class="h-6 w-6 text-blue-400"
            />
          </div>
          <div class="ml-3 flex-1 min-w-0">
            <p class="text-sm font-semibold text-gray-900 break-words">
              {{ toast.title }}
            </p>
            <p v-if="toast.message" class="mt-2 text-sm text-gray-600 break-words">
              {{ toast.message }}
            </p>
          </div>
          <div class="ml-4 flex-shrink-0 flex">
            <button
              @click="removeToast(toast.id)"
              class="inline-flex rounded-md bg-white text-gray-400 hover:text-gray-600 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-blue-500 transition-colors"
            >
              <span class="sr-only">Close</span>
              <XMarkIcon class="h-5 w-5" />
            </button>
          </div>
        </div>
      </div>
    </div>
  </TransitionGroup>
</template>

<script setup>
import { ref, onMounted, onUnmounted } from 'vue'
import {
  CheckCircleIcon,
  XCircleIcon,
  ExclamationTriangleIcon,
  InformationCircleIcon,
  XMarkIcon
} from '@heroicons/vue/24/outline'

const toasts = ref([])
let toastIdCounter = 0

const addToast = (toast) => {
  const id = ++toastIdCounter
  const newToast = {
    id,
    type: toast.type || 'info',
    title: toast.title,
    message: toast.message,
    duration: toast.duration || 5000
  }
  
  toasts.value.push(newToast)
  
  // Auto remove after duration
  if (newToast.duration > 0) {
    setTimeout(() => {
      removeToast(id)
    }, newToast.duration)
  }
  
  return id
}

const removeToast = (id) => {
  const index = toasts.value.findIndex(t => t.id === id)
  if (index > -1) {
    toasts.value.splice(index, 1)
  }
}

// Expose methods globally
const showToast = {
  success: (title, message = '', duration = 5000) => 
    addToast({ type: 'success', title, message, duration }),
  error: (title, message = '', duration = 5000) => 
    addToast({ type: 'error', title, message, duration }),
  warning: (title, message = '', duration = 5000) => 
    addToast({ type: 'warning', title, message, duration }),
  info: (title, message = '', duration = 5000) => 
    addToast({ type: 'info', title, message, duration })
}

// Make available globally
window.showToast = showToast

onMounted(() => {
  // Store in global scope for easy access
  if (typeof window !== 'undefined') {
    window.$toast = showToast
  }
})

onUnmounted(() => {
  if (typeof window !== 'undefined') {
    delete window.showToast
    delete window.$toast
  }
})

defineExpose({
  addToast,
  removeToast,
  showToast
})
</script>

<style scoped>
.toast-enter-active,
.toast-leave-active {
  transition: all 0.3s ease;
}

.toast-enter-from {
  opacity: 0;
  transform: translateX(100%);
}

.toast-leave-to {
  opacity: 0;
  transform: translateX(100%);
}

.toast-move {
  transition: transform 0.3s ease;
}
</style>

