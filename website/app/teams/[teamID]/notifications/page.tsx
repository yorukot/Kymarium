import type { Metadata } from "next";
import { redirect } from "next/navigation";

import { buildCookieHeader } from "@/lib/api/cookies";
import { parseNotifications } from "@/lib/parsers/notifications";
import type { Notification, NotificationRawData } from "@/lib/schemas/notification";
import NotificationsPageClient from "@/components/notification/notifications-page";

type NotificationsResponse = {
  message?: string;
  data?: NotificationRawData[];
};

export const metadata: Metadata = {
  title: "Notifications",
};

async function fetchNotifications(teamID: string): Promise<Notification[]> {
  const apiBase = process.env.NEXT_PUBLIC_BACKEND_URL;
  if (!apiBase) {
    throw new Error("Missing NEXT_PUBLIC_BACKEND_URL");
  }

  const cookieHeader = await buildCookieHeader();
  const res = await fetch(`${apiBase}/api/teams/${teamID}/notifications`, {
    method: "GET",
    headers: cookieHeader ? { cookie: cookieHeader } : undefined,
    cache: "no-store",
  });

  if (res.status === 401) {
    redirect(`/login?next=/teams/${teamID}/notifications`);
  }

  if (!res.ok) {
    throw new Error("Failed to load notifications");
  }

  const body = (await res.json()) as NotificationsResponse;
  if (!Array.isArray(body?.data)) {
    return [];
  }

  return parseNotifications(body.data);
}

export default async function NotificationsPage({
  params,
}: {
  params: Promise<{ teamID: string }>;
}) {
  const { teamID } = await params;
  const notifications = await fetchNotifications(teamID);

  return <NotificationsPageClient teamID={teamID} notifications={notifications} />;
}
