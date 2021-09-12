drop table users if exists;

create table if not exists users (
    id          serial primary key,
    name        varchar(255),
    email       varchar(255) not null unique,
    password    varchar(255) not null unique
);