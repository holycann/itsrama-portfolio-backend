CREATE TABLE public.discussion_participants (
    thread_id uuid NOT NULL,
    user_id uuid NOT NULL,
    created_at timestamp with time zone DEFAULT now(),
    CONSTRAINT discussion_participants_pkey PRIMARY KEY (thread_id, user_id),
    CONSTRAINT discussion_participants_thread_id_fkey FOREIGN KEY (thread_id) REFERENCES public.threads (id),
    CONSTRAINT discussion_participants_user_id_fkey FOREIGN KEY (user_id) REFERENCES auth.users (id)
);

INSERT INTO public.discussion_participants (thread_id, user_id)
VALUES (
    (SELECT id FROM threads LIMIT 1),
    (SELECT id FROM auth.users WHERE email = 'admin@gmail.com' LIMIT 1)
),
(
    (SELECT id FROM threads LIMIT 1 OFFSET 1),
    (SELECT id FROM auth.users WHERE email = 'admin@gmail.com' LIMIT 1)
),
(
    (SELECT id FROM threads LIMIT 1 OFFSET 2),
    (SELECT id FROM auth.users WHERE email = 'admin@gmail.com' LIMIT 1)
);
