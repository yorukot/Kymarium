// ============================================================================
// Team Types
// ============================================================================

export type MemberRole = 'owner' | 'admin' | 'member' | 'viewer';

export interface Team {
	id: string;
	name: string;
	updatedAt: string;
	createdAt: string;
}

export interface TeamMember {
	id: string;
	teamId: string;
	userId: string;
	role: MemberRole;
	updatedAt: string;
	createdAt: string;
}

export interface TeamInvite {
	id: string;
	teamId: string;
	invitedBy: string;
	invitedTo: string;
	invitedEmail: string;
	role: MemberRole;
	status: 'pending' | 'accepted' | 'rejected' | 'canceled';
	expiresAt: string;
	acceptedAt?: string | null;
	rejectedAt?: string | null;
	canceledAt?: string | null;
	updatedAt: string;
	createdAt: string;
}

export interface TeamInviteWithTeam extends TeamInvite {
	teamName: string;
}

export interface TeamMemberWithUser extends TeamMember {
	displayName: string;
	email: string;
	avatar?: string;
}

export interface TeamWithRole extends Team {
	role: MemberRole;
}
