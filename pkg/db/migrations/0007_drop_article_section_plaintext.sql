update article_sections
set body_markdown = case
    when nullif(trim(body_markdown), '') is null then coalesce(body_plaintext, '')
    else body_markdown
end;

alter table article_sections
    drop column if exists body_plaintext;
