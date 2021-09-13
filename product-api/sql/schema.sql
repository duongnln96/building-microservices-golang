DROP TABLE IF EXISTS products;

create table if not exists products (
    id              serial primary key,
    name            varchar(255) not null,
    description     text,
    price           numeric not null,
    sku             varchar(255) not null
);
