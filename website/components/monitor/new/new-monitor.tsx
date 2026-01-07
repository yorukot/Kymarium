"use client";

import React from "react";
import {
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
  CardFooter,
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

import { monitorSchema, MonitorFormValues } from "@/lib/schemas/monitor";

const regionOptions = [
  { id: 1, name: "us-east-1", label: "US East (N. Virginia)" },
  { id: 2, name: "us-west-1", label: "US West (N. California)" },
  { id: 3, name: "eu-west-1", label: "EU West (Ireland)" },
  { id: 4, name: "ap-southeast-1", label: "AP Southeast (Singapore)" },
];

const notificationOptions = [
  { id: 11, name: "email", label: "Email" },
  { id: 12, name: "slack", label: "Slack" },
  { id: 13, name: "discord", label: "Discord" },
  { id: 14, name: "telegram", label: "Telegram" },
];

function BasicSettings() {
  const {
    register,
    setValue,
    control,
    formState: { errors, isSubmitting },
  } = useFormContext<MonitorFormValues>();

  const regions = useWatch({ control, name: "regions", defaultValue: [] });
  const notifications = useWatch({
    control,
    name: "notifications",
    defaultValue: [],
  });

  const toggleListValue = (
    value: number,
    list: number[],
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

            <Field>
              <FieldLabel htmlFor="monitor-interval">Check interval</FieldLabel>
              <Input
                id="monitor-interval"
                type="number"
                min={30}
                max={2592000}
                aria-invalid={!!errors.interval}
                {...register("interval", { valueAsNumber: true })}
              />
              <FieldError errors={[errors.interval]} />
              <FieldDescription>
                Interval is in seconds. Minimum 30s, maximum 30 days.
              </FieldDescription>
            </Field>

            <div className="grid gap-6 md:grid-cols-2">
              <Field>
                <FieldLabel htmlFor="failure-threshold">
                  Failure threshold
                </FieldLabel>
                <Input
                  id="failure-threshold"
                  type="number"
                  min={1}
                  aria-invalid={!!errors.failureThreshold}
                  {...register("failureThreshold", { valueAsNumber: true })}
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
                <Input
                  id="recovery-threshold"
                  type="number"
                  min={1}
                  aria-invalid={!!errors.recoveryThreshold}
                  {...register("recoveryThreshold", { valueAsNumber: true })}
                />
                <FieldError errors={[errors.recoveryThreshold]} />
                <FieldDescription>
                  Consecutive passes required to resolve incidents.
                </FieldDescription>
              </Field>
            </div>

            <Field>
              <FieldTitle>Regions</FieldTitle>
              <FieldDescription>
                Choose at least one region to run checks from.
              </FieldDescription>

              <div className="mt-3 grid gap-2 sm:grid-cols-2">
                {regionOptions.map((region) => (
                  <label
                    key={region.id}
                    className="flex items-center gap-2 rounded-md border px-3 py-2 text-sm"
                  >
                    <Checkbox
                      checked={regions.includes(region.id)}
                      onCheckedChange={() =>
                        toggleListValue(region.id, regions, "regions")
                      }
                    />
                    <span>{region.label}</span>
                  </label>
                ))}
              </div>

              <FieldError errors={[errors.regions]} />
            </Field>

            <Field>
              <FieldTitle>Notifications</FieldTitle>
              <FieldDescription>
                Optional: select channels to notify when incidents occur.
              </FieldDescription>

              <div className="mt-3 grid gap-2 sm:grid-cols-2">
                {notificationOptions.map((option) => (
                  <label
                    key={option.id}
                    className="flex items-center gap-2 rounded-md border px-3 py-2 text-sm"
                  >
                    <Checkbox
                      checked={notifications.includes(option.id)}
                      onCheckedChange={() =>
                        toggleListValue(option.id, notifications, "notifications")
                      }
                    />
                    <span>{option.label}</span>
                  </label>
                ))}
              </div>

              <FieldError errors={[errors.notifications]} />
            </Field>
          </FieldGroup>
        </FieldSet>
      </CardContent>

      <CardFooter className="flex flex-col gap-3">
        <FieldError errors={[errors.root]} />
        <Button type="submit" className="w-full" disabled={isSubmitting}>
          Save monitor
        </Button>
      </CardFooter>
    </Card>
  );
}

export default function NewMonitorForm() {
  const form = useForm<MonitorFormValues>({
    resolver: zodResolver(monitorSchema),
    defaultValues: {
      name: "API - Production",
      interval: 60,
      failureThreshold: 3,
      recoveryThreshold: 2,
      regions: [1, 2],
      notifications: [],
    },
    mode: "onSubmit",
  });

  const onSubmit = async (values: MonitorFormValues) => {
    console.log(values);
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
          <FieldGroup>
            <BasicSettings />
          </FieldGroup>
        </form>
      </FormProvider>
    </div>
  );
}
