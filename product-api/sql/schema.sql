DROP TABLE IF EXISTS users;
DROP TABLE IF EXISTS products;

create table if not exists users (
    id              serial primary key,
    name            varchar(255),
    email           varchar(255) not null unique,
    password        varchar(255) not null unique
);

create table if not exists products (
    id              serial primary key,
    name            varchar(255) not null,
    description     text,
    price           numeric not null,
    sku             varchar(255) not null
);
