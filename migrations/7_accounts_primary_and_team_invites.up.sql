ALTER TABLE "public"."accounts"
    ADD COLUMN "is_primary" boolean NOT NULL DEFAULT false;

WITH ranked_accounts AS (
    SELECT id,
           row_number() OVER (PARTITION BY user_id ORDER BY created_at) AS rn
    FROM accounts
)
UPDATE accounts
SET is_primary = (ranked_accounts.rn = 1)
FROM ranked_accounts
WHERE accounts.id = ranked_accounts.id;

CREATE UNIQUE INDEX "uq_accounts_primary_per_user" ON "public"."accounts" ("user_id") WHERE is_primary;

CREATE TYPE "invite_status" AS ENUM ('pending', 'accepted', 'rejected', 'canceled');

ALTER TABLE "public"."team_invites"
    ADD COLUMN "invited_email" text;

UPDATE "public"."team_invites"
SET invited_email = (
    SELECT email
    FROM accounts
    WHERE accounts.user_id = team_invites.invited_to
    ORDER BY is_primary DESC, created_at ASC
    LIMIT 1
);

ALTER TABLE "public"."team_invites"
    ALTER COLUMN "invited_email" SET NOT NULL,
    ADD COLUMN "role" member_role NOT NULL DEFAULT 'member',
    ADD COLUMN "status" invite_status NOT NULL DEFAULT 'pending',
    ADD COLUMN "token" text,
    ADD COLUMN "expires_at" timestamp NOT NULL DEFAULT (now() + interval '7 days'),
    ADD COLUMN "accepted_at" timestamp,
    ADD COLUMN "rejected_at" timestamp,
    ADD COLUMN "canceled_at" timestamp;

CREATE UNIQUE INDEX "uq_team_invites_token" ON "public"."team_invites" ("token") WHERE token IS NOT NULL;
CREATE INDEX "idx_team_invites_team_id" ON "public"."team_invites" ("team_id");
CREATE INDEX "idx_team_invites_status" ON "public"."team_invites" ("status");
