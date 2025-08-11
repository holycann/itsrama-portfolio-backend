CREATE TABLE public.event_with_views (
    event_id UUID PRIMARY KEY REFERENCES public.events (id) ON DELETE CASCADE,
    views BIGINT DEFAULT 0
);

-- Enable Row Level Security
ALTER TABLE public.event_with_views ENABLE ROW LEVEL SECURITY;

CREATE OR REPLACE FUNCTION public.update_event_views_count()
RETURNS TRIGGER AS $$
BEGIN
    -- Increment or update view counter
    INSERT INTO public.event_with_views (event_id, views)
    VALUES (NEW.event_id, 1)
    ON CONFLICT (event_id) DO UPDATE
        SET views = event_with_views.views + 1;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trg_update_event_views_count
AFTER INSERT ON public.event_views
FOR EACH ROW
EXECUTE FUNCTION public.update_event_views_count();

-- Populate initial view counts for existing events
INSERT INTO public.event_with_views (event_id, views)
SELECT e.id, COALESCE(COUNT(ev.user_id), 0)
FROM public.events e
LEFT JOIN public.event_views ev ON ev.event_id = e.id
GROUP BY e.id
ON CONFLICT (event_id) DO NOTHING;