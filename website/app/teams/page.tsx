"use client";

import Image from "next/image";
import Link from "next/link";
import { GalleryVerticalEnd } from "lucide-react";
import { createAvatar } from "@dicebear/core";
import { shapes } from "@dicebear/collection";

import { useTeams } from "@/components/context/teams-context";
import { Button } from "@/components/ui/button";
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import {
  Field,
  FieldDescription,
  FieldGroup,
  FieldSeparator,
} from "@/components/ui/field";

export default function TeamsPage() {
  const teams = useTeams();
  const avatarFor = (seed: string) =>
    createAvatar(shapes, {
      seed,
      size: 64,
    }).toDataUri();

  return (
    <div className="bg-muted flex min-h-svh flex-col items-center justify-center gap-6 p-6 md:p-10">
      <div className="flex w-full max-w-sm flex-col gap-6">
        <a href="#" className="flex items-center gap-2 self-center font-medium">
          <div className="bg-primary text-primary-foreground flex size-6 items-center justify-center rounded-md">
            <GalleryVerticalEnd className="size-4" />
          </div>
          Acme Inc.
        </a>
        <Card>
          <CardHeader className="text-center">
            <CardTitle className="text-xl">Choose a team</CardTitle>
            <CardDescription>
              Select the team you want to continue with
            </CardDescription>
          </CardHeader>
          <CardContent>
            <FieldGroup>
              <div className="flex flex-col gap-3">
                {teams.length === 0 ? (
                  <Field>
                    <FieldDescription>
                      You don&apos;t belong to any teams yet.
                    </FieldDescription>
                  </Field>
                ) : (
                  teams.map((team) => (
                    <Field key={team.id}>
                      <Link href={`/teams/${team.id}`}>
                        <div className="flex items-center justify-between p-2.5 bg-muted rounded">
                          <div className="flex items-center gap-2">
                            <Image
                              src={avatarFor(team.id)}
                              alt=""
                              unoptimized
                              width={24}
                              height={24}
                              className="rounded"
                            />
                            <span className="truncate font-medium">
                              {team.name}
                            </span>
                          </div>
                          <span className="text-muted-foreground text-sm capitalize">
                            {team.role}
                          </span>
                        </div>
                      </Link>
                    </Field>
                  ))
                )}
              </div>
              <FieldSeparator className="*:data-[slot=field-separator-content]:bg-card">
                Or
              </FieldSeparator>
              <Field>
                <Button asChild className="w-full">
                  <Link href="/teams/new-team">Create a new team</Link>
                </Button>
              </Field>
            </FieldGroup>
          </CardContent>
        </Card>
      </div>
    </div>
  );
}
