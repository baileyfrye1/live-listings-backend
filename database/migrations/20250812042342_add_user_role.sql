-- +goose Up
-- +goose StatementBegin
ALTER TABLE users
ADD COLUMN role TEXT NOT NULL DEFAULT 'user';

UPDATE users
SET role = 'user'
WHERE role IS NULL or role = '';

ALTER TABLE users
ADD CONSTRAINT chk_role
CHECK (role IN ('user', 'agent'));
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE users
DROP COLUMN role;
-- +goose StatementEnd
