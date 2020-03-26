CREATE OR REPLACE FUNCTION stops_tokens_trigger() RETURNS trigger AS
$$
begin
    new.tokens :=
                    setweight(to_tsvector('english', new.id), 'A') ||
                    setweight(to_tsvector('english', coalesce(new.description, '')), 'B') ||
                    setweight(to_tsvector('english', coalesce(new.road, '')), 'C');
    return new;
end
$$ LANGUAGE plpgsql;
