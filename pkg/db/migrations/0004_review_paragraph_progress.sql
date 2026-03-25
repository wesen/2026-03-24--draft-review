create table if not exists review_paragraph_progress (
    review_session_id uuid not null references review_sessions(id) on delete cascade,
    section_id uuid not null references article_sections(id) on delete cascade,
    paragraph_key text not null,
    reached_at timestamptz not null default now(),
    completed_at timestamptz,
    primary key (review_session_id, section_id, paragraph_key)
);

create index if not exists idx_review_paragraph_progress_session
    on review_paragraph_progress(review_session_id, section_id);
