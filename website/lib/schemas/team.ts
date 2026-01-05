import { z } from "zod";

export const teamObjectSchema = z.object({
  name: z
    .string()
    .min(1, "Team name is required")
    .max(255, "Team name is too long"),
});

export const teamSchema = teamObjectSchema;

export type TeamFormValues = z.infer<typeof teamSchema>;

export const teamPayloadSchema = teamObjectSchema;

export type TeamPayload = z.infer<typeof teamPayloadSchema>;
