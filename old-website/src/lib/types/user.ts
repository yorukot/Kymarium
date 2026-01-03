// ============================================================================
// User Types
// ============================================================================

export interface User {
	id: string;
	passwordHash?: string;
	displayName: string;
	avatar?: string;
	verified: boolean;
	createdAt: string;
	updatedAt: string;
}

export interface Session {
	id: string;
	userAgent?: string;
	ip?: string;
	createdAt: string;
}
