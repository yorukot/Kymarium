DROP INDEX IF EXISTS "uq_team_invites_token";
DROP INDEX IF EXISTS "idx_team_invites_team_id";
DROP INDEX IF EXISTS "idx_team_invites_status";

ALTER TABLE "public"."team_invites"
    DROP COLUMN IF EXISTS "canceled_at",
    DROP COLUMN IF EXISTS "rejected_at",
    DROP COLUMN IF EXISTS "accepted_at",
    DROP COLUMN IF EXISTS "expires_at",
    DROP COLUMN IF EXISTS "token",
    DROP COLUMN IF EXISTS "status",
    DROP COLUMN IF EXISTS "role",
    DROP COLUMN IF EXISTS "invited_email";

DROP TYPE IF EXISTS "invite_status";

DROP INDEX IF EXISTS "uq_accounts_primary_per_user";

ALTER TABLE "public"."accounts"
    DROP COLUMN IF EXISTS "is_primary";
