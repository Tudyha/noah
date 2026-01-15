const dashboard = [
  {
    path: "/agent",
    name: "Agent",
    redirect: "/agent/list",
    meta: {
      title: "Agent",
      requiresAuth: true,
      hideSidebar: false,
    },
    children: [
      {
        path: "list",
        name: "AgentList",
        component: () => import("@/views/agent/index.vue"),
        meta: {
          title: "Agent List",
          requiresAuth: true,
          hideSidebar: false,
        },
      },
      {
        path: "console/:id",
        name: "AgentConsole",
        component: () => import("@/views/agent/console.vue"),
        meta: {
          title: "Agent Console",
          requiresAuth: true,
          hideSidebar: false,
        },
      }
    ],
  },
];

export default dashboard;
