-- +goose Up
-- +goose StatementBegin
ALTER TABLE users
DROP CONSTRAINT chk_role;

ALTER TABLE users
ADD CONSTRAINT chk_role
CHECK (role IN ('user', 'agent', 'admin'));
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE users
DROP CONSTRAINT chk_role;

ALTER TABLE users
ADD CONSTRAINT chk_role
CHECK (role IN ('user', 'agent'));
-- +goose StatementEnd
