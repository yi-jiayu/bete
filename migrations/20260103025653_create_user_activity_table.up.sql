create table user_activity
(
    id      bigint generated always as identity,
    user_id bigint                   not null,
    type    text                     not null,
    ts      timestamp with time zone not null default now(),
    data    jsonb                    not null
)