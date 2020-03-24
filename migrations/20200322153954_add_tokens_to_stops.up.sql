begin;

alter table stops
    add column tokens tsvector;
create index stops_tokens_idx on stops using gin (tokens);

CREATE FUNCTION stops_tokens_trigger() RETURNS trigger AS
$$
begin
    new.tokens :=
                    setweight(to_tsvector('english', new.id), 'A') ||
                    setweight(to_tsvector('english', coalesce(new.description, '')), 'B') ||
                    setweight(to_tsvector('english', coalesce(new.road, '')), 'C');
    return new;
end
$$ LANGUAGE plpgsql;

CREATE TRIGGER tsvectorupdate
    BEFORE INSERT OR UPDATE
    ON stops
    FOR EACH ROW
EXECUTE FUNCTION stops_tokens_trigger();

end;