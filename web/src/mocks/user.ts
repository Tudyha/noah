import type { MockMethod } from "vite-plugin-mock";
export default [
  {
    url: "/api/v1/user/login",
    method: "post",
    response: {
      code: 0,
      msg: "success",
      data: {
        token: "mock-token",
      },
    },
  },
  {
    url: "/api/v1/user",
    method: "get",
    response: {
      code: 0,
      msg: "success",
      data: {
        nickname: "admin",
        avatar:
          "https://wpimg.wallstcn.com/f778738c-e4f8-4870-b634-56703b4acafe.gif",
        roles: ["admin"],
      },
    },
  },
] as MockMethod[];
