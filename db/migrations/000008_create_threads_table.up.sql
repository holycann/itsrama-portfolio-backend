CREATE TABLE public.threads (
    id uuid NOT NULL DEFAULT gen_random_uuid(),
    event_id uuid NOT NULL,
    creator_id uuid NOT NULL,
    status character varying(20) DEFAULT 'active' CHECK (status IN ('active', 'closed', 'archived')),
    created_at timestamp with time zone DEFAULT now(),
    updated_at timestamp with time zone NULL,
    CONSTRAINT threads_pkey PRIMARY KEY (id),
    CONSTRAINT threads_event_id_fkey FOREIGN KEY (event_id) REFERENCES public.events (id) ON DELETE CASCADE,
    CONSTRAINT threads_creator_id_fkey FOREIGN KEY (creator_id) REFERENCES auth.users (id) ON DELETE CASCADE,
    CONSTRAINT unique_event_thread UNIQUE (event_id)
);

-- Insert sample threads for initial events
INSERT INTO
    public.threads (id, event_id, creator_id, status)
VALUES (
        gen_random_uuid(),
        (
            SELECT id
            FROM events
            WHERE name = 'Bali Cultural Festival'
            LIMIT 1
        ),
        (
            SELECT id 
            FROM auth.users 
            WHERE email = 'admin@gmail.com' 
            LIMIT 1
        ),
        'active'
    ),
    (
        gen_random_uuid(),
        (
            SELECT id
            FROM events
            WHERE name = 'Jakarta Heritage Walk'
            LIMIT 1
        ),
        (
            SELECT id 
            FROM auth.users 
            WHERE email = 'admin@gmail.com' 
            LIMIT 1
        ),
        'active'
    ),
    (
        gen_random_uuid(),
        (
            SELECT id
            FROM events
            WHERE name = 'Semarang Night Heritage Tour'
            LIMIT 1
        ),
        (
            SELECT id 
            FROM auth.users 
            WHERE email = 'admin@gmail.com' 
            LIMIT 1
        ),
        'active'
    );

-- Verify the insertion
SELECT COUNT(*) as total_threads FROM public.threads;