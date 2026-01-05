"use client";

import dynamic from "next/dynamic";
import type { Control } from "react-hook-form";

const DevTool = dynamic(
  () => import("@hookform/devtools").then((m) => m.DevTool),
  { ssr: false }
);

// eslint-disable-next-line @typescript-eslint/no-explicit-any
export function FormDevTools({ control }: { control: Control<any> }) {
  if (process.env.NODE_ENV !== "development") return null;
  return <DevTool control={control} />;
}
