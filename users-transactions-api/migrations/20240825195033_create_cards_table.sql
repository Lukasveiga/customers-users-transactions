-- +goose Up
-- +goose StatementBegin
CREATE TABLE cards (
    id SERIAL PRIMARY KEY,
    account_id INT REFERENCES accounts(id) ON DELETE CASCADE NOT NULL,
    amount BIGINT DEFAULT 0 NOT NULL,
    created_at timestamptz NOT NULL DEFAULT 'now()',
    updated_at timestamptz,
    deleted_at timestamptz
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS cards CASCADE;
-- +goose StatementEnd
