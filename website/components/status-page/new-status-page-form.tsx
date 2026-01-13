"use client";

import { zodResolver } from "@hookform/resolvers/zod";
import { useRouter } from "next/navigation";
import { useForm } from "react-hook-form";
import { toast } from "sonner";
import { z } from "zod";

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
  FieldError,
  FieldGroup,
  FieldLabel,
  FieldSet,
} from "@/components/ui/field";
import { Input } from "@/components/ui/input";
import { Spinner } from "@/components/ui/spinner";
import { ApiError } from "@/lib/api/client";
import { createStatusPage } from "@/lib/api/status-page";
import Link from "next/link";

const statusPageCreateSchema = z.object({
  title: z.string().min(1, "Name is required.").max(255, "Name is too long."),
  slug: z
    .string()
    .min(3, "Slug must be at least 3 characters.")
    .max(255, "Slug is too long.")
    .regex(/^[a-zA-Z0-9]+$/, "Slug must be alphanumeric (a-z, A-Z, 0-9)."),
});

type StatusPageCreateValues = z.infer<typeof statusPageCreateSchema>;

export default function NewStatusPageForm({ teamID }: { teamID: string }) {
  const router = useRouter();
  const form = useForm<StatusPageCreateValues>({
    resolver: zodResolver(statusPageCreateSchema),
    defaultValues: { title: "", slug: "" },
    mode: "onSubmit",
  });

  const onSubmit = async (values: StatusPageCreateValues) => {
    form.clearErrors();
    const parsed = statusPageCreateSchema.safeParse(values);
    if (!parsed.success) {
      form.setError("root", {
        type: "validate",
        message: "Invalid form data.",
      });
      return;
    }

    try {
      const res = await createStatusPage(teamID, {
        title: parsed.data.title,
        slug: parsed.data.slug,
        elements: [],
      });

      const id = res.data?.data?.status_page?.id;
      if (!id) {
        toast.success("Status page created.");
        router.push(`/teams/${teamID}/status-pages`);
        return;
      }

      toast.success("Status page created.");
      router.push(`/teams/${teamID}/status-pages/${id}`);
    } catch (error) {
      if (error instanceof ApiError) {
        form.setError("root", {
          type: "server",
          message:
            error.status >= 500
              ? "Server error. Please try again later."
              : error.message,
        });
        return;
      }

      form.setError("root", {
        type: "network",
        message: "Network error. Please try again.",
      });
    }
  };

  const isSubmitting = form.formState.isSubmitting;

  return (
    <div className="flex flex-col gap-4">
      <div className="flex gap-1">
        <div className="flex flex-col gap-1">
          <h1 className="text-xl font-semibold">Create new status page</h1>
          <p className="text-sm text-muted-foreground">
            Set a name and a public slug. You can add groups and monitors after
            creation.
          </p>
        </div>
        <Link href={`/teams/${teamID}/status-pages`}>
          <Button variant="outline">Back to list</Button>
        </Link>
      </div>

      <Card>
        <CardHeader>
          <CardTitle>Basics</CardTitle>
          <CardDescription>
            Slug must be alphanumeric (no dashes).
          </CardDescription>
        </CardHeader>
        <CardContent>
          <form
            noValidate
            onSubmit={form.handleSubmit(onSubmit)}
            className="flex flex-col gap-5"
          >
            <FieldSet>
              <FieldGroup>
                <Field>
                  <FieldLabel htmlFor="status-page-title">Name</FieldLabel>
                  <Input
                    id="status-page-title"
                    autoComplete="off"
                    placeholder="Kymarium Status"
                    aria-invalid={!!form.formState.errors.title}
                    {...form.register("title")}
                  />
                  <FieldError errors={[form.formState.errors.title]} />
                </Field>

                <Field>
                  <FieldLabel htmlFor="status-page-slug">Slug</FieldLabel>
                  <Input
                    id="status-page-slug"
                    autoComplete="off"
                    placeholder="kymariumstatus"
                    aria-invalid={!!form.formState.errors.slug}
                    {...form.register("slug")}
                  />
                  <FieldError errors={[form.formState.errors.slug]} />
                </Field>

                <FieldError errors={[form.formState.errors.root]} />
              </FieldGroup>
            </FieldSet>

            <div className="flex items-center justify-end gap-2">
              <Button type="submit" disabled={isSubmitting}>
                {isSubmitting ? (
                  <>
                    <Spinner /> Creating…
                  </>
                ) : (
                  "Create"
                )}
              </Button>
            </div>
          </form>
        </CardContent>
      </Card>
    </div>
  );
}
