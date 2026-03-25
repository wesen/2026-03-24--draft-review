create extension if not exists pgcrypto;

create table if not exists users (
    id uuid primary key default gen_random_uuid(),
    email text not null unique,
    name text not null,
    password_hash text not null,
    email_verified_at timestamptz,
    created_at timestamptz not null default now(),
    updated_at timestamptz not null default now()
);

create table if not exists author_sessions (
    id uuid primary key default gen_random_uuid(),
    user_id uuid not null references users(id) on delete cascade,
    token_hash text not null unique,
    expires_at timestamptz not null,
    created_at timestamptz not null default now(),
    revoked_at timestamptz
);

create table if not exists password_reset_tokens (
    id uuid primary key default gen_random_uuid(),
    user_id uuid not null references users(id) on delete cascade,
    token_hash text not null unique,
    expires_at timestamptz not null,
    consumed_at timestamptz,
    created_at timestamptz not null default now()
);

create table if not exists email_verification_tokens (
    id uuid primary key default gen_random_uuid(),
    user_id uuid not null references users(id) on delete cascade,
    token_hash text not null unique,
    expires_at timestamptz not null,
    consumed_at timestamptz,
    created_at timestamptz not null default now()
);

create table if not exists articles (
    id uuid primary key default gen_random_uuid(),
    owner_user_id uuid not null references users(id) on delete cascade,
    current_version_id uuid,
    title text not null,
    author_display_name text not null,
    status text not null check (status in ('draft', 'in_review', 'complete', 'archived')),
    access_mode text not null check (access_mode in ('invite_link', 'link', 'password')),
    share_token text unique,
    share_password_hash text,
    show_author_note boolean not null default true,
    reader_can_see_reactions boolean not null default true,
    reader_can_see_names boolean not null default false,
    require_note boolean not null default false,
    allow_anonymous boolean not null default true,
    created_at timestamptz not null default now(),
    updated_at timestamptz not null default now()
);

create table if not exists article_versions (
    id uuid primary key default gen_random_uuid(),
    article_id uuid not null references articles(id) on delete cascade,
    version_number integer not null,
    version_label text not null,
    intro text not null,
    author_note text not null,
    import_source text,
    import_metadata jsonb not null default '{}'::jsonb,
    created_at timestamptz not null default now(),
    unique (article_id, version_number)
);

alter table articles
    add constraint articles_current_version_fkey
    foreign key (current_version_id) references article_versions(id)
    deferrable initially deferred;

create table if not exists article_sections (
    id uuid primary key default gen_random_uuid(),
    article_version_id uuid not null references article_versions(id) on delete cascade,
    section_key text not null,
    position integer not null,
    title text not null,
    body_markdown text not null,
    body_plaintext text not null,
    estimated_read_seconds integer not null default 0,
    created_at timestamptz not null default now(),
    updated_at timestamptz not null default now(),
    unique (article_version_id, section_key),
    unique (article_version_id, position)
);

create table if not exists default_reaction_types (
    key text primary key,
    label text not null,
    icon text not null,
    position integer not null unique
);

create table if not exists article_reaction_types (
    id uuid primary key default gen_random_uuid(),
    article_id uuid not null references articles(id) on delete cascade,
    type_key text not null,
    label text not null,
    icon text not null,
    is_default boolean not null default true,
    enabled boolean not null default true,
    position integer not null,
    unique (article_id, type_key)
);

create table if not exists reader_invites (
    id uuid primary key default gen_random_uuid(),
    article_id uuid not null references articles(id) on delete cascade,
    email text not null,
    invite_token text not null unique,
    invite_note text,
    sent_at timestamptz not null default now(),
    opened_at timestamptz,
    revoked_at timestamptz
);

create table if not exists review_sessions (
    id uuid primary key default gen_random_uuid(),
    article_id uuid not null references articles(id) on delete cascade,
    article_version_id uuid not null references article_versions(id) on delete cascade,
    reader_invite_id uuid references reader_invites(id) on delete set null,
    reader_name text,
    reader_email text,
    is_anonymous boolean not null default false,
    started_at timestamptz not null default now(),
    last_active_at timestamptz not null default now(),
    completed_at timestamptz,
    progress_percent integer not null default 0
);

create table if not exists review_section_progress (
    review_session_id uuid not null references review_sessions(id) on delete cascade,
    section_id uuid not null references article_sections(id) on delete cascade,
    reached_at timestamptz not null default now(),
    completed_at timestamptz,
    primary key (review_session_id, section_id)
);

create table if not exists reactions (
    id uuid primary key default gen_random_uuid(),
    review_session_id uuid not null references review_sessions(id) on delete cascade,
    article_id uuid not null references articles(id) on delete cascade,
    article_version_id uuid not null references article_versions(id) on delete cascade,
    section_id uuid not null references article_sections(id) on delete cascade,
    paragraph_key text not null,
    reaction_type_key text not null,
    comment_text text not null,
    created_at timestamptz not null default now(),
    resolved_at timestamptz,
    author_note text
);

create table if not exists review_summaries (
    review_session_id uuid primary key references review_sessions(id) on delete cascade,
    overall_thoughts text,
    recommendability text check (recommendability in ('maybe', 'yes', 'absolutely')),
    notify_new_version boolean not null default false,
    submitted_at timestamptz not null default now()
);

create index if not exists idx_articles_owner_status_updated on articles(owner_user_id, status, updated_at desc);
create index if not exists idx_article_versions_article_version on article_versions(article_id, version_number desc);
create index if not exists idx_article_sections_version_position on article_sections(article_version_id, position);
create index if not exists idx_reader_invites_article_email on reader_invites(article_id, email);
create index if not exists idx_review_sessions_article_last_active on review_sessions(article_id, last_active_at desc);
create index if not exists idx_reactions_article_created on reactions(article_id, created_at desc);
create index if not exists idx_reactions_section_type on reactions(section_id, reaction_type_key);
