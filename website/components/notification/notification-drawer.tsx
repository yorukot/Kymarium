"use client";

import type { ReactNode } from "react";
import { useEffect, useMemo, useState } from "react";
import { useRouter } from "next/navigation";
import { useFieldArray, useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import { toast } from "sonner";
import { Plus, Trash2 } from "lucide-react";

import {
  Sheet,
  SheetContent,
  SheetDescription,
  SheetFooter,
  SheetHeader,
  SheetTitle,
  SheetTrigger,
} from "@/components/ui/sheet";
import { Button } from "@/components/ui/button";
import {
  Field,
  FieldDescription,
  FieldError,
  FieldGroup,
  FieldLabel,
} from "@/components/ui/field";
import { Input } from "@/components/ui/input";
import { Spinner } from "@/components/ui/spinner";

import { ApiError } from "@/lib/api/client";
import { applyServerFieldErrors } from "@/lib/api/error";
import {
  createNotification,
  getNotification,
  updateNotification,
} from "@/lib/api/notification";
import DeleteNotificationAlertDialog from "@/components/notification/delete-notification-alert-dialog";
import {
  defaultNotificationFormValues,
  isNotificationType,
  notificationFormValuesFromApi,
} from "@/lib/parsers/notification-form";
import {
  notificationSchema,
  type NotificationFormValues,
} from "@/lib/schemas/notification";

type Mode = "create" | "update";

type PropsBase = {
  teamID: string;
  initialType: NotificationFormValues["type"];
};

type PropsCreate = PropsBase & {
  mode: "create";
  trigger: ReactNode;
  open?: never;
  onOpenChange?: never;
  notificationID?: never;
};

type PropsUpdate = PropsBase & {
  mode: "update";
  trigger?: ReactNode;
  open?: boolean;
  onOpenChange?: (open: boolean) => void;
  notificationID: string;
};

type Props = PropsCreate | PropsUpdate;

function titleFor(mode: Mode, type: NotificationFormValues["type"]) {
  const action = mode === "create" ? "Create" : "Update";
  if (type === "slack") return `${action} Slack notification`;
  if (type === "discord") return `${action} Discord notification`;
  if (type === "telegram") return `${action} Telegram notification`;
  return `${action} Email notification`;
}

export default function NotificationDrawer(props: Props) {
  const router = useRouter();

  const controlled = props.mode === "update" && props.open !== undefined;
  const [internalOpen, setInternalOpen] = useState(false);
  const open = controlled ? (props.open as boolean) : internalOpen;
  const setOpen = controlled
    ? (props.onOpenChange as (open: boolean) => void)
    : setInternalOpen;

  const [loading, setLoading] = useState(false);
  const [loadedType, setLoadedType] = useState<NotificationFormValues["type"]>(
    props.initialType,
  );
  const type = props.mode === "create" ? props.initialType : loadedType;

  const form = useForm<NotificationFormValues>({
    resolver: zodResolver(notificationSchema),
    defaultValues: defaultNotificationFormValues(props.initialType),
    mode: "onSubmit",
  });

  const notificationID = props.mode === "update" ? props.notificationID : "";
  const { clearErrors, reset, setError } = form;

  const emailAddresses = useFieldArray({
    control: form.control,
    // Note: only used for email notifications; safe to keep registered for other types.
    name: "config.emailAddresses" as never,
  });
  const { replace: replaceEmailAddresses } = emailAddresses;

  useEffect(() => {
    if (!open) return;

    clearErrors();

    // Create flow: reset to clean defaults each time the drawer opens.
    if (props.mode === "create") {
      const defaults = defaultNotificationFormValues(props.initialType);
      reset(defaults);
      if (defaults.type === "email") {
        replaceEmailAddresses(defaults.config.emailAddresses as never);
      }
      return;
    }

    // Update flow: fetch current values and populate the form.
    let cancelled = false;

    Promise.resolve()
      .then(() => {
        if (cancelled) return null;
        setLoading(true);
        return getNotification(props.teamID, notificationID);
      })
      .then((result) => {
        if (cancelled) return;
        if (!result) return;

        const { data } = result;
        if (cancelled) return;

        const notification = data.data;
        if (!notification) {
          setError("root", {
            type: "server",
            message: "Failed to load notification.",
          });
          return;
        }

        if (!isNotificationType(notification.type)) {
          setError("root", {
            type: "server",
            message: `Unsupported notification type: ${String(notification.type)}`,
          });
          return;
        }

        const nextType = notification.type;
        setLoadedType(nextType);

        const values = notificationFormValuesFromApi(
          nextType,
          notification.name,
          notification.config,
        );
        reset(values);
        if (values.type === "email") {
          replaceEmailAddresses(values.config.emailAddresses as never);
        }
      })
      .catch((error) => {
        if (cancelled) return;

        if (error instanceof ApiError) {
          setError("root", {
            type: "server",
            message:
              error.status >= 500
                ? "Server error. Please try again later."
                : error.message,
          });
          return;
        }

        setError("root", {
          type: "network",
          message: "Network error. Please try again.",
        });
      })
      .finally(() => {
        if (cancelled) return;
        setLoading(false);
      });

    return () => {
      cancelled = true;
    };
  }, [
    open,
    props.initialType,
    props.mode,
    props.teamID,
    notificationID,
    clearErrors,
    reset,
    setError,
    replaceEmailAddresses,
  ]);

  const title = useMemo(() => titleFor(props.mode, type), [props.mode, type]);

  const onSubmit = async (values: NotificationFormValues) => {
    clearErrors();

    try {
      if (props.mode === "create") {
        await createNotification(props.teamID, {
          ...values,
          type,
        } as NotificationFormValues);
        toast.success("Notification created");

        setOpen(false);
        reset(defaultNotificationFormValues(type));
        router.refresh();
        return;
      }

      await updateNotification(props.teamID, props.notificationID, {
        ...values,
        type,
      } as NotificationFormValues);

      toast.success("Notification updated");
      setOpen(false);
      router.refresh();
    } catch (error) {
      if (error instanceof ApiError) {
        const hasFieldErrors = applyServerFieldErrors(setError, error.body);
        if (!hasFieldErrors) {
          setError("root", {
            type: "server",
            message:
              error.status >= 500
                ? "Server error. Please try again later."
                : error.message,
          });
        }
        return;
      }

      setError("root", {
        type: "network",
        message: "Network error. Please try again.",
      });
    }
  };

  const {
    register,
    handleSubmit,
    formState: { errors, isSubmitting },
  } = form;

  const disabled = isSubmitting || loading;

  return (
    <Sheet
      open={open}
      onOpenChange={(next) => {
        setOpen(next);
        if (!next) {
          setLoading(false);
          form.reset(defaultNotificationFormValues(props.initialType));
        }
      }}
    >
      {"trigger" in props && props.trigger ? (
        <SheetTrigger asChild>{props.trigger}</SheetTrigger>
      ) : null}

      <SheetContent side="right" className="sm:max-w-lg">
        <SheetHeader className="bg-muted">
          <SheetTitle>{title}</SheetTitle>
          <SheetDescription>
            {props.mode === "create"
              ? "Add a channel so your team can receive incident alerts."
              : "Update where incident alerts are delivered for your team."}
          </SheetDescription>
        </SheetHeader>

        <form
          noValidate
          onSubmit={handleSubmit(onSubmit)}
          className="flex min-h-0 flex-1 flex-col"
        >
          <div className="flex min-h-0 flex-1 flex-col overflow-auto px-4 pb-4">
            <FieldGroup>
              <Field>
                <FieldLabel>Name</FieldLabel>
                <Input
                  placeholder={
                    type === "email"
                      ? "e.g. Email - On call"
                      : "e.g. Slack - #alerts"
                  }
                  aria-invalid={!!errors.name}
                  disabled={disabled}
                  {...register("name")}
                />
                <FieldError errors={[errors.name]} />
                <FieldDescription>
                  Use a clear label so teammates know where alerts go.
                </FieldDescription>
              </Field>

              {type === "email" ? (
                <Field>
                  <FieldLabel>Email recipients</FieldLabel>
                  <FieldDescription>
                    Add one or more email addresses to receive alerts.
                  </FieldDescription>

                  <div className="mt-3 flex flex-col gap-2">
                    {emailAddresses.fields.map((field, index) => (
                      <div key={field.id} className="flex gap-2">
                        <Input
                          placeholder="name@company.com"
                          aria-invalid={
                            !!errors.config?.emailAddresses?.[index]?.value
                          }
                          disabled={disabled}
                          {...register(`config.emailAddresses.${index}.value`)}
                        />
                        <Button
                          type="button"
                          variant="outline"
                          size="icon"
                          aria-label="Remove email"
                          onClick={() => emailAddresses.remove(index)}
                          disabled={
                            disabled || emailAddresses.fields.length <= 1
                          }
                        >
                          <Trash2 className="size-4" />
                        </Button>
                      </div>
                    ))}
                    <Button
                      type="button"
                      variant="outline"
                      onClick={() => emailAddresses.append({ value: "" })}
                      disabled={disabled}
                    >
                      <Plus className="size-4" />
                      Add email
                    </Button>
                  </div>

                  <FieldError
                    errors={[
                      errors.config?.emailAddresses,
                      ...(Array.isArray(errors.config?.emailAddresses)
                        ? errors.config?.emailAddresses.map((e) => e?.value)
                        : []),
                    ]}
                  />
                </Field>
              ) : null}

              {type === "slack" || type === "discord" ? (
                <Field>
                  <FieldLabel>Webhook URL</FieldLabel>
                  <Input
                    placeholder="https://..."
                    aria-invalid={!!errors.config?.webhookUrl}
                    disabled={disabled}
                    {...register("config.webhookUrl")}
                  />
                  <FieldError errors={[errors.config?.webhookUrl]} />
                </Field>
              ) : null}

              {type === "telegram" ? (
                <>
                  <Field>
                    <FieldLabel>Bot token</FieldLabel>
                    <Input
                      placeholder="123456:ABC..."
                      aria-invalid={!!errors.config?.botToken}
                      disabled={disabled}
                      {...register("config.botToken")}
                    />
                    <FieldError errors={[errors.config?.botToken]} />
                  </Field>

                  <Field>
                    <FieldLabel>Chat ID</FieldLabel>
                    <Input
                      placeholder="e.g. -1001234567890"
                      aria-invalid={!!errors.config?.chatId}
                      disabled={disabled}
                      {...register("config.chatId")}
                    />
                    <FieldError errors={[errors.config?.chatId]} />
                  </Field>
                </>
              ) : null}

              <FieldError errors={[errors.root]} />
            </FieldGroup>
          </div>

          <SheetFooter className="flex-col items-center bg-muted">
            {props.mode === "update" ? (
              <DeleteNotificationAlertDialog
                teamID={props.teamID}
                notificationID={props.notificationID}
                notificationName={form.getValues("name") || "this notification"}
                disabled={disabled}
                buttonClassName="w-full"
                onDeleted={() => {
                  setOpen(false);
                  router.refresh();
                }}
              />
            ) : null}

            <Button type="submit" disabled={disabled} className="w-full">
              {disabled ? <Spinner /> : null}
              {props.mode === "create" ? "Create" : "Update"}
            </Button>
          </SheetFooter>
        </form>
      </SheetContent>
    </Sheet>
  );
}
