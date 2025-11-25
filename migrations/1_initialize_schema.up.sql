CREATE TABLE "monitors" (
    "id" bigint NOT NULL,
    "url" text NOT NULL,
    "interval" integer NOT NULL,
    "last_check" timestamp NOT NULL,
    "next_check" timestamp NOT NULL,
    PRIMARY KEY ("id")
);

CREATE TABLE "pings" (
    "time" timestamp NOT NULL,
    "monitor_id" bigint NOT NULL,
    "latency" smallint NOT NULL,
    PRIMARY KEY ("time", "monitor_id")
);

-- Foreign key constraints
-- Schema: public
ALTER TABLE "pings" ADD CONSTRAINT "fk_pings_monitor_id_monitors_id" FOREIGN KEY("monitor_id") REFERENCES "monitors"("id");