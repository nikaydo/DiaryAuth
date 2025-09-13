-- +goose Up
-- +goose StatementBegin
CREATE EXTENSION IF NOT EXISTS pgcrypto;
CREATE TABLE users (
	uuid UUID PRIMARY KEY DEFAULT gen_random_uuid(),
	login TEXT NOT NULL UNIQUE,
	password TEXT NOT NULL,
	refresh_token TEXT);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE users;
-- +goose StatementEnd
