-- Drop the trigger before dropping the table
DROP TRIGGER IF EXISTS trg_update_discussion_participant ON public.discussion_participants;
DROP FUNCTION IF EXISTS update_discussion_participant_timestamp();

-- Drop the table
DROP TABLE IF EXISTS public.discussion_participants;