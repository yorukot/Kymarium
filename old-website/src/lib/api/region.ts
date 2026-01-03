import type { Region } from '$lib/types';
import { apiRequest, type ApiResponse } from './utils';

export type RegionListResponse = ApiResponse<Region[]>;

export function getRegions(): Promise<RegionListResponse> {
	return apiRequest<RegionListResponse>('/regions', {
		defaultError: 'Failed to fetch regions'
	});
}
