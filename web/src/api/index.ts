import axios, {
  type AxiosInstance,
  type AxiosRequestConfig,
  type AxiosResponse,
} from "axios";
import router from "@/router";
import { useUserStore } from "@/stores/auth";
import type { BaseResponse } from "@/types";

const defaultConfig: AxiosRequestConfig = {
  baseURL: import.meta.env.VITE_API_BASE_URL,
  timeout: 5000,
};

class HttpService {
  private instance: AxiosInstance;

  constructor(config: AxiosRequestConfig = defaultConfig) {
    this.instance = axios.create(config);
    this.httpInterceptorRequest();
    this.httpInterceptorResponse();
  }

  private httpInterceptorRequest() {
    this.instance.interceptors.request.use(
      (config) => {
        // 在发送请求之前做些什么
        const userStore = useUserStore();
        // 获取token鉴权
        if (userStore.token) {
          // 有token，在请求头中携带token
          config.headers.Authorization = "Bearer " + userStore.token;
        }
        if (userStore.currentWorkSpace) {
          config.headers.WorkSpaceId = userStore.currentWorkSpace;
        }
        if (userStore.currentWorkApp) {
          config.headers['app-id'] = userStore.currentWorkApp;
        }
        return config;
      },
      (error: any) => {
        // 对请求错误做些什么
        return Promise.reject(error);
      }
    );
  }

  private httpInterceptorResponse() {
    this.instance.interceptors.response.use(
      (response: AxiosResponse<BaseResponse<unknown>>) => {
        const res = response.data;
        if (res.code !== 0) {
          if (res.code == 401) {
            // 登录过期
            const userStore = useUserStore();
            userStore.logout();
            router.push("/login");
          }
          return Promise.reject(res.msg || "Error");
        }
        return res.data as any;
      },
      (error) => {
        return Promise.reject(error);
      }
    );
  }

  public get<T>(
    url: string,
    parame?: any,
    config?: AxiosRequestConfig
  ): Promise<T> {
    return this.instance.get(url, { params: parame, ...config });
  }

  public post<T>(
    url: string,
    data?: any,
    config?: AxiosRequestConfig
  ): Promise<T> {
    return this.instance.post(url, data, config);
  }
}

const httpService = new HttpService();
export default httpService;
