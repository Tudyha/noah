const home = [
  {
    path: "/home",
    name: "Home",
    component: () => import("@/views/home/index.vue"),
    meta: {
      title: "Home",
      requireAuth: false,
      hideSidebar: true, // home页面不显示侧边栏
    },
  },
];

export default home;
