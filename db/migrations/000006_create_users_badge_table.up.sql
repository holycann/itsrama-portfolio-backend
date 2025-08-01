CREATE TABLE public.users_badge (
    id uuid NOT NULL DEFAULT gen_random_uuid (),
    user_id uuid NOT NULL,
    badge_id uuid NOT NULL,
    created_at timestamp
    with
        time zone DEFAULT now(),
        updated_at timestamp
    with
        time zone NULL,
        CONSTRAINT users_badge_pkey PRIMARY KEY (id),
        CONSTRAINT users_badge_user_id_fkey FOREIGN KEY (user_id) REFERENCES auth.users (id),
        CONSTRAINT users_badge_badge_id_fkey FOREIGN KEY (badge_id) REFERENCES public.badges (id),
        CONSTRAINT users_badge_unique_user_badge UNIQUE (user_id, badge_id)
);

-- Dummy data for users badge
WITH
    badges AS (
        SELECT id, name
        FROM public.badges
    )
INSERT INTO
    public.users_badge (id, user_id, badge_id)
VALUES (
        gen_random_uuid (),
        '0244478e-d0d7-4cfe-b868-aa608afc126b',
        (
            SELECT id
            FROM badges
            WHERE
                name = 'Penjelajah'
        )
    ),
    (
        gen_random_uuid (),
        '120060ef-7c2e-4457-a677-c8f839e8e2a7',
        (
            SELECT id
            FROM badges
            WHERE
                name = 'Penjelajah'
        )
    ),
    (
        gen_random_uuid (),
        '34be7296-a530-41b0-872f-f6946441f49f',
        (
            SELECT id
            FROM badges
            WHERE
                name = 'Penjelajah'
        )
    ),
    (
        gen_random_uuid (),
        '609ae64c-78e7-4d49-9300-803ffcab4547',
        (
            SELECT id
            FROM badges
            WHERE
                name = 'Penjelajah'
        )
    ),
    (
        gen_random_uuid (),
        'b0bd223b-d25b-477b-8350-d582b8fb12f1',
        (
            SELECT id
            FROM badges
            WHERE
                name = 'Penjelajah'
        )
    );