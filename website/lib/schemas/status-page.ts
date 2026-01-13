export type StatusPageElementType =
  | "historical_timeline"
  | "current_status_indicator";

export type StatusPageModelRawData = {
  id: string;
  team_id: string;
  title: string;
  slug: string;
  icon?: string | null;
  created_at: string;
  updated_at: string;
};

export type StatusPageMonitorRawData = {
  id: string;
  status_page_id: string;
  monitor_id: string;
  group_id?: string | null;
  name: string;
  type: StatusPageElementType;
  sort_order: number;
};

export type StatusPageElementRawData = {
  id: string;
  status_page_id: string;
  name: string;
  type: StatusPageElementType;
  sort_order: number;
  monitor: boolean;
  monitor_id?: string | null;
  monitors: StatusPageMonitorRawData[];
};

export type StatusPageRawData = {
  status_page: StatusPageModelRawData;
  elements: StatusPageElementRawData[];
};

export type StatusPageListItem = {
  id: string;
  title: string;
  slug: string;
  createdAt: string;
  updatedAt: string;
  elementCount: number;
  monitorCount: number;
};

export type StatusPageDetailItem = {
  id: string;
  teamId: string;
  title: string;
  slug: string;
  createdAt: string;
  updatedAt: string;
  elements: StatusPageElementRawData[];
  elementCount: number;
  monitorCount: number;
};

