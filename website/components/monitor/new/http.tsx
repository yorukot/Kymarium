"use client";

import * as React from "react";
import {
  Controller,
  useFieldArray,
  useFormContext,
  useWatch,
} from "react-hook-form";

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
  FieldTitle,
} from "@/components/ui/field";
import { Input } from "@/components/ui/input";
import { Checkbox } from "@/components/ui/checkbox";
import { Button } from "@/components/ui/button";
import { MultiSelect } from "@/components/ui/multi-selector";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";

import type { MonitorFormValues } from "@/lib/schemas/monitor";
import {
  httpBodyEncodingValues,
  httpMethodValues,
  httpStatusCodes,
} from "@/lib/schemas/monitor";
import { Plus, Trash2 } from "lucide-react";
import { Textarea } from "@/components/ui/textarea";
import {
  Accordion,
  AccordionContent,
  AccordionItem,
  AccordionTrigger,
} from "@/components/ui/accordion";

const methodLabels: Record<(typeof httpMethodValues)[number], string> = {
  GET: "GET",
  POST: "POST",
  PUT: "PUT",
  DELETE: "DELETE",
  PATCH: "PATCH",
  HEAD: "HEAD",
  OPTIONS: "OPTIONS",
};

const bodyEncodingLabels: Record<
  (typeof httpBodyEncodingValues)[number],
  string
> = {
  json: "JSON",
  xml: "XML",
};

const statusCodeGroupDefinitions = [
  { key: "1xx", heading: "1xx Informational", start: 100, end: 199 },
  { key: "2xx", heading: "2xx Success", start: 200, end: 299 },
  { key: "3xx", heading: "3xx Redirection", start: 300, end: 399 },
  { key: "4xx", heading: "4xx Client Error", start: 400, end: 499 },
  { key: "5xx", heading: "5xx Server Error", start: 500, end: 599 },
];

const statusCodeGroups = statusCodeGroupDefinitions
  .map((group) => {
    const groupCodes = httpStatusCodes.filter(
      (status) => status.code >= group.start && status.code <= group.end,
    );

    return {
      heading: group.heading,
      options: [
        { label: `All ${group.key}`, value: `group:${group.key}` },
        ...groupCodes.map((status) => ({
          label: `${status.code} ${status.label}${
            status.deprecated ? " (Deprecated)" : ""
          }`,
          value: String(status.code),
        })),
      ],
    };
  })
  .filter((group) => group.options.length > 1);

type HttpMonitorFormValues = Extract<MonitorFormValues, { type: "http" }>;

export default function HttpMonitorSettings() {
  const {
    register,
    control,
    formState: { errors },
  } = useFormContext<HttpMonitorFormValues>();
  const monitorType = useWatch({ control, name: "type" });
  const selectedStatusGroupKeysRef = React.useRef<string[]>([]);

  const headers = useFieldArray({
    control,
    name: "http.headers",
  });

  const bodyEncoding = useWatch({ control, name: "http.bodyEncoding" });

  if (monitorType !== "http") {
    return null;
  }

  const acceptedStatusCodeErrors = Array.isArray(
    errors.http?.acceptedStatusCodes,
  )
    ? errors.http?.acceptedStatusCodes
    : [errors.http?.acceptedStatusCodes];

  return (
    <Card>
      <CardHeader>
        <CardTitle>HTTP settings</CardTitle>
        <CardDescription>
          Configure the endpoint, request options, and response validation.
        </CardDescription>
      </CardHeader>

      <CardContent>
        <FieldSet>
          <FieldGroup>
            <Field>
              <FieldLabel htmlFor="http-url">Request URL</FieldLabel>
              <Input
                id="http-url"
                placeholder="https://api.kymarium.com/health"
                autoComplete="off"
                aria-invalid={!!errors.http?.url}
                {...register("http.url")}
              />
              <FieldError errors={[errors.http?.url]} />
              <FieldDescription>
                Include the full scheme and path for the target endpoint.
              </FieldDescription>
            </Field>

            <Accordion type="single" collapsible>
              <AccordionItem value="item-1">
                <AccordionTrigger className="text-lg">
                  Advanced Setting
                </AccordionTrigger>
                <AccordionContent>
                  <FieldGroup>
                    <div className="grid gap-6 md:grid-cols-2">
                      <Field>
                        <FieldLabel htmlFor="http-method">
                          HTTP method
                        </FieldLabel>
                        <Controller
                          name="http.method"
                          control={control}
                          render={({ field }) => (
                            <Select
                              name={field.name}
                              value={field.value ?? httpMethodValues[0]}
                              onValueChange={field.onChange}
                            >
                              <SelectTrigger id="http-method">
                                <SelectValue placeholder="Select method" />
                              </SelectTrigger>
                              <SelectContent>
                                {httpMethodValues.map((value) => (
                                  <SelectItem key={value} value={value}>
                                    {methodLabels[value]}
                                  </SelectItem>
                                ))}
                              </SelectContent>
                            </Select>
                          )}
                        />
                        <FieldError errors={[errors.http?.method]} />
                        <FieldDescription>
                          Choose how the uptime check calls the endpoint.
                        </FieldDescription>
                      </Field>

                      <Field>
                        <FieldLabel htmlFor="http-timeout">
                          Request timeout
                        </FieldLabel>
                        <Input
                          id="http-timeout"
                          type="number"
                          min={0}
                          aria-invalid={!!errors.http?.requestTimeout}
                          {...register("http.requestTimeout", {
                            valueAsNumber: true,
                          })}
                        />
                        <FieldError errors={[errors.http?.requestTimeout]} />
                        <FieldDescription>Timeout in seconds.</FieldDescription>
                      </Field>
                    </div>

                    <Field>
                      <FieldLabel htmlFor="http-redirects">
                        Max redirects
                      </FieldLabel>
                      <Input
                        id="http-redirects"
                        type="number"
                        min={0}
                        max={1000}
                        aria-invalid={!!errors.http?.maxRedirects}
                        {...register("http.maxRedirects", {
                          valueAsNumber: true,
                        })}
                      />
                      <FieldError errors={[errors.http?.maxRedirects]} />
                      <FieldDescription>
                        Limit how many redirects the check will follow.
                      </FieldDescription>
                    </Field>

                    <Field>
                      <FieldTitle className="flex items-center justify-between">
                        Headers
                        <Button
                          type="button"
                          variant="outline"
                          size="sm"
                          onClick={() => headers.append({ key: "", value: "" })}
                        >
                          <Plus />
                          Add header
                        </Button>
                      </FieldTitle>
                      <FieldDescription>
                        Optional: include headers such as authorization tokens.
                      </FieldDescription>

                      <div className="mt-3 grid gap-3">
                        <div className="grid gap-3 mr-8 grid-cols-[minmax(0,1fr)_minmax(0,1fr)_auto]">
                          <FieldLabel className="text-xs">
                            Header name
                          </FieldLabel>
                          <FieldLabel className="text-xs">
                            Header value
                          </FieldLabel>
                        </div>
                        {headers.fields.length === 0 ? (
                          <p className="text-sm text-muted-foreground">
                            No headers added yet.
                          </p>
                        ) : (
                          headers.fields.map((field, index) => (
                            <div
                              key={field.id}
                              className="grid gap-3 grid-cols-[minmax(0,1fr)_minmax(0,1fr)_auto]"
                            >
                              <div className="grid gap-2">
                                <Input
                                  placeholder="Authorization"
                                  aria-invalid={
                                    !!errors.http?.headers?.[index]?.key
                                  }
                                  {...register(`http.headers.${index}.key`)}
                                />
                                <FieldError
                                  errors={[errors.http?.headers?.[index]?.key]}
                                />
                              </div>

                              <div className="grid gap-2">
                                <Input
                                  placeholder="Bearer ..."
                                  aria-invalid={
                                    !!errors.http?.headers?.[index]?.value
                                  }
                                  {...register(`http.headers.${index}.value`)}
                                />
                                <FieldError
                                  errors={[
                                    errors.http?.headers?.[index]?.value,
                                  ]}
                                />
                              </div>

                              <div className="flex items-start">
                                <Button
                                  type="button"
                                  variant="ghost"
                                  size="icon"
                                  onClick={() => headers.remove(index)}
                                  aria-label="Remove header"
                                >
                                  <Trash2 className="h-4 w-4" />
                                </Button>
                              </div>
                            </div>
                          ))
                        )}
                      </div>
                    </Field>

                    <Field>
                      <FieldLabel htmlFor="http-body-encoding">Body</FieldLabel>
                      <FieldDescription>
                        Include a request body choose your encoding and type
                        your body.
                      </FieldDescription>
                      <Controller
                        name="http.bodyEncoding"
                        control={control}
                        render={({ field }) => (
                          <Select
                            name={field.name}
                            value={field.value ?? "none"}
                            onValueChange={(value) =>
                              field.onChange(
                                value === "none" ? undefined : value,
                              )
                            }
                          >
                            <SelectTrigger id="http-body-encoding">
                              <SelectValue placeholder="Select encoding" />
                            </SelectTrigger>
                            <SelectContent>
                              <SelectItem value="none">None</SelectItem>
                              {httpBodyEncodingValues.map((value) => (
                                <SelectItem key={value} value={value}>
                                  {bodyEncodingLabels[value]}
                                </SelectItem>
                              ))}
                            </SelectContent>
                          </Select>
                        )}
                      />
                      <FieldError errors={[errors.http?.bodyEncoding]} />
                    </Field>

                    <Field>
                      <div className="grid gap-2">
                        <Textarea
                          id="http-body"
                          rows={4}
                          disabled={!bodyEncoding}
                          aria-invalid={!!errors.http?.body}
                          {...register("http.body")}
                        />
                        <FieldError errors={[errors.http?.body]} />
                      </div>
                    </Field>

                    <Field>
                      <FieldTitle>Accepted status codes</FieldTitle>
                      <FieldDescription>
                        Define which status codes count as an uptime success.
                      </FieldDescription>

                      <div className="mt-3 grid gap-2">
                        <Controller
                          name="http.acceptedStatusCodes"
                          control={control}
                          render={({ field }) => (
                            <MultiSelect
                              hideSelectAll
                              options={statusCodeGroups}
                              defaultValue={(field.value ?? []).map(
                                (value: string) => String(value),
                              )}
                              resetOnDefaultValueChange={false}
                              onValueChange={(values) => {
                                const nextGroupKeys = values
                                  .filter((value) => value.startsWith("group:"))
                                  .map((value) => value.slice("group:".length));
                                const prevGroupKeys =
                                  selectedStatusGroupKeysRef.current;
                                const removedGroupKeys = prevGroupKeys.filter(
                                  (key) => !nextGroupKeys.includes(key),
                                );

                                const expandedValues = values.flatMap(
                                  (value) => {
                                    if (!value.startsWith("group:")) {
                                      return [value];
                                    }

                                    const groupKey = value.slice(
                                      "group:".length,
                                    );
                                    const group =
                                      statusCodeGroupDefinitions.find(
                                        (item) => item.key === groupKey,
                                      );
                                    if (!group) {
                                      return [];
                                    }

                                    return httpStatusCodes
                                      .filter(
                                        (status) =>
                                          status.code >= group.start &&
                                          status.code <= group.end,
                                      )
                                      .map((status) => String(status.code));
                                  },
                                );

                                let nextValues = Array.from(
                                  new Set(expandedValues),
                                );

                                if (removedGroupKeys.length > 0) {
                                  removedGroupKeys.forEach((groupKey) => {
                                    const group =
                                      statusCodeGroupDefinitions.find(
                                        (item) => item.key === groupKey,
                                      );
                                    if (!group) return;

                                    nextValues = nextValues.filter((value) => {
                                      const code = Number(value);
                                      return (
                                        Number.isNaN(code) ||
                                        code < group.start ||
                                        code > group.end
                                      );
                                    });
                                  });
                                }

                                selectedStatusGroupKeysRef.current =
                                  nextGroupKeys;

                                field.onChange(
                                  nextValues.map((value) => Number(value)),
                                );
                              }}
                              placeholder="Select status codes"
                              maxCount={6}
                              searchable
                              singleLine
                            />
                          )}
                        />
                        <FieldError errors={acceptedStatusCodeErrors} />
                        <p className="text-sm text-muted-foreground">
                          Leave empty to accept whatever the endpoint returns.
                        </p>
                      </div>
                    </Field>

                    <Field>
                      <FieldTitle>Response validation</FieldTitle>
                      <FieldDescription>
                        Customize how this monitor evaluates the response.
                      </FieldDescription>

                      <div className="flex flex-col gap-4">
                        <label className="flex items-center gap-2 text-sm">
                          <Controller
                            name="http.upsideDownMode"
                            control={control}
                            render={({ field }) => (
                              <Checkbox
                                checked={field.value}
                                onCheckedChange={(checked) =>
                                  field.onChange(Boolean(checked))
                                }
                              />
                            )}
                          />
                          Treat failures as success (upside-down mode)
                        </label>

                        <label className="flex items-center gap-2 text-sm">
                          <Controller
                            name="http.certificateExpiryNotification"
                            control={control}
                            render={({ field }) => (
                              <Checkbox
                                checked={field.value}
                                onCheckedChange={(checked) =>
                                  field.onChange(Boolean(checked))
                                }
                              />
                            )}
                          />
                          Notify on TLS certificate expiry
                        </label>

                        <label className="flex items-center gap-2 text-sm">
                          <Controller
                            name="http.ignoreTLSError"
                            control={control}
                            render={({ field }) => (
                              <Checkbox
                                checked={field.value}
                                onCheckedChange={(checked) =>
                                  field.onChange(Boolean(checked))
                                }
                              />
                            )}
                          />
                          Ignore TLS errors
                        </label>
                      </div>

                      <FieldError
                        errors={[
                          errors.http?.upsideDownMode,
                          errors.http?.certificateExpiryNotification,
                          errors.http?.ignoreTLSError,
                        ]}
                      />
                    </Field>
                  </FieldGroup>
                </AccordionContent>
              </AccordionItem>
            </Accordion>
          </FieldGroup>
        </FieldSet>
      </CardContent>
    </Card>
  );
}
