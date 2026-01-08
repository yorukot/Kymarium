"use client";

import { useFormContext } from "react-hook-form";

import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import {
  Field,
  FieldDescription,
  FieldError,
  FieldGroup,
  FieldLabel,
  FieldSet,
} from "@/components/ui/field";
import { Input } from "@/components/ui/input";

import type { MonitorFormValues } from "@/lib/schemas/monitor";

export default function PingMonitorSettings() {
  const {
    register,
    formState: { errors },
  } = useFormContext<MonitorFormValues>();

  return (
    <Card>
      <CardHeader>
        <CardTitle>Ping settings</CardTitle>
        <CardDescription>
          Configure the target host and ping parameters.
        </CardDescription>
      </CardHeader>

      <CardContent>
        <FieldSet>
          <FieldGroup>
            <Field>
              <FieldLabel htmlFor="ping-host">Host</FieldLabel>
              <Input
                id="ping-host"
                placeholder="api.kymarium.com or 203.0.113.10"
                autoComplete="off"
                aria-invalid={!!errors.ping?.host}
                {...register("ping.host")}
              />
              <FieldError errors={[errors.ping?.host]} />
              <FieldDescription>
                Use a hostname or IP address reachable from the regions you
                selected.
              </FieldDescription>
            </Field>

            <div className="grid gap-6 md:grid-cols-2">
              <Field>
                <FieldLabel htmlFor="ping-timeout">Timeout</FieldLabel>
                <Input
                  id="ping-timeout"
                  type="number"
                  min={0}
                  aria-invalid={!!errors.ping?.timeoutSeconds}
                  {...register("ping.timeoutSeconds", { valueAsNumber: true })}
                />
                <FieldError errors={[errors.ping?.timeoutSeconds]} />
                <FieldDescription>Timeout in seconds.</FieldDescription>
              </Field>

              <Field>
                <FieldLabel htmlFor="ping-packet-size">
                  Packet size (bytes)
                </FieldLabel>
                <Input
                  id="ping-packet-size"
                  type="number"
                  min={1}
                  max={65000}
                  aria-invalid={!!errors.ping?.packetSize}
                  {...register("ping.packetSize", {
                    setValueAs: (value) =>
                      value === "" ? undefined : Number(value),
                  })}
                />
                <FieldError errors={[errors.ping?.packetSize]} />
                <FieldDescription>
                  Optional: set a custom payload size.
                </FieldDescription>
              </Field>
            </div>
          </FieldGroup>
        </FieldSet>
      </CardContent>
    </Card>
  );
}
