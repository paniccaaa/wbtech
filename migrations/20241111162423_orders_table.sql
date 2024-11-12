-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS orders
(
    order_uid TEXT PRIMARY KEY,
    order_data JSONB
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS orders;
-- +goose StatementEnd
