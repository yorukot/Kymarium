"use client"

import { useRouter } from "next/navigation"

import { zodResolver } from "@hookform/resolvers/zod"
import { useForm } from "react-hook-form"

import { Button } from "@/components/ui/button"
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card"
import {
  Field,
  FieldDescription,
  FieldError,
  FieldGroup,
  FieldLabel,
} from "@/components/ui/field"
import { Input } from "@/components/ui/input"
import { Spinner } from "@/components/ui/spinner"
import { createTeam } from "@/lib/api/teams"
import { ApiError } from "@/lib/api/client"
import { applyServerFieldErrors } from "@/lib/api/error"
import {
  type TeamFormValues,
  teamPayloadSchema,
  teamSchema,
} from "@/lib/schemas/team"

export default function NewTeamPage() {
  const router = useRouter()
  const form = useForm<TeamFormValues>({
    resolver: zodResolver(teamSchema),
    defaultValues: {
      name: "",
    },
    mode: "onSubmit",
  })

  const onSubmit = async (values: TeamFormValues) => {
    form.clearErrors()

    const parsed = teamPayloadSchema.safeParse(values)
    if (!parsed.success) {
      form.setError("root", {
        type: "validate",
        message: "Invalid form data. Please try again.",
      })
      return
    }

    try {
      const response = await createTeam(parsed.data)
      const team = response.data?.data
      form.reset()

      if (team?.id) {
        router.replace(`/teams/${team.id}`)
        return
      }

      router.replace("/teams")
    } catch (error) {
      if (error instanceof ApiError) {
        const hasFieldErrors = applyServerFieldErrors(form.setError, error.body)

        if (!hasFieldErrors) {
          form.setError("root", {
            type: "server",
            message:
              error.status >= 500
                ? "Server error. Please try again later."
                : error.message,
          })
        }

        return
      }

      form.setError("root", {
        type: "network",
        message: "Network error. Please try again.",
      })
    }
  }

  const {
    register,
    handleSubmit,
    formState: { errors, isSubmitting },
  } = form

  return (
    <div className="bg-muted flex min-h-svh flex-col items-center justify-center gap-6 p-6 md:p-10">
      <div className="flex w-full max-w-sm flex-col gap-6">
        <Card>
          <CardHeader className="text-center">
            <CardTitle className="text-xl">Create a team</CardTitle>
            <CardDescription>
              Choose a name for the team you want to manage
            </CardDescription>
          </CardHeader>
          <CardContent>
            <form noValidate onSubmit={handleSubmit(onSubmit)}>
              <FieldGroup>
                <Field>
                  <FieldLabel htmlFor="teamName">Team name</FieldLabel>
                  <Input
                    id="teamName"
                    type="text"
                    placeholder="Acme Inc."
                    aria-invalid={!!errors.name}
                    {...register("name")}
                  />
                  <FieldError errors={[errors.name]} />
                </Field>
                <Field>
                  <FieldDescription>
                    You can change this later in team settings.
                  </FieldDescription>
                </Field>
                <Field>
                  <FieldError errors={[errors.root]} />
                  <Button type="submit" className="w-full" disabled={isSubmitting}>
                    {isSubmitting ? (
                      <>
                        <Spinner className="mr-2 h-4 w-4 animate-spin" />
                        Creating team...
                      </>
                    ) : (
                      "Create team"
                    )}
                  </Button>
                </Field>
              </FieldGroup>
            </form>
          </CardContent>
        </Card>
      </div>
    </div>
  )
}
