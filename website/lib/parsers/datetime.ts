import { format, formatDistanceStrict } from "date-fns";

export function isRFC3339Seconds(value: string): boolean {
  return /^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}(Z|[+-]\d{2}:\d{2})$/.test(
    value,
  );
}

export function toRFC3339Seconds(date: Date): string {
  return format(date, "yyyy-MM-dd'T'HH:mm:ssxxx");
}

export function safeParseDate(value: string | null): Date | null {
  if (!value) return null;
  const date = new Date(value);
  if (Number.isNaN(date.getTime())) return null;
  return date;
}

export function formatRelativeTime(value: string): string {
  if (!value) return "--";

  const date = new Date(value);
  if (Number.isNaN(date.getTime())) return "--";

  return formatDistanceStrict(date, new Date(Date.now()), {
    addSuffix: true,
    roundingMethod: "round",
  });
}
