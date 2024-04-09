-- create table if not exists tags (
--     tag_id bigserial primary key
-- );
-- 
-- create table if not exists features (
--     feature_id bigserial primary key
-- );

create table if not exists banners (
    banner_id bigserial primary key,
    tag_ids integer[] not null,
    feature_id integer not null,
    content json not null,
    is_active boolean not null,
    created_at timestamp(0) with time zone not null default now(),
    updated_at timestamp(6) with time zone not null default now()
);

-- create table if not exists bannertags (
--     banner_id integer not null references banners(banner_id),
--     tag_id integer not null references tags(tag_id),
--     primary key (banner_id, tag_id)
-- );
