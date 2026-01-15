import http from "./index";
import type { PageResponse, AgentResponse, AgentBindResponse, AgentSystemInfoResponse } from "@/types";


export async function getAgentPage(params: Record<string, any>): Promise<PageResponse<AgentResponse>> {
  return http.get("/v1/agent/page", params);
}

export async function getAgentBind(): Promise<AgentBindResponse> {
  return http.get("/v1/agent/bind");
}

export async function deleteAgent(id: number): Promise<void> {
  return http.delete(`/v1/agent/${id}`);
}

export async function getAgentDetail(id: string): Promise<AgentResponse> {
  return http.get(`/v1/agent/${id}`);
}

export async function getAgentSystemInfo(id: number): Promise<AgentSystemInfoResponse[]> {
  return http.get(`/v1/agent/${id}/metric`);
}

export async function getV2raySubscribe(ids: number[]): Promise<string> {
  return http.post(`/v1/agent/v2ray/sub`, {ids: ids});
}
