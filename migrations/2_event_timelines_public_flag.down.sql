-- Reintroduce created_by on event_timelines and remove is_public flag
ALTER TABLE public.event_timelines ADD COLUMN IF NOT EXISTS created_by bigint;
ALTER TABLE public.event_timelines ADD CONSTRAINT fk_event_timelines_created_by_users_id FOREIGN KEY (created_by) REFERENCES public.users(id) ON DELETE SET NULL;
ALTER TABLE public.event_timelines DROP COLUMN IF EXISTS is_public;
