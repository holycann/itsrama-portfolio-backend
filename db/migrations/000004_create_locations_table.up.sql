CREATE TABLE public.locations (
    id uuid NOT NULL DEFAULT gen_random_uuid (),
    name character varying NOT NULL,
    city_id uuid NOT NULL,
    latitude double precision NOT NULL,
    longitude double precision NOT NULL,
    created_at timestamp,
    updated_at timestamp NULL
    with
        time zone NOT NULL DEFAULT now(),
        CONSTRAINT locations_pkey PRIMARY KEY (id),
        CONSTRAINT locations_city_id_fkey FOREIGN KEY (city_id) REFERENCES public.cities (id)
);

-- Dummy data for locations
WITH cities AS (
  SELECT id, name FROM public.cities
)
INSERT INTO public.locations (id, name, city_id, latitude, longitude) VALUES 
  (gen_random_uuid(), 'Tanah Lot Temple', (SELECT id FROM cities WHERE name = 'Denpasar'), -8.6225, 115.0872),
  (gen_random_uuid(), 'Monas', (SELECT id FROM cities WHERE name = 'Jakarta Selatan'), -6.1753, 106.8272),
  (gen_random_uuid(), 'Lawang Sewu', (SELECT id FROM cities WHERE name = 'Semarang'), -6.9824, 110.4168),
  (gen_random_uuid(), 'Surabaya North Quay', (SELECT id FROM cities WHERE name = 'Surabaya'), -7.2459, 112.7381),
  (gen_random_uuid(), 'Istana Maimun', (SELECT id FROM cities WHERE name = 'Medan'), 3.5952, 98.6722);