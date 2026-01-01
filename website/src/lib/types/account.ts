// ============================================================================
// Account Types
// ============================================================================

export type Provider = 'email' | 'google';

export interface Account {
	id: string;
	provider: Provider;
	providerUserId: string;
	userId: string;
	email: string;
	isPrimary: boolean;
	createdAt: string;
	updatedAt: string;
}

export interface OAuthToken {
	accountId: string;
	accessToken: string;
	refreshToken?: string;
	expiry: string;
	tokenType: string;
	provider: Provider;
	createdAt: string;
	updatedAt: string;
}
