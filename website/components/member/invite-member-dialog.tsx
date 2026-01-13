"use client";

import * as React from "react";
import { zodResolver } from "@hookform/resolvers/zod";
import { useForm, Controller } from "react-hook-form";
import { toast } from "sonner";
import { UserPlus } from "lucide-react";

import { Button } from "@/components/ui/button";
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
  DialogTrigger,
} from "@/components/ui/dialog";
import { Field, FieldError, FieldGroup, FieldLabel } from "@/components/ui/field";
import { Input } from "@/components/ui/input";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";

import { ApiError } from "@/lib/api/client";
import { createTeamInvite } from "@/lib/api/team-invites";
import {
  createTeamInviteSchema,
  inviteRoleValues,
  type CreateTeamInviteValues,
} from "@/lib/schemas/team-invite";

export default function InviteMemberDialog({
  teamID,
  disabled,
  onInvited,
  trigger,
}: {
  teamID: string;
  disabled?: boolean;
  onInvited: () => void;
  trigger?: React.ReactNode;
}) {
  const [open, setOpen] = React.useState(false);

  const {
    control,
    register,
    handleSubmit,
    reset,
    formState: { errors, isSubmitting },
  } = useForm<CreateTeamInviteValues>({
    resolver: zodResolver(createTeamInviteSchema),
    defaultValues: {
      email: "",
      role: "member",
    },
  });

  const onSubmit = React.useCallback(
    async (values: CreateTeamInviteValues) => {
      const toastId = toast.loading("Sending invite...");
      try {
        await createTeamInvite(teamID, values);
        toast.success("Invite sent", { id: toastId });
        onInvited();
        reset();
        setOpen(false);
      } catch (error) {
        if (error instanceof ApiError) {
          toast.error(error.message, { id: toastId });
          return;
        }
        toast.error("Failed to send invite", { id: toastId });
      }
    },
    [onInvited, reset, teamID],
  );

  return (
    <Dialog open={open} onOpenChange={setOpen}>
      <DialogTrigger asChild disabled={disabled}>
        {trigger ?? (
          <Button disabled={disabled}>
            <UserPlus />
            Invite member
          </Button>
        )}
      </DialogTrigger>
      <DialogContent>
        <DialogHeader>
          <DialogTitle>Invite a member</DialogTitle>
          <DialogDescription>
            The invited user must already have an account.
          </DialogDescription>
        </DialogHeader>

        <form noValidate onSubmit={handleSubmit(onSubmit)}>
          <FieldGroup>
            <Field>
              <FieldLabel htmlFor="email">Email</FieldLabel>
              <Input
                id="email"
                type="email"
                placeholder="user@example.com"
                aria-invalid={!!errors.email}
                {...register("email")}
              />
              <FieldError errors={[errors.email]} />
            </Field>

            <Field>
              <FieldLabel>Role</FieldLabel>
              <Controller
                name="role"
                control={control}
                render={({ field }) => (
                  <Select
                    name={field.name}
                    value={field.value}
                    onValueChange={field.onChange}
                  >
                    <SelectTrigger aria-invalid={!!errors.role}>
                      <SelectValue placeholder="Select role" />
                    </SelectTrigger>
                    <SelectContent>
                      {inviteRoleValues.map((role) => (
                        <SelectItem key={role} value={role}>
                          {role}
                        </SelectItem>
                      ))}
                    </SelectContent>
                  </Select>
                )}
              />
              <FieldError errors={[errors.role]} />
            </Field>

            <DialogFooter>
              <Button
                type="button"
                variant="outline"
                onClick={() => setOpen(false)}
              >
                Cancel
              </Button>
              <Button type="submit" disabled={isSubmitting}>
                Send invite
              </Button>
            </DialogFooter>
          </FieldGroup>
        </form>
      </DialogContent>
    </Dialog>
  );
}

