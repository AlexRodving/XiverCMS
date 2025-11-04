<template>
  <div class="relative">
    <button
      @click="toggleMenu"
      class="flex items-center space-x-2 px-3 py-2 text-sm font-medium text-gray-700 hover:text-gray-900 hover:bg-gray-100 rounded-lg transition-colors"
      :title="$t('language.changeLanguage')"
    >
      <GlobeAltIcon class="w-5 h-5" />
      <span class="hidden sm:inline">{{ currentLanguageName }}</span>
    </button>

    <!-- Dropdown menu -->
    <div
      v-if="isOpen"
      ref="dropdownRef"
      class="absolute left-0 bottom-full mb-2 w-48 bg-white rounded-md shadow-lg py-1 z-50 border border-gray-200"
    >
      <button
        @click="setLanguage('ru')"
        class="w-full text-left px-4 py-2 text-sm text-gray-700 hover:bg-gray-100 flex items-center justify-between"
        :class="{ 'bg-blue-50 text-blue-700': i18nStore.currentLocale === 'ru' }"
      >
        <span>🇷🇺 {{ $t('language.russian') }}</span>
        <CheckIcon v-if="i18nStore.currentLocale === 'ru'" class="w-4 h-4 text-blue-600" />
      </button>
      <button
        @click="setLanguage('en')"
        class="w-full text-left px-4 py-2 text-sm text-gray-700 hover:bg-gray-100 flex items-center justify-between"
        :class="{ 'bg-blue-50 text-blue-700': i18nStore.currentLocale === 'en' }"
      >
        <span>🇬🇧 {{ $t('language.english') }}</span>
        <CheckIcon v-if="i18nStore.currentLocale === 'en'" class="w-4 h-4 text-blue-600" />
      </button>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, onUnmounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { useI18nStore } from '../stores/i18n'
import { GlobeAltIcon, CheckIcon } from '@heroicons/vue/24/outline'

const { t } = useI18n()
const i18nStore = useI18nStore()

const isOpen = ref(false)
const dropdownRef = ref(null)

const currentLanguageName = computed(() => {
  return i18nStore.currentLocale === 'ru' ? t('language.russian') : t('language.english')
})

const toggleMenu = () => {
  isOpen.value = !isOpen.value
}

const closeMenu = () => {
  isOpen.value = false
}

const setLanguage = (locale) => {
  i18nStore.setLocale(locale)
  closeMenu()
}

// Handle click outside
const handleClickOutside = (event) => {
  if (dropdownRef.value && !dropdownRef.value.contains(event.target) && !event.target.closest('button')) {
    closeMenu()
  }
}

onMounted(() => {
  document.addEventListener('click', handleClickOutside)
})

onUnmounted(() => {
  document.removeEventListener('click', handleClickOutside)
})
</script>

