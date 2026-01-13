import { apiRequest } from "@/lib/api/client";
import type { NotificationFormValues } from "@/lib/schemas/notification";

export type NotificationApiData = {
  id: string;
  teamId: string;
  type: string;
  name: string;
  config?: unknown;
  updatedAt?: string;
  createdAt?: string;
};

type CreateNotificationPayload = {
  type: NotificationFormValues["type"];
  name: string;
  config: unknown;
};

function uniqueStrings(values: string[]) {
  if (values.length <= 1) return values;
  return Array.from(new Set(values));
}

export function buildCreateNotificationPayload(
  values: NotificationFormValues,
): CreateNotificationPayload {
  if (values.type === "email") {
    const emails = uniqueStrings(
      values.config.emailAddresses
        .map((entry) => entry.value.trim())
        .filter(Boolean),
    );

    return {
      type: values.type,
      name: values.name.trim(),
      config: {
        emailAddress: emails,
      },
    };
  }

  if (values.type === "telegram") {
    return {
      type: values.type,
      name: values.name.trim(),
      config: {
        botToken: values.config.botToken.trim(),
        chatId: values.config.chatId.trim(),
      },
    };
  }

  return {
    type: values.type,
    name: values.name.trim(),
    config: {
      webhookUrl: values.config.webhookUrl.trim(),
    },
  };
}

export function createNotification(teamID: string, values: NotificationFormValues) {
  const payload = buildCreateNotificationPayload(values);

  return apiRequest<{ message?: string; data?: unknown }>(
    `/api/teams/${teamID}/notifications`,
    {
      method: "POST",
      body: payload,
      defaultError: "Create notification failed",
      redirectOn401: true,
    },
  );
}

export function getNotification(teamID: string, notificationID: string) {
  return apiRequest<{ message?: string; data?: NotificationApiData }>(
    `/api/teams/${teamID}/notifications/${notificationID}`,
    {
      method: "GET",
      defaultError: "Load notification failed",
      redirectOn401: true,
    },
  );
}

export function updateNotification(
  teamID: string,
  notificationID: string,
  values: NotificationFormValues,
) {
  const payload = buildCreateNotificationPayload(values);

  return apiRequest<{ message?: string; data?: unknown }>(
    `/api/teams/${teamID}/notifications/${notificationID}`,
    {
      method: "PATCH",
      body: payload,
      defaultError: "Update notification failed",
      redirectOn401: true,
    },
  );
}

export function deleteNotification(teamID: string, notificationID: string) {
  return apiRequest<{ message?: string; data?: unknown }>(
    `/api/teams/${teamID}/notifications/${notificationID}`,
    {
      method: "DELETE",
      defaultError: "Delete notification failed",
      redirectOn401: true,
    },
  );
}

export function testNotification(teamID: string, notificationID: string) {
  return apiRequest<{ message?: string; data?: unknown }>(
    `/api/teams/${teamID}/notifications/${notificationID}/test`,
    {
      method: "POST",
      defaultError: "Test notification failed",
      redirectOn401: true,
    },
  );
}
