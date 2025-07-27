CREATE TABLE public.provinces (
    id uuid NOT NULL DEFAULT gen_random_uuid (),
    name character varying NOT NULL UNIQUE,
    description text,
    created_at timestamp
    with
        time zone NOT NULL DEFAULT now(),
        CONSTRAINT provinces_pkey PRIMARY KEY (id)
);

-- Dummy data for provinces
INSERT INTO public.provinces (id, name, description) VALUES 
  (gen_random_uuid(), 'Bali', 'Beautiful island province known for its culture and tourism'),
  (gen_random_uuid(), 'Jakarta', 'Capital city province, center of Indonesian economy and politics'),
  (gen_random_uuid(), 'Jawa Tengah', 'Central Java province with rich Javanese cultural heritage'),
  (gen_random_uuid(), 'Jawa Timur', 'East Java province with diverse landscapes and industries'),
  (gen_random_uuid(), 'Sumatera Utara', 'North Sumatra province with varied ethnic groups and natural resources');