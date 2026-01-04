import http from "./index";
import type { PageResponse, ClientResponse, ClientBindResponse, ClientSystemInfoResponse } from "@/types";


export async function getClientPage(): Promise<PageResponse<ClientResponse>> {
  return http.get("/v1/agent/page");
}

export async function getClientBind(): Promise<ClientBindResponse> {
  return http.get("/v1/agent/bind");
}

export async function deleteClient(id: number): Promise<void> {
  return http.delete(`/v1/agent/${id}`);
}

export async function getClientDetail(id: string): Promise<ClientResponse> {
  return http.get(`/v1/agent/${id}`);
}

export async function getClientSystemInfo(id: number): Promise<ClientSystemInfoResponse[]> {
  return http.get(`/v1/agent/${id}/metric`);
}

export async function getV2raySubscribe(ids: number[]): Promise<string> {
  return http.post(`/v1/agent/v2ray/sub`, {ids: ids});
}