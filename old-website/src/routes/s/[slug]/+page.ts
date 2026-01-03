import { error } from '@sveltejs/kit';
import { getPublicStatusPage } from '$lib/api/status-page';
import type { PublicStatusPageData } from '$lib/types';
import type { PageLoad } from './$types';

export type PublicStatusPageLoad = {
	statusPage: PublicStatusPageData;
};

export const load: PageLoad<PublicStatusPageLoad> = async ({ params }) => {
	const { slug } = params;

	try {
		const response = await getPublicStatusPage(slug);
		return { statusPage: response.data };
	} catch (err) {
		const message = err instanceof Error ? err.message : 'Status page not found';
		throw error(404, message);
	}
};
