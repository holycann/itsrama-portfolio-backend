CREATE OR REPLACE VIEW event_with_views AS
SELECT e.*, COALESCE(ev.views, 0) AS views
FROM events e
    LEFT JOIN event_views ev ON ev.event_id = e.id;