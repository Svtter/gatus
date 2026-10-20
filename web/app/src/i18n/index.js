import { createI18n } from 'vue-i18n'
import en from './locales/en'
import zhCN from './locales/zh-CN'

const STORAGE_KEY = 'gatus-language'
const SUPPORTED = ['en', 'zh-CN']

function detectLanguage() {
  const stored = localStorage.getItem(STORAGE_KEY)
  if (stored && SUPPORTED.includes(stored)) {
    return stored
  }
  if (navigator.language && navigator.language.toLowerCase().startsWith('zh')) {
    return 'zh-CN'
  }
  return 'en'
}

const i18n = createI18n({
  legacy: false,
  locale: detectLanguage(),
  fallbackLocale: 'en',
  messages: {
    en,
    'zh-CN': zhCN,
  },
})

function setLanguage(locale) {
  i18n.global.locale.value = locale
  localStorage.setItem(STORAGE_KEY, locale)
  document.documentElement.setAttribute('lang', locale)
}

setLanguage(i18n.global.locale.value)

export { setLanguage }
export default i18n
