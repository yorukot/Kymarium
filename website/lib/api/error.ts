import type { FieldValues, Path, UseFormSetError } from "react-hook-form";

function isRecord(value: unknown): value is Record<string, unknown> {
  return typeof value === "object" && value !== null;
}

export function applyServerFieldErrors<T extends FieldValues>(
  setError: UseFormSetError<T>,
  body: unknown
): boolean {
  if (!isRecord(body)) return false;

  const errors = body["errors"];
  if (!isRecord(errors)) return false;

  let applied = false;

  for (const [key, value] of Object.entries(errors)) {
    if (typeof value !== "string" || !value) continue;

    setError(key as Path<T>, { type: "server", message: value });
    applied = true;
  }

  return applied;
}
