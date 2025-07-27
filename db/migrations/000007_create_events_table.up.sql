CREATE TABLE public.events (
    id uuid NOT NULL DEFAULT gen_random_uuid (),
    created_at timestamp
    with
        time zone NOT NULL DEFAULT now(),
        name character varying NOT NULL,
        location_id uuid NOT NULL,
        description text,
        start_date timestamp
    with
        time zone,
        end_date timestamp
    with
        time zone,
        is_kid_friendly boolean DEFAULT false,
        image_url character varying,
        city_id uuid NOT NULL,
        province_id uuid NOT NULL,
        user_id uuid NOT NULL,
        CONSTRAINT events_pkey PRIMARY KEY (id),
        CONSTRAINT events_city_id_fkey FOREIGN KEY (city_id) REFERENCES public.cities (id),
        CONSTRAINT events_user_id_fkey FOREIGN KEY (user_id) REFERENCES auth.users (id),
        CONSTRAINT events_location_id_fkey FOREIGN KEY (location_id) REFERENCES public.locations (id),
        CONSTRAINT events_province_id_fkey FOREIGN KEY (province_id) REFERENCES public.provinces (id)
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
        name,
        location_id,
        description,
        start_date,
        end_date,
        is_kid_friendly,
        image_url,
        city_id,
        province_id,
        user_id
    )
VALUES (
        gen_random_uuid (),
        'Bali Cultural Festival',
        (
            SELECT id
            FROM locations
            WHERE
                name = 'Tanah Lot Temple'
        ),
        'Annual festival celebrating Balinese culture and traditions',
        NOW() + INTERVAL '1 month',
        NOW() + INTERVAL '1 month 3 days',
        true,
        'https://placehold.co/400x300.png?text=Bali+Festival',
        (
            SELECT id
            FROM cities
            WHERE
                name = 'Denpasar'
        ),
        (
            SELECT id
            FROM provinces
            WHERE
                name = 'Bali'
        ),
        '0244478e-d0d7-4cfe-b868-aa608afc126b'
    ),
    (
        gen_random_uuid (),
        'Jakarta Heritage Walk',
        (
            SELECT id
            FROM locations
            WHERE
                name = 'Monas'
        ),
        'Guided walking tour exploring Jakarta''s historical landmarks',
        NOW() + INTERVAL '2 weeks',
        NOW() + INTERVAL '2 weeks 1 day',
        false,
        'https://placehold.co/400x300.png?text=Jakarta+Walk',
        (
            SELECT id
            FROM cities
            WHERE
                name = 'Jakarta Selatan'
        ),
        (
            SELECT id
            FROM provinces
            WHERE
                name = 'Jakarta'
        ),
        '120060ef-7c2e-4457-a677-c8f839e8e2a7'
    ),
    (
        gen_random_uuid (),
        'Semarang Night Heritage Tour',
        (
            SELECT id
            FROM locations
            WHERE
                name = 'Lawang Sewu'
        ),
        'Evening tour of historical sites in Semarang',
        NOW() + INTERVAL '3 weeks',
        NOW() + INTERVAL '3 weeks 1 day',
        false,
        'https://placehold.co/400x300.png?text=Semarang+Tour',
        (
            SELECT id
            FROM cities
            WHERE
                name = 'Semarang'
        ),
        (
            SELECT id
            FROM provinces
            WHERE
                name = 'Jawa Tengah'
        ),
        '34be7296-a530-41b0-872f-f6946441f49f'
    );