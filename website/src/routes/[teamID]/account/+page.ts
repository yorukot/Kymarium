import { getAccounts, getUser, listSessions } from '$lib/api/user';
import type { Account, Session, User } from '$lib/types';
import type { PageLoad } from './$types';

export type AccountPageData = {
	user: User;
	accounts: Account[];
	sessions: Session[];
};

export const load: PageLoad<AccountPageData> = async () => {
	const [userResponse, accountResponse, sessionResponse] = await Promise.all([
		getUser(),
		getAccounts(),
		listSessions()
	]);

	return {
		user: userResponse.data,
		accounts: accountResponse.data,
		sessions: sessionResponse.data
	};
};

export const ssr = false;
