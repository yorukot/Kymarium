import type {
	PublicStatusPageData,
	StatusPage,
	StatusPageElementType,
	StatusPageWithElements
} from '$lib/types';
import { apiRequest, publicApiRequest } from './utils';

export type StatusPageListResponse = {
	message: string;
	data: StatusPageWithElements[];
};

export type StatusPageSingleResponse = {
	message: string;
	data: StatusPageWithElements;
};

export type StatusPageUpsertRequest = {
	title: string;
	slug: string;
	icon?: string | null;
	elements?: {
		id?: string;
		name: string;
		type: StatusPageElementType;
		sort_order: number;
		monitor: boolean;
		monitor_id?: string | null;
		monitors?: {
			id?: string;
			monitor_id: string;
			group_id?: string | null;
			name: string;
			type: StatusPageElementType;
			sort_order: number;
		}[];
	}[];
	groups?: {
		id?: string;
		name: string;
		type: StatusPageElementType;
		sort_order: number;
	}[];
	monitors?: {
		id?: string;
		monitor_id: string;
		group_id?: string | null;
		name: string;
		type: StatusPageElementType;
		sort_order: number;
	}[];
};

export type StatusPageCreateResponse = {
	message: string;
	data: StatusPageWithElements;
};

export type StatusPageUpdateResponse = {
	message: string;
	data: StatusPageWithElements;
};

export type StatusPageDeleteResponse = {
	message: string;
};

export type PublicStatusPageResponse = {
	message: string;
	data: PublicStatusPageData;
};

export type StatusPageModelResponse = {
	message: string;
	data: StatusPage;
};

export function getStatusPages(teamID: string): Promise<StatusPageListResponse> {
	return apiRequest<StatusPageListResponse>(`/teams/${teamID}/status-pages`, {
		defaultError: 'Failed to fetch status pages'
	});
}

export function getStatusPage(teamID: string, statusPageID: string): Promise<StatusPageSingleResponse> {
	return apiRequest<StatusPageSingleResponse>(
		`/teams/${teamID}/status-pages/${statusPageID}`,
		{
			defaultError: 'Failed to fetch status page'
		}
	);
}

export function createStatusPage(
	teamID: string,
	payload: StatusPageUpsertRequest
): Promise<StatusPageCreateResponse> {
	return apiRequest<StatusPageCreateResponse>(`/teams/${teamID}/status-pages`, {
		method: 'POST',
		body: payload,
		defaultError: 'Failed to create status page'
	});
}

export function updateStatusPage(
	teamID: string,
	statusPageID: string,
	payload: StatusPageUpsertRequest
): Promise<StatusPageUpdateResponse> {
	return apiRequest<StatusPageUpdateResponse>(`/teams/${teamID}/status-pages/${statusPageID}`, {
		method: 'PUT',
		body: payload,
		defaultError: 'Failed to update status page'
	});
}

export function deleteStatusPage(
	teamID: string,
	statusPageID: string
): Promise<StatusPageDeleteResponse> {
	return apiRequest<StatusPageDeleteResponse>(`/teams/${teamID}/status-pages/${statusPageID}`, {
		method: 'DELETE',
		defaultError: 'Failed to delete status page'
	});
}

export function getPublicStatusPage(slug: string): Promise<PublicStatusPageResponse> {
	return publicApiRequest<PublicStatusPageResponse>(`/status-pages/${slug}`, {
		defaultError: 'Failed to fetch status page'
	});
}
