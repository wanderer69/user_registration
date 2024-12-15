-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS users
(
	uuid                        TEXT      NOT NULL PRIMARY KEY,
	login                       TEXT      NOT NULL DEFAULT '',
	email                       TEXT      NOT NULL DEFAULT '',
    registration_code           TEXT      NULL,
    hash                        TEXT      NOT NULL DEFAULT '',
	access_code                 TEXT      NOT NULL DEFAULT '',

	created_at                  TIMESTAMP NOT NULL,
	confirmed_at                TIMESTAMP NULL,
	updated_at                  TIMESTAMP NOT NULL,
	deleted_at                  TIMESTAMP NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS users CASCADE;
-- +goose StatementEnd
