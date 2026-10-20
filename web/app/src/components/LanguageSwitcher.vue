<template>
  <select
    :value="locale"
    :aria-label="t('app.language.label')"
    class="h-9 rounded-md border border-border bg-transparent px-2 text-sm font-medium text-foreground transition-colors hover:bg-accent hover:text-accent-foreground cursor-pointer"
    @change="onChange"
  >
    <option value="en">English</option>
    <option value="zh-CN">中文</option>
  </select>
</template>

<script setup>
import { computed } from 'vue'
import { useI18n } from 'vue-i18n'
import i18n, { setLanguage } from '@/i18n'

const { t } = useI18n()
const locale = computed(() => i18n.global.locale.value)

const onChange = (event) => setLanguage(event.target.value)
</script>
