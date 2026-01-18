"use client";

import { useEffect } from "react";
import { useRouter } from "next/navigation";
import { Card, CardContent } from "@/components/ui/card";
import type { IncidentListItem } from "@/lib/schemas/incident";
import { formatRelativeTime } from "@/lib/parsers/datetime";
import { humanizeIdentifier } from "@/lib/parsers/strings";
import { Clock, Eye, EyeOff, MoreVertical } from "lucide-react";
import Link from "next/link";
import { IncidentStatusDot } from "./status-dot";
import { Button } from "@/components/ui/button";
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuGroup,
  DropdownMenuItem,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu";
import { updateIncidentSettings, updateIncidentStatus } from "@/lib/api/incident";
import type { IncidentStatus } from "@/lib/schemas/incident";
import { toast } from "sonner";
import { ApiError } from "@/lib/api/client";
import { useTeamEntities } from "@/components/context/team-entities-context";

export default function IncidentList({
  teamID,
  incidents,
}: {
  teamID: string;
  incidents: IncidentListItem[];
}) {
  const router = useRouter();
  const { setIncidents } = useTeamEntities();

  useEffect(() => {
    setIncidents(
      incidents.map((incident) => ({
        id: incident.id,
        title: incident.title ?? "",
      })),
    );
  }, [incidents, setIncidents]);

  if (!incidents.length) return null;

  return (
    <>
      {incidents.map((incident) => (
        <Card key={incident.id} className="p-3.5">
          <CardContent className="flex items-center justify-between gap-3 p-0">
            <div className="flex min-w-0 flex-col gap-2">
              <div className="flex items-center gap-2 text-md min-w-0">
                <IncidentStatusDot status={incident.status} />
                <span className="min-w-0 flex-1 truncate hover:underline">
                  <Link href={`/teams/${teamID}/incidents/${incident.id}`}>
                    {incident.title || `Incident ${incident.id}`}
                  </Link>
                </span>
              </div>

              <div className="flex items-center gap-2 text-sm text-muted-foreground min-w-0">
                <Clock
                  size={14}
                  className="flex-none shrink-0"
                  aria-label="Incident time"
                />
                <span className="min-w-0 flex-1 truncate">
                  {incident.resolvedAt
                    ? `Resolved ${formatRelativeTime(incident.resolvedAt)}`
                    : `Started ${formatRelativeTime(incident.startedAt)}`}
                </span>
              </div>
            </div>

            <div className="flex flex-none gap-3">
              <div className="hidden md:flex flex-col items-center justify-start gap-1">
                <span className="text-muted-foreground text-sm">Status</span>
                <span className="capitalize">
                  {humanizeIdentifier(incident.status)}
                </span>
              </div>

              <div className="flex flex-col items-center justify-start gap-1">
                <span className="text-muted-foreground text-sm">Severity</span>
                <span className="capitalize">
                  {humanizeIdentifier(incident.severity)}
                </span>
              </div>

              <div className="flex items-center justify-center">
                <DropdownMenu>
                  <DropdownMenuTrigger asChild>
                    <Button size="icon" variant="ghost" aria-label="Actions">
                      <MoreVertical size={16} />
                    </Button>
                  </DropdownMenuTrigger>
                  <DropdownMenuContent align="start">
                    <DropdownMenuGroup>
                      <DropdownMenuItem asChild>
                        <Link href={`/teams/${teamID}/incidents/${incident.id}`}>
                          <Eye /> View
                        </Link>
                      </DropdownMenuItem>
                      <DropdownMenuSeparator />
                      <DropdownMenuItem
                        onSelect={async () => {
                          try {
                            await updateIncidentSettings(teamID, incident.id, {
                              public: !incident.isPublic,
                            });
                            toast.success(
                              incident.isPublic ? "Marked as private" : "Made public",
                            );
                            router.refresh();
                          } catch (error) {
                            if (error instanceof ApiError) {
                              toast.error(error.message);
                              return;
                            }
                            toast.error("Failed to update visibility");
                          }
                        }}
                      >
                        {incident.isPublic ? <EyeOff /> : <Eye />}
                        {incident.isPublic ? "Make private" : "Make public"}
                      </DropdownMenuItem>
                      <DropdownMenuSeparator />
                      {(
                        [
                          "detected",
                          "investigating",
                          "identified",
                          "monitoring",
                          "resolved",
                        ] as const
                      ).map((status) => (
                        <DropdownMenuItem
                          key={status}
                          onSelect={async () => {
                            try {
                              await updateIncidentStatus(teamID, incident.id, {
                                status: status as IncidentStatus,
                              });
                              toast.success(`Status set to ${humanizeIdentifier(status)}`);
                              router.refresh();
                            } catch (error) {
                              if (error instanceof ApiError) {
                                toast.error(error.message);
                                return;
                              }
                              toast.error("Failed to update status");
                            }
                          }}
                        >
                          <IncidentStatusDot status={status} />
                          Set {humanizeIdentifier(status)}
                        </DropdownMenuItem>
                      ))}
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
