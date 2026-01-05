"use client";

import React from "react";
import { cn } from "@/lib/utils";

import { Button } from "@/components/ui/button";
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import {
  Field,
  FieldDescription,
  FieldError,
  FieldGroup,
  FieldLabel,
  FieldSeparator,
} from "@/components/ui/field";
import { Input } from "@/components/ui/input";

import { useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";

import { FormDevTools } from "@/components/devtools/form-dev-tools";
import {
  SignupFormValues,
  signupPayloadSchema,
  signupSchema,
} from "@/lib/schemas/auth";
import { signup } from "@/lib/api/auth";
import { ApiError } from "@/lib/api/client";
import { applyServerFieldErrors } from "@/lib/api/error";
import { Spinner } from "@/components/ui/spinner";
import { useRouter } from "next/navigation";

export function SignupForm({
  className,
  ...props
}: React.ComponentProps<"div">) {
  const router = useRouter();
  const form = useForm<SignupFormValues>({
    resolver: zodResolver(signupSchema),
    defaultValues: {
      displayName: "",
      email: "",
      password: "",
      confirmPassword: "",
    },
    mode: "onSubmit",
  });

  const onSubmit = async (values: SignupFormValues) => {
    form.clearErrors();

    const parsed = signupPayloadSchema.safeParse(values);
    if (!parsed.success) {
      form.setError("root", {
        type: "validate",
        message: "Invalid form data. Please try again.",
      });
      return;
    }

    try {
      const response = await signup(parsed.data);
      form.reset();
      if (response.data.message.includes("Verification email sent")) {
        router.push(`/email-sent?email=${encodeURIComponent(values.email)}`);
      }
    } catch (e) {
      if (e instanceof ApiError) {
        const hasFieldErrors = applyServerFieldErrors(form.setError, e.body);

        if (!hasFieldErrors) {
          form.setError("root", {
            type: "server",
            message:
              e.status >= 500
                ? "Server error. Please try again later."
                : e.message,
          });
        }

        return;
      }

      form.setError("root", {
        type: "network",
        message: "Network error. Please try again.",
      });
    }
  };

  const {
    register,
    handleSubmit,
    formState: { errors, isSubmitting },
  } = form;

  return (
    <div className={cn("flex flex-col gap-6", className)} {...props}>
      <Card>
        <CardHeader className="text-center">
          <CardTitle className="text-xl">Create your account</CardTitle>
          <CardDescription>
            Enter your email below to create your account
          </CardDescription>
        </CardHeader>

        <CardContent>
          {/* noValidate：不要讓瀏覽器 required 擋住你的 zod 訊息 */}
          <form noValidate onSubmit={handleSubmit(onSubmit)}>
            <FieldGroup>
              {/* OAuth button */}
              <Field>
                <Button variant="outline" type="button" className="w-full">
                  <svg
                    xmlns="http://www.w3.org/2000/svg"
                    viewBox="0 0 24 24"
                    aria-hidden="true"
                  >
                    <path
                      d="M12.48 10.92v3.28h7.84c-.24 1.84-.853 3.187-1.787 4.133-1.147 1.147-2.933 2.4-6.053 2.4-4.827 0-8.6-3.893-8.6-8.72s3.773-8.72 8.6-8.72c2.6 0 4.507 1.027 5.907 2.347l2.307-2.307C18.747 1.44 16.133 0 12.48 0 5.867 0 .307 5.387.307 12s5.56 12 12.173 12c3.573 0 6.267-1.173 8.373-3.36 2.16-2.16 2.84-5.213 2.84-7.667 0-.76-.053-1.467-.173-2.053H12.48z"
                      fill="currentColor"
                    />
                  </svg>
                  Create account with Google
                </Button>
              </Field>

              <FieldSeparator className="*:data-[slot=field-separator-content]:bg-card">
                Or continue with
              </FieldSeparator>

              {/* Display Name */}
              <Field>
                <FieldLabel htmlFor="displayName">Display Name</FieldLabel>
                <Input
                  id="displayName"
                  type="text"
                  placeholder="Kymarium Kot"
                  aria-invalid={!!errors.displayName}
                  {...register("displayName")}
                />
                <FieldError errors={[errors.displayName]} />
              </Field>

              {/* Email */}
              <Field>
                <FieldLabel htmlFor="email">Email</FieldLabel>
                <Input
                  id="email"
                  type="email"
                  placeholder="m@example.com"
                  aria-invalid={!!errors.email}
                  autoComplete="email"
                  {...register("email")}
                />
                <FieldError errors={[errors.email]} />
              </Field>

              {/* Password + Confirm */}
              <Field>
                <div className="grid grid-cols-1 gap-4 md:grid-cols-2">
                  <Field>
                    <FieldLabel htmlFor="password">Password</FieldLabel>
                    <Input
                      id="password"
                      type="password"
                      aria-invalid={!!errors.password}
                      autoComplete="new-password"
                      {...register("password")}
                    />
                    <FieldError errors={[errors.password]} />
                  </Field>

                  <Field>
                    <FieldLabel htmlFor="confirmPassword">
                      Confirm Password
                    </FieldLabel>
                    <Input
                      id="confirmPassword"
                      type="password"
                      aria-invalid={!!errors.confirmPassword}
                      autoComplete="new-password"
                      {...register("confirmPassword")}
                    />
                    <FieldError errors={[errors.confirmPassword]} />
                  </Field>
                </div>

                <FieldDescription>
                  Must be at least 8 characters long.
                </FieldDescription>
              </Field>

              {/* Root error + submit */}
              <Field>
                <FieldError errors={[errors.root]} />

                <Button
                  type="submit"
                  className="w-full"
                  disabled={isSubmitting}
                >
                  {isSubmitting ? (
                    <>
                      <Spinner className="mr-2 h-4 w-4 animate-spin" />
                      Creating...
                    </>
                  ) : (
                    "Create Account"
                  )}
                </Button>

                <FieldDescription className="text-center">
                  Already have an account? <a href="/login">Login in</a>
                </FieldDescription>
              </Field>
            </FieldGroup>
          </form>
        </CardContent>
      </Card>

      {process.env.NODE_ENV === "development" && (
        <FormDevTools control={form.control} />
      )}
    </div>
  );
}
