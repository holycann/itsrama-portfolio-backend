CREATE TABLE event_views (
    event_id UUID REFERENCES events (id) ON DELETE CASCADE,
    user_id UUID REFERENCES auth.users (id) ON DELETE CASCADE,
    PRIMARY KEY (event_id, user_id)
);

-- Dummy data for event_views based on events in 000007_create_events_table.up.sql
WITH
    events AS (
        SELECT id
        FROM public.events
    ),
    users AS (
        SELECT id
        FROM auth.users
        LIMIT 3
    )
INSERT INTO
    event_views (event_id, user_id)
SELECT events.id, users.id
FROM events, users;