import { defineStore } from "pinia";
import type { LoginRequest, UserResponse } from "@/types";
import { login as loginApi, getUser as getUserApi } from "@/api/user";

export const useUserStore = defineStore(
  "user",
  () => {
    const token = ref<string | null>(null);
    const user = ref<UserResponse | null>(null);
    const currentWorkSpace = ref<string | null>(null);
    const currentWorkApp = ref<string | null>(null);

    const isLogined = computed(() => !!token.value);

    const workSpaceList = computed(() => user.value?.work_space_list || []);
    const workAppList = computed(
      () =>
        workSpaceList.value.find(
          (workSpace) => workSpace.id === currentWorkSpace.value
        )?.app_list
    );

    const login = async (data: LoginRequest): Promise<void> => {
      const res = await loginApi(data);
      token.value = res.token;

      const userInfoRes = await getUserApi();
      user.value = userInfoRes;

      currentWorkSpace.value = userInfoRes.work_space_list[0]?.id || null;
      currentWorkApp.value =
        userInfoRes.work_space_list[0]?.app_list[0]?.id || null;
    };

    const logout = () => {
      token.value = null;
      user.value = null;
      currentWorkSpace.value = null;
      currentWorkApp.value = null;
    };

    const changeWorkSpace = (workSpaceId?: string, appId?: string) => {
      if (workSpaceId) {
        currentWorkSpace.value = workSpaceId;
      }
      if (appId) {
        currentWorkApp.value = appId;
      }
    };

    return {
      token,
      user,
      isLogined,
      login,
      logout,
      currentWorkSpace,
      currentWorkApp,
      workSpaceList,
      workAppList,
      changeWorkSpace,
    };
  },
  {
    persist: true,
  }
);
