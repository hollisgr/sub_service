-- +goose Up
-- +goose StatementBegin
CREATE TABLE subscriptions (
    id SERIAL PRIMARY KEY,
    service_name TEXT,
    price INTEGER,
    user_id TEXT,
    start_date DATE,
    end_date DATE
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS subscriptions
-- +goose StatementEnd
