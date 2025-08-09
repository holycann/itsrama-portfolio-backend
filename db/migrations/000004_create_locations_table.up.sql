CREATE EXTENSION IF NOT EXISTS postgis;

CREATE TABLE public.locations (
    id uuid NOT NULL DEFAULT gen_random_uuid(),
    name character varying NOT NULL,
    city_id uuid NOT NULL,
    latitude double precision NOT NULL,
    longitude double precision NOT NULL,
    location geography(Point, 4326),
    created_at timestamp with time zone DEFAULT now(),
    updated_at timestamp with time zone,
    CONSTRAINT locations_pkey PRIMARY KEY (id),
    CONSTRAINT locations_city_id_fkey FOREIGN KEY (city_id) REFERENCES public.cities (id) ON DELETE CASCADE
);

-- Create GIST index for location column
CREATE INDEX IF NOT EXISTS locations_location_idx ON public.locations USING GIST(location);

-- Create function to automatically update location geometry
CREATE OR REPLACE FUNCTION update_location_geom()
RETURNS trigger AS $$
BEGIN
  NEW.location := ST_SetSRID(ST_MakePoint(NEW.longitude, NEW.latitude), 4326)::geography;
  RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Create trigger to call the function before insert or update
CREATE TRIGGER trg_update_location_geom
BEFORE INSERT OR UPDATE ON public.locations
FOR EACH ROW
EXECUTE FUNCTION update_location_geom();

-- Dummy data for locations
WITH
    cities AS (
        SELECT id, name
        FROM public.cities
    )
INSERT INTO public.locations (id, name, city_id, latitude, longitude) 
SELECT 
  gen_random_uuid(), 'Tanah Lot Temple', (SELECT id FROM cities WHERE name = 'Denpasar'), -8.6225, 115.0872
WHERE NOT EXISTS (SELECT 1 FROM public.locations WHERE name = 'Tanah Lot Temple')
UNION ALL
SELECT 
  gen_random_uuid(), 'Monas', (SELECT id FROM cities WHERE name = 'Jakarta Selatan'), -6.1753, 106.8272
WHERE NOT EXISTS (SELECT 1 FROM public.locations WHERE name = 'Monas')
UNION ALL
SELECT 
  gen_random_uuid(), 'Lawang Sewu', (SELECT id FROM cities WHERE name = 'Semarang'), -6.9824, 110.4168
WHERE NOT EXISTS (SELECT 1 FROM public.locations WHERE name = 'Lawang Sewu')
UNION ALL
SELECT 
  gen_random_uuid(), 'Surabaya North Quay', (SELECT id FROM cities WHERE name = 'Surabaya'), -7.2459, 112.7381
WHERE NOT EXISTS (SELECT 1 FROM public.locations WHERE name = 'Surabaya North Quay')
UNION ALL
SELECT 
  gen_random_uuid(), 'Istana Maimun', (SELECT id FROM cities WHERE name = 'Medan'), 3.5952, 98.6722
WHERE NOT EXISTS (SELECT 1 FROM public.locations WHERE name = 'Istana Maimun');

-- Verify the insertion
SELECT COUNT(*) as total_locations FROM public.locations;