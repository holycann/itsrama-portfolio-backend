CREATE TYPE public.message_type AS ENUM ('discussion', 'ai');

CREATE TABLE public.messages (
    id uuid NOT NULL DEFAULT gen_random_uuid (),
    thread_id uuid NOT NULL,
    user_id uuid NOT NULL,
    content text NOT NULL,
    type public.message_type NOT NULL DEFAULT 'discussion',
    created_at timestamp with time zone DEFAULT now(),
    updated_at timestamp with time zone NULL,
    CONSTRAINT messages_pkey PRIMARY KEY (id),
    CONSTRAINT messages_thread_id_fkey FOREIGN KEY (thread_id) REFERENCES public.threads (id),
    CONSTRAINT messages_user_id_fkey FOREIGN KEY (user_id) REFERENCES auth.users (id)
);

-- Dummy data for messages
WITH
    threads AS (
        SELECT id
        FROM public.threads
    )
INSERT INTO
    public.messages (
        id,
        thread_id,
        user_id,
        content,
        type
    )
VALUES (
        gen_random_uuid (),
        (
            SELECT id
            FROM threads
            LIMIT 1
        ),
        '0244478e-d0d7-4cfe-b868-aa608afc126b',
        'The Bali Cultural Festival was an incredible experience! The traditional dances were mesmerizing.',
        'discussion'
    ),
    (
        gen_random_uuid (),
        (
            SELECT id
            FROM threads
            LIMIT 1 OFFSET 1
        ),
        '120060ef-7c2e-4457-a677-c8f839e8e2a7',
        'I learned so much about Jakarta''s history during this walking tour. The guide was incredibly knowledgeable.',
        'discussion'
    ),
    (
        gen_random_uuid (),
        (
            SELECT id
            FROM threads
            LIMIT 1 OFFSET 2
        ),
        '34be7296-a530-41b0-872f-f6946441f49f',
        'Lawang Sewu at night is both beautiful and haunting. Such a rich historical site!',
        'discussion'
    );