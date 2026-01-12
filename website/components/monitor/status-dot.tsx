import { cn } from "@/lib/utils";

type StatusDotTone = "successed" | "destructive";

function toneFromStatus(status?: string | null): StatusDotTone {
  return status === "up" ? "successed" : "destructive";
}

export function StatusDot({
  status,
  className,
}: {
  status?: string | null;
  className?: string;
  }) {
  const tone = toneFromStatus(status);
  const outer =
    tone === "successed" ? "bg-successed/40" : "bg-destructive/40";
  const inner =
    tone === "successed" ? "bg-successed/70" : "bg-destructive/70";

  return (
    <span
      className={cn(
        "inline-flex h-4 w-4 flex-none items-center justify-center rounded-full",
        outer,
        className,
      )}
      aria-label={status ?? undefined}
      title={status ?? undefined}
    >
      <span className={cn("h-2 w-2 rounded-full", inner)} />
    </span>
  );
}
