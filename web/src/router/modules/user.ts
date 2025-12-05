
const user = [
  {
    path: "/user",
    name: "User",
    meta: {
      title: "User",
    },
    children: [
      {
        path: "login",
        name: "Login",
        component: () => import("@/views/user/login.vue"),
        meta: {
          title: "Login",
          requireAuth: false,
          hideSidebar: true,
        },
      },
      {
        path: "list",
        name: "UserList",
        component: () => import("@/views/user/list.vue"),
        meta: {
          title: "UserList",
          requireAuth: true,
        },
      }
    ],
  },
];

export default user;
