-- =====================================================
-- SUPABASE ROW LEVEL SECURITY (RLS) POLICIES
-- =====================================================
-- This file contains RLS policies for all tables in the application
-- Assumes auth.users table has a 'role' column or uses user metadata

-- =====================================================
-- USERS_PROFILE TABLE POLICIES
-- =====================================================

-- Enable Row Level Security
ALTER TABLE public.users_profile ENABLE ROW LEVEL SECURITY;

-- Policy: Users can view all profiles (public information)
CREATE POLICY "Users can view all profiles" ON public.users_profile 
FOR SELECT 
USING (true);

-- Policy: Users can only update their own profile
CREATE POLICY "Users can update their own profile" ON public.users_profile 
FOR UPDATE 
USING (auth.uid() = user_id);

-- Policy: Users can only create their own profile during registration
CREATE POLICY "Users can create their own profile" ON public.users_profile 
FOR INSERT 
WITH CHECK (auth.uid() = user_id);

-- Policy: Admins can manage all profiles
-- Note: This assumes 'role' is stored in auth.users metadata or a custom column
CREATE POLICY "Admins can manage all profiles" ON public.users_profile 
FOR ALL 
USING (
    COALESCE(
        (auth.jwt() -> 'user_metadata' ->> 'role'),
        (SELECT raw_user_meta_data ->> 'role' FROM auth.users WHERE id = auth.uid())
    ) = 'admin'
);

-- =====================================================
-- PROVINCES TABLE POLICIES
-- =====================================================

-- Enable Row Level Security
ALTER TABLE public.provinces ENABLE ROW LEVEL SECURITY;

-- Policy: All authenticated users can view provinces (reference data)
CREATE POLICY "All authenticated users can view provinces" ON public.provinces 
FOR SELECT 
USING (auth.uid() IS NOT NULL);

-- Policy: Only admins can modify provinces
CREATE POLICY "Only admins can modify provinces" ON public.provinces 
FOR ALL 
USING (
    COALESCE(
        (auth.jwt() -> 'user_metadata' ->> 'role'),
        (SELECT raw_user_meta_data ->> 'role' FROM auth.users WHERE id = auth.uid())
    ) = 'admin'
);

-- =====================================================
-- CITIES TABLE POLICIES
-- =====================================================

-- Enable Row Level Security
ALTER TABLE public.cities ENABLE ROW LEVEL SECURITY;

-- Policy: All authenticated users can view cities (reference data)
CREATE POLICY "All authenticated users can view cities" ON public.cities 
FOR SELECT 
USING (auth.uid() IS NOT NULL);

-- Policy: Only admins can modify cities
CREATE POLICY "Only admins can modify cities" ON public.cities 
FOR ALL 
USING (
    COALESCE(
        (auth.jwt() -> 'user_metadata' ->> 'role'),
        (SELECT raw_user_meta_data ->> 'role' FROM auth.users WHERE id = auth.uid())
    ) = 'admin'
);

-- =====================================================
-- LOCATIONS TABLE POLICIES
-- =====================================================

-- Enable Row Level Security
ALTER TABLE public.locations ENABLE ROW LEVEL SECURITY;

-- Policy: All authenticated users can view locations
CREATE POLICY "All authenticated users can view locations" ON public.locations 
FOR SELECT 
USING (auth.uid() IS NOT NULL);

-- Policy: Users can create locations (general creation allowed)
-- Note: Original policy had circular dependency with events table
CREATE POLICY "Users can create locations" ON public.locations 
FOR INSERT 
WITH CHECK (auth.uid() IS NOT NULL);

-- Policy: Users can update locations for their events
CREATE POLICY "Users can update locations for their events" ON public.locations 
FOR UPDATE 
USING (
    EXISTS (
        SELECT 1 
        FROM public.events e 
        WHERE e.user_id = auth.uid() 
        AND e.location_id = locations.id
    )
);

-- Policy: Only admins can delete locations
CREATE POLICY "Only admins can delete locations" ON public.locations 
FOR DELETE 
USING (
    COALESCE(
        (auth.jwt() -> 'user_metadata' ->> 'role'),
        (SELECT raw_user_meta_data ->> 'role' FROM auth.users WHERE id = auth.uid())
    ) = 'admin'
);

-- =====================================================
-- BADGES TABLE POLICIES
-- =====================================================

-- Enable Row Level Security
ALTER TABLE public.badges ENABLE ROW LEVEL SECURITY;

-- Policy: All authenticated users can view badges (reference data)
CREATE POLICY "All authenticated users can view badges" ON public.badges 
FOR SELECT 
USING (auth.uid() IS NOT NULL);

-- Policy: Only admins can modify badges
CREATE POLICY "Only admins can modify badges" ON public.badges 
FOR ALL 
USING (
    COALESCE(
        (auth.jwt() -> 'user_metadata' ->> 'role'),
        (SELECT raw_user_meta_data ->> 'role' FROM auth.users WHERE id = auth.uid())
    ) = 'admin'
);

-- =====================================================
-- USERS_BADGE TABLE POLICIES
-- =====================================================

-- Enable Row Level Security
ALTER TABLE public.users_badge ENABLE ROW LEVEL SECURITY;

-- Policy: Users can view their own badges
CREATE POLICY "Users can view their own badges" ON public.users_badge 
FOR SELECT 
USING (auth.uid() = user_id);

-- Policy: Admins can view all user badges
CREATE POLICY "Admins can view all user badges" ON public.users_badge 
FOR SELECT 
USING (
    COALESCE(
        (auth.jwt() -> 'user_metadata' ->> 'role'),
        (SELECT raw_user_meta_data ->> 'role' FROM auth.users WHERE id = auth.uid())
    ) = 'admin'
);

-- Policy: Prevent manual badge assignment by users
-- Note: Badges should be assigned through database functions or admin actions
CREATE POLICY "Users cannot manually assign badges" ON public.users_badge 
FOR INSERT 
WITH CHECK (false);

-- Policy: Only admins can manage user badges
CREATE POLICY "Only admins can manage user badges" ON public.users_badge 
FOR ALL 
USING (
    COALESCE(
        (auth.jwt() -> 'user_metadata' ->> 'role'),
        (SELECT raw_user_meta_data ->> 'role' FROM auth.users WHERE id = auth.uid())
    ) = 'admin'
);

-- =====================================================
-- EVENTS TABLE POLICIES
-- =====================================================

-- Enable Row Level Security
ALTER TABLE public.events ENABLE ROW LEVEL SECURITY;

-- Policy: All authenticated users can view events
CREATE POLICY "All authenticated users can view events" ON public.events 
FOR SELECT 
USING (auth.uid() IS NOT NULL);

-- Policy: Users can create their own events
CREATE POLICY "Users can create their own events" ON public.events 
FOR INSERT 
WITH CHECK (auth.uid() = user_id);

-- Policy: Users can update their own events
CREATE POLICY "Users can update their own events" ON public.events 
FOR UPDATE 
USING (auth.uid() = user_id);

-- Policy: Users can delete their own events
CREATE POLICY "Users can delete their own events" ON public.events 
FOR DELETE 
USING (auth.uid() = user_id);

-- Policy: Admins can manage all events
CREATE POLICY "Admins can manage all events" ON public.events 
FOR ALL 
USING (
    COALESCE(
        (auth.jwt() -> 'user_metadata' ->> 'role'),
        (SELECT raw_user_meta_data ->> 'role' FROM auth.users WHERE id = auth.uid())
    ) = 'admin'
);

-- =====================================================
-- THREADS TABLE POLICIES
-- =====================================================

-- Enable Row Level Security
ALTER TABLE public.threads ENABLE ROW LEVEL SECURITY;

-- Policy: All authenticated users can view threads
CREATE POLICY "All authenticated users can view threads" ON public.threads 
FOR SELECT 
USING (auth.uid() IS NOT NULL);

-- Policy: Users can create threads for their events
-- Note: Simplified to avoid circular dependency
CREATE POLICY "Users can create threads for their events" ON public.threads 
FOR INSERT 
WITH CHECK (
    auth.uid() = creator_id
    AND EXISTS (
        SELECT 1 
        FROM public.events e 
        WHERE e.id = event_id 
        AND e.user_id = auth.uid()
    )
);

-- Policy: Thread creators and event owners can update threads
CREATE POLICY "Users can update their threads" ON public.threads 
FOR UPDATE 
USING (
    auth.uid() = creator_id
    OR EXISTS (
        SELECT 1 
        FROM public.events e 
        WHERE e.id = event_id 
        AND e.user_id = auth.uid()
    )
);

-- Policy: Thread creators can delete their threads
CREATE POLICY "Users can delete their threads" ON public.threads 
FOR DELETE 
USING (
    auth.uid() = creator_id
    OR EXISTS (
        SELECT 1 
        FROM public.events e 
        WHERE e.id = event_id 
        AND e.user_id = auth.uid()
    )
);

-- =====================================================
-- MESSAGES TABLE POLICIES
-- =====================================================

-- Enable Row Level Security
ALTER TABLE public.messages ENABLE ROW LEVEL SECURITY;

-- Policy: Users can view messages in threads they participate in
CREATE POLICY "Users can view messages in their threads" ON public.messages 
FOR SELECT 
USING (
    -- User is a participant in the thread
    EXISTS (
        SELECT 1 
        FROM public.discussion_participants dp 
        WHERE dp.thread_id = messages.thread_id 
        AND dp.user_id = auth.uid()
    )
    -- User is the thread creator
    OR EXISTS (
        SELECT 1 
        FROM public.threads t 
        WHERE t.id = messages.thread_id 
        AND t.creator_id = auth.uid()
    )
    -- User is the event owner
    OR EXISTS (
        SELECT 1 
        FROM public.events e
        JOIN public.threads t ON t.event_id = e.id 
        WHERE t.id = messages.thread_id 
        AND e.user_id = auth.uid()
    )
);

-- Policy: Users can create messages in authorized threads
CREATE POLICY "Users can create messages in their threads" ON public.messages 
FOR INSERT 
WITH CHECK (
    auth.uid() = sender_id
    AND (
        -- User is a participant in the thread
        EXISTS (
            SELECT 1 
            FROM public.discussion_participants dp 
            WHERE dp.thread_id = thread_id 
            AND dp.user_id = auth.uid()
        )
        -- User is the thread creator
        OR EXISTS (
            SELECT 1 
            FROM public.threads t 
            WHERE t.id = thread_id 
            AND t.creator_id = auth.uid()
        )
        -- User is the event owner
        OR EXISTS (
            SELECT 1 
            FROM public.events e
            JOIN public.threads t ON t.event_id = e.id 
            WHERE t.id = thread_id 
            AND e.user_id = auth.uid()
        )
    )
);

-- Policy: Users can update their own messages
CREATE POLICY "Users can update their own messages" ON public.messages 
FOR UPDATE 
USING (auth.uid() = sender_id);

-- Policy: Users can delete their own messages
CREATE POLICY "Users can delete their own messages" ON public.messages 
FOR DELETE 
USING (auth.uid() = sender_id);

-- Policy: Admins can manage all messages
CREATE POLICY "Admins can manage all messages" ON public.messages 
FOR ALL 
USING (
    COALESCE(
        (auth.jwt() -> 'user_metadata' ->> 'role'),
        (SELECT raw_user_meta_data ->> 'role' FROM auth.users WHERE id = auth.uid())
    ) = 'admin'
);

-- =====================================================
-- EVENT_VIEWS TABLE POLICIES
-- =====================================================

-- Enable Row Level Security
ALTER TABLE public.event_views ENABLE ROW LEVEL SECURITY;

-- Policy: Users can view their own event views
CREATE POLICY "Users can view their own event views" ON public.event_views 
FOR SELECT 
USING (auth.uid() = user_id);

-- Policy: Users can create event views for themselves
CREATE POLICY "Users can create event views for themselves" ON public.event_views 
FOR INSERT 
WITH CHECK (auth.uid() = user_id);

-- Policy: Users can update their own event views (if needed)
CREATE POLICY "Users can update their own event views" ON public.event_views 
FOR UPDATE 
USING (auth.uid() = user_id);

-- =====================================================
-- DISCUSSION_PARTICIPANTS TABLE POLICIES
-- =====================================================

-- Enable Row Level Security
ALTER TABLE public.discussion_participants ENABLE ROW LEVEL SECURITY;

-- Policy: Users can view participants in threads they are part of
CREATE POLICY "Users can view participants in their threads" ON public.discussion_participants 
FOR SELECT 
USING (
    -- User is viewing their own participation
    auth.uid() = user_id
    -- User is the thread creator
    OR EXISTS (
        SELECT 1 
        FROM public.threads t 
        WHERE t.id = thread_id 
        AND t.creator_id = auth.uid()
    )
    -- User is the event owner
    OR EXISTS (
        SELECT 1 
        FROM public.events e
        JOIN public.threads t ON t.event_id = e.id 
        WHERE t.id = thread_id 
        AND e.user_id = auth.uid()
    )
    -- User is also a participant in the same thread
    OR EXISTS (
        SELECT 1 
        FROM public.discussion_participants dp2 
        WHERE dp2.thread_id = discussion_participants.thread_id 
        AND dp2.user_id = auth.uid()
    )
);

-- Policy: Users can join threads for events or open discussions
CREATE POLICY "Users can join thread discussions" ON public.discussion_participants 
FOR INSERT 
WITH CHECK (
    auth.uid() = user_id
    AND (
        -- Thread belongs to an event (anyone can join event discussions)
        EXISTS (
            SELECT 1 
            FROM public.threads t 
            WHERE t.id = thread_id 
            AND t.event_id IS NOT NULL
        )
        -- Thread creator allows participation
        OR EXISTS (
            SELECT 1 
            FROM public.threads t 
            WHERE t.id = thread_id 
            AND t.creator_id = auth.uid()
        )
        -- Thread is marked as public/open (if you have such a field)
        OR EXISTS (
            SELECT 1 
            FROM public.threads t 
            WHERE t.id = thread_id 
            AND COALESCE(t.status, 'active') = 'active'
        )
    )
);

-- Policy: Users can leave threads they have joined
CREATE POLICY "Users can leave threads they joined" ON public.discussion_participants 
FOR DELETE 
USING (auth.uid() = user_id);

-- Policy: Admins can manage all discussion participants
CREATE POLICY "Admins can manage all discussion participants" ON public.discussion_participants 
FOR ALL 
USING (
    COALESCE(
        (auth.jwt() -> 'user_metadata' ->> 'role'),
        (SELECT raw_user_meta_data ->> 'role' FROM auth.users WHERE id = auth.uid())
    ) = 'admin'
);