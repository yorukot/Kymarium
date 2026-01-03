// ============================================================================
// Auth Types
// ============================================================================

export type Provider = 'email' | 'google';

export interface RefreshToken {
	id: string;
	userId: string;
	token: string;
	userAgent?: string;
	ip: string;
	usedAt?: string;
	createdAt: string;
}

export const CookieName = {
	OAuthSession: 'oauth_session',
	RefreshToken: 'refresh_token',
	AccessToken: 'access_token'
} as const;
