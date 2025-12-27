import type { Incident } from '$lib/types';
import type { IncidentEvent, IncidentEventType } from '$lib/types/incident';
import { apiRequest } from './utils';

export type IncidentsResponse = {
	message: string;
	data: Incident[];
};

export function getIncidents(teamID: string): Promise<IncidentsResponse> {
	return apiRequest<IncidentsResponse>(`/teams/${teamID}/incidents`, {
		defaultError: 'Failed to fetch incidents'
	});
}

export type IncidentResponse = {
	message: string;
	data: Incident;
};

export function getIncident(teamID: string, incidentID: string): Promise<IncidentResponse> {
	return apiRequest<IncidentResponse>(`/teams/${teamID}/incidents/${incidentID}`, {
		defaultError: 'Failed to fetch incident'
	});
}

export type IncidentEventsResponse = {
	message: string;
	data: IncidentEvent[];
};

export function getIncidentEvents(
	teamID: string,
	incidentID: string
): Promise<IncidentEventsResponse> {
	return apiRequest<IncidentEventsResponse>(`/teams/${teamID}/incidents/${incidentID}/events`, {
		defaultError: 'Failed to fetch incident timeline'
	});
}

export type IncidentEventCreateRequest = {
	message: string;
	event_type?: IncidentEventType;
};

export type IncidentEventCreateResponse = {
	message: string;
	data: IncidentEvent;
};

export function createIncidentEvent(
	teamID: string,
	incidentID: string,
	payload: IncidentEventCreateRequest
): Promise<IncidentEventCreateResponse> {
	return apiRequest<IncidentEventCreateResponse>(`/teams/${teamID}/incidents/${incidentID}/events`, {
		method: 'POST',
		body: payload,
		defaultError: 'Failed to add timeline update'
	});
}

export type IncidentStatusUpdateRequest = {
	status: Incident['status'];
	message?: string;
	public?: boolean;
};

export type IncidentStatusUpdateResponse = {
	message: string;
	data: {
		incident: Incident;
		event: IncidentEvent;
	};
};

export function updateIncidentStatus(
	teamID: string,
	incidentID: string,
	payload: IncidentStatusUpdateRequest
): Promise<IncidentStatusUpdateResponse> {
	return apiRequest<IncidentStatusUpdateResponse>(`/teams/${teamID}/incidents/${incidentID}/status`, {
		method: 'POST',
		body: payload,
		defaultError: 'Failed to update incident status'
	});
}

export type IncidentUpdateRequest = {
	title?: string | null;
	public?: boolean;
	auto_resolve?: boolean;
};

export type IncidentUpdateResponse = {
	message: string;
	data: Incident;
};

export function updateIncident(
	teamID: string,
	incidentID: string,
	payload: IncidentUpdateRequest
): Promise<IncidentUpdateResponse> {
	return apiRequest<IncidentUpdateResponse>(`/teams/${teamID}/incidents/${incidentID}`, {
		method: 'PATCH',
		body: payload,
		defaultError: 'Failed to update incident'
	});
}

export type IncidentCreateRequest = {
	title?: string | null;
	status?: Incident['status'];
	severity?: Incident['severity'];
	message?: string;
	started_at?: string;
	public?: boolean;
	auto_resolve?: boolean;
	monitor_ids: string[];
};

export type IncidentCreateResponse = {
	message: string;
	data: {
		incident: Incident;
		event: IncidentEvent;
	};
};

export function createIncident(
	teamID: string,
	payload: IncidentCreateRequest
): Promise<IncidentCreateResponse> {
	return apiRequest<IncidentCreateResponse>(`/teams/${teamID}/incidents`, {
		method: 'POST',
		body: payload,
		defaultError: 'Failed to create incident'
	});
}
