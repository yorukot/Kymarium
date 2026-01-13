import { apiRequest } from "@/lib/api/client";
import type { StatusPageRawData } from "@/lib/schemas/status-page";

export type StatusPageUpsertPayload = {
  title: string;
  slug: string;
  elements?: Array<{
    name: string;
    type: "historical_timeline" | "current_status_indicator";
    sortOrder: number;
    monitor: boolean;
    monitorId?: string;
    monitors?: Array<{
      monitorId: string;
      name: string;
      type: "historical_timeline" | "current_status_indicator";
      sortOrder: number;
    }>;
  }>;
};

export type StatusPageUpsertResponse = {
  message: string;
  data?: StatusPageRawData;
};

export function createStatusPage(teamID: string, payload: StatusPageUpsertPayload) {
  return apiRequest<StatusPageUpsertResponse>(`/api/teams/${teamID}/status-pages`, {
    method: "POST",
    body: payload,
    defaultError: "Create status page failed",
    redirectOn401: true,
    camelcaseResponse: false,
  });
}

export function updateStatusPage(
  teamID: string,
  statusPageID: string,
  payload: StatusPageUpsertPayload,
) {
  return apiRequest<StatusPageUpsertResponse>(
    `/api/teams/${teamID}/status-pages/${statusPageID}`,
    {
      method: "PUT",
      body: payload,
      defaultError: "Update status page failed",
      redirectOn401: true,
      camelcaseResponse: false,
    },
  );
}

