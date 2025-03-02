-- +goose Up
-- +goose StatementBegin
alter table if exists lenses
    add column if not exists quantity int not null default 0;

comment on column lenses.quantity is 'Количество доступных пар линз';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
select 'down SQL query';
-- +goose StatementEnd