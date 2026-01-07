ALTER TABLE "public"."regions" ADD COLUMN "display_name" text;
UPDATE "public"."regions" SET "display_name" = "name" WHERE "display_name" IS NULL;
ALTER TABLE "public"."regions" ALTER COLUMN "display_name" SET NOT NULL;
CREATE UNIQUE INDEX "uq_regions_display_name" ON "public"."regions" ("display_name");
