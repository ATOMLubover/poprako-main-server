INSERT INTO user_tbl(
    id,
    qq,
    nickname,
    password_hash,
    is_admin,
    assigned_translator_at,
    assigned_proofreader_at,
    assigned_typesetter_at,
    assigned_redrawer_at,
    assigned_reviewer_at,
    assigned_uploader_at,
    created_at,
    updated_at
) VALUES (
    '019bbf6f-fa6b-7119-b1cd-c961a808c864',
    '3384539248',
    'TestAdmin',
    '$2a$10$EA6AJEOBejXnJFmEzfbopezyBjh4FDhmlmA5XyqoMGWijL.vuh1K2',
    TRUE,
    NOW(),
    NOW(),
    NOW(),
    NOW(),
    NOW(),
    NOW(),
    NOW(),
    NOW()
) ON CONFLICT (id) DO NOTHING;