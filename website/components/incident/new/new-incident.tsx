"use client";

import { useMemo } from "react";
import { useRouter } from "next/navigation";
import { Controller, useForm, useWatch } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import { toast } from "sonner";

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
import { StatusDot } from "@/components/monitor/status-dot";

import { ApiError } from "@/lib/api/client";
import { applyServerFieldErrors } from "@/lib/api/error";
import { createIncident } from "@/lib/api/incident";
import { humanizeIdentifier } from "@/lib/parsers/strings";
import type { MonitorListItem } from "@/lib/schemas/monitor";
import {
  createIncidentPayloadSchema,
  createIncidentSchema,
  incidentSeverityValues,
  incidentStatusValues,
  type CreateIncidentFormValues,
} from "@/lib/schemas/incident";

function uniqueStrings(values: string[]) {
  if (values.length <= 1) return values;
  return Array.from(new Set(values));
}

export default function NewIncidentForm({
  teamID,
  monitors,
}: {
  teamID: string;
  monitors: MonitorListItem[];
}) {
  const router = useRouter();
  const form = useForm<CreateIncidentFormValues>({
    resolver: zodResolver(createIncidentSchema),
    defaultValues: {
      title: "",
      message: "",
      status: "detected",
      severity: "major",
      public: true,
      autoResolve: false,
      monitorIds: [],
    },
    mode: "onSubmit",
  });

  const monitorOptions = useMemo(
    () =>
      monitors.map((monitor) => ({
        id: monitor.id,
        name: monitor.name,
        status: monitor.status,
        targetLabel: monitor.targetLabel,
        targetValue: monitor.targetValue,
      })),
    [monitors],
  );

  const toggleMonitor = (monitorID: string) => {
    const current = form.getValues("monitorIds");
    const next = current.includes(monitorID)
      ? current.filter((id) => id !== monitorID)
      : [...current, monitorID];

    form.setValue("monitorIds", next, { shouldDirty: true, shouldValidate: true });
  };

  const onSubmit = async (values: CreateIncidentFormValues) => {
    form.clearErrors();

    const normalized = {
      ...values,
      monitorIds: uniqueStrings(values.monitorIds),
      title: values.title.trim() || undefined,
      message: values.message.trim() || undefined,
    };

    const parsed = createIncidentPayloadSchema.safeParse(normalized);
    if (!parsed.success) {
      form.setError("root", {
        type: "validate",
        message: "Invalid form data. Please try again.",
      });
      return;
    }

    try {
      await createIncident(teamID, parsed.data);

      toast.success("Incident created");
      router.replace(`/teams/${teamID}/incidents`);
    } catch (error) {
      if (error instanceof ApiError) {
        const hasFieldErrors = applyServerFieldErrors(form.setError, error.body);
        if (!hasFieldErrors) {
          form.setError("root", {
            type: "server",
            message:
              error.status >= 500
                ? "Server error. Please try again later."
                : error.message,
          });
        }
        return;
      }

      form.setError("root", {
        type: "network",
        message: "Network error. Please try again.",
      });
    }
  };

  const {
    register,
    handleSubmit,
    formState: { errors, isSubmitting },
    control,
  } = form;

  const selectedMonitorIDs =
    useWatch({ control, name: "monitorIds" }) ?? [];
  const isPublic = !!useWatch({ control, name: "public" });
  const autoResolve = !!useWatch({ control, name: "autoResolve" });

  return (
    <div className="flex flex-col gap-4">
      <div className="flex justify-between items-center">
        <div className="flex flex-col gap-1">
          <span className="text-xl font-bold">Create Incident</span>
          <span className="text-sm text-muted-foreground">
            Manually open an incident for one or more monitors
          </span>
        </div>
      </div>

      <Card>
        <CardHeader>
          <CardTitle className="text-xl">Incident details</CardTitle>
        </CardHeader>
        <CardContent>
          <form noValidate onSubmit={handleSubmit(onSubmit)}>
            <FieldGroup>
              <Field>
                <FieldLabel htmlFor="title">Title</FieldLabel>
                <Input
                  id="title"
                  type="text"
                  placeholder="Optional incident title"
                  aria-invalid={!!errors.title}
                  {...register("title")}
                />
                <FieldError errors={[errors.title]} />
              </Field>

              <div className="grid gap-6 md:grid-cols-2">
                <Field>
                  <FieldLabel>Status</FieldLabel>
                  <Controller
                    name="status"
                    control={control}
                    render={({ field }) => (
                      <Select
                        name={field.name}
                        value={field.value}
                        onValueChange={field.onChange}
                      >
                        <SelectTrigger aria-invalid={!!errors.status}>
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
                    )}
                  />
                  <FieldError errors={[errors.status]} />
                </Field>

                <Field>
                  <FieldLabel>Severity</FieldLabel>
                  <Controller
                    name="severity"
                    control={control}
                    render={({ field }) => (
                      <Select
                        name={field.name}
                        value={field.value}
                        onValueChange={field.onChange}
                      >
                        <SelectTrigger aria-invalid={!!errors.severity}>
                          <SelectValue placeholder="Select severity" />
                        </SelectTrigger>
                        <SelectContent>
                          {incidentSeverityValues.map((severity) => (
                            <SelectItem key={severity} value={severity}>
                              {humanizeIdentifier(severity)}
                            </SelectItem>
                          ))}
                        </SelectContent>
                      </Select>
                    )}
                  />
                  <FieldError errors={[errors.severity]} />
                </Field>
              </div>

              <Field>
                <FieldLabel htmlFor="message">Message</FieldLabel>
                <Textarea
                  id="message"
                  placeholder="Optional: add context for this incident"
                  aria-invalid={!!errors.message}
                  {...register("message")}
                />
                <FieldError errors={[errors.message]} />
              </Field>

              <Field>
                <FieldTitle>Monitors</FieldTitle>
                <FieldDescription>
                  Choose at least one monitor to attach this incident to.
                </FieldDescription>

                {monitorOptions.length === 0 ? (
                  <div className="rounded-md border border-dashed px-3 py-2 text-sm text-muted-foreground">
                    No monitors found.
                  </div>
                ) : (
                  <div className="mt-3 grid gap-2 sm:grid-cols-2">
                    {monitorOptions.map((monitor) => (
                      <label
                        key={monitor.id}
                        className="flex items-center gap-2 rounded-md border px-3 py-2 text-sm"
                      >
                        <Checkbox
                          checked={selectedMonitorIDs.includes(monitor.id)}
                          onCheckedChange={() => toggleMonitor(monitor.id)}
                        />
                        <div className="flex min-w-0 flex-col leading-tight">
                          <div className="flex items-center gap-2 min-w-0">
                            <StatusDot status={monitor.status} />
                            <span className="truncate font-medium">
                              {monitor.name}
                            </span>
                          </div>
                          <span className="text-xs text-muted-foreground truncate">
                            {monitor.targetLabel}: {monitor.targetValue || "--"}
                          </span>
                        </div>
                      </label>
                    ))}
                  </div>
                )}

                <FieldError errors={[errors.monitorIds]} />
              </Field>

              <div className="grid gap-6 md:grid-cols-2">
                <Field orientation="horizontal">
                  <Checkbox
                    checked={isPublic}
                    onCheckedChange={(checked) =>
                      form.setValue("public", !!checked, {
                        shouldDirty: true,
                        shouldValidate: true,
                      })
                    }
                  />
                  <div className="flex flex-col gap-1">
                    <FieldTitle>Public</FieldTitle>
                    <FieldDescription>
                      Show this incident on your status page.
                    </FieldDescription>
                  </div>
                </Field>

                <Field orientation="horizontal">
                  <Checkbox
                    checked={autoResolve}
                    onCheckedChange={(checked) =>
                      form.setValue("autoResolve", !!checked, {
                        shouldDirty: true,
                        shouldValidate: true,
                      })
                    }
                  />
                  <div className="flex flex-col gap-1">
                    <FieldTitle>Auto-resolve</FieldTitle>
                    <FieldDescription>
                      Allow the system to resolve it automatically.
                    </FieldDescription>
                  </div>
                </Field>
              </div>

              <Field>
                <FieldError errors={[errors.root]} />
                <Button type="submit" disabled={isSubmitting}>
                  {isSubmitting ? (
                    <>
                      <Spinner className="mr-2 h-4 w-4 animate-spin" />
                      Creating incident...
                    </>
                  ) : (
                    "Create incident"
                  )}
                </Button>
              </Field>
            </FieldGroup>
          </form>
        </CardContent>
      </Card>
    </div>
  );
}
