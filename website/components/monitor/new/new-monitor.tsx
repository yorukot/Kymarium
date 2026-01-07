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
import { ExternalLink } from "lucide-react";

const notificationOptions = [
  { id: "11", name: "email", label: "Email", icon: <SiGmail /> },
  { id: "12", name: "slack", label: "Slack", icon: <SiSlack /> },
  { id: "13", name: "discord", label: "Discord", icon: <SiDiscord /> },
  { id: "14", name: "telegram", label: "Telegram", icon: <SiTelegram /> },
];

function BasicSettings({ regions }: { regions: Region[] }) {
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
  const notifications = useWatch({
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

            <Field>
              <FieldLabel htmlFor="monitor-interval">Check interval</FieldLabel>
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
              <FieldTitle className="flex items-center justify-between">
                Notifications
                <Button variant="outline" size="sm">
                  Create new notifications <ExternalLink />
                </Button>
              </FieldTitle>
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
                        toggleListValue(
                          option.id,
                          notifications,
                          "notifications",
                        )
                      }
                    />
                    <div className="flex gap-1 items-center">
                      {option.icon} {option.label}
                    </div>
                  </label>
                ))}
              </div>

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

export default function NewMonitorForm({ regions }: { regions: Region[] }) {
  const defaultRegions = regions.map((region) => region.id);
  const form = useForm<MonitorFormValues>({
    resolver: zodResolver(monitorSchema),
    defaultValues: {
      name: "API - Production",
      type: "http",
      interval: 60,
      failureThreshold: 3,
      recoveryThreshold: 2,
      regions: defaultRegions,
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
          <div className="flex flex-col gap-6">
            <FieldGroup>
              <BasicSettings regions={regions} />
            </FieldGroup>
            <div className="flex justify-end">
              <Button type="submit">Create new monitor</Button>
            </div>
          </div>
        </form>
      </FormProvider>
    </div>
  );
}
