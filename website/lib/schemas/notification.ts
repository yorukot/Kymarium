import { z } from "zod";

export type NotificationType = "email" | "slack" | "discord" | "telegram";

export type NotificationRawData = {
  id: string;
  team_id: string;
  type: string;
  name: string;
  config?: unknown;
  updated_at?: string;
  created_at?: string;
};

export type Notification = {
  id: string;
  type: NotificationType | string;
  name: string;
  typeLabel: string;
  detail?: string;
  createdAt?: string;
  updatedAt?: string;
};

export const notificationTypeValues = [
  "email",
  "slack",
  "discord",
  "telegram",
] as const;

const emailConfigSchema = z.object({
  emailAddresses: z
    .array(
      z.object({
        value: z
          .string()
          .trim()
          .email("Enter a valid email address."),
      }),
    )
    .min(1, "Add at least one email address."),
});

const webhookConfigSchema = z.object({
  webhookUrl: z.string().trim().url("Enter a valid webhook URL."),
});

const telegramConfigSchema = z.object({
  botToken: z.string().trim().min(1, "Bot token is required."),
  chatId: z.string().trim().min(1, "Chat ID is required."),
});

const baseNotificationSchema = z.object({
  type: z.enum(notificationTypeValues),
  name: z.string().trim().min(1, "Notification name is required."),
});

export const notificationSchema = z.discriminatedUnion("type", [
  baseNotificationSchema.extend({
    type: z.literal("email"),
    config: emailConfigSchema,
  }),
  baseNotificationSchema.extend({
    type: z.literal("slack"),
    config: webhookConfigSchema,
  }),
  baseNotificationSchema.extend({
    type: z.literal("discord"),
    config: webhookConfigSchema,
  }),
  baseNotificationSchema.extend({
    type: z.literal("telegram"),
    config: telegramConfigSchema,
  }),
]);

export type NotificationFormValues = z.infer<typeof notificationSchema>;
