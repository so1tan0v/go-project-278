-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS links (
    id BIGSERIAL PRIMARY KEY,
    original_url TEXT NOT NULL,
    short_name TEXT NOT NULL UNIQUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS links;
-- +goose StatementEnd

