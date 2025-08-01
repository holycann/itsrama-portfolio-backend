CREATE OR REPLACE FUNCTION insert_event_view_once(_event_id UUID, _user_id UUID)
RETURNS VOID AS $$
BEGIN
    INSERT INTO event_views (event_id, user_id)
    VALUES (_event_id, _user_id)
    ON CONFLICT DO NOTHING;
END;
$$ LANGUAGE plpgsql;