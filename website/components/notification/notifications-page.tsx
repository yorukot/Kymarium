"use client";

import type { ComponentProps } from "react";
import { useState } from "react";
import { useRouter } from "next/navigation";
import Link from "next/link";
import { toast } from "sonner";
import { Bell, MoreVertical, Pencil, Send, Trash } from "lucide-react";

import { Card, CardContent } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuGroup,
  DropdownMenuItem,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu";

import { ApiError } from "@/lib/api/client";
import { testNotification } from "@/lib/api/notification";
import { getNotificationTypeMeta } from "@/lib/parsers/notification-meta";
import {
  notificationTypeValues,
  type Notification,
  type NotificationType,
} from "@/lib/schemas/notification";
import CreateNotificationDrawer from "@/components/notification/create-notification-drawer";
import DeleteNotificationAlertDialog from "@/components/notification/delete-notification-alert-dialog";
import NotificationDrawer from "@/components/notification/notification-drawer";
import type { NotificationTypeMeta } from "@/lib/parsers/notification-meta";
import { cn } from "@/lib/utils";

const CREATE_NOTIFICATION_TYPES: NotificationType[] = [
  "slack",
  "discord",
  "telegram",
  "email",
];

export default function NotificationsPageClient({
  teamID,
  notifications,
}: {
  teamID: string;
  notifications: Notification[];
}) {
  const router = useRouter();

  return (
    <div className="flex flex-col gap-4">
      <div className="flex flex-col gap-1">
        <span className="text-xl font-bold">Notifications</span>
        <span className="text-sm text-muted-foreground">
          Manage channels used for incident alerts
        </span>
      </div>

      <div className="flex flex-col gap-3">
        <span className="text-sm font-medium">Create new notification</span>
        <div className="grid gap-2 grid-cols-2 lg:grid-cols-4">
          {CREATE_NOTIFICATION_TYPES.map((type) => {
            const meta = getNotificationTypeMeta(type);
            return (
              <CreateNotificationDrawer
                key={type}
                teamID={teamID}
                initialType={type}
                trigger={<CreateNotificationButton meta={meta} />}
              />
            );
          })}
        </div>
      </div>

      <div className="flex items-center justify-between">
        <span className="text-sm font-medium">Notifications</span>
      </div>

      {notifications.length === 0 ? (
        <div className="rounded-md border border-dashed px-3 py-2 text-sm text-muted-foreground">
          No notifications yet. Create one to receive alerts.
        </div>
      ) : (
        <div className="flex flex-col gap-3">
          {notifications.map((notification) => (
            <NotificationCard
              key={notification.id}
              teamID={teamID}
              notification={notification}
              onChanged={() => router.refresh()}
            />
          ))}
        </div>
      )}

      <div className="text-sm text-muted-foreground">
        Tip: You can attach these channels to monitors when creating or editing
        a monitor.{" "}
        <Link
          href={`/teams/${teamID}/monitors/new`}
          className="underline underline-offset-4"
        >
          Create a monitor
        </Link>
        .
      </div>
    </div>
  );
}

function NotificationCard({
  teamID,
  notification,
  onChanged,
}: {
  teamID: string;
  notification: Notification;
  onChanged: () => void;
}) {
  const displayName = notification.name?.trim() || notification.typeLabel;
  const canEdit = notificationTypeValues.includes(
    notification.type as NotificationType,
  );
  const [editOpen, setEditOpen] = useState(false);
  const [deleteOpen, setDeleteOpen] = useState(false);
  const meta = canEdit
    ? getNotificationTypeMeta(notification.type as NotificationType)
    : null;
  const Icon = meta?.Icon;
  const subtitleParts = [
    notification.name?.trim() ? notification.typeLabel : null,
    notification.detail ?? null,
  ].filter(Boolean);
  const subtitle = subtitleParts.join(" · ");

  return (
    <Card className="py-1.5 px-3">
      <CardContent className="flex items-center justify-between gap-3 p-0">
        <div className="flex min-w-0 flex-col gap-2">
          <div className="flex items-center gap-2 text-md min-w-0">
            <span className="text-base">
              {Icon ? <Icon /> : <Bell className="size-4" />}
            </span>
            <div className="min-w-0 flex-1">
              <div className="truncate font-medium">{displayName}</div>
              {subtitle ? (
                <div className="truncate text-xs text-muted-foreground">
                  {subtitle}
                </div>
              ) : null}
            </div>
          </div>
        </div>

        <div className="flex flex-none gap-3">
          <div className="flex items-center justify-center">
            <DropdownMenu>
              <DropdownMenuTrigger asChild>
                <Button size="icon" variant="ghost" aria-label="Actions">
                  <MoreVertical size={16} />
                </Button>
              </DropdownMenuTrigger>
              <DropdownMenuContent align="start">
                <DropdownMenuGroup>
                  <DropdownMenuItem
                    onSelect={async () => {
                      const toastId = toast.loading(
                        "Sending test notification...",
                      );
                      try {
                        await testNotification(teamID, notification.id);
                        toast.success("Test notification sent", {
                          id: toastId,
                        });
                      } catch (error) {
                        if (error instanceof ApiError) {
                          toast.error(error.message, { id: toastId });
                          return;
                        }
                        toast.error("Failed to send test notification", {
                          id: toastId,
                        });
                      }
                    }}
                  >
                    <Send />
                    Send test
                  </DropdownMenuItem>

                  {canEdit ? (
                    <DropdownMenuItem onSelect={() => setEditOpen(true)}>
                      <Pencil />
                      Edit
                    </DropdownMenuItem>
                  ) : null}

                  <DropdownMenuSeparator />

                  <DropdownMenuItem
                    variant="destructive"
                    onSelect={() => setDeleteOpen(true)}
                  >
                    <Trash />
                    Delete
                  </DropdownMenuItem>
                </DropdownMenuGroup>
              </DropdownMenuContent>
            </DropdownMenu>
          </div>
        </div>

        {canEdit ? (
          <NotificationDrawer
            mode="update"
            teamID={teamID}
            notificationID={notification.id}
            initialType={notification.type as NotificationType}
            open={editOpen}
            onOpenChange={setEditOpen}
          />
        ) : null}

        <DeleteNotificationAlertDialog
          teamID={teamID}
          notificationID={notification.id}
          notificationName={displayName}
          open={deleteOpen}
          onOpenChange={setDeleteOpen}
          showDefaultTrigger={false}
          onDeleted={onChanged}
        />
      </CardContent>
    </Card>
  );
}

function CreateNotificationButton({
  meta,
  className,
  variant = "outline",
  ...props
}: { meta: NotificationTypeMeta } & ComponentProps<typeof Button>) {
  const Icon = meta.Icon;

  return (
    <Button
      variant={variant}
      className={cn(
        "h-auto items-start justify-start gap-3 py-3",
        className,
      )}
      {...props}
    >
      <span className="text-base">
        <Icon />
      </span>
      <span className="flex min-w-0 flex-col items-start text-left leading-tight">
        <span className="font-medium">{meta.label}</span>
      </span>
    </Button>
  );
}
