export function formatUptime(value?: number | null) {
  if (value === undefined || value === null || Number.isNaN(value)) return "--";
  return `${value.toFixed(2)}%`;
}

export function formatCount(value?: number | null) {
  if (value === undefined || value === null || Number.isNaN(value)) return "--";
  return new Intl.NumberFormat("en", { notation: "compact" }).format(value);
}

export function formatLatencyMs(value?: number | null) {
  if (value === undefined || value === null || Number.isNaN(value)) return "--";
  return `${Math.round(value)}ms`;
}

