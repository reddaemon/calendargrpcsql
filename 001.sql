create table events
(
    id          serial not null
        constraint events_pk
            primary key,
    title        text   not null,
    description text,
    date        timestamp
);
