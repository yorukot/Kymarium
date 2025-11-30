CREATE SCHEMA IF NOT EXISTS "public";

CREATE TYPE "auth_provider" AS ENUM ('email', 'google');
CREATE TYPE "member_role" AS ENUM ('owner', 'admin', 'member', 'viewer');
CREATE TYPE "monitor_type" AS ENUM ('http');
CREATE TYPE "notification_type" AS ENUM ('discord', 'telegram', 'email');
CREATE TYPE "ping_status" AS ENUM ('successful', 'failed', 'timeout');

CREATE TABLE "public"."oauth_tokens" (
    "account_id" bigint NOT NULL,
    "access_token" text NOT NULL,
    "refresh_token" text,
    "expiry" timestamp NOT NULL,
    "token_type" character varying(50) NOT NULL,
    "provider" auth_provider NOT NULL,
    "created_at" timestamp NOT NULL,
    "updated_at" timestamp NOT NULL,
    CONSTRAINT "pk_oauth_tokens_account_id" PRIMARY KEY ("account_id")
);
-- Indexes
CREATE INDEX "oauth_tokens_idx_oauth_tokens_provider" ON "public"."oauth_tokens" ("provider");

CREATE TABLE "public"."users" (
    "id" bigint NOT NULL,
    "password_hash" text,
    "avatar" text,
    "display_name" text NOT NULL,
    "created_at" timestamp NOT NULL,
    "updated_at" timestamp NOT NULL,
    PRIMARY KEY ("id")
);

CREATE TABLE "public"."accounts" (
    "id" bigint NOT NULL,
    "provider" auth_provider NOT NULL,
    "provider_user_id" character varying(255) NOT NULL,
    "user_id" bigint NOT NULL,
    "email" character varying(255) NOT NULL,
    "created_at" timestamp NOT NULL,
    "updated_at" timestamp NOT NULL,
    PRIMARY KEY ("id")
);
-- Indexes
CREATE UNIQUE INDEX "accounts_accounts_provider_provider_user_id_key" ON "public"."accounts" ("provider", "provider_user_id");
CREATE INDEX "accounts_idx_accounts_provider" ON "public"."accounts" ("provider");
CREATE UNIQUE INDEX "accounts_accounts_provider_email_key" ON "public"."accounts" ("provider", "email");
CREATE INDEX "accounts_idx_accounts_email" ON "public"."accounts" ("email");
CREATE INDEX "accounts_idx_accounts_user_id" ON "public"."accounts" ("user_id");

CREATE TABLE "public"."team_members" (
    "id" bigint NOT NULL,
    "team_id" bigint NOT NULL,
    "user_id" bigint NOT NULL,
    "role" member_role NOT NULL,
    "updated_at" timestamp NOT NULL,
    "created_at" timestamp NOT NULL,
    CONSTRAINT "pk_team_users_id" PRIMARY KEY ("id")
);

CREATE TABLE "public"."refresh_tokens" (
    "id" bigint NOT NULL,
    "user_id" bigint NOT NULL,
    "token" text NOT NULL UNIQUE,
    "user_agent" text,
    "ip" inet,
    "used_at" timestamp,
    "created_at" timestamp NOT NULL,
    PRIMARY KEY ("id")
);
-- Indexes
CREATE UNIQUE INDEX "refresh_tokens_refresh_tokens_token_key" ON "public"."refresh_tokens" ("token");
CREATE INDEX "refresh_tokens_idx_refresh_tokens_token" ON "public"."refresh_tokens" ("token");
CREATE INDEX "refresh_tokens_idx_refresh_tokens_created_at" ON "public"."refresh_tokens" ("created_at");
CREATE INDEX "refresh_tokens_idx_refresh_tokens_user_id" ON "public"."refresh_tokens" ("user_id");

CREATE TABLE "public"."teams" (
    "id" bigint NOT NULL,
    "name" text NOT NULL,
    "updated_at" timestamp NOT NULL,
    "created_at" timestamp NOT NULL,
    CONSTRAINT "pk_teams_id" PRIMARY KEY ("id")
);

CREATE TABLE "public"."team_invites" (
    "id" bigint NOT NULL,
    "team_id" bigint NOT NULL,
    "Invited_by" bigint NOT NULL,
    "Invited_to" bigint NOT NULL,
    "updated_at" timestamp NOT NULL,
    "created_at" timestamp NOT NULL,
    PRIMARY KEY ("id")
);

CREATE TABLE "public"."monitors" (
    "id" bigint NOT NULL,
    "team_id" bigint NOT NULL,
    "name" text NOT NULL,
    "type" monitor_type NOT NULL,
    "interval" integer NOT NULL,
    "config" jsonb NOT NULL,
    "last_checked" timestamp NOT NULL,
    "next_check" timestamp NOT NULL,
    "notification" bigint[] NOT NULL,
    "updated_at" timestamp NOT NULL,
    "creted_at" timestamp NOT NULL,
    "group" bigint,
    CONSTRAINT "pk_monitors_id" PRIMARY KEY ("id")
);

CREATE TABLE "public"."notifications" (
    "id" bigint NOT NULL,
    "team_id" bigint NOT NULL,
    "type" notification_type NOT NULL,
    "name" text NOT NULL,
    "config" jsonb NOT NULL,
    "updated_at" timestamp NOT NULL,
    "created_at" timestamp NOT NULL,
    CONSTRAINT "pk_table_10_id" PRIMARY KEY ("id")
);

CREATE TABLE "public"."pings" (
    "time" timestamp NOT NULL,
    "monitor_id" bigint NOT NULL,
    "latency" smallint NOT NULL,
    "status" ping_status NOT NULL,
    "data" jsonb,
    CONSTRAINT "pk_table_9_id" PRIMARY KEY ("time", "monitor_id")
);

-- Foreign key constraints
-- Schema: public
ALTER TABLE "public"."accounts" ADD CONSTRAINT "fk_accounts_user_id_users_id" FOREIGN KEY("user_id") REFERENCES "public"."users"("id");
ALTER TABLE "public"."notifications" ADD CONSTRAINT "fk_notifications_team_id_teams_id" FOREIGN KEY("team_id") REFERENCES "public"."teams"("id");
ALTER TABLE "public"."oauth_tokens" ADD CONSTRAINT "fk_oauth_tokens_account_id_accounts_id" FOREIGN KEY("account_id") REFERENCES "public"."accounts"("id");
ALTER TABLE "public"."pings" ADD CONSTRAINT "fk_pings_monitor_id_monitors_id" FOREIGN KEY("monitor_id") REFERENCES "public"."monitors"("id");
ALTER TABLE "public"."refresh_tokens" ADD CONSTRAINT "fk_refresh_tokens_user_id_users_id" FOREIGN KEY("user_id") REFERENCES "public"."users"("id");
ALTER TABLE "public"."team_invites" ADD CONSTRAINT "fk_team_invites_Invited_by_users_id" FOREIGN KEY("Invited_by") REFERENCES "public"."users"("id");
ALTER TABLE "public"."team_invites" ADD CONSTRAINT "fk_team_invites_Invited_to_users_id" FOREIGN KEY("Invited_to") REFERENCES "public"."users"("id");
ALTER TABLE "public"."team_invites" ADD CONSTRAINT "fk_team_invites_team_id_teams_id" FOREIGN KEY("team_id") REFERENCES "public"."teams"("id");
ALTER TABLE "public"."team_members" ADD CONSTRAINT "fk_team_members_team_id_teams_id" FOREIGN KEY("team_id") REFERENCES "public"."teams"("id");
ALTER TABLE "public"."team_members" ADD CONSTRAINT "fk_team_members_user_id_users_id" FOREIGN KEY("user_id") REFERENCES "public"."users"("id");