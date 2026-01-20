import type { Metadata } from "next";
import type { ReactNode } from "react";

export const metadata: Metadata = {
  title: {
    default: "Email Sent",
    template: "%s | Kymarium",
  },
};

export default function EmailSentLayout({
  children,
}: {
  children: ReactNode;
}) {
  return children;
}
