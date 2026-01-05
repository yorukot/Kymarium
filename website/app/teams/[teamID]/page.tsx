import { redirect } from "next/navigation";

export default async function TeamPage({
  params,
}: {
  params: Promise<{ teamID: string }>;
}) {
  const { teamID } = await params;
  redirect(`/teams/${teamID}/monitors`);
}
