create table event_views (
    event_id uuid primary key references events (id) on delete cascade,
    views integer not null default 0
);