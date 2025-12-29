import { defineStore } from "pinia";
import { ref } from "vue";

export type ToastType = "info" | "success" | "warning" | "error";

export interface ToastMessage {
  id: number;
  message: string;
  type: ToastType;
  duration?: number;
}

export const useUIStore = defineStore("ui", () => {
  const toasts = ref<ToastMessage[]>([]);
  let nextId = 0;

  const showToast = (message: string, type: ToastType = "info", duration = 3000) => {
    const id = nextId++;
    const toast: ToastMessage = { id, message, type, duration };
    toasts.value.push(toast);

    if (duration > 0) {
      setTimeout(() => {
        removeToast(id);
      }, duration);
    }
  };

  const removeToast = (id: number) => {
    const index = toasts.value.findIndex((t) => t.id === id);
    if (index !== -1) {
      toasts.value.splice(index, 1);
    }
  };

  return {
    toasts,
    showToast,
    removeToast,
  };
});
