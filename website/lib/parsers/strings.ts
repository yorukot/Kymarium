function capitalize(part: string): string {
  const first = part[0];
  if (!first) return "";
  return first.toUpperCase() + part.slice(1);
}

export function humanizeIdentifier(value: string): string {
  if (!value) return "";
  return value
    .split(/[-_]/)
    .filter(Boolean)
    .map(capitalize)
    .join(" ");
}

