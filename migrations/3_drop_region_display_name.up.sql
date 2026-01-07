DROP INDEX IF EXISTS "uq_regions_display_name";
ALTER TABLE "public"."regions" DROP COLUMN IF EXISTS "display_name";
