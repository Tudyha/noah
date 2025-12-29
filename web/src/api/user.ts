import http from "./index";
import type { LoginRequest, LoginResponse, UserResponse } from "@/types";

export async function login(data: LoginRequest): Promise<LoginResponse> {
  return http.post("/v1/auth/login", data);
}

export async function sendCode(data: { type: number, target: string }): Promise<any> {
  return http.post("/v1/auth/send_code", data);
}

export async function getUser(): Promise<UserResponse> {
  return http.get("/v1/user");
}
