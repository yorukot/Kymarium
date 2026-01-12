import { z } from "zod";

export const incidentStatusValues = [
  "detected",
  "investigating",
  "identified",
  "monitoring",
  "resolved",
] as const;

export const incidentSeverityValues = [
  "emergency",
  "critical",
  "major",
  "minor",
  "info",
] as const;

export const incidentEventTypeValues = [
  "detected",
  "notification_sent",
  "manually_resolved",
  "auto_resolved",
  "unpublished",
  "published",
  "investigating",
  "identified",
  "update",
  "monitoring",
] as const;

export type IncidentStatus = (typeof incidentStatusValues)[number];
export type IncidentSeverity = (typeof incidentSeverityValues)[number];
export type IncidentEventType = (typeof incidentEventTypeValues)[number];

const createIncidentObjectSchema = z.object({
  title: z.string().max(255, "Title must be 255 characters or fewer").optional(),
  status: z.enum(incidentStatusValues),
  severity: z.enum(incidentSeverityValues),
  message: z.string().optional(),
  public: z.boolean(),
  autoResolve: z.boolean(),
  monitorIds: z.array(z.string()).min(1, "Select at least one monitor."),
});

export const createIncidentSchema = createIncidentObjectSchema;

export type CreateIncidentFormValues = z.infer<typeof createIncidentSchema>;

export const createIncidentPayloadSchema = createIncidentObjectSchema;

export type CreateIncidentPayload = z.infer<typeof createIncidentPayloadSchema>;

export type CreateIncidentResponse = {
  message?: string;
  data?: {
    incident?: IncidentRawData;
    event?: unknown;
  };
};

export type IncidentRawData = {
  id: string;
  title?: string | null;
  status: string;
  severity: string;
  is_public: boolean;
  auto_resolve: boolean;
  started_at: string;
  resolved_at?: string | null;
  created_at: string;
  updated_at: string;
};

export type IncidentDetailItem = {
  id: string;
  title?: string;
  status: string;
  severity: string;
  isPublic: boolean;
  autoResolve: boolean;
  startedAt: string;
  resolvedAt?: string;
  createdAt: string;
  updatedAt: string;
};

export type IncidentListItem = {
  id: string;
  title?: string;
  status: string;
  severity: string;
  startedAt: string;
  resolvedAt?: string;
};

export type IncidentEventRawData = {
  id: string;
  incident_id: string;
  created_by?: string | null;
  message: string;
  event_type: string;
  created_at: string;
  updated_at: string;
};

export type IncidentEventItem = {
  id: string;
  incidentId: string;
  createdBy?: string;
  message: string;
  eventType: string;
  createdAt: string;
  updatedAt: string;
};

const incidentSettingsObjectSchema = z.object({
  title: z
    .string()
    .max(255, "Title must be 255 characters or fewer")
    .optional(),
  public: z.boolean(),
  autoResolve: z.boolean(),
});

export const incidentSettingsSchema = incidentSettingsObjectSchema;
export type IncidentSettingsFormValues = z.infer<typeof incidentSettingsSchema>;

export type UpdateIncidentSettingsPayload = {
  title?: string | null;
  public?: boolean;
  autoResolve?: boolean;
};

const incidentStatusUpdateObjectSchema = z.object({
  status: z.enum(incidentStatusValues),
  message: z.string().optional(),
});

export const incidentStatusUpdateSchema = incidentStatusUpdateObjectSchema;
export type IncidentStatusUpdateFormValues = z.infer<
  typeof incidentStatusUpdateSchema
>;

export type UpdateIncidentStatusPayload = {
  status: IncidentStatus;
  message?: string;
  public?: boolean;
};

const incidentEventCreateObjectSchema = z.object({
  message: z.string().min(1, "Message is required"),
  eventType: z.enum(incidentEventTypeValues).optional(),
});

export const incidentEventCreateSchema = incidentEventCreateObjectSchema;
export type IncidentEventCreateFormValues = z.infer<
  typeof incidentEventCreateSchema
>;

export type CreateIncidentEventPayload = {
  message: string;
  eventType?: IncidentEventType;
};
