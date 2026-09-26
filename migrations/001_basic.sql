-- +goose Up
CREATE TABLE links (
    id          BIGSERIAL PRIMARY KEY,
    code        VARCHAR(10) NOT NULL UNIQUE,
    url         TEXT NOT NULL,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT now(),
    click_count   BIGINT NOT NULL DEFAULT 0

);

-- +goose Down
DROP TABLE links;