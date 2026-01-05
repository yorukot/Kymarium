import { z } from "zod";

export const signupObjectSchema = z.object({
  displayName: z.string().min(1, "Full name is required"),
  email: z.string().email("Please enter a valid email address"),
  password: z
    .string()
    .min(8, "At least 8 characters are required")
    .max(255, "Password is too long"),
  confirmPassword: z.string().min(1, "Please confirm your password"),
});

export const signupSchema = signupObjectSchema.superRefine((val, ctx) => {
  if (val.password !== val.confirmPassword) {
    ctx.addIssue({
      code: "custom",
      path: ["confirmPassword"],
      message: "Passwords do not match",
    });
  }
});

export type SignupFormValues = z.infer<typeof signupSchema>;

export const signupPayloadSchema = signupObjectSchema.omit({
  confirmPassword: true,
});

export type SignupPayload = z.infer<typeof signupPayloadSchema>;
