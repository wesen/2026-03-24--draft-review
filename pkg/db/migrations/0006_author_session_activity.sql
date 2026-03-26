alter table author_sessions
    add column if not exists last_used_at timestamptz;

update author_sessions
set last_used_at = coalesce(last_used_at, created_at)
where last_used_at is null;

alter table author_sessions
    alter column last_used_at set default now();

alter table author_sessions
    alter column last_used_at set not null;

create index if not exists idx_author_sessions_user_last_used
    on author_sessions(user_id, last_used_at desc);
