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
                name = 'Lhokseumawe City Center'
        ),
        'Aceh Perkusi',
        'Pertunjukan musik perkusi tradisional Aceh',
        '2025-08-22 00:00:00+07',
        '2025-08-24 23:59:59+07',
        true,
        'https://rhfhplcxngijmfanrxzo.supabase.co/storage/v1/object/public/cultour/cultour/images/events/aceh.jpg'
    ),
    (
        gen_random_uuid(),
        '120060ef-7c2e-4457-a677-c8f839e8e2a7',
        (
            SELECT id
            FROM locations
            WHERE
                name = 'Jalan Asia Afrika'
        ),
        'Asia Afrika Festival',
        'Perayaan sejarah Konferensi Asia-Afrika dengan parade budaya, pertunjukan seni tradisional dan kontemporer sepanjang Jalan Asia-Afrika, menghadirkan rasa persatuan antarnegara Asia dan Afrika.',
        '2025-08-10 10:00:00+07',
        '2025-08-10 18:00:00+07',
        true,
        'https://rhfhplcxngijmfanrxzo.supabase.co/storage/v1/object/public/cultour/cultour/images/events/Asia.jpg'
    ),
    (
        gen_random_uuid(),
        '34be7296-a530-41b0-872f-f6946441f49f',
        (
            SELECT id
            FROM locations
            WHERE
                name = 'Dieng Plateau'
        ),
        'Dieng Culture Festival 2025',
        'Festival budaya tahunan di dataran tinggi Dieng—menampilkan seni tradisi, ritual pencukuran rambut gimbal anak-anak Dieng, konser jazz, pelepasan lampion, dan bazar UMKM.',
        '2025-08-23 00:00:00+07',
        '2025-08-24 23:59:59+07',
        true,
        'https://rhfhplcxngijmfanrxzo.supabase.co/storage/v1/object/public/cultour/cultour/images/events/Dieng.jpg'
    ),
    (
        gen_random_uuid(),
        '120060ef-7c2e-4457-a677-c8f839e8e2a7',
        (
            SELECT id
            FROM locations
            WHERE
                name = 'Gandoriah Beach'
        ),
        'Festival Tabuik 2026',
        'Perayaan tradisi masyarakat Pariaman yang memperingati Asyura dengan prosesi mengarak Tabuik ke pantai dan pelepasan ke laut.',
        '2026-08-01 00:00:00+07',
        '2026-08-02 23:59:59+07',
        true,
        'https://rhfhplcxngijmfanrxzo.supabase.co/storage/v1/object/public/cultour/cultour/images/events/tabuik.jpg'
    ),
    (
        gen_random_uuid(),
        '0244478e-d0d7-4cfe-b868-aa608afc126b',
        (
            SELECT id
            FROM locations
            WHERE
                name = 'Ubud Cultural Center'
        ),
        'Ubud Writers & Readers Festival 2025',
        'Festival sastra internasional menghadirkan penulis, pembaca, dan seniman dari berbagai negara untuk diskusi, pertunjukan, dan lokakarya.',
        '2025-10-29 00:00:00+07',
        '2025-11-02 23:59:59+07',
        false,
        'https://rhfhplcxngijmfanrxzo.supabase.co/storage/v1/object/public/cultour/cultour/images/events/ubud.jpg'
    ),
    (
        gen_random_uuid(),
        '0244478e-d0d7-4cfe-b868-aa608afc126b',
        (
            SELECT id
            FROM locations
            WHERE
                name = 'Tanah Lot Temple'
        ),
        'Bali Cultural Festival',
        'Festival budaya tahunan di Bali yang menampilkan berbagai kesenian tradisional dari seluruh Pulau Bali.',
        '2025-08-09 14:00:00+07',
        '2025-08-11 22:00:00+07',
        true,
        'https://rhfhplcxngijmfanrxzo.supabase.co/storage/v1/object/public/cultour/cultour/images/events/bali.jpg'
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
        'Tur sejarah berkeliling pusat kota Jakarta, mengunjungi landmark bersejarah dan museum.',
        '2025-09-15 08:00:00+07',
        '2025-09-15 16:00:00+07',
        true,
        'https://rhfhplcxngijmfanrxzo.supabase.co/storage/v1/object/public/cultour/cultour/images/events/jakarta.jpg'
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
        'Tur malam menelusuri bangunan bersejarah di Semarang, dengan pemandu yang berpakaian kostum era kolonial.',
        '2025-10-18 19:00:00+07',
        '2025-10-18 23:00:00+07',
        false,
        'https://rhfhplcxngijmfanrxzo.supabase.co/storage/v1/object/public/cultour/cultour/images/events/semarang.jpg'
    );

-- Verify the insertion
SELECT COUNT(*) as total_events FROM public.events;