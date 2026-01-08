import { apiRequest } from "@/lib/api/client";
import type {
  MonitorFormValues,
  HttpMonitorFormConfig,
  PingMonitorFormConfig,
  CreateMonitorResponse,
} from "@/lib/schemas/monitor";

function normalizeHeaders(headers: Array<{ key: string; value: string }> = []) {
  const entries = headers
    .map((header) => ({
      key: header.key.trim(),
      value: header.value.trim(),
    }))
    .filter((header) => header.key.length > 0);

  if (entries.length === 0) {
    return undefined;
  }

  return entries.reduce<Record<string, string>>((acc, header) => {
    acc[header.key] = header.value;
    return acc;
  }, {});
}

function uniqueStrings(values: string[]) {
  if (values.length <= 1) {
    return values;
  }

  return Array.from(new Set(values));
}

type HeaderKV = { key: string; value: string };

type CreateMonitorPayloadBase = {
  name: string;
  interval: number;
  failureThreshold: number;
  recoveryThreshold: number;
  regions: string[];
  notifications: string[];
};

type CreateMonitorPayload =
  | (CreateMonitorPayloadBase & {
      type: "http";
      config: HttpMonitorFormConfig;
    })
  | (CreateMonitorPayloadBase & {
      type: "ping";
      config: PingMonitorFormConfig;
    });

function recordToHeaderArray(headers?: Record<string, string>): HeaderKV[] {
  return headers
    ? Object.entries(headers).map(([key, value]) => ({ key, value }))
    : [];
}

function buildHttpConfig(values: MonitorFormValues): HttpMonitorFormConfig {
  const http = values.http;
  if (!http) {
    throw new Error("HTTP config is required for HTTP monitors.");
  }

  const headers = normalizeHeaders(http.headers);

  const config: HttpMonitorFormConfig = {
    url: http.url,
    method: http.method,
    maxRedirects: http.maxRedirects,
    requestTimeout: http.requestTimeout,
    upsideDownMode: http.upsideDownMode,
    certificateExpiryNotification: http.certificateExpiryNotification,
    headers: recordToHeaderArray(headers),
    acceptedStatusCodes: http.acceptedStatusCodes ?? [],
    body: http.body ?? "",
    ignoreTLSError: http.ignoreTLSError,
    ...(http.bodyEncoding ? { bodyEncoding: http.bodyEncoding } : {}),
  };

  return config;
}

function buildPingConfig(values: MonitorFormValues): PingMonitorFormConfig {
  const ping = values.ping;
  if (!ping) {
    throw new Error("Ping config is required for ping monitors.");
  }

  return {
    host: ping.host,
    timeoutSeconds: ping.timeoutSeconds,
    ...(ping.packetSize ? { packetSize: ping.packetSize } : {}),
  };
}

export function buildCreateMonitorPayload(
  values: MonitorFormValues,
): CreateMonitorPayload {
  const basePayload: CreateMonitorPayloadBase = {
    name: values.name,
    interval: values.interval,
    failureThreshold: values.failureThreshold,
    recoveryThreshold: values.recoveryThreshold,
    regions: uniqueStrings(values.regions),
    notifications: uniqueStrings(values.notifications),
  };

  if (values.type === "http") {
    return {
      ...basePayload,
      type: "http",
      config: buildHttpConfig(values),
    };
  }

  return {
    ...basePayload,
    type: "ping",
    config: buildPingConfig(values),
  };
}

export function createMonitor(teamID: string, values: MonitorFormValues) {
  const payload = buildCreateMonitorPayload(values);

  return apiRequest<CreateMonitorResponse>(`/api/teams/${teamID}/monitors`, {
    method: "POST",
    body: payload,
    defaultError: "Create monitor failed",
    redirectOn401: true,
  });
}
