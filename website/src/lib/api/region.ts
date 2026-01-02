import type { Region } from '$lib/types';
import { apiRequest } from './utils';

export type RegionListResponse = {
	message: string;
	data: Region[];
};

export function getRegions(): Promise<RegionListResponse> {
	return apiRequest<RegionListResponse>('/regions', {
		defaultError: 'Failed to fetch regions'
	});
}
