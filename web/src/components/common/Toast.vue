<script setup lang="ts">
import { useUIStore } from '@/stores/ui';
import { storeToRefs } from 'pinia';

const uiStore = useUIStore();
const { toasts } = storeToRefs(uiStore);

const getAlertClass = (type: string) => {
  switch (type) {
    case 'success': return 'alert-success';
    case 'error': return 'alert-error';
    case 'warning': return 'alert-warning';
    case 'info':
    default: return 'alert-info';
  }
};

const getIcon = (type: string) => {
  switch (type) {
    case 'success': return 'mdi:check-circle-outline';
    case 'error': return 'mdi:alert-circle-outline';
    case 'warning': return 'mdi:alert-outline';
    case 'info':
    default: return 'mdi:information-outline';
  }
};
</script>

<template>
  <div class="toast toast-top toast-end z-[9999]">
    <TransitionGroup name="toast">
      <div v-for="toast in toasts" :key="toast.id" class="alert shadow-lg mb-2 py-3 px-4 min-w-[300px]"
        :class="getAlertClass(toast.type)">
        <div class="flex items-center gap-2">
          <Icon :icon="getIcon(toast.type)" class="w-5 h-5 flex-shrink-0" />
          <span class="text-sm font-medium">{{ toast.message }}</span>
        </div>
      </div>
    </TransitionGroup>
  </div>
</template>

<style scoped>
.toast-enter-active,
.toast-leave-active {
  transition: all 0.3s ease;
}

.toast-enter-from {
  opacity: 0;
  transform: translateX(30px);
}

.toast-leave-to {
  opacity: 0;
  transform: translateX(30px);
}
</style>
