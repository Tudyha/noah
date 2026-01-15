import { defineStore } from "pinia";

export const useThemeStore = defineStore(
  "theme",
  () => {
    const currentTheme = ref("light");

    function setTheme(theme: string) {
      currentTheme.value = theme;
      document.documentElement.setAttribute("data-theme", theme);
    }

    const initTheme = () => {
      setTheme(currentTheme.value)
    }

    return {
      currentTheme,
      setTheme,
      initTheme,
    };
  },
  {
    persist: true,
  }
);
