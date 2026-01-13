import { z } from "zod";

export const userProfileSchema = z.object({
  displayName: z
    .string()
    .trim()
    .min(1, "Display name is required")
    .max(255, "Display name is too long"),
  avatar: z.string().max(2048, "Avatar URL is too long"),
});

export type UserProfileFormValues = z.infer<typeof userProfileSchema>;

export const userProfilePayloadSchema = userProfileSchema;

export type UserProfilePayload = z.infer<typeof userProfilePayloadSchema>;

const passwordUpdateObjectSchema = z.object({
  currentPassword: z
    .string()
    .min(8, "Password must be at least 8 characters")
    .max(255, "Password is too long"),
  newPassword: z
    .string()
    .min(8, "Password must be at least 8 characters")
    .max(255, "Password is too long"),
  confirmNewPassword: z.string(),
});

export const passwordUpdateSchema = passwordUpdateObjectSchema.superRefine(
  ({ newPassword, confirmNewPassword }, ctx) => {
    if (newPassword !== confirmNewPassword) {
      ctx.addIssue({
        code: z.ZodIssueCode.custom,
        message: "Passwords do not match",
        path: ["confirmNewPassword"],
      });
    }
  },
);

export type PasswordUpdateFormValues = z.infer<typeof passwordUpdateSchema>;

export const passwordUpdatePayloadSchema = passwordUpdateObjectSchema.pick({
  currentPassword: true,
  newPassword: true,
});

export type PasswordUpdatePayload = z.infer<typeof passwordUpdatePayloadSchema>;
