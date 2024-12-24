-- +goose Up
-- +goose StatementBegin
CREATE TABLE accounts (
  id SERIAL PRIMARY KEY,
  userid UUID NOT NULL,
  username TEXT NOT NULL
);
-- +goose StatementEnd
