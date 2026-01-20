import type { Metadata } from "next";
import NewStatusPageForm from "@/components/status-page/new-status-page-form";

export const metadata: Metadata = {
  title: "New Status Page",
};

export default async function NewStatusPage({
  params,
}: {
  params: Promise<{ teamID: string }>;
}) {
  const { teamID } = await params;

  return (
    <div className="flex flex-col gap-4">
      <NewStatusPageForm teamID={teamID} />
    </div>
  );
}
