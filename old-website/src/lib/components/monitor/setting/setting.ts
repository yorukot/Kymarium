import type { BodyEncoding, HTTPMethod } from '../../$lib/type/monitor-config';
import type { MonitorType } from '$lib/types';
import type { MultiSelectOption } from '$lib/components/ui/multi-select';

export type MonitorTypeSelect = {
	title: string;
	description: string;
	value: MonitorType;
};

export const monitorTypeSelectData: readonly MonitorTypeSelect[] = [
	{
		title: 'HTTP',
		description: 'Monitor an HTTP or HTTPS endpoint and check response status and latency.',
		value: 'http'
	},
	{
		title: 'Ping',
		description: 'Monitor a host using ICMP ping to check network availability and latency.',
		value: 'ping'
	}
] as const;

export type IntervalOption = {
	label: string;
	seconds: number;
};

export const intervalOptions: readonly IntervalOption[] = [
	{ label: '30s', seconds: 30 },
	{ label: '45s', seconds: 45 },
	{ label: '1m', seconds: 60 },
	{ label: '3m', seconds: 180 },
	{ label: '5m', seconds: 300 },
	{ label: '10m', seconds: 600 },
	{ label: '15m', seconds: 900 },
	{ label: '30m', seconds: 1800 },
	{ label: '1h', seconds: 3600 },
	{ label: '2h', seconds: 7200 }
] as const;

export type ThresholdOption = {
	label: string;
	value: number;
};

export const thresholdOptions: readonly ThresholdOption[] = [
	{ label: 'Immediate (after 1 check)', value: 1 },
	{ label: 'After 2 checks', value: 2 },
	{ label: 'After 3 checks', value: 3 },
	{ label: 'After 4 checks', value: 4 },
	{ label: 'After 5 checks', value: 5 }
] as const;

export const httpMethods: readonly HTTPMethod[] = [
	'GET',
	'POST',
	'PUT',
	'DELETE',
	'PATCH',
	'HEAD',
	'OPTIONS'
] as const;

export const bodyEncodingOptions: Array<{ label: string; value: BodyEncoding | '' }> = [
	{ label: 'None', value: '' },
	{ label: 'JSON', value: 'json' },
	{ label: 'XML', value: 'xml' }
];

const statusRangeOptions: MultiSelectOption[] = [
	{ label: 'Any 1xx', value: '1xx', keywords: ['1xx', 'informational'] },
	{ label: 'Any 2xx', value: '2xx', keywords: ['2xx', 'success'] },
	{ label: 'Any 3xx', value: '3xx', keywords: ['3xx', 'redirect'] },
	{ label: 'Any 4xx', value: '4xx', keywords: ['4xx', 'client'] },
	{ label: 'Any 5xx', value: '5xx', keywords: ['5xx', 'server'] }
];

const statusCodeEntries: Array<{ code: number; label: string; keywords?: string[] }> = [
	// 1xx Informational
	{ code: 100, label: 'Continue' },
	{ code: 101, label: 'Switching Protocols' },
	{ code: 102, label: 'Processing' },
	{ code: 103, label: 'Early Hints' },
	// 2xx Success
	{ code: 200, label: 'OK' },
	{ code: 201, label: 'Created' },
	{ code: 202, label: 'Accepted' },
	{ code: 203, label: 'Non-Authoritative Information' },
	{ code: 204, label: 'No Content' },
	{ code: 205, label: 'Reset Content' },
	{ code: 206, label: 'Partial Content' },
	{ code: 207, label: 'Multi-Status' },
	{ code: 208, label: 'Already Reported' },
	{ code: 226, label: 'IM Used' },
	// 3xx Redirection
	{ code: 300, label: 'Multiple Choices' },
	{ code: 301, label: 'Moved Permanently' },
	{ code: 302, label: 'Found' },
	{ code: 303, label: 'See Other' },
	{ code: 304, label: 'Not Modified' },
	{ code: 305, label: 'Use Proxy' },
	{ code: 306, label: 'Switch Proxy' },
	{ code: 307, label: 'Temporary Redirect' },
	{ code: 308, label: 'Permanent Redirect' },
	// 4xx Client Error
	{ code: 400, label: 'Bad Request' },
	{ code: 401, label: 'Unauthorized' },
	{ code: 402, label: 'Payment Required' },
	{ code: 403, label: 'Forbidden' },
	{ code: 404, label: 'Not Found' },
	{ code: 405, label: 'Method Not Allowed' },
	{ code: 406, label: 'Not Acceptable' },
	{ code: 407, label: 'Proxy Authentication Required' },
	{ code: 408, label: 'Request Timeout' },
	{ code: 409, label: 'Conflict' },
	{ code: 410, label: 'Gone' },
	{ code: 411, label: 'Length Required' },
	{ code: 412, label: 'Precondition Failed' },
	{ code: 413, label: 'Payload Too Large' },
	{ code: 414, label: 'URI Too Long' },
	{ code: 415, label: 'Unsupported Media Type' },
	{ code: 416, label: 'Range Not Satisfiable' },
	{ code: 417, label: 'Expectation Failed' },
	{ code: 418, label: "I'm a teapot" },
	{ code: 421, label: 'Misdirected Request' },
	{ code: 422, label: 'Unprocessable Content' },
	{ code: 423, label: 'Locked' },
	{ code: 424, label: 'Failed Dependency' },
	{ code: 425, label: 'Too Early' },
	{ code: 426, label: 'Upgrade Required' },
	{ code: 428, label: 'Precondition Required' },
	{ code: 429, label: 'Too Many Requests' },
	{ code: 431, label: 'Request Header Fields Too Large' },
	{ code: 451, label: 'Unavailable For Legal Reasons' },
	// 5xx Server Error
	{ code: 500, label: 'Internal Server Error' },
	{ code: 501, label: 'Not Implemented' },
	{ code: 502, label: 'Bad Gateway' },
	{ code: 503, label: 'Service Unavailable' },
	{ code: 504, label: 'Gateway Timeout' },
	{ code: 505, label: 'HTTP Version Not Supported' },
	{ code: 506, label: 'Variant Also Negotiates' },
	{ code: 507, label: 'Insufficient Storage' },
	{ code: 508, label: 'Loop Detected' },
	{ code: 510, label: 'Not Extended' },
	{ code: 511, label: 'Network Authentication Required' }
];

const commonStatusOptions: MultiSelectOption[] = statusCodeEntries.map(({ code, label, keywords }) => ({
	label: `${code} ${label}`,
	value: code.toString(),
	keywords
}));

export const successStatusCodes = statusCodeEntries
	.filter(({ code }) => code >= 200 && code < 300)
	.map(({ code }) => code.toString());

export const acceptedStatusOptions: MultiSelectOption[] = [...statusRangeOptions, ...commonStatusOptions];
