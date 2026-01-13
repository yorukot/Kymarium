import { z } from "zod";

import type { TeamRole } from "@/lib/schemas/team-member";

export const inviteRoleValues = ["admin", "member", "viewer"] as const;

export const createTeamInviteSchema = z.object({
  email: z
    .string()
    .trim()
    .email("Enter a valid email address.")
    .max(255, "Email is too long."),
  role: z.enum(inviteRoleValues),
});

export type CreateTeamInviteValues = z.infer<typeof createTeamInviteSchema>;

export type TeamInviteRawData = {
  id: string;
  team_id: string;
  invited_by: string;
  invited_to: string;
  invited_email: string;
  role: string;
  status: string;
  expires_at: string;
  accepted_at?: string | null;
  rejected_at?: string | null;
  canceled_at?: string | null;
  updated_at?: string;
  created_at?: string;
};

export type TeamInviteWithTeamRawData = TeamInviteRawData & {
  team_name: string;
};

export type TeamInvite = {
  id: string;
  teamID: string;
  invitedBy: string;
  invitedTo: string;
  invitedEmail: string;
  role: TeamRole;
  status: string;
  expiresAt: string;
  acceptedAt?: string | null;
  rejectedAt?: string | null;
  canceledAt?: string | null;
  updatedAt?: string;
  createdAt?: string;
};

export type TeamInviteWithTeam = TeamInvite & {
  teamName: string;
};
