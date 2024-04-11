create table if not exists tokens (
    token_hash bytea primary key,
    user_id bigint not null references users(user_id) on delete cascade
);
