"use client";

import { useEffect } from "react";
import { Edit, Eye, Link as LinkIcon, MoreVertical, Trash } from "lucide-react";
import { Card, CardContent } from "../ui/card";
import { Button } from "../ui/button";
import type { MonitorListItem } from "@/lib/schemas/monitor";
import Link from "next/link";
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuGroup,
  DropdownMenuItem,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu";
import { useTeamEntities } from "@/components/context/team-entities-context";
import { StatusDot } from "./status-dot";
import { formatRelativeTime } from "@/lib/parsers/datetime";
import { formatUptime } from "@/lib/parsers/format";

export default function Monitor({
  teamID,
  monitors,
}: {
  teamID: string;
  monitors: MonitorListItem[];
}) {
  const { setMonitors } = useTeamEntities();

  useEffect(() => {
    setMonitors(
      monitors.map((monitor) => ({ id: monitor.id, name: monitor.name })),
    );
  }, [monitors, setMonitors]);

  if (!monitors.length) return null;

  return (
    <>
      {monitors.map((monitor) => (
        <Card key={monitor.id} className="p-3.5">
          <CardContent className="flex items-center justify-between gap-3 p-0">
            <div className="flex min-w-0 flex-col gap-2">
              <div className="flex items-center gap-2 text-md min-w-0">
                <StatusDot status={monitor.status} />
                <span className="min-w-0 flex-1 truncate hover:underline">
                  <Link href={`/teams/${teamID}/monitors/${monitor.id}`}>
                    {monitor.name}
                  </Link>
                </span>
              </div>

              <div className="flex items-center gap-2 text-sm text-muted-foreground min-w-0">
                <LinkIcon
                  size={14}
                  className="flex-none shrink-0"
                  aria-label={monitor.targetLabel}
                />
                <span className="min-w-0 flex-1 truncate">
                  {monitor.targetValue || "--"}
                </span>
              </div>
            </div>

            <div className="flex flex-none gap-3">
              <div className="hidden md:flex flex-col items-center justify-start gap-1">
                <span className="text-muted-foreground text-sm">
                  Last checked
                </span>
                <span suppressHydrationWarning>
                  {formatRelativeTime(monitor.lastChecked)}
                </span>
              </div>

              <div className="flex flex-col items-center justify-start gap-1">
                <span className="text-muted-foreground text-sm">
                  30d uptime
                </span>
                <span>{formatUptime(monitor.uptimeSLI30)}</span>
              </div>

              <div className="flex items-center justify-center">
                <DropdownMenu>
                  <DropdownMenuTrigger asChild>
                    <Button size="icon" variant="ghost">
                      <MoreVertical size={16} />
                    </Button>
                  </DropdownMenuTrigger>
                  <DropdownMenuContent align="start">
                    <DropdownMenuGroup>
                      <DropdownMenuItem asChild>
                        <Link
                          href={`/teams/${teamID}/monitors/${monitor.id}`}
                        >
                          <Eye /> View
                        </Link>
                      </DropdownMenuItem>
                      <DropdownMenuItem asChild>
                        <Link
                          href={`/teams/${teamID}/monitors/${monitor.id}/edit`}
                        >
                          <Edit /> Edit
                        </Link>
                      </DropdownMenuItem>
                      <DropdownMenuSeparator />
                      <DropdownMenuItem variant="destructive">
                        <Trash /> Delete
                      </DropdownMenuItem>
                    </DropdownMenuGroup>
                  </DropdownMenuContent>
                </DropdownMenu>
              </div>
            </div>
          </CardContent>
        </Card>
      ))}
    </>
  );
}
