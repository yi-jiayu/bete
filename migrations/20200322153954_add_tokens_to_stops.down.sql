begin;

alter table stops
    drop column tokens;
drop trigger tsvectorupdate on stops;
drop function stops_tokens_trigger;

end;
