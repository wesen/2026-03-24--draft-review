alter table reader_invites
    alter column email drop not null;

alter table reader_invites
    add column if not exists display_name text,
    add column if not exists identity_mode text not null default 'email',
    add column if not exists is_preview boolean not null default false;

update reader_invites
set identity_mode = case
        when coalesce(identity_mode, '') <> '' then identity_mode
        when email is not null and btrim(email) <> '' then 'email'
        else 'anonymous'
    end
where identity_mode is null or btrim(identity_mode) = '';

create index if not exists idx_reader_invites_article_identity_mode
    on reader_invites(article_id, identity_mode);
