CREATE TABLE public.users_profile (
    id uuid NOT NULL DEFAULT gen_random_uuid (),
    user_id uuid NOT NULL,
    fullname character varying,
    bio text,
    avatar_url character varying,
    identity_image_url character varying,
    created_at timestamp
    with
        time zone NOT NULL DEFAULT now(),
        CONSTRAINT users_profile_pkey PRIMARY KEY (id),
        CONSTRAINT users_profile_user_id_fkey FOREIGN KEY (user_id) REFERENCES auth.users (id),
        CONSTRAINT users_profile_unique_user_id UNIQUE (user_id)
);

-- Dummy data for users profile using provided user IDs
INSERT INTO
    public.users_profile (
        id,
        user_id,
        fullname,
        bio,
        avatar_url,
        identity_image_url
    )
VALUES (
        gen_random_uuid (),
        '0244478e-d0d7-4cfe-b868-aa608afc126b',
        'Budi Santoso',
        'Passionate about Indonesian culture and travel',
        'https://placehold.co/100x100.png?text=Budi',
        'https://placehold.co/200x300.png?text=ID'
    ),
    (
        gen_random_uuid (),
        '120060ef-7c2e-4457-a677-c8f839e8e2a7',
        'Ani Widya',
        'Local storyteller and culture enthusiast',
        'https://placehold.co/100x100.png?text=Ani',
        'https://placehold.co/200x300.png?text=ID'
    ),
    (
        gen_random_uuid (),
        '34be7296-a530-41b0-872f-f6946441f49f',
        'Rudi Hartono',
        'Explorer of hidden gems across Indonesia',
        'https://placehold.co/100x100.png?text=Rudi',
        'https://placehold.co/200x300.png?text=ID'
    ),
    (
        gen_random_uuid (),
        '609ae64c-78e7-4d49-9300-803ffcab4547',
        'Siti Nurhaliza',
        'Cultural researcher and event organizer',
        'https://placehold.co/100x100.png?text=Siti',
        'https://placehold.co/200x300.png?text=ID'
    ),
    (
        gen_random_uuid (),
        'b0bd223b-d25b-477b-8350-d582b8fb12f1',
        'Eko Prasetyo',
        'Travel blogger documenting Indonesian heritage',
        'https://placehold.co/100x100.png?text=Eko',
        'https://placehold.co/200x300.png?text=ID'
    );