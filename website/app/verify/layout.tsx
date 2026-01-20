import type { Metadata } from "next";
import type { ReactNode } from "react";

export const metadata: Metadata = {
  title: {
    default: "Verify Email",
    template: "%s | Kymarium",
  },
};

export default function VerifyLayout({
  children,
}: {
  children: ReactNode;
}) {
  return children;
}
