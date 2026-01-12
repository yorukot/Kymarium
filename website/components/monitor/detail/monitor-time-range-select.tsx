"use client";

import * as React from "react";
import { format, subDays, subHours } from "date-fns";
import { usePathname, useRouter, useSearchParams } from "next/navigation";

import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";
import { cn } from "@/lib/utils";
import { safeParseDate, toRFC3339Seconds } from "@/lib/parsers/datetime";

type PresetKey = "12h" | "24h" | "3d" | "7d" | "14d" | "30d" | "90d";
const defaultPreset: PresetKey = "24h";

const presets: Array<{ key: PresetKey; label: string }> = [
  { key: "12h", label: "12hour" },
  { key: "24h", label: "24hour" },
  { key: "3d", label: "3d" },
  { key: "7d", label: "7d" },
  { key: "14d", label: "14d" },
  { key: "30d", label: "30d" },
  { key: "90d", label: "90d" },
];

function isPresetKey(value: string): value is PresetKey {
  return presets.some((p) => p.key === value);
}

function applyPresetRange(key: PresetKey): { start: Date; end: Date } {
  const end = new Date();
  end.setMilliseconds(0);

  switch (key) {
    case "12h": {
      const start = subHours(end, 12);
      start.setMilliseconds(0);
      return { start, end };
    }
    case "24h": {
      const start = subHours(end, 24);
      start.setMilliseconds(0);
      return { start, end };
    }
    case "3d":
    case "7d":
    case "14d":
    case "30d":
    case "90d": {
      const days = Number(key.replace("d", ""));
      const start = subDays(end, days);
      start.setMilliseconds(0);
      return { start, end };
    }
  }
}

export function MonitorTimeRangeSelect({ className }: { className?: string }) {
  const router = useRouter();
  const pathname = usePathname();
  const searchParams = useSearchParams();

  const presetRaw = searchParams.get("preset") || "";
  const preset = isPresetKey(presetRaw) ? presetRaw : "";
  const start = safeParseDate(searchParams.get("start"));
  const end = safeParseDate(searchParams.get("end"));

  const setQueryRange = React.useCallback(
    ({
      nextStart,
      nextEnd,
      nextPreset,
    }: {
      nextStart: Date | null;
      nextEnd: Date | null;
      nextPreset?: PresetKey | null;
    }) => {
      const params = new URLSearchParams(searchParams.toString());

      if (nextStart && nextEnd) {
        params.set("start", toRFC3339Seconds(nextStart));
        params.set("end", toRFC3339Seconds(nextEnd));
      } else {
        params.delete("start");
        params.delete("end");
      }

      if (nextPreset) {
        params.set("preset", nextPreset);
      } else {
        params.delete("preset");
      }

      const nextUrl =
        params.size > 0 ? `${pathname}?${params.toString()}` : pathname;
      router.replace(nextUrl, { scroll: false });
      router.refresh();
    },
    [pathname, router, searchParams],
  );

  React.useEffect(() => {
    // Keep frontend default aligned with backend (default last 24h) so the UI is not empty
    // and shareable URLs always include a concrete window.
    if (preset) {
      if (!start || !end) {
        const nextRange = applyPresetRange(preset);
        setQueryRange({
          nextStart: nextRange.start,
          nextEnd: nextRange.end,
          nextPreset: preset,
        });
      }
      return;
    }

    if (!start && !end) {
      const nextRange = applyPresetRange(defaultPreset);
      setQueryRange({
        nextStart: nextRange.start,
        nextEnd: nextRange.end,
        nextPreset: defaultPreset,
      });
    }
  }, [end, preset, setQueryRange, start]);

  return (
    <div className={cn("flex flex-wrap items-center gap-2", className)}>
      <Select
        value={preset || undefined}
        onValueChange={(value) => {
          const key = value as PresetKey;
          const { start, end } = applyPresetRange(key);
          setQueryRange({ nextStart: start, nextEnd: end, nextPreset: key });
        }}
      >
        <SelectTrigger className="min-w-40" size="sm">
          <SelectValue placeholder="Select range" />
        </SelectTrigger>
        <SelectContent align="start">
          {presets.map((p) => (
            <SelectItem key={p.key} value={p.key}>
              {p.label}
            </SelectItem>
          ))}
        </SelectContent>
      </Select>
      <span className="text-xs text-muted-foreground">
        {start && end ? `(${format(start, "PP p")} → ${format(end, "PP p")})` : ""}
      </span>
    </div>
  );
}
