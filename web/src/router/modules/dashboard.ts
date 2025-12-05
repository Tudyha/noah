const dashboard = [
  {
    path: "/dashboard",
    name: "Dashboard",
    component: () => import("@/views/dashboard/index.vue"),
    meta: {
      title: "Dashboard",
      requiresAuth: true,
      hideSidebar: false, // dashboard页面显示侧边栏
    },
  },
];

export default dashboard;
