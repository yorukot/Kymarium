"use client";

import * as React from "react";
import { createAvatar } from "@dicebear/core";
import { thumbs } from "@dicebear/collection";
import { zodResolver } from "@hookform/resolvers/zod";
import { KeyRound, Save } from "lucide-react";
import { useRouter } from "next/navigation";
import { useForm, useWatch } from "react-hook-form";
import { toast } from "sonner";

import { useUser } from "@/components/context/user-context";
import { Avatar, AvatarImage } from "@/components/ui/avatar";
import { Button } from "@/components/ui/button";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import {
  Field,
  FieldDescription,
  FieldError,
  FieldGroup,
  FieldLabel,
} from "@/components/ui/field";
import { Input } from "@/components/ui/input";
import { Spinner } from "@/components/ui/spinner";

import { ApiError } from "@/lib/api/client";
import { updatePassword, updateUserProfile } from "@/lib/api/users";
import {
  passwordUpdateSchema,
  type PasswordUpdateFormValues,
  userProfileSchema,
  type UserProfileFormValues,
} from "@/lib/schemas/user";

export default function AccountSettingsPageClient() {
  const user = useUser();
  const router = useRouter();

  const avatarSeed = user.id;
  const generatedAvatarSrc = React.useMemo(
    () =>
      createAvatar(thumbs, {
        seed: avatarSeed,
        size: 64,
      }).toDataUri(),
    [avatarSeed],
  );

  const profileForm = useForm<UserProfileFormValues>({
    resolver: zodResolver(userProfileSchema),
    defaultValues: {
      displayName: user.displayName,
      avatar: user.avatar ?? "",
    },
  });

  React.useEffect(() => {
    profileForm.reset({
      displayName: user.displayName,
      avatar: user.avatar ?? "",
    });
  }, [profileForm, user.avatar, user.displayName]);

  const passwordForm = useForm<PasswordUpdateFormValues>({
    resolver: zodResolver(passwordUpdateSchema),
    defaultValues: {
      currentPassword: "",
      newPassword: "",
      confirmNewPassword: "",
    },
  });

  const watchedAvatar = useWatch({
    control: profileForm.control,
    name: "avatar",
  });
  const avatarSrc = watchedAvatar?.trim()
    ? watchedAvatar.trim()
    : user.avatar || generatedAvatarSrc;

  return (
    <div className="flex flex-col gap-4">
      <div className="flex flex-col gap-1">
        <span className="text-xl font-bold">Account</span>
        <span className="text-sm text-muted-foreground">
          Manage your profile and security
        </span>
      </div>

      <Card>
        <CardHeader>
          <CardTitle className="text-xl">Profile</CardTitle>
        </CardHeader>
        <CardContent>
          <form
            noValidate
            onSubmit={profileForm.handleSubmit(async (values) => {
              const payload: Partial<UserProfileFormValues> = {};

              if (profileForm.formState.dirtyFields.displayName) {
                payload.displayName = values.displayName;
              }

              if (profileForm.formState.dirtyFields.avatar) {
                payload.avatar = values.avatar;
              }

              if (!payload.displayName && payload.avatar === undefined) {
                toast.message("No changes to save");
                return;
              }

              const toastId = toast.loading("Saving profile...");
              try {
                await updateUserProfile(payload);
                toast.success("Profile updated", { id: toastId });
                profileForm.reset(values);
                router.refresh();
              } catch (error) {
                if (error instanceof ApiError) {
                  toast.error(error.message, { id: toastId });
                  return;
                }
                toast.error("Failed to update profile", { id: toastId });
              }
            })}
          >
            <FieldGroup>
              <div className="flex items-center gap-4">
                <Avatar className="h-14 w-14 rounded-lg">
                  <AvatarImage src={avatarSrc} alt={user.displayName} />
                </Avatar>
                <div className="min-w-0">
                  <div className="truncate text-sm font-medium">
                    {user.displayName}
                  </div>
                  <div className="truncate text-xs text-muted-foreground">
                    {user.email}
                  </div>
                </div>
              </div>

              <Field>
                <FieldLabel htmlFor="displayName">Display name</FieldLabel>
                <Input
                  id="displayName"
                  type="text"
                  placeholder="Your name"
                  aria-invalid={!!profileForm.formState.errors.displayName}
                  {...profileForm.register("displayName")}
                />
                <FieldError errors={[profileForm.formState.errors.displayName]} />
              </Field>

              <Field>
                <FieldLabel htmlFor="avatar">Avatar URL</FieldLabel>
                <Input
                  id="avatar"
                  type="url"
                  placeholder="https://..."
                  aria-invalid={!!profileForm.formState.errors.avatar}
                  {...profileForm.register("avatar")}
                />
                <FieldError errors={[profileForm.formState.errors.avatar]} />
                <FieldDescription>
                  Leave blank to use the generated avatar.
                </FieldDescription>
              </Field>

              <Field>
                <FieldLabel htmlFor="email">Email</FieldLabel>
                <Input id="email" type="email" value={user.email} disabled readOnly />
              </Field>

              <div className="flex justify-end">
                <Button
                  type="submit"
                  disabled={
                    profileForm.formState.isSubmitting ||
                    !profileForm.formState.isDirty
                  }
                >
                  {profileForm.formState.isSubmitting ? <Spinner /> : <Save />}
                  Save changes
                </Button>
              </div>
            </FieldGroup>
          </form>
        </CardContent>
      </Card>

      <Card>
        <CardHeader>
          <CardTitle className="text-xl">Security</CardTitle>
        </CardHeader>
        <CardContent>
          <form
            noValidate
            onSubmit={passwordForm.handleSubmit(async (values) => {
              const toastId = toast.loading("Updating password...");
              try {
                await updatePassword({
                  currentPassword: values.currentPassword,
                  newPassword: values.newPassword,
                });
                toast.success("Password updated", { id: toastId });
                passwordForm.reset();
              } catch (error) {
                if (error instanceof ApiError) {
                  toast.error(error.message, { id: toastId });
                  return;
                }
                toast.error("Failed to update password", { id: toastId });
              }
            })}
          >
            <FieldGroup>
              <Field>
                <FieldLabel htmlFor="currentPassword">Current password</FieldLabel>
                <Input
                  id="currentPassword"
                  type="password"
                  autoComplete="current-password"
                  aria-invalid={!!passwordForm.formState.errors.currentPassword}
                  {...passwordForm.register("currentPassword")}
                />
                <FieldError errors={[passwordForm.formState.errors.currentPassword]} />
              </Field>

              <Field>
                <div className="grid grid-cols-1 gap-4 md:grid-cols-2">
                  <Field>
                    <FieldLabel htmlFor="newPassword">New password</FieldLabel>
                    <Input
                      id="newPassword"
                      type="password"
                      autoComplete="new-password"
                      aria-invalid={!!passwordForm.formState.errors.newPassword}
                      {...passwordForm.register("newPassword")}
                    />
                    <FieldError errors={[passwordForm.formState.errors.newPassword]} />
                  </Field>

                  <Field>
                    <FieldLabel htmlFor="confirmNewPassword">
                      Confirm new password
                    </FieldLabel>
                    <Input
                      id="confirmNewPassword"
                      type="password"
                      autoComplete="new-password"
                      aria-invalid={
                        !!passwordForm.formState.errors.confirmNewPassword
                      }
                      {...passwordForm.register("confirmNewPassword")}
                    />
                    <FieldError
                      errors={[passwordForm.formState.errors.confirmNewPassword]}
                    />
                  </Field>
                </div>
                <FieldDescription>Must be at least 8 characters long.</FieldDescription>
              </Field>

              <div className="flex justify-end">
                <Button type="submit" disabled={passwordForm.formState.isSubmitting}>
                  {passwordForm.formState.isSubmitting ? <Spinner /> : <KeyRound />}
                  Update password
                </Button>
              </div>
            </FieldGroup>
          </form>
        </CardContent>
      </Card>
    </div>
  );
}
