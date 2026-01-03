import type {
	PublicStatusPageData,
	StatusPage,
	StatusPageElementType,
	StatusPageWithElements
} from '$lib/types';
import { apiRequest, publicApiRequest, type ApiDefaultBody, type ApiResponse } from './utils';

export type StatusPageListResponse = ApiResponse<StatusPageWithElements[]>;

export type StatusPageSingleResponse = ApiResponse<StatusPageWithElements>;

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

export type StatusPageCreateResponse = ApiResponse<StatusPageWithElements>;

export type StatusPageUpdateResponse = ApiResponse<StatusPageWithElements>;

export type StatusPageDeleteResponse = ApiDefaultBody;

export type PublicStatusPageResponse = ApiResponse<PublicStatusPageData>;

export type StatusPageModelResponse = ApiResponse<StatusPage>;

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
