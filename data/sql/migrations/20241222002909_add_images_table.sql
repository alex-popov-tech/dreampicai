-- +goose Up
-- +goose StatementBegin
DROP TYPE IF EXISTS image_status;

CREATE TYPE image_status AS ENUM('started', 'cancelled', 'failed', 'succeeded');

CREATE TABLE IF NOT EXISTS images (
  id SERIAL PRIMARY KEY,
  provider_id TEXT NOT NULL,
  owner_id INTEGER REFERENCES accounts (id),
  status image_status NOT NULL DEFAULT 'started',
  prompt TEXT NOT NULL,
  negative_prompt TEXT NOT NULL,
  model TEXT NOT NULL,
  url TEXT,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- +goose StatementEnd
