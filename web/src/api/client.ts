import http from "./index";
import type { PageResponse, ClientResponse, ClientBindResponse } from "@/types";


export async function getClientPage(): Promise<PageResponse<ClientResponse>> {
  return http.get("/v1/client/page");
}

export async function getClientBind(): Promise<ClientBindResponse> {
  return http.get("/v1/client/bind");
}
