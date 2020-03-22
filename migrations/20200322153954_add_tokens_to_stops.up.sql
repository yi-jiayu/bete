alter table stops
    add column tokens tsvector generated always as (setweight(to_tsvector('english', id), 'A') ||
                                                    setweight(to_tsvector('english', coalesce(description, '')), 'B') ||
                                                    setweight(to_tsvector('english', coalesce(road, '')), 'C')) stored;
create index stops_tokens_idx on stops using gin (tokens);
