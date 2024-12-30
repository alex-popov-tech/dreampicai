-- +goose Up
-- +goose StatementBegin
CREATE TABLE accounts (
  id SERIAL PRIMARY KEY,
  user_id UUID NOT NULL,
  username TEXT NOT NULL
);

-- +goose StatementEnd
