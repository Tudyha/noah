import { createApp } from "vue";
import { createPinia } from "pinia";
import piniaPluginPersistedstate from "pinia-plugin-persistedstate";
import router from "@/router/index";
import "./style.css";
import App from "./App.vue";
import i18n from "@/locales/index";
import { Icon } from "@iconify/vue"; // 导入 Icon 组件

const app = createApp(App);
const pinia = createPinia();
pinia.use(piniaPluginPersistedstate);
app.use(pinia);
app.use(router);
app.use(i18n);
app.component("Icon", Icon); // 全局注册 Icon 组件
app.mount("#app");
