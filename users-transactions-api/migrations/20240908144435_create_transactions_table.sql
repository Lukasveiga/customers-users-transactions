-- +goose Up
-- +goose StatementBegin
CREATE TABLE transactions (
    id SERIAL PRIMARY KEY,
    card_id INT REFERENCES cards(id) ON DELETE CASCADE NOT NULL,
    kind VARCHAR(145) NOT NULL,
    value BIGINT NOT NULL,
    created_at timestamptz NOT NULL DEFAULT 'now()',
    updated_at timestamptz,
    deleted_at timestamptz
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS transactions CASCADE;
-- +goose StatementEnd
