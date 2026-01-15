<template>
  <RouterView />
  <Toast />
</template>

<script setup lang="ts">
import Toast from '@/components/common/Toast.vue';
import { useThemeStore} from '@/stores/theme';
const themeStore = useThemeStore();

onMounted(() => {
  themeStore.initTheme();
});

</script>
