export type TeamRole = "owner" | "admin" | "member" | "viewer" | (string & {});

export type TeamMemberRawData = {
  id: string;
  team_id: string;
  user_id: string;
  role: string;
  display_name?: string;
  email?: string;
  avatar?: string | null;
  updated_at?: string;
  created_at?: string;
};

export type TeamMember = {
  id: string;
  teamID: string;
  userID: string;
  role: TeamRole;
  displayName: string;
  email: string;
  avatar?: string | null;
  updatedAt?: string;
  createdAt?: string;
};

