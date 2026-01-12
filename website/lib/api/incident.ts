import { apiRequest } from "@/lib/api/client";
import type {
  CreateIncidentEventPayload,
  CreateIncidentPayload,
  CreateIncidentResponse,
  IncidentEventRawData,
  IncidentRawData,
  UpdateIncidentSettingsPayload,
  UpdateIncidentStatusPayload,
} from "@/lib/schemas/incident";

export function createIncident(teamID: string, payload: CreateIncidentPayload) {
  return apiRequest<CreateIncidentResponse>(`/api/teams/${teamID}/incidents`, {
    method: "POST",
    body: payload,
    defaultError: "Create incident failed",
    redirectOn401: true,
  });
}

export type GetIncidentResponse = {
  message?: string;
  data?: IncidentRawData;
};

export type ListIncidentEventsResponse = {
  message?: string;
  data?: IncidentEventRawData[];
};

export function updateIncidentStatus(
  teamID: string,
  incidentID: string,
  payload: UpdateIncidentStatusPayload,
) {
  return apiRequest<{ message?: string; data?: unknown }>(
    `/api/teams/${teamID}/incidents/${incidentID}/status`,
    {
      method: "POST",
      body: payload,
      defaultError: "Update incident status failed",
      redirectOn401: true,
    },
  );
}

export function createIncidentEvent(
  teamID: string,
  incidentID: string,
  payload: CreateIncidentEventPayload,
) {
  return apiRequest<{ message?: string; data?: unknown }>(
    `/api/teams/${teamID}/incidents/${incidentID}/events`,
    {
      method: "POST",
      body: payload,
      defaultError: "Create incident update failed",
      redirectOn401: true,
    },
  );
}

export function updateIncidentSettings(
  teamID: string,
  incidentID: string,
  payload: UpdateIncidentSettingsPayload,
) {
  return apiRequest<{ message?: string; data?: unknown }>(
    `/api/teams/${teamID}/incidents/${incidentID}`,
    {
      method: "PATCH",
      body: payload,
      defaultError: "Update incident failed",
      redirectOn401: true,
    },
  );
}
