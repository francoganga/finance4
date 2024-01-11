-- +goose Up
-- +goose StatementBegin
CREATE TABLE transactions (
    id   INTEGER PRIMARY KEY,
    date TEXT NOT NULL,
    code TEXT,
    description TEXT NOT NULL,
    amount INTEGER NOT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE transactions;
-- +goose StatementEnd
