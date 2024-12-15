-- +goose NO TRANSACTION
-- +goose Up
-- +goose StatementBegin
CREATE DATABASE user_registration;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP DATABASE user_registration;
-- +goose StatementEnd
