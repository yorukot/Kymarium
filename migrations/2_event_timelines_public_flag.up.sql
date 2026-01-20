-- Add is_public flag to event_timelines and drop creator reference
ALTER TABLE public.event_timelines DROP CONSTRAINT IF EXISTS fk_event_timelines_created_by_users_id;
ALTER TABLE public.event_timelines DROP COLUMN IF EXISTS created_by;
ALTER TABLE public.event_timelines ADD COLUMN IF NOT EXISTS is_public boolean NOT NULL DEFAULT false;
