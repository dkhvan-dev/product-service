create table if not exists lenses(
    id bigserial primary key not null,
    created_at timestamp with time zone default now() not null,
    created_by int not null,
    updated_at timestamp with time zone,
    name varchar(255) not null,
    model_id bigint not null references lenses_models(id),
    brand varchar(255) not null,
    description text,
    actual_price_id bigint references lenses_price_history(id),
    color varchar(255) not null,
    optical_power numeric(4, 2) not null,
    diameter numeric(4, 2) not null,
    curvature_radius numeric(4, 2) not null,
    is_available boolean default true,
    is_deleted boolean default false,
    deleted_at timestamp with time zone
);

comment on table lenses is 'Линзы';
comment on column lenses.id is 'Идентификатор линзы';
comment on column lenses.created_at is 'Дата создания';
comment on column lenses.created_by is 'Автор добавления';
comment on column lenses.updated_at is 'Дата последнего обновления';
comment on column lenses.name is 'Наименование';
comment on column lenses.model_id is 'Ссылка на модель линз';
comment on column lenses.brand is 'Бренд';
comment on column lenses.description is 'Описание';
comment on column lenses.actual_price_id is 'Ссылка на актуальную цену';
comment on column lenses.color is 'Цвет';
comment on column lenses.optical_power is 'Оптическая сила';
comment on column lenses.diameter is 'Диаметр';
comment on column lenses.curvature_radius is 'Радиус кривизны';
comment on column lenses.is_available is 'Доступна?';
comment on column lenses.is_deleted is 'Удалена?';
comment on column lenses.deleted_at is 'Дата удаления';

alter table if exists lenses_price_history
    add constraint fk_lenses_price_history_to_lenses
        foreign key (lens_id) references lenses(id) on delete cascade;