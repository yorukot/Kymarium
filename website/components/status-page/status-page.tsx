"use client";

import Link from "next/link";
import { Eye, Link as LinkIcon, MoreVertical } from "lucide-react";

import { Card, CardContent } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuGroup,
  DropdownMenuItem,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu";
import { formatRelativeTime } from "@/lib/parsers/datetime";
import type { StatusPageListItem } from "@/lib/schemas/status-page";

export default function StatusPageList({
  teamID,
  statusPages,
}: {
  teamID: string;
  statusPages: StatusPageListItem[];
}) {
  if (!statusPages.length) return null;

  return (
    <>
      {statusPages.map((page) => (
        <Card key={page.id} className="p-3.5">
          <CardContent className="flex items-center justify-between gap-3 p-0">
            <div className="flex min-w-0 flex-col gap-2">
              <div className="flex items-center gap-2 text-md min-w-0">
                <span className="min-w-0 flex-1 truncate hover:underline">
                  <Link href={`/teams/${teamID}/status-pages/${page.id}`}>
                    {page.title}
                  </Link>
                </span>
              </div>

              <div className="flex items-center gap-2 text-sm text-muted-foreground min-w-0">
                <LinkIcon
                  size={14}
                  className="flex-none shrink-0"
                  aria-label="Status page slug"
                />
                <span className="min-w-0 flex-1 truncate">/{page.slug}</span>
              </div>
            </div>

            <div className="flex flex-none gap-3">
              <div className="hidden md:flex flex-col items-center justify-start gap-1">
                <span className="text-muted-foreground text-sm">Updated</span>
                <span suppressHydrationWarning>
                  {formatRelativeTime(page.updatedAt)}
                </span>
              </div>

              <div className="flex flex-col items-center justify-start gap-1">
                <span className="text-muted-foreground text-sm">Monitors</span>
                <span>{page.monitorCount}</span>
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
                        <Link href={`/teams/${teamID}/status-pages/${page.id}`}>
                          <Eye /> View
                        </Link>
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

