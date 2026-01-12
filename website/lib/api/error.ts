import type { FieldValues, Path, UseFormSetError } from "react-hook-form";
import { isPlainObject } from "@/lib/parsers/guards";

export function applyServerFieldErrors<T extends FieldValues>(
  setError: UseFormSetError<T>,
  body: unknown
): boolean {
  if (!isPlainObject(body)) return false;

  const errors = body["errors"];
  if (!isPlainObject(errors)) return false;

  let applied = false;

  for (const [key, value] of Object.entries(errors)) {
    if (typeof value !== "string" || !value) continue;

    setError(key as Path<T>, { type: "server", message: value });
    applied = true;
  }

  return applied;
}
