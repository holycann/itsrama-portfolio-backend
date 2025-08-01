CREATE OR REPLACE VIEW event_with_views AS
SELECT 
    e.*, 
    COUNT(ev.user_id) AS views
FROM 
    events e
LEFT JOIN 
    event_views ev ON ev.event_id = e.id
GROUP BY 
    e.id;