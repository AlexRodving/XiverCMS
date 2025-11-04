import { defineStore } from 'pinia'
import { useI18n } from 'vue-i18n'
import { setLocale, getLocale } from '../locales'

export const useI18nStore = defineStore('i18n', {
  state: () => ({
    locale: getLocale()
  }),
  
  actions: {
    setLocale(locale) {
      this.locale = locale
      setLocale(locale)
    },
    
    toggleLocale() {
      const newLocale = this.locale === 'ru' ? 'en' : 'ru'
      this.setLocale(newLocale)
    }
  },
  
  getters: {
    currentLocale: (state) => state.locale,
    isRussian: (state) => state.locale === 'ru',
    isEnglish: (state) => state.locale === 'en'
  }
})

