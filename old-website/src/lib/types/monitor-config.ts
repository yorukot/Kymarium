// ============================================================================
// Monitor Configuration Types
// ============================================================================

export type HTTPMethod = 'GET' | 'POST' | 'PUT' | 'DELETE' | 'PATCH' | 'HEAD' | 'OPTIONS';

export type BodyEncoding = 'json' | 'xml';

export interface HTTPMonitorConfig {
	url: string;
	method: HTTPMethod;
	maxRedirects: number;
	requestTimeoutSeconds: number;
	headers?: Record<string, string>;
	bodyEncoding?: BodyEncoding | '';
	body?: string;
	upsideDownMode: boolean;
	certificateExpiryNotification: boolean;
	ignoreTlsError: boolean;
	acceptedStatusCodes?: string[];
}

export interface PingMonitorConfig {
	host: string;
	timeoutSeconds: number;
	packetSize?: number;
}
