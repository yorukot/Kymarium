import { z } from "zod";

export const monitorSchema = z.object({
  name: z.string().min(1, "Monitor name is required."),
  interval: z
    .number()
    .int()
    .min(30, "Interval must be at least 30 seconds.")
    .max(2592000, "Interval must be 30 days or less."),
  failureThreshold: z
    .number()
    .int()
    .min(1, "Failure threshold must be at least 1."),
  recoveryThreshold: z
    .number()
    .int()
    .min(1, "Recovery threshold must be at least 1."),
  regions: z.array(z.number()).min(1, "Select at least one region."),
  notifications: z.array(z.number()),
});

export type MonitorFormValues = z.infer<typeof monitorSchema>;
