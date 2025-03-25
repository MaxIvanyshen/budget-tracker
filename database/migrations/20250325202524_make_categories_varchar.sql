-- +goose Up
-- +goose StatementBegin
ALTER TABLE transactions DROP COLUMN category_id;
DROP TABLE IF EXISTS categories;
ALTER TABLE transactions ADD COLUMN category VARCHAR(255);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
