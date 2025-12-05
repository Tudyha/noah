<template>
  <div class="dropdown dropdown-end">
    <div tabindex="0" role="button">
      <div class="btn group btn-sm btn-ghost p-1">
        <div class="bg-base-100">
          <Icon icon="mdi:language" class="w-5 h-5" />
        </div>
      </div>
    </div>
    <ul tabindex="0" class="dropdown-content menu p-2 shadow bg-base-200 rounded-box w-25 mt-4">
      <li v-for="lang in languages" :key="lang.code">
        <a @click="changeLanguage(lang.code)" :class="{ 'active': currentLocale === lang.code }">
          {{ lang.name }}
        </a>
      </li>
    </ul>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { useI18n } from 'vue-i18n'

const { locale } = useI18n()

const languages = [
  { code: 'zh-CN', name: '简体中文' },
  { code: 'en', name: 'English' }
]

const currentLocale = computed(() => locale.value)

const changeLanguage = (lang: string) => {
  locale.value = lang
  localStorage.setItem('locale', lang)
}
</script>
