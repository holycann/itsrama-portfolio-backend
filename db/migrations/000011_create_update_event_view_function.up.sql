create or replace function increment_or_create_event_views(event_id uuid)
returns void as $$
begin
  insert into event_views (event_id, views)
  values (event_id, 1)
  on conflict (event_id)
  do update set views = event_views.views + 1;
end;
$$ language plpgsql;