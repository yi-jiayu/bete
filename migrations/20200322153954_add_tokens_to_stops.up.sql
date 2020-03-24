alter table stops
    add column tokens tsvector;
create index stops_tokens_idx on stops using gin (tokens);
