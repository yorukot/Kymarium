"use client";

import * as React from "react";
import Link from "next/link";
import { useParams, useRouter, useSearchParams } from "next/navigation";

import { Button } from "@/components/ui/button";
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import { Field, FieldDescription, FieldGroup } from "@/components/ui/field";
import { GalleryVerticalEnd } from "lucide-react";

const REDIRECT_SECONDS = 3;

type VerifyStatus = "success" | "expired" | "unknown";

export default function VerifyResultPage() {
  const params = useParams<{ result?: string }>();
  const searchParams = useSearchParams();
  const router = useRouter();

  const email = searchParams?.get("email")?.trim() || "";
  const rawResult = params?.result;
  const normalizedResult = Array.isArray(rawResult)
    ? rawResult[0]?.toLowerCase()
    : rawResult?.toLowerCase();

  const status: VerifyStatus =
    normalizedResult === "success"
      ? "success"
      : normalizedResult === "expired"
        ? "expired"
        : "unknown";

  const shouldRedirect = status === "success" || status === "expired";
  const [countdown, setCountdown] = React.useState(
    shouldRedirect ? REDIRECT_SECONDS : 0
  );

  React.useEffect(() => {
    setCountdown(shouldRedirect ? REDIRECT_SECONDS : 0);
  }, [shouldRedirect, status]);

  React.useEffect(() => {
    if (!shouldRedirect) return;
    if (countdown <= 0) {
      router.push("/login");
      return;
    }

    const timer = window.setTimeout(() => {
      setCountdown((prev) => Math.max(0, prev - 1));
    }, 1000);

    return () => window.clearTimeout(timer);
  }, [countdown, router, shouldRedirect]);

  const title =
    status === "success"
      ? "Email verified"
      : status === "expired"
        ? "Verification link expired"
        : "Verification link missing";

  const description =
    status === "success"
      ? email
        ? `Your email ${email} is verified. Please log in to continue.`
        : "Your email is verified. Please log in to continue."
      : status === "expired"
        ? "Your email verification link has expired. Please log in and request a new one."
        : "The verification link is missing. Please log in and request a new one.";

  return (
    <div className="bg-muted flex min-h-svh flex-col items-center justify-center gap-6 p-6 md:p-10">
      <div className="flex w-full max-w-sm flex-col gap-6">
        <a href="#" className="flex items-center gap-2 self-center font-medium">
          <div className="bg-primary text-primary-foreground flex size-6 items-center justify-center rounded-md">
            <GalleryVerticalEnd className="size-4" />
          </div>
          Acme Inc.
        </a>
        <Card>
          <CardHeader className="text-center">
            <CardTitle className="text-xl">{title}</CardTitle>
            <CardDescription>{description}</CardDescription>
          </CardHeader>
          <CardContent>
            <FieldGroup>
              {shouldRedirect && (
                <Field>
                  <FieldDescription>
                    Redirecting to login in {countdown}s...
                  </FieldDescription>
                </Field>
              )}
              <Field>
                <Button className="w-full" asChild>
                  <Link href="/login">Go to login</Link>
                </Button>
              </Field>
            </FieldGroup>
          </CardContent>
        </Card>
      </div>
    </div>
  );
}
