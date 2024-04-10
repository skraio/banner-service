create table if not exists banners (
    banner_id bigserial primary key,
    tag_ids integer[] not null,
    feature_id integer not null,
    content json not null,
    is_active boolean not null,
    created_at timestamp(3) with time zone not null default now(),
    updated_at timestamp(3) with time zone not null default now()
);
