import type {
  StatusPageDetailItem,
  StatusPageListItem,
  StatusPageRawData,
} from "@/lib/schemas/status-page";

function countMonitors(raw: StatusPageRawData): number {
  let total = 0;
  for (const element of raw.elements ?? []) {
    if (element.monitor) {
      total += 1;
      continue;
    }
    total += element.monitors?.length ?? 0;
  }
  return total;
}

export function parseStatusPage(raw: StatusPageRawData): StatusPageListItem {
  return {
    id: raw.status_page.id,
    title: raw.status_page.title,
    slug: raw.status_page.slug,
    createdAt: raw.status_page.created_at,
    updatedAt: raw.status_page.updated_at,
    elementCount: raw.elements?.length ?? 0,
    monitorCount: countMonitors(raw),
  };
}

export function parseStatusPages(rawList: StatusPageRawData[]): StatusPageListItem[] {
  return rawList.map(parseStatusPage);
}

export function parseStatusPageDetail(raw: StatusPageRawData): StatusPageDetailItem {
  return {
    id: raw.status_page.id,
    teamId: raw.status_page.team_id,
    title: raw.status_page.title,
    slug: raw.status_page.slug,
    createdAt: raw.status_page.created_at,
    updatedAt: raw.status_page.updated_at,
    elements: raw.elements ?? [],
    elementCount: raw.elements?.length ?? 0,
    monitorCount: countMonitors(raw),
  };
}

