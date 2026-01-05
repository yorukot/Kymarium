"use client"

import { useSelectedLayoutSegments } from "next/navigation"
import { BreadcrumbPage } from "@/components/ui/breadcrumb"

const LABELS: Record<string, string> = {
  monitors: "Monitors",
  incidents: "Incidents",
}

function toTitle(segment: string) {
  return segment
    .split(/[-_]/)
    .filter(Boolean)
    .map((part) => part[0]?.toUpperCase() + part.slice(1))
    .join(" ")
}

export function TeamBreadcrumbs() {
  const segments = useSelectedLayoutSegments()
  const section = segments[0]
  const label = section ? (LABELS[section] ?? toTitle(section)) : "Monitors"

  return <BreadcrumbPage>{label}</BreadcrumbPage>
}
