-- +goose Up
-- +goose StatementBegin
CREATE TABLE label (
    id INTEGER PRIMARY KEY,
    name TEXT NOT NULL
);
-- +goose StatementEnd

-- +goose StatementBegin
CREATE TABLE transaction_label (
    id INTEGER PRIMARY KEY,
    transaction_id INTEGER NOT NULL,
    label_id INTEGER NOT NULL,
    FOREIGN KEY(transaction_id) REFERENCES transactions(id),
    FOREIGN KEY(label_id) REFERENCES label(id),
    UNIQUE(transaction_id, label_id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS label;
-- +goose StatementEnd

-- +goose StatementBegin
DROP TABLE IF EXISTS transaction_label;
-- +goose StatementEnd
