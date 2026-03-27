create table article_assets (
    id uuid primary key,
    article_id uuid not null references articles(id) on delete cascade,
    storage_key text not null unique,
    original_filename text not null,
    content_type text not null,
    byte_size bigint not null,
    created_at timestamptz not null default now()
);

create index idx_article_assets_article_id on article_assets(article_id, created_at desc);
