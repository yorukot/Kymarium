"use client";
import "flag-icons/css/flag-icons.min.css";

import {
  Controller,
  FormProvider,
  useForm,
  useFormContext,
  useWatch,
} from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";

import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import {
  Field,
  FieldDescription,
  FieldError,
  FieldGroup,
  FieldLabel,
  FieldSet,
  FieldTitle,
} from "@/components/ui/field";
import { Checkbox } from "@/components/ui/checkbox";

import {
  monitorSchema,
  monitorTypes,
  MonitorFormValues,
} from "@/lib/schemas/monitor";
import HttpMonitorSettings from "@/components/monitor/new/http";
import PingMonitorSettings from "@/components/monitor/new/ping";
import { SiGmail, SiSlack, SiDiscord, SiTelegram } from "react-icons/si";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";
import { RadioGroup, RadioGroupItem } from "@/components/ui/radio-group";
import type { Region } from "@/lib/schemas/region";
import type {
  Notification,
  NotificationType,
} from "@/lib/schemas/notification";
import { ExternalLink } from "lucide-react";
import { FormDevTools } from "@/components/devtools/form-dev-tools";
import { createMonitor } from "@/lib/api/monitor";
import { ApiError } from "@/lib/api/client";
import { applyServerFieldErrors } from "@/lib/api/error";
import { Spinner } from "@/components/ui/spinner";
import { useRouter } from "next/navigation";
import { IconType } from "react-icons/lib";

const NOTIFICATION_TYPE_ICONS: Partial<Record<NotificationType, IconType>> = {
  email: SiGmail,
  slack: SiSlack,
  discord: SiDiscord,
  telegram: SiTelegram,
};

function BasicSettings({
  regions,
  notificationOptions,
}: {
  regions: Region[];
  notificationOptions: Notification[];
}) {
  const {
    register,
    setValue,
    control,
    formState: { errors },
  } = useFormContext<MonitorFormValues>();

  const selectedRegions = useWatch({
    control,
    name: "regions",
  });
  const selectedNotifications = useWatch({
    control,
    name: "notifications",
  });

  const toggleListValue = (
    value: string,
    list: string[],
    field: "regions" | "notifications",
  ) => {
    const next = list.includes(value)
      ? list.filter((item) => item !== value)
      : [...list, value];

    setValue(field, next, { shouldDirty: true, shouldValidate: true });
  };

  return (
    <Card>
      <CardHeader>
        <CardTitle>Basic settings</CardTitle>
        <CardDescription>
          Define the monitor name, schedule, and alert thresholds.
        </CardDescription>
      </CardHeader>

      <CardContent>
        <FieldSet>
          <FieldGroup>
            <Field>
              <FieldTitle>Monitor type</FieldTitle>
              <FieldDescription>
                Choose whether this monitor checks an HTTP endpoint or a ping
                target.
              </FieldDescription>

              <Controller
                name="type"
                control={control}
                render={({ field }) => (
                  <RadioGroup
                    value={field.value}
                    onValueChange={field.onChange}
                    className="mt-3 grid gap-2 sm:grid-cols-2"
                    aria-invalid={!!errors.type}
                  >
                    {monitorTypes.map((type) => (
                      <FieldLabel
                        key={type.value}
                        className="flex items-center gap-2 rounded-md border px-3 py-2 text-sm w-full"
                      >
                        <RadioGroupItem
                          id={`monitor-type-${type.value}`}
                          value={type.value}
                          aria-invalid={!!errors.type}
                        />
                        <span>{type.label}</span>
                      </FieldLabel>
                    ))}
                  </RadioGroup>
                )}
              />

              <FieldError errors={[errors.type]} />
            </Field>

            <div className="flex flex-col gap-2 sm:flex-row">
              <Field className="md:w-xs">
                <FieldLabel htmlFor="monitor-interval">
                  Check interval
                </FieldLabel>
                <Input
                  id="monitor-interval"
                  type="number"
                  min={2}
                  max={2592000}
                  aria-invalid={!!errors.interval}
                  {...register("interval", { valueAsNumber: true })}
                />
                <FieldError errors={[errors.interval]} />
                <FieldDescription>
                  Interval is in seconds. Minimum 2s, maximum 30 days.
                </FieldDescription>
              </Field>

              <Field>
                <FieldLabel htmlFor="monitor-name">Monitor name</FieldLabel>
                <Input
                  id="monitor-name"
                  placeholder="API - Production"
                  autoComplete="off"
                  aria-invalid={!!errors.name}
                  {...register("name")}
                />
                <FieldError errors={[errors.name]} />
                <FieldDescription>
                  Use a clear label so teammates know what is being checked.
                </FieldDescription>
              </Field>
            </div>

            <Field>
              <FieldTitle className="flex items-center justify-between">
                Notifications
                <Button variant="outline" size="sm">
                  Create new notifications <ExternalLink />
                </Button>
              </FieldTitle>
              <FieldDescription>
                Optional: select channels to notify when incidents occur.
              </FieldDescription>

              {notificationOptions.length === 0 ? (
                <div className="rounded-md border border-dashed px-3 py-2 text-sm text-muted-foreground">
                  No notifications yet. Create one to receive alerts.
                </div>
              ) : (
                <div className="mt-3 grid gap-2 sm:grid-cols-2">
                  {notificationOptions.map((option) => {
                    const Icon =
                      NOTIFICATION_TYPE_ICONS[option.type as NotificationType];
                    const displayName =
                      option.name?.trim() || option.typeLabel;
                    return (
                      <label
                        key={option.id}
                        className="flex items-center gap-2 rounded-md border px-3 py-2 text-sm"
                      >
                        <Checkbox
                          checked={selectedNotifications.includes(option.id)}
                          onCheckedChange={() =>
                            toggleListValue(
                              option.id,
                              selectedNotifications,
                              "notifications",
                            )
                          }
                        />
                        <div className="flex items-center gap-2">
                          <span className="text-base">
                            {Icon ? <Icon /> : null}
                          </span>
                          <div className="flex flex-col leading-tight">
                            <span className="font-medium">{displayName}</span>
                            {option.name?.trim() ? (
                              <span className="text-xs text-muted-foreground">
                                {option.typeLabel}
                              </span>
                            ) : null}
                          </div>
                        </div>
                      </label>
                    );
                  })}
                </div>
              )}

              <FieldError errors={[errors.notifications]} />
            </Field>

            <Field>
              <FieldTitle>Regions</FieldTitle>
              <FieldDescription>
                Choose at least one region to run checks from.
              </FieldDescription>

              <div className="mt-3 grid gap-2 sm:grid-cols-2">
                {regions.map((region) => {
                  return (
                    <label
                      key={region.id}
                      className="flex items-center gap-2 rounded-md border px-3 py-2 text-sm"
                    >
                      <Checkbox
                        checked={selectedRegions.includes(region.id)}
                        onCheckedChange={() =>
                          toggleListValue(region.id, selectedRegions, "regions")
                        }
                      />
                      <span>
                        <span className={`fi fi-${region.flag}`}></span>{" "}
                        {region.displayName}
                      </span>
                    </label>
                  );
                })}
              </div>

              <FieldError errors={[errors.regions]} />
            </Field>

            <div className="grid gap-6 md:grid-cols-2">
              <Field>
                <FieldLabel htmlFor="failure-threshold">
                  Failure threshold
                </FieldLabel>

                <Controller
                  name="failureThreshold"
                  control={control}
                  render={({ field }) => (
                    <Select
                      name={field.name}
                      value={String(field.value ?? 2)}
                      onValueChange={(v) => field.onChange(Number(v))}
                    >
                      <SelectTrigger aria-invalid={!!errors.failureThreshold}>
                        <SelectValue placeholder="Select threshold" />
                      </SelectTrigger>

                      <SelectContent>
                        {[1, 2, 3, 4, 5, 6, 7, 8, 9, 10].map((n) => (
                          <SelectItem key={n} value={String(n)}>
                            {n}
                          </SelectItem>
                        ))}
                      </SelectContent>
                    </Select>
                  )}
                />

                <FieldError errors={[errors.failureThreshold]} />
                <FieldDescription>
                  Consecutive failures before an incident opens.
                </FieldDescription>
              </Field>

              <Field>
                <FieldLabel htmlFor="recovery-threshold">
                  Recovery threshold
                </FieldLabel>
                <Controller
                  name="recoveryThreshold"
                  control={control}
                  render={({ field }) => (
                    <Select
                      name={field.name}
                      value={String(field.value ?? 2)}
                      onValueChange={(v) => field.onChange(Number(v))}
                    >
                      <SelectTrigger
                        id="recovery-threshold"
                        aria-invalid={!!errors.recoveryThreshold}
                      >
                        <SelectValue placeholder="Select recovery threshold" />
                      </SelectTrigger>

                      <SelectContent>
                        {[1, 2, 3, 4, 5, 6, 7, 8, 9, 10].map((n) => (
                          <SelectItem key={n} value={String(n)}>
                            {n}
                          </SelectItem>
                        ))}
                      </SelectContent>
                    </Select>
                  )}
                />
                <FieldError errors={[errors.recoveryThreshold]} />
                <FieldDescription>
                  Consecutive passes required to resolve incidents.
                </FieldDescription>
              </Field>
            </div>
          </FieldGroup>
        </FieldSet>
      </CardContent>
    </Card>
  );
}

export default function NewMonitorForm({
  teamID,
  regions,
  notifications,
}: {
  teamID: string;
  regions: Region[];
  notifications: Notification[];
}) {
  const router = useRouter();
  const defaultRegions = regions.map((region) => region.id);
  const form = useForm<MonitorFormValues>({
    resolver: zodResolver(monitorSchema),
    defaultValues: {
      name: "API - Production",
      type: "http",
      interval: 60,
      failureThreshold: 2,
      recoveryThreshold: 2,
      regions: defaultRegions,
      notifications: [],
      http: {
        url: "",
        method: "GET",
        maxRedirects: 10,
        requestTimeout: 30,
        headers: [],
        bodyEncoding: undefined,
        body: "",
        acceptedStatusCodes: [200],
        upsideDownMode: false,
        certificateExpiryNotification: true,
        ignoreTLSError: false,
      },
      ping: {
        host: "",
        timeoutSeconds: 5,
        packetSize: undefined,
      },
    },
    mode: "onSubmit",
  });

  const monitorType = useWatch({
    control: form.control,
    name: "type",
  });

  const onSubmit = async (values: MonitorFormValues) => {
    form.clearErrors();

    const parsed = monitorSchema.safeParse(values);
    if (!parsed.success) {
      form.setError("root", {
        type: "validate",
        message: "Invalid form data. Please review the fields.",
      });
      return;
    }

    try {
      await createMonitor(teamID, parsed.data);
      form.reset();
      router.push(`/teams/${teamID}/monitors`);
    } catch (error) {
      if (error instanceof ApiError) {
        const hasFieldErrors = applyServerFieldErrors(
          form.setError,
          error.body,
        );

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

  return (
    <div className="flex flex-col gap-6">
      <div>
        <h1 className="text-2xl font-semibold">Create new monitor</h1>
        <p className="text-sm text-muted-foreground">
          Configure the monitor, then wire it to the API when you are ready.
        </p>
      </div>

      <FormProvider {...form}>
        <form noValidate onSubmit={form.handleSubmit(onSubmit)}>
          <div className="flex flex-col gap-6">
            <FieldGroup>
              <BasicSettings
                regions={regions}
                notificationOptions={notifications}
              />
            </FieldGroup>
            {monitorType === "http" ? (
              <FieldGroup>
                <HttpMonitorSettings />
              </FieldGroup>
            ) : null}
            {monitorType === "ping" ? (
              <FieldGroup>
                <PingMonitorSettings />
              </FieldGroup>
            ) : null}
            <div className="flex flex-col gap-3">
              <FieldError errors={[form.formState.errors.root]} />
              <div className="flex justify-end">
                <Button type="submit" disabled={form.formState.isSubmitting}>
                  {form.formState.isSubmitting ? (
                    <>
                      <Spinner className="mr-2 h-4 w-4 animate-spin" />
                      Creating...
                    </>
                  ) : (
                    "Create new monitor"
                  )}
                </Button>
              </div>
            </div>
          </div>
        </form>
      </FormProvider>

      {process.env.NODE_ENV === "development" && (
        <FormDevTools control={form.control} />
      )}
    </div>
  );
}
