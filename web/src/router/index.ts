import { createRouter, createWebHistory } from "vue-router";
import type { RouteRecordRaw } from "vue-router";
import NProgress from "@/config/nprogress";
import layout from "./dynamic-route";
import { useUserStore } from "@/stores/auth";

// 根路由
export const RootRoute = {
  path: "/",
  name: "Root",
  redirect: "/home",
  meta: {
    title: "Root",
  },
};

// Basic routing without permission
// 无需认证的基本路由
export const basicRoutes = [RootRoute, layout];

const router = createRouter({
  history: createWebHistory(),
  routes: basicRoutes as RouteRecordRaw[],
  scrollBehavior: () => ({ left: 0, top: 0 }),
});

// Injection Progress
router.beforeEach((to, _, next) => {
  NProgress.start(); // 开启进度条
  const requireAuth = to.meta.requireAuth ?? true;
  const userStore = useUserStore();
  if (to.name === "Login") {
    if (userStore.isLogined) {
      return next("/");
    } else {
      return next();
    }
  }

  if (requireAuth && !userStore.isLogined) {
    return next({ name: "Login" });
  }
  next();
});

router.afterEach(() => {
  NProgress.done();
});

export default router;
