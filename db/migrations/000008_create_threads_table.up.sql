CREATE TABLE public.threads (
    id uuid NOT NULL DEFAULT gen_random_uuid (),
    created_at timestamp
    with
        time zone NOT NULL DEFAULT now(),
        title character varying NOT NULL,
        event_id uuid NOT NULL,
        status character varying DEFAULT 'active',
        CONSTRAINT threads_pkey PRIMARY KEY (id),
        CONSTRAINT threads_event_id_fkey FOREIGN KEY (event_id) REFERENCES public.events (id)
);

-- Dummy data for threads
WITH
    events AS (
        SELECT id, name
        FROM public.events
    )
INSERT INTO
    public.threads (id, title, event_id, status)
VALUES (
        gen_random_uuid (),
        'Sharing Experiences at Bali Cultural Festival',
        (
            SELECT id
            FROM events
            WHERE
                name = 'Bali Cultural Festival'
        ),
        'active'
    ),
    (
        gen_random_uuid (),
        'Hidden Stories of Jakarta Heritage Walk',
        (
            SELECT id
            FROM events
            WHERE
                name = 'Jakarta Heritage Walk'
        ),
        'active'
    ),
    (
        gen_random_uuid (),
        'Exploring Semarang''s Night Secrets',
        (
            SELECT id
            FROM events
            WHERE
                name = 'Semarang Night Heritage Tour'
        ),
        'active'
    );