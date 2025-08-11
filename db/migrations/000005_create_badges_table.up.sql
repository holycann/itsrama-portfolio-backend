CREATE TABLE public.badges (
  id uuid NOT NULL DEFAULT gen_random_uuid(),
  name character varying(100) NOT NULL UNIQUE,
  description character varying(500),
  icon_url character varying,
  created_at timestamp with time zone DEFAULT now(),
  updated_at timestamp with time zone NULL,
  CONSTRAINT badges_pkey PRIMARY KEY (id)
);

-- Dummy data for badges
INSERT INTO public.badges (name, description, icon_url) VALUES 
  ('Warlok', 'Terverifikasi sebagai warga lokal', 'https://rhfhplcxngijmfanrxzo.supabase.co/storage/v1/object/public/cultour/cultour/images/badges/warlok.png'),
  ('Penjelajah', 'Menjelajahi situs budaya', 'https://rhfhplcxngijmfanrxzo.supabase.co/storage/v1/object/public/cultour/cultour/images/badges/eksproler.png'),
  ('Event Enthusiast', 'Attended 3 local events', 'https://rhfhplcxngijmfanrxzo.supabase.co/storage/v1/object/public/cultour/cultour/images/badges/event.png'),
  ('Province Traveler', 'Explored locations in 3 different provinces', 'https://rhfhplcxngijmfanrxzo.supabase.co/storage/v1/object/public/cultour/cultour/images/badges/province.png'),
  ('Community Contributor', 'Actively participated in discussions', 'https://rhfhplcxngijmfanrxzo.supabase.co/storage/v1/object/public/cultour/cultour/images/badges/community.png');

-- Verify the insertion
SELECT COUNT(*) as total_badges FROM public.badges;
