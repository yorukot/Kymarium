"use client";

import type { ReactNode } from "react";

import NotificationDrawer from "@/components/notification/notification-drawer";
import type { NotificationFormValues } from "@/lib/schemas/notification";

export default function CreateNotificationDrawer({
  teamID,
  trigger,
  initialType,
}: {
  teamID: string;
  trigger: ReactNode;
  initialType: NotificationFormValues["type"];
}) {
  return (
    <NotificationDrawer
      mode="create"
      teamID={teamID}
      trigger={trigger}
      initialType={initialType}
    />
  );
}

