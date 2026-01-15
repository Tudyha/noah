import http from "./index";
import type { DashboardResponse } from "@/types";

export async function getDashboard(): Promise<DashboardResponse> {
  return http.get("/v1/dashboard");
}
