DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'user_role') THEN
        CREATE TYPE user_role AS ENUM ('user', 'admin');
    END IF;
END
$$;

create table if not exists users (
    user_id bigserial primary key,
    username text unique not null,
    role user_role not null,
    password_hash bytea not null,
    created_at timestamp(0) with time zone not null default now()
);
