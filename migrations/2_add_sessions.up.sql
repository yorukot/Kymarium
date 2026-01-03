CREATE TABLE "public"."sessions" (
    "id" bigint NOT NULL,
    "user_id" bigint NOT NULL,
    "token" text NOT NULL UNIQUE,
    "user_agent" text,
    "ip" inet,
    "expires_at" timestamp NOT NULL,
    "created_at" timestamp NOT NULL,
    PRIMARY KEY ("id")
);

CREATE UNIQUE INDEX "uq_sessions_token" ON "public"."sessions" ("token");
CREATE INDEX "idx_sessions_user_id" ON "public"."sessions" ("user_id");
CREATE INDEX "idx_sessions_expires_at" ON "public"."sessions" ("expires_at");

ALTER TABLE "public"."sessions" ADD CONSTRAINT "fk_sessions_user_id_users_id" FOREIGN KEY("user_id") REFERENCES "public"."users"("id") ON DELETE CASCADE;
