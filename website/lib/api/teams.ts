import { apiRequest } from "@/lib/api/client";
import type { TeamPayload } from "@/lib/schemas/team";

export type TeamResponse = {
  id: string;
  name: string;
  updatedAt: string;
  createdAt: string;
};

export type CreateTeamResponse = {
  message: string;
  data?: TeamResponse;
};

export function createTeam(payload: TeamPayload) {
  return apiRequest<CreateTeamResponse>("/api/teams", {
    method: "POST",
    body: payload,
    defaultError: "Create team failed",
  });
}
