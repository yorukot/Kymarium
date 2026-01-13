import Link from "next/link";

import { Button } from "@/components/ui/button";

export default async function NewStatusPage({
  params,
}: {
  params: Promise<{ teamID: string }>;
}) {
  const { teamID } = await params;

  return (
    <div className="flex flex-col gap-4">
      <div className="flex justify-between items-center">
        <div className="flex flex-col gap-1">
          <span className="text-xl font-bold">New Status Page</span>
          <span className="text-sm text-muted-foreground">
            Status page creation UI is coming soon.
          </span>
        </div>

        <Link href={`/teams/${teamID}/status-pages`}>
          <Button variant="outline">Back to list</Button>
        </Link>
      </div>
    </div>
  );
}

