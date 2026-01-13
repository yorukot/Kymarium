import { apiRequest } from "@/lib/api/client";
import type { PasswordUpdatePayload, UserProfilePayload } from "@/lib/schemas/user";

type UserResponse = {
  id: string;
  displayName: string;
  avatar?: string | null;
};

type UpdateMeResponse = {
  message?: string;
  data?: UserResponse;
};

export function updateUserProfile(payload: Partial<UserProfilePayload>) {
  return apiRequest<UpdateMeResponse>("/api/users/me", {
    method: "PATCH",
    body: payload,
    defaultError: "Update profile failed",
    redirectOn401: true,
  });
}

type MessageResponse = {
  message?: string;
};

export function updatePassword(payload: PasswordUpdatePayload) {
  return apiRequest<MessageResponse>("/api/users/me/password", {
    method: "PATCH",
    body: payload,
    defaultError: "Update password failed",
    redirectOn401: true,
  });
}

