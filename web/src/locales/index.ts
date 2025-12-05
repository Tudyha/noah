import { createI18n } from "vue-i18n";
import en from "./en";
import zhCN from "./zh-CN";

const messages = {
  en,
  "zh-CN": zhCN,
};

const savedLocale = localStorage.getItem("locale") || "zh-CN";

const i18n = createI18n({
  legacy: false,
  locale: savedLocale,
  fallbackLocale: "en",
  messages,
});

export default i18n;
