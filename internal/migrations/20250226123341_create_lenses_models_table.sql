-- +goose Up
-- +goose StatementBegin
create table if not exists lenses_models(
    id bigserial primary key not null,
    created_at timestamp with time zone default now() not null,
    created_by int not null,
    updated_at timestamp with time zone,
    name varchar(255) not null unique,
    is_available boolean default true,
    is_deleted boolean default false,
    deleted_at timestamp with time zone
);

comment on table lenses_models is 'Модели линз';
comment on column lenses_models.id is 'Идентификатор';
comment on column lenses_models.created_at is 'Дата создания';
comment on column lenses_models.created_by is 'Автор добавления';
comment on column lenses_models.updated_at is 'Дата последнего обновления';
comment on column lenses_models.name is 'Наименование';
comment on column lenses_models.is_available is 'Доступна?';
comment on column lenses_models.is_deleted is 'Удалена?';
comment on column lenses_models.deleted_at is 'Дата удаления';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
select 'down SQL query';
-- +goose StatementEnd