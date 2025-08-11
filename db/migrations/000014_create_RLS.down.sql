-- Drop Row Level Security Policies and Disable RLS for all tables

-- Drop policies for users_profile table
DROP POLICY "Users can view all profiles" ON public.users_profile;
DROP POLICY "Users can update their own profile" ON public.users_profile;
DROP POLICY "Users can create their own profile" ON public.users_profile;
DROP POLICY "Admins can manage all profiles" ON public.users_profile;
ALTER TABLE public.users_profile DISABLE ROW LEVEL SECURITY;

-- Drop policies for provinces table
DROP POLICY "All authenticated users can view provinces" ON public.provinces;
DROP POLICY "Only admins can modify provinces" ON public.provinces;
ALTER TABLE public.provinces DISABLE ROW LEVEL SECURITY;

-- Drop policies for cities table
DROP POLICY "All authenticated users can view cities" ON public.cities;
DROP POLICY "Only admins can modify cities" ON public.cities;
ALTER TABLE public.cities DISABLE ROW LEVEL SECURITY;

-- Drop policies for locations table
DROP POLICY "All authenticated users can view locations" ON public.locations;
DROP POLICY "Users can create locations" ON public.locations;
DROP POLICY "Users can update locations for their events" ON public.locations;
DROP POLICY "Only admins can delete locations" ON public.locations;
ALTER TABLE public.locations DISABLE ROW LEVEL SECURITY;

-- Drop policies for badges table
DROP POLICY "All authenticated users can view badges" ON public.badges;
DROP POLICY "Only admins can modify badges" ON public.badges;
ALTER TABLE public.badges DISABLE ROW LEVEL SECURITY;

-- Drop policies for users_badge table
DROP POLICY "Users can view their own badges" ON public.users_badge;
DROP POLICY "Admins can view all user badges" ON public.users_badge;
DROP POLICY "Users cannot manually assign badges" ON public.users_badge;
DROP POLICY "Only admins can manage user badges" ON public.users_badge;
ALTER TABLE public.users_badge DISABLE ROW LEVEL SECURITY;

-- Drop policies for events table
DROP POLICY "All authenticated users can view events" ON public.events;
DROP POLICY "Users can create their own events" ON public.events;
DROP POLICY "Users can update their own events" ON public.events;
DROP POLICY "Users can delete their own events" ON public.events;
DROP POLICY "Admins can manage all events" ON public.events;
ALTER TABLE public.events DISABLE ROW LEVEL SECURITY;

-- Drop policies for threads table
DROP POLICY "All authenticated users can view threads" ON public.threads;
DROP POLICY "Users can create threads for their events" ON public.threads;
DROP POLICY "Users can update their threads" ON public.threads;
DROP POLICY "Users can delete their threads" ON public.threads;
ALTER TABLE public.threads DISABLE ROW LEVEL SECURITY;

-- Drop policies for messages table
DROP POLICY "Users can view messages in their threads" ON public.messages;
DROP POLICY "Users can create messages in their threads" ON public.messages;
DROP POLICY "Users can update their own messages" ON public.messages;
DROP POLICY "Users can delete their own messages" ON public.messages;
DROP POLICY "Admins can manage all messages" ON public.messages;
ALTER TABLE public.messages DISABLE ROW LEVEL SECURITY;

-- Drop policies for event_views table
DROP POLICY "Users can view their own event views" ON public.event_views;
DROP POLICY "Users can create event views for themselves" ON public.event_views;
DROP POLICY "Users can update their own event views" ON public.event_views;
ALTER TABLE public.event_views DISABLE ROW LEVEL SECURITY;

-- Drop policies for discussion_participants table
DROP POLICY "Users can view participants in their threads" ON public.discussion_participants;
DROP POLICY "Users can join thread discussions" ON public.discussion_participants;
DROP POLICY "Users can leave threads they joined" ON public.discussion_participants;
DROP POLICY "Admins can manage all discussion participants" ON public.discussion_participants;
ALTER TABLE public.discussion_participants DISABLE ROW LEVEL SECURITY;
