CREATE TABLE public.discussion_participants (
    thread_id uuid NOT NULL,
    user_id uuid NOT NULL,
    joined_at timestamp with time zone DEFAULT now(),
    updated_at timestamp with time zone DEFAULT now(),
    CONSTRAINT discussion_participants_pkey PRIMARY KEY (thread_id, user_id),
    CONSTRAINT discussion_participants_thread_id_fkey FOREIGN KEY (thread_id) REFERENCES public.threads (id),
    CONSTRAINT discussion_participants_user_id_fkey FOREIGN KEY (user_id) REFERENCES auth.users (id)
);

-- Trigger to automatically update the updated_at timestamp
CREATE OR REPLACE FUNCTION update_discussion_participant_timestamp()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trg_update_discussion_participant
BEFORE UPDATE ON public.discussion_participants
FOR EACH ROW
EXECUTE FUNCTION update_discussion_participant_timestamp();

-- Initial seed data for discussion participants
INSERT INTO public.discussion_participants (thread_id, user_id, joined_at, updated_at)
SELECT 
    t.id, 
    u.id, 
    NOW(), 
    NOW()
FROM 
    threads t
CROSS JOIN 
    (SELECT id FROM auth.users WHERE email = 'admin@gmail.com' LIMIT 3) u
LIMIT 3;

-- Verify the insertion
SELECT COUNT(*) as total_discussion_participants FROM public.discussion_participants;
