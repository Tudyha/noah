import { defineStore } from "pinia";

export const useThemeStore = defineStore(
  "theme",
  () => {
    const currentTheme = ref("light");

    function setTheme(theme: string) {
      currentTheme.value = theme;
      document.documentElement.setAttribute("data-theme", theme);
    }

    return {
      currentTheme,
      setTheme,
    };
  },
  {
    persist: true,
  }
);
