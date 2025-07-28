CREATE TABLE public.cities (
    id uuid NOT NULL DEFAULT gen_random_uuid (),
    name character varying NOT NULL,
    province_id uuid NOT NULL,
    image_url varchar NULL,
    created_at timestamp with time zone NOT NULL DEFAULT now(),
    CONSTRAINT cities_pkey PRIMARY KEY (id),
    CONSTRAINT cities_province_id_fkey FOREIGN KEY (province_id) REFERENCES public.provinces (id),
    CONSTRAINT cities_unique_name_in_province UNIQUE (name, province_id)
);

-- Dummy data for cities
WITH provinces AS (
  SELECT id, name FROM public.provinces
)
INSERT INTO public.cities (id, name, province_id, image_url) VALUES 
  (gen_random_uuid(), 'Denpasar', (SELECT id FROM provinces WHERE name = 'Bali'), 'https://plcahold.co/600x400?text=Denpasar'),
  (gen_random_uuid(), 'Jakarta Selatan', (SELECT id FROM provinces WHERE name = 'Jakarta'), 'https://plcahold.co/600x400?text=Jakarta+Selatan'),
  (gen_random_uuid(), 'Semarang', (SELECT id FROM provinces WHERE name = 'Jawa Tengah'), 'https://plcahold.co/600x400?text=Semarang'),
  (gen_random_uuid(), 'Surabaya', (SELECT id FROM provinces WHERE name = 'Jawa Timur'), 'https://plcahold.co/600x400?text=Surabaya'),
  (gen_random_uuid(), 'Medan', (SELECT id FROM provinces WHERE name = 'Sumatera Utara'), 'https://plcahold.co/600x400?text=Medan');