import { isPlainObject } from "@/lib/parsers/guards";
import {
  notificationTypeValues,
  type NotificationFormValues,
} from "@/lib/schemas/notification";

export function defaultNotificationFormValues(
  type: NotificationFormValues["type"],
): NotificationFormValues {
  if (type === "email") {
    return {
      type,
      name: "",
      config: { emailAddresses: [{ value: "" }] },
    };
  }

  if (type === "telegram") {
    return {
      type,
      name: "",
      config: { botToken: "", chatId: "" },
    };
  }

  return {
    type,
    name: "",
    config: { webhookUrl: "" },
  };
}

export function isNotificationType(
  value: unknown,
): value is NotificationFormValues["type"] {
  return notificationTypeValues.includes(value as NotificationFormValues["type"]);
}

export function notificationFormValuesFromApi(
  type: NotificationFormValues["type"],
  name: unknown,
  config: unknown,
): NotificationFormValues {
  const safeName = typeof name === "string" ? name : "";

  if (!isPlainObject(config)) {
    return { ...defaultNotificationFormValues(type), name: safeName };
  }

  if (type === "email") {
    const rawEmails = config.emailAddress;
    const emails = Array.isArray(rawEmails)
      ? rawEmails.map((value) => String(value)).filter(Boolean)
      : [];

    return {
      type,
      name: safeName,
      config: {
        emailAddresses:
          emails.length > 0 ? emails.map((value) => ({ value })) : [{ value: "" }],
      },
    };
  }

  if (type === "telegram") {
    return {
      type,
      name: safeName,
      config: {
        botToken: typeof config.botToken === "string" ? config.botToken : "",
        chatId: typeof config.chatId === "string" ? config.chatId : "",
      },
    };
  }

  return {
    type,
    name: safeName,
    config: {
      webhookUrl: typeof config.webhookUrl === "string" ? config.webhookUrl : "",
    },
  };
}

