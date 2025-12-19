import http from "./index";
import type { PageResponse, ClientResponse, ClientBindResponse, ClientSystemInfoResponse } from "@/types";


export async function getClientPage(): Promise<PageResponse<ClientResponse>> {
  return http.get("/v1/client/page");
}

export async function getClientBind(): Promise<ClientBindResponse> {
  return http.get("/v1/client/bind");
}

export async function deleteClient(id: number): Promise<void> {
  return http.delete(`/v1/client/${id}`);
}

export async function getClientSystemInfo(id: number): Promise<ClientSystemInfoResponse[]> {
  return http.get(`/v1/client/${id}/stat`);
}