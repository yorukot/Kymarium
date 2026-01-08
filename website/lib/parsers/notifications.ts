import type {
  Notification,
  NotificationRawData,
  NotificationType,
} from "@/lib/schemas/notification";

const NOTIFICATION_TYPE_LABELS: Record<NotificationType, string> = {
  email: "Email",
  slack: "Slack",
  discord: "Discord",
  telegram: "Telegram",
};

export function parseNotification(raw: NotificationRawData): Notification {
  const type = raw.type as NotificationType;
  const typeLabel = NOTIFICATION_TYPE_LABELS[type] ?? raw.type;

  return {
    id: raw.id,
    type: raw.type,
    name: raw.name,
    typeLabel,
  };
}

export function parseNotifications(
  rawList: NotificationRawData[],
): Notification[] {
  return rawList.map(parseNotification);
}
