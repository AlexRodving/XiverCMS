import { createI18n } from 'vue-i18n'
import ru from './ru.js'
import en from './en.js'

// Get saved language from localStorage or default to Russian
const getSavedLocale = () => {
  const saved = localStorage.getItem('xivercms-locale')
  return saved || 'ru'
}

const i18n = createI18n({
  locale: getSavedLocale(),
  fallbackLocale: 'ru',
  messages: {
    ru,
    en
  },
  legacy: false // Use Composition API mode
})

export default i18n

// Helper function to change locale
export const setLocale = (locale) => {
  i18n.global.locale.value = locale
  localStorage.setItem('xivercms-locale', locale)
}

// Helper function to get current locale
export const getLocale = () => {
  return i18n.global.locale.value
}

