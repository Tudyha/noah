<template>
  <div class="dropdown dropdown-end">
    <div tabindex="0" role="button" aria-label="Change Theme">
      <ThemePreview :theme="currentTheme" />
    </div>
    <div tabindex="0"
      class="dropdown-content bg-base-200 text-base-content rounded-box h-122 max-h-[calc(100vh-8.6rem)] overflow-y-auto  border-white/5 shadow-2xl  outline-black/5">
      <ul class="menu">
        <li v-for="theme in themes" :key="theme" class="w-40">
          <a @click="setTheme(theme)" :class="{ 'active': currentTheme === theme }">
            <ThemePreview :theme="theme" />
            {{ theme }}
          </a>
        </li>
      </ul>
    </div>
  </div>
</template>

<script setup lang="ts">
import { useThemeStore } from '@/stores/theme';
import { storeToRefs } from 'pinia';
import ThemePreview from '@/components/theme/preview.vue';

const themeStore = useThemeStore()
const { currentTheme } = storeToRefs(themeStore)

const themes = ['light', 'dark', 'cupcake', 'bumblebee', 'emerald', 'corporate', 'synthwave', 'retro', 'cyberpunk',
  'valentine', 'halloween', 'garden', 'forest', 'aqua', 'lofi', 'pastel', 'fantasy', 'wireframe', 'black', 'luxury', 'dracula', 'cmyk',
  'autumn', 'business', 'acid', 'lemonade', 'night', 'coffee', 'winter', 'dim', 'nord', 'sunset', 'caramellatte', 'abyss', 'silk']

const setTheme = (theme: string) => {
  themeStore.setTheme(theme as any)
}
</script>
