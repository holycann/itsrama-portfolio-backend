CREATE TABLE public.threads (
    id uuid NOT NULL DEFAULT gen_random_uuid (),
    created_at timestamp
    with
        time zone DEFAULT now(),
        updated_at timestamp
    with
        time zone NULL,
        event_id uuid NOT NULL,
        creator_id uuid NOT NULL,
        status character varying DEFAULT 'active',
        CONSTRAINT threads_pkey PRIMARY KEY (id),
        CONSTRAINT threads_event_id_fkey FOREIGN KEY (event_id) REFERENCES public.events (id),
        CONSTRAINT threads_user_id_fkey FOREIGN KEY (creator_id) REFERENCES auth.users (id)
);

INSERT INTO
    public.threads (id, event_id, status, creator_id)
VALUES (
        gen_random_uuid (),
        (
            SELECT id
            FROM events
            WHERE
                name = 'Bali Cultural Festival'
        ),
        'active',
        (
            SELECT id 
            FROM auth.users 
            WHERE email = 'admin@gmail.com' 
            LIMIT 1
        )
    ),
    (
        gen_random_uuid (),
        (
            SELECT id
            FROM events
            WHERE
                name = 'Jakarta Heritage Walk'
        ),
        'active',
        (
            SELECT id 
            FROM auth.users 
            WHERE email = 'admin@gmail.com' 
            LIMIT 1
        )
    ),
    (
        gen_random_uuid (),
        (
            SELECT id
            FROM events
            WHERE
                name = 'Semarang Night Heritage Tour'
        ),
        'active',
        (
            SELECT id 
            FROM auth.users 
            WHERE email = 'admin@gmail.com' 
            LIMIT 1
        )
    );