const dashboard = [
  {
    path: "/client",
    name: "Client",
    component: () => import("@/views/client/index.vue"),
    meta: {
      title: "Client",
      requiresAuth: true,
      hideSidebar: false,
    },
  },
];

export default dashboard;
