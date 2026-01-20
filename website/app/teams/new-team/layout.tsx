import type { Metadata } from "next";
import type { ReactNode } from "react";

export const metadata: Metadata = {
  title: {
    default: "Create Team",
    template: "%s | Kymarium",
  },
};

export default function NewTeamLayout({
  children,
}: {
  children: ReactNode;
}) {
  return children;
}
