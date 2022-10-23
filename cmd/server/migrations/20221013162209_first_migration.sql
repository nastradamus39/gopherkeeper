-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS urls
(
    hash       varchar(256) not null,
    uuid       varchar(256) not null,
    url        text         not null,
    short_url  varchar(256) not null,
    deleted_at date         null
    )
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE urls
-- +goose StatementEnd
