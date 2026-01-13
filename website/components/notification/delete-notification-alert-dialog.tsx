"use client";

import { cloneElement, useState, type ReactElement } from "react";
import { toast } from "sonner";
import { Trash } from "lucide-react";

import { Button } from "@/components/ui/button";
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
} from "@/components/ui/dialog";
import { Spinner } from "@/components/ui/spinner";

import { ApiError } from "@/lib/api/client";
import { deleteNotification } from "@/lib/api/notification";

export default function DeleteNotificationAlertDialog({
  teamID,
  notificationID,
  notificationName,
  disabled,
  onDeleted,
  buttonClassName,
  trigger,
}: {
  teamID: string;
  notificationID: string;
  notificationName: string;
  disabled?: boolean;
  onDeleted: () => void;
  buttonClassName?: string;
  trigger?: ReactElement;
}) {
  const [open, setOpen] = useState(false);
  const [deleting, setDeleting] = useState(false);

  const handleDelete = async () => {
    const toastId = toast.loading("Deleting notification...");
    setDeleting(true);

    try {
      await deleteNotification(teamID, notificationID);
      toast.success("Notification deleted", { id: toastId });
      setOpen(false);
      onDeleted();
    } catch (error) {
      if (error instanceof ApiError) {
        toast.error(error.message, { id: toastId });
        return;
      }
      toast.error("Failed to delete notification", { id: toastId });
    } finally {
      setDeleting(false);
    }
  };

  const triggerNode = trigger
    ? cloneElement(trigger, {
        ...trigger.props,
        disabled: Boolean(trigger.props.disabled || disabled || deleting),
        onClick: (event: unknown) => {
          trigger.props.onClick?.(event);
          if ((event as { defaultPrevented?: boolean } | null)?.defaultPrevented)
            return;
          setOpen(true);
        },
        onSelect: (event: unknown) => {
          trigger.props.onSelect?.(event);
          if ((event as { defaultPrevented?: boolean } | null)?.defaultPrevented)
            return;
          setOpen(true);
        },
      })
    : null;

  return (
    <Dialog open={open} onOpenChange={setOpen}>
      {triggerNode ?? (
        <Button
          type="button"
          variant="destructive"
          onClick={() => setOpen(true)}
          disabled={disabled || deleting}
          className={buttonClassName}
        >
          {deleting ? <Spinner /> : <Trash />}
          Delete
        </Button>
      )}

      <DialogContent showCloseButton={!deleting}>
        <DialogHeader>
          <DialogTitle>Delete notification?</DialogTitle>
          <DialogDescription>
            This will permanently delete{" "}
            <span className="font-medium">{notificationName}</span>. Monitors
            using it may stop sending alerts until updated.
          </DialogDescription>
        </DialogHeader>

        <DialogFooter>
          <Button
            type="button"
            variant="outline"
            onClick={() => setOpen(false)}
            disabled={deleting}
          >
            Cancel
          </Button>
          <Button
            type="button"
            variant="destructive"
            onClick={handleDelete}
            disabled={deleting}
          >
            {deleting ? <Spinner /> : null}
            Delete notification
          </Button>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  );
}
