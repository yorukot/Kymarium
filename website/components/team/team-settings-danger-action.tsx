"use client";

import * as React from "react";
import { useRouter } from "next/navigation";
import { toast } from "sonner";

import { Button } from "@/components/ui/button";
import {
  AlertDialog,
  AlertDialogContent,
  AlertDialogDescription,
  AlertDialogFooter,
  AlertDialogHeader,
  AlertDialogTitle,
} from "@/components/ui/alert-dialog";
import { Spinner } from "@/components/ui/spinner";
import { ApiError } from "@/lib/api/client";

export default function TeamSettingsDangerAction({
  title,
  description,
  buttonText,
  buttonIcon,
  buttonVariant = "outline",
  disabled = false,
  confirmMessage,
  loadingMessage,
  successMessage,
  defaultErrorMessage,
  action,
  redirectTo,
  refreshAfterSuccess = true,
}: {
  title: string;
  description: string;
  buttonText: string;
  buttonIcon?: React.ReactNode;
  buttonVariant?: React.ComponentProps<typeof Button>["variant"];
  disabled?: boolean;
  confirmMessage: string;
  loadingMessage: string;
  successMessage: string;
  defaultErrorMessage: string;
  action: () => Promise<unknown>;
  redirectTo?: string;
  refreshAfterSuccess?: boolean;
}) {
  const router = useRouter();
  const [pending, setPending] = React.useState(false);
  const [open, setOpen] = React.useState(false);

  const confirmButtonVariant: React.ComponentProps<typeof Button>["variant"] =
    buttonVariant === "destructive" ? "destructive" : "default";

  const handleConfirm = async () => {
    if (disabled || pending) return;

    setPending(true);
    const toastId = toast.loading(loadingMessage);
    try {
      await action();
      toast.success(successMessage, { id: toastId });
      setOpen(false);
      if (redirectTo) router.replace(redirectTo);
      if (refreshAfterSuccess) router.refresh();
    } catch (error) {
      if (error instanceof ApiError) {
        toast.error(error.message, { id: toastId });
      } else {
        toast.error(defaultErrorMessage, { id: toastId });
      }
    } finally {
      setPending(false);
    }
  };

  return (
    <div className="flex gap-1 justify-between">
      <div className="flex flex-col gap-1">
        <span className="text-sm font-medium">{title}</span>
        <span className="text-sm text-muted-foreground">{description}</span>
      </div>
      <div className="flex justify-end">
        <AlertDialog
          open={open}
          onOpenChange={(nextOpen) => {
            if (pending) return;
            setOpen(nextOpen);
          }}
        >
          <Button
            type="button"
            variant={buttonVariant}
            disabled={disabled || pending}
            onClick={() => setOpen(true)}
          >
            {buttonIcon}
            {buttonText}
          </Button>

          <AlertDialogContent
            showCloseButton={!pending}
            onEscapeKeyDown={(event) => {
              if (!pending) return;
              event.preventDefault();
            }}
            onInteractOutside={(event) => {
              if (!pending) return;
              event.preventDefault();
            }}
          >
            <AlertDialogHeader>
              <AlertDialogTitle>{title}?</AlertDialogTitle>
              <AlertDialogDescription>{confirmMessage}</AlertDialogDescription>
            </AlertDialogHeader>

            <AlertDialogFooter>
              <Button
                type="button"
                variant="outline"
                onClick={() => setOpen(false)}
                disabled={pending}
              >
                Cancel
              </Button>
              <Button
                type="button"
                variant={confirmButtonVariant}
                onClick={handleConfirm}
                disabled={pending}
              >
                {pending ? <Spinner /> : buttonIcon}
                {buttonText}
              </Button>
            </AlertDialogFooter>
          </AlertDialogContent>
        </AlertDialog>
      </div>
    </div>
  );
}
