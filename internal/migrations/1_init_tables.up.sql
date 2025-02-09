create table if not exists products_price_history(
    id bigserial primary key not null,
    created_at timestamp with time zone default now() not null,
    value numeric not null,
    product_id bigint not null references products(id)
);

comment on table products_price_history is 'История изменения цены товара';
comment on column products_price_history.id is 'Идентификатор цены';
comment on column products_price_history.created_at is 'Дата создания цены';
comment on column products_price_history.value is 'Цена';
comment on column products_price_history.product_id is 'Ссылка на товар';

create table if not exists products(
    id bigserial primary key not null,
    created_at timestamp with time zone default now() not null,
    updated_at timestamp with time zone,
    name varchar(255) not null,
    description text,
    actual_price_id bigint not null references products_price_history(id),
    quantity int not null
);

comment on table products is 'Товары';
comment on column products.id is 'Идентификатор товара';
comment on column products.created_at is 'Дата создания товара';
comment on column products.updated_at is 'Дата последнего обновления товара';
comment on column products.name is 'Наименование товара';
comment on column products.description is 'Описание товара';
comment on column products.actual_price_id is 'Ссылка на актуальную цену товара';
comment on column products.quantity is 'Количество товара';