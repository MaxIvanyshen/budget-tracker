-- +goose Up
-- +goose StatementBegin
ALTER TABLE categories DROP COLUMN color;
ALTER TABLE categories ADD COLUMN color VARCHAR(64) NOT NULL DEFAULT 'yellow-400';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
