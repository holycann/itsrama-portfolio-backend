-- Create message type enum
CREATE TYPE public.message_type AS ENUM ('discussion', 'ai');

-- Create messages table
CREATE TABLE public.messages (
    id uuid NOT NULL DEFAULT gen_random_uuid(),
    thread_id uuid NOT NULL,
    sender_id uuid NOT NULL,
    content text NOT NULL,
    type public.message_type NOT NULL DEFAULT 'discussion',
    created_at timestamp with time zone DEFAULT now(),
    updated_at timestamp with time zone NULL,
    CONSTRAINT messages_pkey PRIMARY KEY (id),
    CONSTRAINT messages_thread_id_fkey FOREIGN KEY (thread_id) REFERENCES public.threads (id) ON DELETE CASCADE,
    CONSTRAINT messages_sender_id_fkey FOREIGN KEY (sender_id) REFERENCES auth.users (id) ON DELETE CASCADE
);

-- Dummy data for messages
WITH
    threads AS (
        SELECT id
        FROM public.threads
        LIMIT 10  -- Ensure we have enough threads
    ),
    users AS (
        SELECT id
        FROM auth.users
        LIMIT 3
    )
INSERT INTO
    public.messages (
        thread_id,
        sender_id,
        content,
        type
    )
SELECT 
    threads.id,
    (SELECT id FROM users ORDER BY RANDOM() LIMIT 1),
    'Hei, ada yang tahu jam mulai acara ini?',
    'discussion'
FROM threads
LIMIT 8;