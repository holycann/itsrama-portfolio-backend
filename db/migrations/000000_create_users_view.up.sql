CREATE VIEW public.users_view AS
SELECT
    id,
    email,
    phone,
    role,
    last_sign_in_at,
    created_at,
    updated_at
FROM auth.users;