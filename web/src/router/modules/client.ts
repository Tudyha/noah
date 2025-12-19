const dashboard = [
  {
    path: "/client",
    name: "Client",
    redirect: "/client/list",
    meta: {
      title: "Client",
      requiresAuth: true,
      hideSidebar: false,
    },
    children: [
      {
        path: "list",
        name: "ClientList",
        component: () => import("@/views/client/index.vue"),
        meta: {
          title: "Client List",
          requiresAuth: true,
          hideSidebar: false,
        },
      },
      {
        path: "console/:id",
        name: "ClientConsole",
        component: () => import("@/views/client/console.vue"),
        meta: {
          title: "Client Console",
          requiresAuth: true,
          hideSidebar: false,
        },
      }
    ],
  },
];

export default dashboard;
