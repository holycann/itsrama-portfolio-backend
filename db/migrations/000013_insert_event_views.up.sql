create or replace function insert_event_views()
returns trigger as $$
begin
  insert into event_views (event_id, user_id)
  values (NEW.id, NEW.user_id);
  return NEW;
end;
$$ language plpgsql;

create trigger trg_insert_event_views
after insert on events
for each row
execute procedure insert_event_views();