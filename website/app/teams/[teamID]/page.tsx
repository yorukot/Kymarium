import type { Metadata } from "next";
import { redirect } from "next/navigation";

export const metadata: Metadata = {
  title: "Team Overview",
};

export default async function TeamPage({
  params,
}: {
  params: Promise<{ teamID: string }>;
}) {
  const { teamID } = await params;
  redirect(`/teams/${teamID}/monitors`);
}
