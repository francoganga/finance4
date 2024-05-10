-- +goose Up
-- +goose StatementBegin
CREATE TABLE label (
    id INTEGER PRIMARY KEY,
    name TEXT NOT NULL
);
-- +goose StatementEnd

-- +goose StatementBegin
ALTER TABLE transactions ADD COLUMN label_id INTEGER;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS label;
-- +goose StatementEnd

-- +goose StatementBegin
ALTER TABLE transactions DROP COLUMN label_id;
-- +goose StatementEnd
