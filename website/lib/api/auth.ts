import { apiRequest } from "@/lib/api/client";
import type { SignupPayload } from "@/lib/schemas/auth";

export type SignupResponse = { message: string };
export type ResendVerifyResponse = { message: string };

export function signup(payload: SignupPayload) {
  return apiRequest<SignupResponse>("/api/auth/register", {
    method: "POST",
    body: payload,
    defaultError: "Signup failed",
  });
}

export function resendVerifyEmail(email: string) {
  return apiRequest<ResendVerifyResponse>("/api/auth/verify/resend", {
    method: "POST",
    body: { email },
    defaultError: "Resend verification failed",
  });
}
