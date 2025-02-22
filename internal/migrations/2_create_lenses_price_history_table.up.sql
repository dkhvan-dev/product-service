create table if not exists lenses_price_history(
    id bigserial primary key not null,
    created_at timestamp with time zone default now() not null,
    price numeric not null,
    lens_id bigint not null
);

comment on table lenses_price_history is 'История изменения цены линзы';
comment on column lenses_price_history.id is 'Идентификатор цены';
comment on column lenses_price_history.created_at is 'Дата создания цены';
comment on column lenses_price_history.price is 'Цена';
comment on column lenses_price_history.lens_id is 'Ссылка на линзы';