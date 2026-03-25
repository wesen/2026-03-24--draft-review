insert into default_reaction_types (key, label, icon, position)
values
    ('useful', 'Useful', '★', 1),
    ('confused', 'Confused', '?', 2),
    ('slow', 'Slow', '◎', 3),
    ('favorite', 'Favorite', '♥', 4)
on conflict (key) do update
set
    label = excluded.label,
    icon = excluded.icon,
    position = excluded.position;
