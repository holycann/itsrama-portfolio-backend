CREATE TABLE public.events (
    id uuid NOT NULL DEFAULT gen_random_uuid(),
    user_id uuid NOT NULL,
    location_id uuid NOT NULL,
    name character varying(100) NOT NULL,
    description character varying(500),
    image_url character varying,
    start_date timestamp with time zone NOT NULL,
    end_date timestamp with time zone NOT NULL,
    is_kid_friendly boolean DEFAULT false,
    created_at timestamp with time zone DEFAULT now(),
    updated_at timestamp with time zone NULL,
    CONSTRAINT events_pkey PRIMARY KEY (id),
    CONSTRAINT events_user_id_fkey FOREIGN KEY (user_id) REFERENCES auth.users (id) ON DELETE CASCADE,
    CONSTRAINT events_location_id_fkey FOREIGN KEY (location_id) REFERENCES public.locations (id) ON DELETE CASCADE
);

-- Dummy data for events
WITH
    locations AS (
        SELECT id, name, city_id
        FROM public.locations
    ),
    cities AS (
        SELECT id, name, province_id
        FROM public.cities
    ),
    provinces AS (
        SELECT id, name
        FROM public.provinces
    )
INSERT INTO
    public.events (
        id,
        user_id,
        location_id,
        name,
        description,
        start_date,
        end_date,
        is_kid_friendly,
        image_url
    )
VALUES (
        gen_random_uuid(),
        '0244478e-d0d7-4cfe-b868-aa608afc126b',
        (
            SELECT id
            FROM locations
            WHERE
                name = 'Tanah Lot Temple'
        ),
        'Bali Cultural Festival',
        'A vibrant festival celebrating Balinese culture and traditions',
        NOW() + INTERVAL '1 month',
        NOW() + INTERVAL '1 month 3 days',
        true,
        'https://placehold.co/400x300.png?text=Bali+Festival'
    ),
    (
        gen_random_uuid(),
        '120060ef-7c2e-4457-a677-c8f839e8e2a7',
        (
            SELECT id
            FROM locations
            WHERE
                name = 'Monas'
        ),
        'Jakarta Heritage Walk',
        'Guided walking tour exploring Jakarta''s historical landmarks',
        NOW() + INTERVAL '2 weeks',
        NOW() + INTERVAL '2 weeks 1 day',
        false,
        'https://placehold.co/400x300.png?text=Jakarta+Walk'
    ),
    (
        gen_random_uuid(),
        '34be7296-a530-41b0-872f-f6946441f49f',
        (
            SELECT id
            FROM locations
            WHERE
                name = 'Lawang Sewu'
        ),
        'Semarang Night Heritage Tour',
        'Evening tour of historical sites in Semarang',
        NOW() + INTERVAL '3 weeks',
        NOW() + INTERVAL '3 weeks 1 day',
        false,
        'https://placehold.co/400x300.png?text=Semarang+Tour'
    );

-- Verify the insertion
SELECT COUNT(*) as total_events FROM public.events;