import type { Monitor, MonitorAnalytics, MonitorWithIncidents } from '../../types';
import { apiRequest } from './utils';

export type MonitorResponse = {
	message: string;
	data: MonitorWithIncidents[];
};

export type MonitorSingleResponse = {
	message: string;
	data: Monitor;
};

export type MonitorCreateRequest = {
	name: string;
	type: 'http' | 'ping';
	interval: number;
	failure_threshold: number;
	recovery_threshold: number;
	regions: string[];
	notification: string[];
	config:
		| {
			url: string;
			method: string;
			max_redirects: number;
			request_timeout: number;
			headers?: Record<string, string>;
			body_encoding?: 'json' | 'xml';
			body?: string;
			upside_down_mode: boolean;
			certificate_expiry_notification: boolean;
			ignore_tls_error: boolean;
			accepted_status_codes: number[];
		}
		| {
			host: string;
			timeout_seconds: number;
			packet_size?: number;
		};
};

export type MonitorCreateResponse = {
	message: string;
	data: Monitor;
};

export type MonitorUpdateResponse = {
	message: string;
	data: Monitor;
};

export type MonitorDeleteResponse = {
	message: string;
};

export type MonitorAnalyticsResponse = {
	message: string;
	data: MonitorAnalytics;
};

export function getMonitors(teamID: string): Promise<MonitorResponse> {
	return apiRequest<MonitorResponse>(`/teams/${teamID}/monitors`, {
		defaultError: 'Failed to fetch monitors'
	});
}

export function getMonitor(teamID: string, monitorID: string): Promise<MonitorSingleResponse> {
	return apiRequest<MonitorSingleResponse>(`/teams/${teamID}/monitors/${monitorID}`, {
		defaultError: 'Failed to fetch monitor'
	});
}

export function createMonitor(
	teamID: string,
	payload: MonitorCreateRequest
): Promise<MonitorCreateResponse> {
	return apiRequest<MonitorCreateResponse>(`/teams/${teamID}/monitors`, {
		method: 'POST',
		body: payload,
		defaultError: 'Failed to create monitor'
	});
}

export function updateMonitor(
	teamID: string,
	monitorID: string,
	payload: MonitorCreateRequest
): Promise<MonitorUpdateResponse> {
	return apiRequest<MonitorUpdateResponse>(`/teams/${teamID}/monitors/${monitorID}`, {
		method: 'PUT',
		body: payload,
		defaultError: 'Failed to update monitor'
	});
}

export function deleteMonitor(teamID: string, monitorID: string): Promise<MonitorDeleteResponse> {
	return apiRequest<MonitorDeleteResponse>(`/teams/${teamID}/monitors/${monitorID}`, {
		method: 'DELETE',
		defaultError: 'Failed to delete monitor'
	});
}

export function getMonitorAnalytics(
	teamID: string,
	monitorID: string
): Promise<MonitorAnalyticsResponse> {
	return apiRequest<MonitorAnalyticsResponse>(
		`/teams/${teamID}/monitors/${monitorID}/analytics`,
		{
			defaultError: 'Failed to fetch monitor analytics'
		}
	);
}
