import { SiDiscord, SiGmail, SiSlack, SiTelegram } from "react-icons/si";
import type { IconType } from "react-icons/lib";

import type { NotificationType } from "@/lib/schemas/notification";

export type NotificationTypeMeta = {
  type: NotificationType;
  label: string;
  description: string;
  Icon: IconType;
};

const NOTIFICATION_TYPE_META: Record<NotificationType, NotificationTypeMeta> = {
  email: {
    type: "email",
    label: "Email",
    description: "Send incident alerts to one or more email recipients.",
    Icon: SiGmail,
  },
  slack: {
    type: "slack",
    label: "Slack",
    description: "Send incident alerts to a Slack Incoming Webhook.",
    Icon: SiSlack,
  },
  discord: {
    type: "discord",
    label: "Discord",
    description: "Send incident alerts to a Discord webhook.",
    Icon: SiDiscord,
  },
  telegram: {
    type: "telegram",
    label: "Telegram",
    description: "Send incident alerts to a Telegram chat via a bot.",
    Icon: SiTelegram,
  },
};

export function getNotificationTypeMeta(type: NotificationType): NotificationTypeMeta {
  return NOTIFICATION_TYPE_META[type];
}

