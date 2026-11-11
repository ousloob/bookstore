-- Demo users. The password hash below is bcrypt of "gophers!" for both rows:
-- demo data only, never load this into an environment that matters.
--
-- ON CONFLICT DO NOTHING makes the file idempotent, so `admin seed` can run
-- as often as you like. The ids are fixed for the same reason.
INSERT INTO users (
    user_id,
    name,
    email,
    roles,
    password_hash,
    date_created,
    date_updated
)
VALUES
    (
        '01a1cb60-1a33-7b12-8a2f-a78266fa47a4',
        'Staff Gopher',
        'staff@bookstore.test',
        ARRAY['staff'],
        '$2a$10$a/C5joZFTCfFid5PM4Wy/eYni7tscsIpLacGcMbY7PPmdZ/WBzEMe'::bytea,
        '2026-01-01 00:00:00+00',
        '2026-01-01 00:00:00+00'
    ),
    (
        '01a1cb60-1a33-7b60-9b71-a1f19f7e2af6',
        'Customer Gopher',
        'customer@bookstore.test',
        ARRAY['customer'],
        '$2a$10$a/C5joZFTCfFid5PM4Wy/eYni7tscsIpLacGcMbY7PPmdZ/WBzEMe'::bytea,
        '2026-01-01 00:00:00+00',
        '2026-01-01 00:00:00+00'
    )
ON CONFLICT (email) DO NOTHING;
