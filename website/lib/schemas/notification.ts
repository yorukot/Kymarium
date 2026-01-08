export type NotificationType = "email" | "slack" | "discord" | "telegram";

export type NotificationRawData = {
  id: string;
  team_id: string;
  type: string;
  name: string;
};

export type Notification = {
  id: string;
  type: NotificationType | string;
  name: string;
  typeLabel: string;
};
