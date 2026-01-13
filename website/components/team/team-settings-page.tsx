"use client";

import * as React from "react";
import { useRouter } from "next/navigation";
import { toast } from "sonner";
import { zodResolver } from "@hookform/resolvers/zod";
import { useForm } from "react-hook-form";
import { LogOut, Save, Trash2 } from "lucide-react";

import { Button } from "@/components/ui/button";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import {
  Field,
  FieldError,
  FieldGroup,
  FieldLabel,
} from "@/components/ui/field";
import { Input } from "@/components/ui/input";
import TeamSettingsDangerAction from "@/components/team/team-settings-danger-action";

import { ApiError } from "@/lib/api/client";
import { deleteTeam, leaveTeam, updateTeam } from "@/lib/api/teams";
import { teamSchema, type TeamFormValues } from "@/lib/schemas/team";
import type { TeamRole } from "@/lib/schemas/team-member";

export default function TeamSettingsPageClient({
  teamID,
  teamName,
  teamRole,
}: {
  teamID: string;
  teamName: string;
  teamRole: TeamRole;
}) {
  const router = useRouter();

  const canEditTeam = teamRole === "owner" || teamRole === "admin";
  const canDeleteTeam = teamRole === "owner";
  const canLeaveTeam = teamRole !== "owner";

  const {
    register,
    handleSubmit,
    reset,
    formState: { errors, isSubmitting, isDirty },
  } = useForm<TeamFormValues>({
    resolver: zodResolver(teamSchema),
    defaultValues: { name: teamName },
  });

  React.useEffect(() => {
    reset({ name: teamName });
  }, [reset, teamName]);

  return (
    <div className="flex flex-col gap-4">
      <div className="flex flex-col gap-1">
        <span className="text-xl font-bold">Settings</span>
        <span className="text-sm text-muted-foreground">
          Manage team details and access
        </span>
      </div>

      <Card>
        <CardHeader>
          <CardTitle className="text-xl">Team</CardTitle>
        </CardHeader>
        <CardContent>
          <form
            noValidate
            onSubmit={handleSubmit(async (values) => {
              if (!canEditTeam) return;

              const toastId = toast.loading("Saving team...");
              try {
                await updateTeam(teamID, values);
                toast.success("Team updated", { id: toastId });
                router.refresh();
              } catch (error) {
                if (error instanceof ApiError) {
                  toast.error(error.message, { id: toastId });
                } else {
                  toast.error("Failed to update team", { id: toastId });
                }
              }
            })}
          >
            <FieldGroup>
              <Field>
                <FieldLabel htmlFor="name">Team name</FieldLabel>
                <Input
                  id="name"
                  type="text"
                  placeholder="Team name"
                  aria-invalid={!!errors.name}
                  disabled={!canEditTeam}
                  {...register("name")}
                />
                <FieldError errors={[errors.name]} />
              </Field>

              <div className="flex justify-end">
                <Button
                  type="submit"
                  disabled={!canEditTeam || isSubmitting || !isDirty}
                >
                  <Save />
                  Save changes
                </Button>
              </div>
            </FieldGroup>
          </form>
        </CardContent>
      </Card>

      <Card>
        <CardHeader>
          <CardTitle className="text-xl">Danger zone</CardTitle>
        </CardHeader>
        <CardContent>
          <div className="flex flex-col gap-3">
            <TeamSettingsDangerAction
              title="Leave team"
              description="You can leave this team at any time (owners cannot leave)."
              buttonText="Leave team"
              buttonIcon={<LogOut />}
              buttonVariant="outline"
              disabled={!canLeaveTeam}
              confirmMessage="Are you sure you want to leave this team?"
              loadingMessage="Leaving team..."
              successMessage="Left team"
              defaultErrorMessage="Failed to leave team"
              action={() => leaveTeam(teamID)}
              redirectTo="/teams"
            />

            <div className="h-px bg-border" />

            <TeamSettingsDangerAction
              title="Delete team"
              description="Permanently deletes this team and all its data. Owner only."
              buttonText="Delete team"
              buttonIcon={<Trash2 />}
              buttonVariant="destructive"
              disabled={!canDeleteTeam}
              confirmMessage="Delete this team permanently? This cannot be undone."
              loadingMessage="Deleting team..."
              successMessage="Team deleted"
              defaultErrorMessage="Failed to delete team"
              action={() => deleteTeam(teamID)}
              redirectTo="/teams"
            />
          </div>
        </CardContent>
      </Card>
    </div>
  );
}
