CREATE SCHEMA IF NOT EXISTS "public";

-- Required for hypertables and continuous aggregates
CREATE EXTENSION IF NOT EXISTS timescaledb;

CREATE TYPE "auth_provider" AS ENUM ('email', 'google');
CREATE TYPE "event_type" AS ENUM ('detected', 'notification_sent', 'manually_resolved', 'auto_resolved', 'unpublished', 'published', 'investigating', 'identified', 'update', 'monitoring');
CREATE TYPE "incident_status" AS ENUM ('detected', 'investigating', 'identified', 'monitoring', 'resolved');
CREATE TYPE "incident_severity" AS ENUM ('emergency', 'critical', 'major', 'minor', 'info');
CREATE TYPE "member_role" AS ENUM ('owner', 'admin', 'member', 'viewer');
CREATE TYPE "monitor_status" AS ENUM ('up', 'down');
CREATE TYPE "monitor_type" AS ENUM ('http', 'ping');
CREATE TYPE "notification_type" AS ENUM ('discord', 'telegram', 'email', 'slack');
CREATE TYPE "ping_status" AS ENUM ('successful', 'failed', 'timeout');
CREATE TYPE "status_page_element_type" AS ENUM ('historical_timeline', 'current_status_indicator');
CREATE TYPE "invite_status" AS ENUM ('pending', 'accepted', 'rejected', 'canceled');

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
CREATE INDEX "idx_oauth_tokens_provider" ON "public"."oauth_tokens" ("provider");

CREATE TABLE "public"."users" (
    "id" bigint NOT NULL,
    "password_hash" text,
    "avatar" text,
    "display_name" text NOT NULL,
    "verified" boolean NOT NULL DEFAULT false,
    "verify_code" text,
    "created_at" timestamp NOT NULL,
    "updated_at" timestamp NOT NULL,
    PRIMARY KEY ("id")
);

CREATE TABLE "public"."status_pages" (
    "id" bigint NOT NULL,
    "team_id" bigint NOT NULL,
    "slug" text NOT NULL,
    "title" text NOT NULL,
    "icon" bytea,
    "created_at" timestamp NOT NULL,
    "updated_at" timestamp NOT NULL,
    CONSTRAINT "pk_status_pages_id" PRIMARY KEY ("id")
);
-- Indexes
CREATE UNIQUE INDEX "uq_status_pages_slug" ON "public"."status_pages" ("slug");

CREATE TABLE "public"."accounts" (
    "id" bigint NOT NULL,
    "provider" auth_provider NOT NULL,
    "provider_user_id" character varying(255) NOT NULL,
    "user_id" bigint NOT NULL,
    "email" character varying(255) NOT NULL,
    "is_primary" boolean NOT NULL DEFAULT false,
    "created_at" timestamp NOT NULL,
    "updated_at" timestamp NOT NULL,
    CONSTRAINT "pk_accounts_id" PRIMARY KEY ("id")
);
-- Indexes
CREATE UNIQUE INDEX "uq_accounts_provider_user_id" ON "public"."accounts" ("provider", "provider_user_id");
CREATE UNIQUE INDEX "uq_accounts_provider_email" ON "public"."accounts" ("provider", "email");
CREATE UNIQUE INDEX "uq_accounts_primary_per_user" ON "public"."accounts" ("user_id") WHERE is_primary;
CREATE INDEX "idx_accounts_provider" ON "public"."accounts" ("provider");
CREATE INDEX "idx_accounts_email" ON "public"."accounts" ("email");
CREATE INDEX "idx_accounts_user_id" ON "public"."accounts" ("user_id");

CREATE TABLE "public"."monitor_regions" (
    "id" bigint NOT NULL,
    "monitor_id" bigint NOT NULL,
    "region_id" bigint NOT NULL,
    CONSTRAINT "pk_monitor_regions_id" PRIMARY KEY ("id")
);
-- Indexes
CREATE UNIQUE INDEX "uq_monitor_regions_monitor_id_region_id" ON "public"."monitor_regions" ("monitor_id", "region_id");

CREATE TABLE "public"."team_members" (
    "id" bigint NOT NULL,
    "team_id" bigint NOT NULL,
    "user_id" bigint NOT NULL,
    "role" member_role NOT NULL,
    "updated_at" timestamp NOT NULL,
    "created_at" timestamp NOT NULL,
    CONSTRAINT "pk_team_members_id" PRIMARY KEY ("id")
);
-- Indexes
CREATE UNIQUE INDEX "uq_team_members_team_id_user_id" ON "public"."team_members" ("team_id", "user_id");

CREATE TABLE "public"."monitor_notifications" (
    "id" bigint NOT NULL,
    "monitor_id" bigint NOT NULL,
    "notification_id" bigint NOT NULL,
    CONSTRAINT "pk_monitor_notifications_id" PRIMARY KEY ("id")
);
-- Indexes
CREATE UNIQUE INDEX "uq_monitor_notifications_monitor_id_notification_id" ON "public"."monitor_notifications" ("monitor_id", "notification_id");

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
CREATE UNIQUE INDEX "uq_refresh_tokens_token" ON "public"."refresh_tokens" ("token");
CREATE INDEX "idx_refresh_tokens_created_at" ON "public"."refresh_tokens" ("created_at");
CREATE INDEX "idx_refresh_tokens_user_id" ON "public"."refresh_tokens" ("user_id");

CREATE TABLE "public"."incidents" (
    "id" bigint NOT NULL,
    "status" incident_status NOT NULL,
    "severity" incident_severity NOT NULL DEFAULT 'major',
    "is_public" boolean NOT NULL,
    "auto_resolve" boolean NOT NULL DEFAULT false,
    "title" text,
    "started_at" timestamp NOT NULL,
    "resolved_at" timestamp,
    "created_at" timestamp NOT NULL,
    "updated_at" timestamp NOT NULL,
    CONSTRAINT "pk_incidents_id" PRIMARY KEY ("id")
);

CREATE TABLE "public"."incident_monitors" (
    "id" bigint NOT NULL,
    "incident_id" bigint NOT NULL,
    "monitor_id" bigint NOT NULL,
    CONSTRAINT "pk_incident_monitors_id" PRIMARY KEY ("id")
);
-- Indexes
CREATE UNIQUE INDEX "uq_incident_monitors_incident_id_monitor_id" ON "public"."incident_monitors" ("incident_id", "monitor_id");

CREATE TABLE "public"."status_page_groups" (
    "id" bigint NOT NULL,
    "status_page_id" bigint NOT NULL,
    "name" text NOT NULL,
    "type" status_page_element_type NOT NULL,
    "sort_order" integer NOT NULL,
    CONSTRAINT "pk_status_page_groups_id" PRIMARY KEY ("id")
);
-- Indexes
CREATE UNIQUE INDEX "uq_status_page_groups_status_page_id_sort_order" ON "public"."status_page_groups" ("status_page_id", "sort_order");

CREATE TABLE "public"."teams" (
    "id" bigint NOT NULL,
    "name" text NOT NULL,
    "updated_at" timestamp NOT NULL,
    "created_at" timestamp NOT NULL,
    CONSTRAINT "pk_teams_id" PRIMARY KEY ("id")
);

CREATE TABLE "public"."status_page_monitors" (
    "id" bigint NOT NULL,
    "status_page_id" bigint NOT NULL,
    "monitor_id" bigint NOT NULL,
    "group_id" bigint,
    "name" text NOT NULL,
    "type" status_page_element_type NOT NULL,
    "sort_order" integer NOT NULL,
    CONSTRAINT "pk_status_page_monitors_id" PRIMARY KEY ("id")
);
-- Indexes
CREATE UNIQUE INDEX "uq_status_page_monitors_status_page_id_group_id_sort_order" ON "public"."status_page_monitors" ("status_page_id", "group_id", "sort_order");

CREATE TABLE "public"."event_timelines" (
    "id" bigint NOT NULL,
    "event_id" bigint NOT NULL,
    "created_by" bigint,
    "message" text NOT NULL,
    "event_type" event_type NOT NULL,
    "created_at" timestamp NOT NULL,
    "updated_at" timestamp NOT NULL,
    CONSTRAINT "pk_event_timelines_id" PRIMARY KEY ("id")
);

CREATE TABLE "public"."team_invites" (
    "id" bigint NOT NULL,
    "team_id" bigint NOT NULL,
    "invited_by" bigint NOT NULL,
    "invited_to" bigint NOT NULL,
    "invited_email" text NOT NULL,
    "role" member_role NOT NULL DEFAULT 'member',
    "status" invite_status NOT NULL DEFAULT 'pending',
    "token" text,
    "expires_at" timestamp NOT NULL DEFAULT (now() + interval '7 days'),
    "accepted_at" timestamp,
    "rejected_at" timestamp,
    "canceled_at" timestamp,
    "updated_at" timestamp NOT NULL,
    "created_at" timestamp NOT NULL,
    CONSTRAINT "pk_team_invites_id" PRIMARY KEY ("id")
);
-- Indexes
CREATE UNIQUE INDEX "uq_team_invites_token" ON "public"."team_invites" ("token") WHERE token IS NOT NULL;
CREATE INDEX "idx_team_invites_team_id" ON "public"."team_invites" ("team_id");
CREATE INDEX "idx_team_invites_status" ON "public"."team_invites" ("status");

CREATE TABLE "public"."regions" (
    "id" bigint NOT NULL,
    "name" text NOT NULL,
    "display_name" text NOT NULL,
    CONSTRAINT "pk_regions_id" PRIMARY KEY ("id")
);
-- Indexes
CREATE UNIQUE INDEX "uq_regions_name" ON "public"."regions" ("name");
CREATE UNIQUE INDEX "uq_regions_display_name" ON "public"."regions" ("display_name");

CREATE TABLE "public"."monitors" (
    "id" bigint NOT NULL,
    "team_id" bigint NOT NULL,
    "name" text NOT NULL,
    "type" monitor_type NOT NULL,
    "interval" integer NOT NULL,
    "config" jsonb NOT NULL,
    "last_checked" timestamp NOT NULL,
    "next_check" timestamp NOT NULL,
    "status" monitor_status NOT NULL,
    "failure_threshold" smallint NOT NULL,
    "recovery_threshold" smallint NOT NULL,
    "updated_at" timestamp NOT NULL,
    "created_at" timestamp NOT NULL,
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
    CONSTRAINT "pk_notifications_id" PRIMARY KEY ("id")
);

CREATE TABLE "public"."pings" (
    "time" timestamp NOT NULL,
    "monitor_id" bigint NOT NULL,
    "region_id" bigint NOT NULL,
    "latency" integer NOT NULL,
    "status" ping_status NOT NULL,
    CONSTRAINT "pk_pings_time_monitor_id_region_id" PRIMARY KEY ("time", "monitor_id", "region_id")
);
-- Indexes
CREATE INDEX "idx_pings_monitor_id" ON "public"."pings" ("monitor_id");

-- TimescaleDB hypertable and continuous aggregate for ping rollups
SELECT create_hypertable('pings', 'time', if_not_exists => TRUE);
SELECT add_retention_policy('pings', INTERVAL '7 days');

CREATE MATERIALIZED VIEW monitor_30min_summary
WITH (timescaledb.continuous) AS
SELECT
    monitor_id,
    region_id,
    time_bucket('30 minutes', time) AS bucket,
    count(*) AS total_count,
    count(*) FILTER (
        WHERE status = 'successful' AND latency <= 5000
    ) AS good_count,
    percentile_cont(0.50) WITHIN GROUP (ORDER BY latency) AS p50_ms,
    percentile_cont(0.75) WITHIN GROUP (ORDER BY latency) AS p75_ms,
    percentile_cont(0.90) WITHIN GROUP (ORDER BY latency) AS p90_ms,
    percentile_cont(0.95) WITHIN GROUP (ORDER BY latency) AS p95_ms,
    percentile_cont(0.99) WITHIN GROUP (ORDER BY latency) AS p99_ms
FROM pings
GROUP BY monitor_id, region_id, bucket
WITH NO DATA;

SELECT add_continuous_aggregate_policy(
    'monitor_30min_summary',
    start_offset => INTERVAL '24 hours',
    end_offset   => INTERVAL '30 minutes',
    schedule_interval => INTERVAL '15 minutes'
);

SELECT add_retention_policy('monitor_30min_summary', INTERVAL '1 year');

ALTER MATERIALIZED VIEW monitor_30min_summary
SET (timescaledb.materialized_only = false);

-- Foreign key constraints
-- Schema: public
ALTER TABLE "public"."accounts" ADD CONSTRAINT "fk_accounts_user_id_users_id" FOREIGN KEY("user_id") REFERENCES "public"."users"("id") ON DELETE CASCADE;
ALTER TABLE "public"."event_timelines" ADD CONSTRAINT "fk_event_timelines_created_by_users_id" FOREIGN KEY("created_by") REFERENCES "public"."users"("id") ON DELETE SET NULL;
ALTER TABLE "public"."event_timelines" ADD CONSTRAINT "fk_event_timelines_event_id_incidents_id" FOREIGN KEY("event_id") REFERENCES "public"."incidents"("id") ON DELETE CASCADE;
ALTER TABLE "public"."incident_monitors" ADD CONSTRAINT "fk_incident_monitors_incident_id_incidents_id" FOREIGN KEY("incident_id") REFERENCES "public"."incidents"("id") ON DELETE CASCADE;
ALTER TABLE "public"."incident_monitors" ADD CONSTRAINT "fk_incident_monitors_monitor_id_monitors_id" FOREIGN KEY("monitor_id") REFERENCES "public"."monitors"("id") ON DELETE CASCADE;
ALTER TABLE "public"."monitor_notifications" ADD CONSTRAINT "fk_monitor_notifications_monitor_id_monitors_id" FOREIGN KEY("monitor_id") REFERENCES "public"."monitors"("id") ON DELETE CASCADE;
ALTER TABLE "public"."monitor_notifications" ADD CONSTRAINT "fk_monitor_notifications_notification_id_notifications_id" FOREIGN KEY("notification_id") REFERENCES "public"."notifications"("id") ON DELETE CASCADE;
ALTER TABLE "public"."monitor_regions" ADD CONSTRAINT "fk_monitor_regions_monitor_id_monitors_id" FOREIGN KEY("monitor_id") REFERENCES "public"."monitors"("id") ON DELETE CASCADE;
ALTER TABLE "public"."monitor_regions" ADD CONSTRAINT "fk_monitor_regions_region_id_regions_id" FOREIGN KEY("region_id") REFERENCES "public"."regions"("id") ON DELETE CASCADE;
ALTER TABLE "public"."monitors" ADD CONSTRAINT "fk_monitors_team_id_teams_id" FOREIGN KEY("team_id") REFERENCES "public"."teams"("id") ON DELETE CASCADE;
ALTER TABLE "public"."notifications" ADD CONSTRAINT "fk_notifications_team_id_teams_id" FOREIGN KEY("team_id") REFERENCES "public"."teams"("id") ON DELETE CASCADE;
ALTER TABLE "public"."oauth_tokens" ADD CONSTRAINT "fk_oauth_tokens_account_id_accounts_id" FOREIGN KEY("account_id") REFERENCES "public"."accounts"("id") ON DELETE CASCADE;
ALTER TABLE "public"."pings" ADD CONSTRAINT "fk_pings_monitor_id_monitors_id" FOREIGN KEY("monitor_id") REFERENCES "public"."monitors"("id") ON DELETE CASCADE;
ALTER TABLE "public"."pings" ADD CONSTRAINT "fk_pings_region_id_regions_id" FOREIGN KEY("region_id") REFERENCES "public"."regions"("id") ON DELETE CASCADE;
ALTER TABLE "public"."refresh_tokens" ADD CONSTRAINT "fk_refresh_tokens_user_id_users_id" FOREIGN KEY("user_id") REFERENCES "public"."users"("id") ON DELETE CASCADE;
ALTER TABLE "public"."status_page_monitors" ADD CONSTRAINT "fk_status_page_monitors_group_id_status_page_groups_id" FOREIGN KEY("group_id") REFERENCES "public"."status_page_groups"("id") ON DELETE CASCADE;
ALTER TABLE "public"."status_page_monitors" ADD CONSTRAINT "fk_status_page_monitors_monitor_id_monitors_id" FOREIGN KEY("monitor_id") REFERENCES "public"."monitors"("id") ON DELETE CASCADE;
ALTER TABLE "public"."status_pages" ADD CONSTRAINT "fk_status_pages_team_id_teams_id" FOREIGN KEY("team_id") REFERENCES "public"."teams"("id") ON DELETE CASCADE;
ALTER TABLE "public"."team_invites" ADD CONSTRAINT "fk_team_invites_invited_by_users_id" FOREIGN KEY("invited_by") REFERENCES "public"."users"("id") ON DELETE CASCADE;
ALTER TABLE "public"."team_invites" ADD CONSTRAINT "fk_team_invites_invited_to_users_id" FOREIGN KEY("invited_to") REFERENCES "public"."users"("id") ON DELETE CASCADE;
ALTER TABLE "public"."team_invites" ADD CONSTRAINT "fk_team_invites_team_id_teams_id" FOREIGN KEY("team_id") REFERENCES "public"."teams"("id") ON DELETE CASCADE;
ALTER TABLE "public"."team_members" ADD CONSTRAINT "fk_team_members_team_id_teams_id" FOREIGN KEY("team_id") REFERENCES "public"."teams"("id") ON DELETE CASCADE;
ALTER TABLE "public"."team_members" ADD CONSTRAINT "fk_team_members_user_id_users_id" FOREIGN KEY("user_id") REFERENCES "public"."users"("id") ON DELETE CASCADE;
ALTER TABLE "public"."status_page_groups" ADD CONSTRAINT "fk_status_page_groups_status_page_id_status_pages_id" FOREIGN KEY("status_page_id") REFERENCES "public"."status_pages"("id") ON DELETE CASCADE;
ALTER TABLE "public"."status_page_monitors" ADD CONSTRAINT "fk_status_page_monitors_status_page_id_status_pages_id" FOREIGN KEY("status_page_id") REFERENCES "public"."status_pages"("id") ON DELETE CASCADE;
