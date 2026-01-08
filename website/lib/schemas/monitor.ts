import { z } from "zod";

export const monitorTypeValues = ["http", "ping"] as const;
export const httpMethodValues = [
  "GET",
  "POST",
  "PUT",
  "DELETE",
  "PATCH",
  "HEAD",
  "OPTIONS",
] as const;

export const httpBodyEncodingValues = ["json", "xml"] as const;

export type HttpStatusCode = {
  code: number;
  label: string;
  deprecated?: boolean;
};

export const httpStatusCodes: HttpStatusCode[] = [
  { code: 100, label: "Continue" },
  { code: 101, label: "Switching Protocols" },
  { code: 103, label: "Early Hints" },
  { code: 104, label: "Upload Resumption Supported (Temporary)" },
  { code: 200, label: "OK" },
  { code: 201, label: "Created" },
  { code: 202, label: "Accepted" },
  { code: 203, label: "Non-Authoritative Information" },
  { code: 204, label: "No Content" },
  { code: 205, label: "Reset Content" },
  { code: 206, label: "Partial Content" },
  { code: 207, label: "Multi-Status" },
  { code: 208, label: "Already Reported" },
  { code: 226, label: "IM Used" },
  { code: 300, label: "Multiple Choices" },
  { code: 301, label: "Moved Permanently" },
  { code: 302, label: "Found" },
  { code: 303, label: "See Other" },
  { code: 304, label: "Not Modified" },
  { code: 307, label: "Temporary Redirect" },
  { code: 308, label: "Permanent Redirect" },
  { code: 400, label: "Bad Request" },
  { code: 401, label: "Unauthorized" },
  { code: 402, label: "Payment Required" },
  { code: 403, label: "Forbidden" },
  { code: 404, label: "Not Found" },
  { code: 405, label: "Method Not Allowed" },
  { code: 406, label: "Not Acceptable" },
  { code: 407, label: "Proxy Authentication Required" },
  { code: 408, label: "Request Timeout" },
  { code: 409, label: "Conflict" },
  { code: 410, label: "Gone" },
  { code: 411, label: "Length Required" },
  { code: 412, label: "Precondition Failed" },
  { code: 413, label: "Content Too Large" },
  { code: 414, label: "URI Too Long" },
  { code: 415, label: "Unsupported Media Type" },
  { code: 416, label: "Range Not Satisfiable" },
  { code: 417, label: "Expectation Failed" },
  { code: 418, label: "Unused" },
  { code: 421, label: "Misdirected Request" },
  { code: 422, label: "Unprocessable Content" },
  { code: 423, label: "Locked" },
  { code: 424, label: "Failed Dependency" },
  { code: 425, label: "Too Early" },
  { code: 426, label: "Upgrade Required" },
  { code: 428, label: "Precondition Required" },
  { code: 429, label: "Too Many Requests" },
  { code: 431, label: "Request Header Fields Too Large" },
  { code: 451, label: "Unavailable For Legal Reasons" },
  { code: 500, label: "Internal Server Error" },
  { code: 501, label: "Not Implemented" },
  { code: 502, label: "Bad Gateway" },
  { code: 503, label: "Service Unavailable" },
  { code: 504, label: "Gateway Timeout" },
  { code: 505, label: "HTTP Version Not Supported" },
  { code: 506, label: "Variant Also Negotiates" },
  { code: 507, label: "Insufficient Storage" },
  { code: 508, label: "Loop Detected" },
  { code: 510, label: "Not Extended (Obsoleted)" },
  { code: 511, label: "Network Authentication Required" },
];

export const monitorTypes = [
  { value: monitorTypeValues[0], label: "HTTP" },
  { value: monitorTypeValues[1], label: "Ping" },
] as const;

export const httpConfigSchema = z.object({
  url: z.string().url("Enter a valid URL."),
  method: z.enum(httpMethodValues),
  maxRedirects: z
    .number()
    .int()
    .min(0, "Max redirects must be 0 or higher.")
    .max(1000, "Max redirects must be 1000 or less."),
  requestTimeout: z
    .number()
    .int()
    .min(0, "Request timeout must be 0 or higher."),
  headers: z.array(
    z.object({
      key: z
        .string()
        .min(1, "Header name is required.")
        .max(128, "Header name is too long.")
        .regex(
          /^[!#$%&'*+\-.^_`|~0-9A-Za-z]+$/,
          "Invalid header name (must be a valid HTTP token).",
        ),
      value: z.string().min(1, "Header value is required."),
    }),
  ),
  bodyEncoding: z.enum(httpBodyEncodingValues).optional(),
  body: z.string().max(1000000, "Body must be 1,000,000 characters or less."),
  acceptedStatusCodes: z.array(
    z
      .number()
      .int()
      .min(100, "Status codes must be between 100 and 599.")
      .max(599, "Status codes must be between 100 and 599."),
  ),
  upsideDownMode: z.boolean(),
  certificateExpiryNotification: z.boolean(),
  ignoreTLSError: z.boolean(),
});

export const DEFAULT_HTTP: MonitorFormValues["http"] = {
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
};

export const pingConfigSchema = z.object({
  host: z.string().min(1, "Host is required."),
  timeoutSeconds: z.number().int().min(0, "Timeout must be 0 or higher."),
  packetSize: z
    .number()
    .int()
    .min(1, "Packet size must be at least 1.")
    .max(65000, "Packet size must be 65000 or less.")
    .optional(),
});

export const DEFAULT_PING: MonitorFormValues["ping"] = {
  host: "",
  timeoutSeconds: 5,
  packetSize: undefined,
};

const baseSchema = z.object({
  name: z.string().min(1, "Monitor name is required."),
  interval: z
    .number()
    .int()
    .min(2, "Interval must be at least 2 seconds.")
    .max(2592000),
  failureThreshold: z
    .number()
    .int()
    .min(1, "Failure threshold must be at least 1."),
  recoveryThreshold: z
    .number()
    .int()
    .min(1, "Recovery threshold must be at least 1."),
  regions: z.array(z.string()).min(1, "Select at least one region."),
  notifications: z.array(z.string()),
  config: z.any().optional(),
});

export const monitorSchema = z.discriminatedUnion("type", [
  baseSchema.extend({
    type: z.literal("http"),
    http: httpConfigSchema,
    ping: z.undefined().optional(),
  }),
  baseSchema.extend({
    type: z.literal("ping"),
    ping: pingConfigSchema,
    http: z.undefined().optional(),
  }),
]);

export type MonitorFormValues = z.infer<typeof monitorSchema>;

export type HttpMonitorFormConfig = z.infer<typeof httpConfigSchema>;
export type PingMonitorFormConfig = z.infer<typeof pingConfigSchema>;

export type MonitorResponse = {
  id: string;
  teamId: string;
  updatedAt: string;
  createdAt: string;
};
export type CreateMonitorResponse = {
  message: string;
  data?: MonitorResponse;
};

export type MonitorRawData = {
  id: string;
  team_id: string;
  name: string;
  type: string;
  config?: Record<string, unknown> | null;
  interval: number;
  status: string;
  uptime_sli_30?: number | null;
  last_checked: string;
  next_check: string;
  failure_threshold: number;
  recovery_threshold: number;
  regions: string[];
  notification: string[];
  incidents?: Array<Record<string, unknown>>;
  updated_at: string;
  created_at: string;
};

export type MonitorListItem = {
  id: string;
  name: string;
  type: string;
  status: string;
  targetLabel: string;
  targetValue: string;
  lastChecked: string;
  uptimeSLI30?: number;
};
