-- Drop the trigger for updating event views count
DROP TRIGGER IF EXISTS trg_update_event_views_count ON public.event_views;

-- Drop the function for updating event views count
DROP FUNCTION IF EXISTS public.update_event_views_count();

-- Drop the table for event views count
DROP TABLE IF EXISTS public.event_with_views;