-- Drop the trigger first to avoid potential errors
DROP TRIGGER IF EXISTS trg_update_location_geom ON public.locations;

-- Drop the function associated with the trigger
DROP FUNCTION IF EXISTS update_location_geom();

-- Drop the GIST index
DROP INDEX IF EXISTS locations_location_idx;

-- Drop the locations table
DROP TABLE IF EXISTS public.locations;

-- Drop the PostGIS extension
DROP EXTENSION IF EXISTS postgis;