import type {
  Notification,
  NotificationRawData,
  NotificationType,
} from "@/lib/schemas/notification";
import { isPlainObject } from "@/lib/parsers/guards";

const NOTIFICATION_TYPE_LABELS: Record<NotificationType, string> = {
  email: "Email",
  slack: "Slack",
  discord: "Discord",
  telegram: "Telegram",
};

function safeUrlHostname(url: string): string | null {
  try {
    return new URL(url).hostname;
  } catch {
    return null;
  }
}

function buildNotificationDetail(
  type: string,
  config: unknown,
): string | undefined {
  if (!isPlainObject(config)) return undefined;

  if (type === "email") {
    const emails = config.email_address;
    const count = Array.isArray(emails) ? emails.filter(Boolean).length : 0;
    return count > 0 ? `${count} recipient${count === 1 ? "" : "s"}` : undefined;
  }

  if (type === "slack" || type === "discord") {
    const webhookUrl = config.webhook_url;
    if (typeof webhookUrl !== "string" || webhookUrl.trim().length === 0) {
      return undefined;
    }
    const hostname = safeUrlHostname(webhookUrl);
    return hostname ? `Webhook: ${hostname}` : "Webhook configured";
  }

  if (type === "telegram") {
    const chatId = config.chat_id;
    if (typeof chatId !== "string" || chatId.trim().length === 0) {
      return "Telegram configured";
    }
    const trimmed = chatId.trim();
    const suffix = trimmed.length > 4 ? trimmed.slice(-4) : trimmed;
    return `Chat ID: …${suffix}`;
  }

  return undefined;
}

export function parseNotification(raw: NotificationRawData): Notification {
  const type = raw.type as NotificationType;
  const typeLabel = NOTIFICATION_TYPE_LABELS[type] ?? raw.type;

  return {
    id: raw.id,
    type: raw.type,
    name: raw.name,
    typeLabel,
    detail: buildNotificationDetail(raw.type, raw.config),
    createdAt: raw.created_at,
    updatedAt: raw.updated_at,
  };
}

export function parseNotifications(
  rawList: NotificationRawData[],
): Notification[] {
  return rawList.map(parseNotification);
}
