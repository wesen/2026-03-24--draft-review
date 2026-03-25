alter table users
    add column if not exists auth_subject text,
    add column if not exists auth_issuer text;

create unique index if not exists idx_users_auth_identity
    on users(auth_issuer, auth_subject)
    where auth_subject is not null;

update users
set auth_issuer = 'dev',
    auth_subject = 'local-author',
    updated_at = now()
where id = '11111111-1111-1111-1111-111111111111'
  and (auth_issuer is null or auth_subject is null);
