-- +goose Up
-- +goose StatementBegin
ALTER TABLE listings
    ADD COLUMN views INT NOT NULL DEFAULT 0;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE listings
DROP COLUMN views;
-- +goose StatementEnd
