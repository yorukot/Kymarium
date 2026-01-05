"use client";

import { useState, use, useEffect } from "react";
import Link from "next/link";

import { Button } from "@/components/ui/button";
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import { resendVerifyEmail } from "@/lib/api/auth";
import { ApiError } from "@/lib/api/client";
import { GalleryVerticalEnd } from "lucide-react";

const RESEND_COOLDOWN_SECONDS = 60;

export default function EmailSentPage({
  searchParams,
}: {
  searchParams: Promise<{ email?: string }>;
}) {
  const params = use(searchParams);
  const email = params.email || "";

  const [isResending, setIsResending] = useState(false);
  const [cooldownSeconds, setCooldownSeconds] = useState(0);
  const [status, setStatus] = useState<
    { type: "success" | "error"; message: string } | undefined
  >();

  useEffect(() => {
    if (cooldownSeconds <= 0) return;

    const timer = window.setTimeout(() => {
      setCooldownSeconds((prev) => Math.max(0, prev - 1));
    }, 1000);

    return () => window.clearTimeout(timer);
  }, [cooldownSeconds]);

  const handleResend = async () => {
    if (!email || isResending || cooldownSeconds > 0) {
      return;
    }

    setIsResending(true);
    setStatus(undefined);

    try {
      await resendVerifyEmail(email);
      setStatus({
        type: "success",
        message: "Verification email resent. Please check your inbox.",
      });
      setCooldownSeconds(RESEND_COOLDOWN_SECONDS);
    } catch (error) {
      if (error instanceof ApiError) {
        setStatus({ type: "error", message: error.message });
      } else {
        setStatus({
          type: "error",
          message: "Network error. Please try again.",
        });
      }
    } finally {
      setIsResending(false);
    }
  };

  const resendLabel =
    cooldownSeconds > 0 ? `Resend email (${cooldownSeconds}s)` : "Resend email";

  const emailCopy = email
    ? `A verification email has been sent to ${email}. Please check your inbox.`
    : "A verification email has been sent. Please check your inbox.";

  return (
    <div className="bg-muted flex min-h-svh flex-col items-center justify-center gap-6 p-6 md:p-10">
      <div className="flex w-full max-w-sm flex-col gap-6">
        <a href="#" className="flex items-center gap-2 self-center font-medium">
          <div className="bg-primary text-primary-foreground flex size-6 items-center justify-center rounded-md">
            <GalleryVerticalEnd className="size-4" />
          </div>
          Acme Inc.
        </a>
        <div className="flex flex-col gap-6">
          <Card>
            <CardHeader className="text-center">
              <CardTitle className="text-xl">Check your inbox</CardTitle>
              <CardDescription>{emailCopy}</CardDescription>
            </CardHeader>
            <CardContent>
              <div className="flex flex-col gap-3">
                <Button
                  type="button"
                  onClick={handleResend}
                  disabled={isResending || !email || cooldownSeconds > 0}
                >
                  {isResending ? "Resending..." : resendLabel}
                </Button>
                <Button variant="outline" asChild>
                  <Link href="/login">Go to login</Link>
                </Button>
              </div>
              {status?.message && (
                <p
                  className={`mt-4 text-sm ${
                    status.type === "error"
                      ? "text-destructive"
                      : "text-muted-foreground"
                  }`}
                >
                  {status.message}
                </p>
              )}
            </CardContent>
          </Card>
        </div>
      </div>
    </div>
  );
}
