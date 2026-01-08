import type { MonitorListItem, MonitorRawData } from "@/lib/schemas/monitor";

type MonitorConfig = Record<string, unknown> | null | undefined;

type MonitorTarget = {
  label: string;
  value: string;
};

function getConfigValue(config: MonitorConfig, key: string): string {
  if (!config || typeof config !== "object") {
    return "";
  }

  const value = config[key as keyof typeof config];
  return typeof value === "string" ? value : "";
}

function parseMonitorTarget(raw: MonitorRawData): MonitorTarget {
  const config = raw.config as MonitorConfig;

  if (raw.type === "http") {
    return {
      label: "URL",
      value: getConfigValue(config, "url"),
    };
  }

  if (raw.type === "ping") {
    return {
      label: "Host",
      value: getConfigValue(config, "host"),
    };
  }

  return {
    label: "Target",
    value: "",
  };
}

export function parseMonitor(raw: MonitorRawData): MonitorListItem {
  const target = parseMonitorTarget(raw);

  return {
    id: raw.id,
    name: raw.name,
    type: raw.type,
    status: raw.status,
    targetLabel: target.label,
    targetValue: target.value,
    lastChecked: raw.last_checked,
    uptimeSLI30: raw.uptime_sli_30 ?? undefined,
  };
}

export function parseMonitors(rawList: MonitorRawData[]): MonitorListItem[] {
  return rawList.map(parseMonitor);
}
