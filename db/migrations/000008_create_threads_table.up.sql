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
WITH 
    events_with_ids AS (
        SELECT id, name 
        FROM events 
        WHERE name IN ('Bali Cultural Festival', 'Jakarta Heritage Walk', 'Semarang Night Heritage Tour')
    ),
    admin_user AS (
        SELECT id 
        FROM auth.users 
        WHERE email = 'admin@gmail.com' 
        LIMIT 1
    )
INSERT INTO
    public.threads (id, event_id, creator_id, status)
SELECT 
    gen_random_uuid(),
    id,
    (SELECT id FROM admin_user),
    'active'
FROM 
    events_with_ids;

-- Verify the insertion
SELECT COUNT(*) as total_threads FROM public.threads;