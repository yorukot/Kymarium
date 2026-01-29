"use client";

import { useMemo, useState } from "react";
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
  UpdateIncidentEventPayload,
} from "@/lib/schemas/incident";
import {
  incidentEventCreateSchema,
  incidentEventTypeValues,
  incidentSettingsSchema,
  incidentStatusUpdateSchema,
  incidentStatusValues,
  updateIncidentEventSchema,
} from "@/lib/schemas/incident";
import {
  createIncidentEvent,
  updateIncidentEvent,
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
import { Timeline } from "@/components/ui/timeline";
import {
  Tooltip,
  TooltipContent,
  TooltipTrigger,
} from "@/components/ui/tooltip";

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

  const [editingEventId, setEditingEventId] = useState<string | null>(null);
  const [editMessage, setEditMessage] = useState("");
  const [editIsPublic, setEditIsPublic] = useState(false);
  const [savingEventId, setSavingEventId] = useState<string | null>(null);
  const [editError, setEditError] = useState<string | null>(null);

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

  const handleToggleVisibility = async (
    event: IncidentEventItem,
    target: boolean,
  ) => {
    setSavingEventId(event.id);
    setEditError(null);
    try {
      await updateIncidentEvent(teamID, incident.id, event.id, {
        public: target,
      });
      toast.success(target ? "Event made public" : "Event made private");
      router.refresh();
    } catch (error) {
      toast.error(
        error instanceof ApiError ? error.message : "Failed to update event",
      );
    } finally {
      setSavingEventId(null);
    }
  };

  const startEdit = (event: IncidentEventItem) => {
    setEditingEventId(event.id);
    setEditMessage(event.message);
    setEditIsPublic(event.isPublic);
    setEditError(null);
  };

  const cancelEdit = () => {
    setEditingEventId(null);
    setEditMessage("");
    setEditIsPublic(false);
    setEditError(null);
  };

  const saveEdit = async (event: IncidentEventItem) => {
    const payload: UpdateIncidentEventPayload = {};
    const trimmed = editMessage.trim();
    if (trimmed !== event.message) {
      payload.message = trimmed;
    }
    if (editIsPublic !== event.isPublic) {
      payload.public = editIsPublic;
    }

    const validation = updateIncidentEventSchema.safeParse(payload);
    if (!validation.success) {
      setEditError(validation.error.issues[0]?.message ?? "Invalid input");
      return;
    }

    setSavingEventId(event.id);
    setEditError(null);
    try {
      await updateIncidentEvent(teamID, incident.id, event.id, validation.data);
      toast.success("Event updated");
      cancelEdit();
      router.refresh();
    } catch (error) {
      setEditError(
        error instanceof ApiError ? error.message : "Failed to update event",
      );
    } finally {
      setSavingEventId(null);
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

      <div className="flex flex-col gap-4">
        <Card>
          <CardHeader>
            <CardTitle className="text-base flex items-center gap-2">
              <Settings className="size-4 text-muted-foreground" />
              Settings
            </CardTitle>
          </CardHeader>
          <CardContent>
            <form
              noValidate
              onSubmit={settingsForm.handleSubmit(submitSettings)}
            >
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
                      <FieldDescription>Show on status page.</FieldDescription>
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
                  <div className="w-fit">
                    <Button
                      type="submit"
                      disabled={settingsForm.formState.isSubmitting}
                      className="w-auto inline-flex"
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
                  </div>
                </Field>
              </FieldGroup>
            </form>
          </CardContent>
        </Card>

        <Card>
          <CardHeader>
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
                  <div className="w-fit">
                    <Button
                      type="submit"
                      disabled={statusForm.formState.isSubmitting}
                      className="w-auto inline-flex"
                    >
                      {statusForm.formState.isSubmitting ? (
                        <>
                          <Spinner className="mr-2 h-4 w-4 animate-spin" />
                          Updating...
                        </>
                      ) : (
                        "Update status"
                      )}
                    </Button>
                  </div>
                </Field>
              </FieldGroup>
            </form>
          </CardContent>
        </Card>

        <Card>
          <CardHeader>
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
                      eventForm.setValue(
                        "eventType",
                        value as IncidentEventType,
                        {
                          shouldDirty: true,
                          shouldValidate: true,
                        },
                      )
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
                  <div className="w-fit">
                    <Button
                      type="submit"
                      disabled={eventForm.formState.isSubmitting}
                      className="w-auto inline-flex"
                    >
                      {eventForm.formState.isSubmitting ? (
                        <>
                          <Spinner className="mr-2 h-4 w-4 animate-spin" />
                          Posting...
                        </>
                      ) : (
                        "Post update"
                      )}
                    </Button>
                  </div>
                </Field>
              </FieldGroup>
            </form>
          </CardContent>
        </Card>
      </div>

      <Card>
        <CardHeader>
          <CardTitle className="text-base">Timeline</CardTitle>
        </CardHeader>
        <CardContent>
          {sortedEvents.length === 0 ? (
            <div className="rounded-md border border-dashed px-3 py-2 text-sm text-muted-foreground">
              No updates yet.
            </div>
          ) : (
            <Timeline
              items={sortedEvents.map((event) => ({
                id: event.id,
                label: (
                  <Tooltip>
                    <TooltipTrigger asChild>
                      <span className="cursor-default">
                        {formatRelativeTime(event.createdAt)}
                      </span>
                    </TooltipTrigger>
                    <TooltipContent side="bottom">
                      {new Date(event.createdAt).toLocaleString()}
                    </TooltipContent>
                  </Tooltip>
                ),
                meta: (
                  <div className="flex items-center gap-2">
                    <EventBadge eventType={event.eventType} />
                    <Badge
                      variant={event.isPublic ? "outline" : "secondary"}
                      className="text-[11px]"
                    >
                      {event.isPublic ? "Public" : "Private"}
                    </Badge>
                  </div>
                ),
                description: (
                  <span className="whitespace-pre-wrap text-foreground">
                    {event.message}
                  </span>
                ),
                children: (
                  <div className="space-y-2">
                    {editingEventId === event.id ? (
                      <div className="space-y-3 rounded-md border p-3">
                        <Field>
                          <FieldLabel className="text-xs">Message</FieldLabel>
                          <Textarea
                            value={editMessage}
                            onChange={(e) => setEditMessage(e.target.value)}
                            minLength={1}
                            maxLength={1000}
                            aria-label="Edit event message"
                          />
                        </Field>

                        <Field orientation="horizontal">
                          <Checkbox
                            checked={editIsPublic}
                            onCheckedChange={(checked) =>
                              setEditIsPublic(!!checked)
                            }
                          />
                          <div className="flex flex-col">
                            <FieldTitle className="text-sm">Public</FieldTitle>
                            <FieldDescription className="text-xs">
                              Show this update on the status page.
                            </FieldDescription>
                          </div>
                        </Field>

                        {editError ? (
                          <div className="text-sm text-destructive">
                            {editError}
                          </div>
                        ) : null}

                        <div className="flex items-center gap-2">
                          <Button
                            size="sm"
                            disabled={savingEventId === event.id}
                            onClick={() => saveEdit(event)}
                          >
                            {savingEventId === event.id ? "Saving..." : "Save"}
                          </Button>
                          <Button
                            size="sm"
                            variant="outline"
                            onClick={cancelEdit}
                            disabled={savingEventId === event.id}
                          >
                            Cancel
                          </Button>
                        </div>
                      </div>
                    ) : (
                      <div className="flex flex-wrap items-center gap-2 text-xs">
                        <Button
                          size="sm"
                          variant="outline"
                          onClick={() =>
                            handleToggleVisibility(event, !event.isPublic)
                          }
                          disabled={savingEventId === event.id}
                        >
                          {savingEventId === event.id
                            ? "Updating..."
                            : event.isPublic
                              ? "Make private"
                              : "Make public"}
                        </Button>
                        <Button
                          size="sm"
                          variant="outline"
                          onClick={() => startEdit(event)}
                          disabled={savingEventId === event.id}
                        >
                          Edit
                        </Button>
                      </div>
                    )}
                  </div>
                ),
              }))}
            />
          )}
        </CardContent>
      </Card>
    </div>
  );
}
