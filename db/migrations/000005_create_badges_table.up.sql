CREATE TABLE public.badges (
  id uuid NOT NULL DEFAULT gen_random_uuid(),
  name character varying NOT NULL UNIQUE,
  description text,
  icon_url character varying,
  created_at timestamp with time zone NOT NULL DEFAULT now(),
  updated_at timestamp with time zone,
  CONSTRAINT badges_pkey PRIMARY KEY (id)
);

-- Dummy data for badges
INSERT INTO public.badges (id, name, description, icon_url) VALUES 
  (gen_random_uuid(), 'Cultural Explorer', 'Visited 5 unique cultural sites', 'https://placehold.co/100x100.png?text=Explorer'),
  (gen_random_uuid(), 'Local Story Master', 'Contributed 10 local stories', 'https://placehold.co/100x100.png?text=Story'),
  (gen_random_uuid(), 'Event Enthusiast', 'Attended 3 local events', 'https://placehold.co/100x100.png?text=Event'),
  (gen_random_uuid(), 'Province Traveler', 'Explored locations in 3 different provinces', 'https://placehold.co/100x100.png?text=Travel'),
  (gen_random_uuid(), 'Community Contributor', 'Actively participated in discussions', 'https://placehold.co/100x100.png?text=Community');
