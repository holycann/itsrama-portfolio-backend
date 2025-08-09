CREATE TABLE public.event_views (
    event_id UUID NOT NULL REFERENCES public.events (id) ON DELETE CASCADE,
    user_id UUID NOT NULL REFERENCES auth.users (id) ON DELETE CASCADE,
    viewed_at timestamp with time zone DEFAULT now(),
    PRIMARY KEY (event_id, user_id)
);

-- Dummy data for event views based on events in 000007_create_events_table.up.sql
WITH
    events AS (
        SELECT id, user_id
        FROM public.events
    ),
    users AS (
        SELECT id
        FROM auth.users
    )
INSERT INTO
    event_views (event_id, user_id)
SELECT events.id, users.id
FROM events, users;

-- Verify the insertion
SELECT COUNT(*) as total_event_views FROM event_views;