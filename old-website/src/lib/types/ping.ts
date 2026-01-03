// ============================================================================
// Ping Types
// ============================================================================

export type PingStatus = 'successful' | 'failed' | 'timeout';

export interface Ping {
	time: string;
	monitorId: string;
	regionId: string;
	latency: number;
	status: PingStatus;
}
