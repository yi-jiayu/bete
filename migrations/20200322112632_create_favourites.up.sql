create table favourites
(
    user_id integer,
    name    text,
    query   text not null,
    primary key (user_id, name)
);