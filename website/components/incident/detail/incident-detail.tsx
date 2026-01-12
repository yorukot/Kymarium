"use client";

import { useMemo } from "react";
import { useRouter } from "next/navigation";
import { useForm, useWatch } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import { toast } from "sonner";
import {
  CalendarClock,
  Dot,
  FileText,
  Megaphone,
  Settings,
  Shield,
} from "lucide-react";

import { cn } from "@/lib/utils";
import { humanizeIdentifier } from "@/lib/parsers/strings";
import { formatRelativeTime } from "@/lib/parsers/datetime";
import type {
  IncidentDetailItem,
  IncidentEventItem,
  IncidentEventType,
  IncidentEventCreateFormValues,
  IncidentSettingsFormValues,
  IncidentStatus,
  IncidentStatusUpdateFormValues,
} from "@/lib/schemas/incident";
import {
  incidentEventCreateSchema,
  incidentEventTypeValues,
  incidentSettingsSchema,
  incidentStatusUpdateSchema,
  incidentStatusValues,
} from "@/lib/schemas/incident";
import {
  createIncidentEvent,
  updateIncidentSettings,
  updateIncidentStatus,
} from "@/lib/api/incident";
import { ApiError } from "@/lib/api/client";
import { applyServerFieldErrors } from "@/lib/api/error";

import { IncidentStatusDot } from "@/components/incident/status-dot";
import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Checkbox } from "@/components/ui/checkbox";
import {
  Field,
  FieldDescription,
  FieldError,
  FieldGroup,
  FieldLabel,
  FieldTitle,
} from "@/components/ui/field";
import { Input } from "@/components/ui/input";
import { Spinner } from "@/components/ui/spinner";
import { Textarea } from "@/components/ui/textarea";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";

function toneForEvent(eventType: string) {
  if (eventType === "manually_resolved" || eventType === "auto_resolved") {
    return "success";
  }

  if (
    eventType === "detected" ||
    eventType === "investigating" ||
    eventType === "identified" ||
    eventType === "monitoring"
  ) {
    return "danger";
  }

  return "neutral";
}

function EventBadge({ eventType }: { eventType: string }) {
  const tone = toneForEvent(eventType);
  if (tone === "success") {
    return (
      <Badge className="bg-successed/15 text-successed border-successed/30">
        {humanizeIdentifier(eventType)}
      </Badge>
    );
  }

  if (tone === "danger") {
    return <Badge variant="destructive">{humanizeIdentifier(eventType)}</Badge>;
  }

  return <Badge variant="secondary">{humanizeIdentifier(eventType)}</Badge>;
}

function labelForStatus(status: string) {
  return humanizeIdentifier(status);
}

function labelForSeverity(severity: string) {
  return humanizeIdentifier(severity);
}

export function IncidentDetail({
  teamID,
  incident,
  events,
}: {
  teamID: string;
  incident: IncidentDetailItem;
  events: IncidentEventItem[];
}) {
  const router = useRouter();

  const title = incident.title?.trim() || `Incident ${incident.id}`;
  const isResolved = incident.status === "resolved";

  const sortedEvents = useMemo(() => {
    // API returns ASC; showing newest first reads better.
    return [...events].sort((a, b) => (a.createdAt < b.createdAt ? 1 : -1));
  }, [events]);

  const settingsForm = useForm<IncidentSettingsFormValues>({
    resolver: zodResolver(incidentSettingsSchema),
    defaultValues: {
      title: incident.title ?? "",
      public: incident.isPublic,
      autoResolve: incident.autoResolve,
    },
    mode: "onSubmit",
  });

  const settingsPublic = !!useWatch({
    control: settingsForm.control,
    name: "public",
  });
  const settingsAutoResolve = !!useWatch({
    control: settingsForm.control,
    name: "autoResolve",
  });

  const statusForm = useForm<IncidentStatusUpdateFormValues>({
    resolver: zodResolver(incidentStatusUpdateSchema),
    defaultValues: {
      status: incident.status as IncidentStatus,
      message: "",
    },
    mode: "onSubmit",
  });

  const selectedStatus =
    useWatch({
      control: statusForm.control,
      name: "status",
    }) ?? (incident.status as IncidentStatus);

  const eventForm = useForm<IncidentEventCreateFormValues>({
    resolver: zodResolver(incidentEventCreateSchema),
    defaultValues: {
      eventType: "update" as IncidentEventType,
      message: "",
    },
    mode: "onSubmit",
  });

  const selectedEventType =
    useWatch({
      control: eventForm.control,
      name: "eventType",
    }) ?? ("update" as IncidentEventType);

  const submitSettings = async (values: IncidentSettingsFormValues) => {
    settingsForm.clearErrors();

    const trimmedTitle = values.title?.trim() ?? "";

    try {
      await updateIncidentSettings(teamID, incident.id, {
        title: trimmedTitle ? trimmedTitle : null,
        public: values.public,
        autoResolve: values.autoResolve,
      });

      toast.success("Incident updated");
      router.refresh();
    } catch (error) {
      if (error instanceof ApiError) {
        const hasFieldErrors = applyServerFieldErrors(
          settingsForm.setError,
          error.body,
        );

        if (!hasFieldErrors) {
          settingsForm.setError("root", {
            type: "server",
            message:
              error.status >= 500
                ? "Server error. Please try again later."
                : error.message,
          });
        }

        return;
      }

      settingsForm.setError("root", {
        type: "network",
        message: "Network error. Please try again.",
      });
    }
  };

  const submitStatus = async (values: IncidentStatusUpdateFormValues) => {
    statusForm.clearErrors();

    try {
      await updateIncidentStatus(teamID, incident.id, {
        status: values.status as IncidentStatus,
        message: values.message?.trim() || undefined,
      });

      toast.success("Status updated");
      statusForm.reset({
        status: values.status,
        message: "",
      });
      router.refresh();
    } catch (error) {
      if (error instanceof ApiError) {
        const hasFieldErrors = applyServerFieldErrors(
          statusForm.setError,
          error.body,
        );

        if (!hasFieldErrors) {
          statusForm.setError("root", {
            type: "server",
            message:
              error.status >= 500
                ? "Server error. Please try again later."
                : error.message,
          });
        }

        return;
      }

      statusForm.setError("root", {
        type: "network",
        message: "Network error. Please try again.",
      });
    }
  };

  const submitEvent = async (values: IncidentEventCreateFormValues) => {
    eventForm.clearErrors();

    try {
      await createIncidentEvent(teamID, incident.id, {
        message: values.message.trim(),
        eventType: values.eventType as IncidentEventType,
      });

      toast.success("Update posted");
      eventForm.reset({ eventType: values.eventType, message: "" });
      router.refresh();
    } catch (error) {
      if (error instanceof ApiError) {
        const hasFieldErrors = applyServerFieldErrors(
          eventForm.setError,
          error.body,
        );

        if (!hasFieldErrors) {
          eventForm.setError("root", {
            type: "server",
            message:
              error.status >= 500
                ? "Server error. Please try again later."
                : error.message,
          });
        }

        return;
      }

      eventForm.setError("root", {
        type: "network",
        message: "Network error. Please try again.",
      });
    }
  };

  return (
    <div className="flex flex-col gap-4">
      <div className="flex flex-col gap-1">
        <div className="flex items-center gap-2 min-w-0">
          <IncidentStatusDot status={incident.status} />
          <span className="text-xl font-bold min-w-0 truncate">{title}</span>
        </div>
        <div className="flex flex-wrap items-center gap-2 text-sm text-muted-foreground">
          <Badge variant={isResolved ? "secondary" : "destructive"}>
            {labelForStatus(incident.status)}
          </Badge>
          <Badge variant="outline">{labelForSeverity(incident.severity)}</Badge>
          <span className="inline-flex items-center">
            <CalendarClock className="mr-1.5 size-4" />
            Started {formatRelativeTime(incident.startedAt)}
          </span>
          {incident.resolvedAt ? (
            <span className="inline-flex items-center">
              <Dot className="-mx-1 size-4" />
              Resolved {formatRelativeTime(incident.resolvedAt)}
            </span>
          ) : null}
        </div>
      </div>

      <div className="grid gap-4 lg:grid-cols-5">
        <div className="flex flex-col gap-4 lg:col-span-2">
          <Card>
            <CardHeader className="pb-3">
              <CardTitle className="text-base flex items-center gap-2">
                <Settings className="size-4 text-muted-foreground" />
                Settings
              </CardTitle>
            </CardHeader>
            <CardContent>
              <form noValidate onSubmit={settingsForm.handleSubmit(submitSettings)}>
                <FieldGroup>
                  <Field>
                    <FieldLabel htmlFor="incidentTitle">Title</FieldLabel>
                    <Input
                      id="incidentTitle"
                      type="text"
                      placeholder="Optional incident title"
                      aria-invalid={!!settingsForm.formState.errors.title}
                      {...settingsForm.register("title")}
                    />
                    <FieldError errors={[settingsForm.formState.errors.title]} />
                  </Field>

                  <div className="grid gap-4 md:grid-cols-2">
                    <Field orientation="horizontal">
                      <Checkbox
                        checked={settingsPublic}
                        onCheckedChange={(checked) =>
                          settingsForm.setValue("public", !!checked, {
                            shouldDirty: true,
                            shouldValidate: true,
                          })
                        }
                      />
                      <div className="flex flex-col gap-1">
                        <FieldTitle className="flex items-center gap-2">
                          <Shield className="size-4 text-muted-foreground" />
                          Public
                        </FieldTitle>
                        <FieldDescription>
                          Show on status page.
                        </FieldDescription>
                      </div>
                    </Field>

                    <Field orientation="horizontal">
                      <Checkbox
                        checked={settingsAutoResolve}
                        onCheckedChange={(checked) =>
                          settingsForm.setValue("autoResolve", !!checked, {
                            shouldDirty: true,
                            shouldValidate: true,
                          })
                        }
                      />
                      <div className="flex flex-col gap-1">
                        <FieldTitle>Auto-resolve</FieldTitle>
                        <FieldDescription>
                          Allow automatic resolution.
                        </FieldDescription>
                      </div>
                    </Field>
                  </div>

                  <Field>
                    <FieldError errors={[settingsForm.formState.errors.root]} />
                    <Button
                      type="submit"
                      disabled={settingsForm.formState.isSubmitting}
                    >
                      {settingsForm.formState.isSubmitting ? (
                        <>
                          <Spinner className="mr-2 h-4 w-4 animate-spin" />
                          Saving...
                        </>
                      ) : (
                        "Save changes"
                      )}
                    </Button>
                  </Field>
                </FieldGroup>
              </form>
            </CardContent>
          </Card>

          <Card>
            <CardHeader className="pb-3">
              <CardTitle className="text-base flex items-center gap-2">
                <Megaphone className="size-4 text-muted-foreground" />
                Update status
              </CardTitle>
            </CardHeader>
            <CardContent>
              <form noValidate onSubmit={statusForm.handleSubmit(submitStatus)}>
                <FieldGroup>
                  <Field>
                    <FieldLabel>Status</FieldLabel>
                    <Select
                      value={selectedStatus}
                      onValueChange={(value) =>
                        statusForm.setValue("status", value as IncidentStatus, {
                          shouldDirty: true,
                          shouldValidate: true,
                        })
                      }
                    >
                      <SelectTrigger
                        aria-invalid={!!statusForm.formState.errors.status}
                      >
                        <SelectValue placeholder="Select status" />
                      </SelectTrigger>
                      <SelectContent>
                        {incidentStatusValues.map((status) => (
                          <SelectItem key={status} value={status}>
                            {humanizeIdentifier(status)}
                          </SelectItem>
                        ))}
                      </SelectContent>
                    </Select>
                    <FieldError errors={[statusForm.formState.errors.status]} />
                  </Field>

                  <Field>
                    <FieldLabel htmlFor="statusMessage">Message</FieldLabel>
                    <Textarea
                      id="statusMessage"
                      placeholder="Optional message"
                      aria-invalid={!!statusForm.formState.errors.message}
                      {...statusForm.register("message")}
                    />
                    <FieldError errors={[statusForm.formState.errors.message]} />
                  </Field>

                  <Field>
                    <FieldError errors={[statusForm.formState.errors.root]} />
                    <Button type="submit" disabled={statusForm.formState.isSubmitting}>
                      {statusForm.formState.isSubmitting ? (
                        <>
                          <Spinner className="mr-2 h-4 w-4 animate-spin" />
                          Updating...
                        </>
                      ) : (
                        "Update status"
                      )}
                    </Button>
                  </Field>
                </FieldGroup>
              </form>
            </CardContent>
          </Card>

          <Card>
            <CardHeader className="pb-3">
              <CardTitle className="text-base flex items-center gap-2">
                <FileText className="size-4 text-muted-foreground" />
                Post update
              </CardTitle>
            </CardHeader>
            <CardContent>
              <form noValidate onSubmit={eventForm.handleSubmit(submitEvent)}>
                <FieldGroup>
                  <Field>
                    <FieldLabel>Type</FieldLabel>
                    <Select
                      value={selectedEventType}
                      onValueChange={(value) =>
                        eventForm.setValue("eventType", value as IncidentEventType, {
                          shouldDirty: true,
                          shouldValidate: true,
                        })
                      }
                    >
                      <SelectTrigger
                        aria-invalid={!!eventForm.formState.errors.eventType}
                      >
                        <SelectValue placeholder="Select event type" />
                      </SelectTrigger>
                      <SelectContent>
                        {incidentEventTypeValues.map((eventType) => (
                          <SelectItem key={eventType} value={eventType}>
                            {humanizeIdentifier(eventType)}
                          </SelectItem>
                        ))}
                      </SelectContent>
                    </Select>
                    <FieldError errors={[eventForm.formState.errors.eventType]} />
                  </Field>

                  <Field>
                    <FieldLabel htmlFor="eventMessage">Message</FieldLabel>
                    <Textarea
                      id="eventMessage"
                      placeholder="Write an update…"
                      aria-invalid={!!eventForm.formState.errors.message}
                      {...eventForm.register("message")}
                    />
                    <FieldError errors={[eventForm.formState.errors.message]} />
                  </Field>

                  <Field>
                    <FieldError errors={[eventForm.formState.errors.root]} />
                    <Button type="submit" disabled={eventForm.formState.isSubmitting}>
                      {eventForm.formState.isSubmitting ? (
                        <>
                          <Spinner className="mr-2 h-4 w-4 animate-spin" />
                          Posting...
                        </>
                      ) : (
                        "Post update"
                      )}
                    </Button>
                  </Field>
                </FieldGroup>
              </form>
            </CardContent>
          </Card>
        </div>

        <div className="lg:col-span-3">
          <Card>
            <CardHeader className="pb-3">
              <CardTitle className="text-base">Timeline</CardTitle>
            </CardHeader>
            <CardContent>
              {sortedEvents.length === 0 ? (
                <div className="rounded-md border border-dashed px-3 py-2 text-sm text-muted-foreground">
                  No updates yet.
                </div>
              ) : (
                <div className="relative">
                  <div className="absolute left-2 top-0 bottom-0 w-px bg-border" />
                  <div className="flex flex-col gap-4">
                    {sortedEvents.map((event) => (
                      <div key={event.id} className="relative pl-7">
                        <span
                          className={cn(
                            "absolute left-[5px] top-1.5 h-4 w-4 rounded-full border bg-card",
                            toneForEvent(event.eventType) === "success" &&
                              "border-successed/40",
                            toneForEvent(event.eventType) === "danger" &&
                              "border-destructive/40",
                          )}
                        >
                          <span
                            className={cn(
                              "absolute inset-1 rounded-full",
                              toneForEvent(event.eventType) === "success" &&
                                "bg-successed/60",
                              toneForEvent(event.eventType) === "danger" &&
                                "bg-destructive/60",
                              toneForEvent(event.eventType) === "neutral" &&
                                "bg-muted-foreground/40",
                            )}
                          />
                        </span>

                        <div className="flex flex-col gap-1">
                          <div className="flex flex-wrap items-center gap-2">
                            <EventBadge eventType={event.eventType} />
                            <span
                              className="text-sm text-muted-foreground"
                              title={event.createdAt}
                              suppressHydrationWarning
                            >
                              {formatRelativeTime(event.createdAt)}
                            </span>
                          </div>
                          <div className="text-sm whitespace-pre-wrap">
                            {event.message}
                          </div>
                        </div>
                      </div>
                    ))}
                  </div>
                </div>
              )}
            </CardContent>
          </Card>
        </div>
      </div>
    </div>
  );
}
