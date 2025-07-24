-- +goose Up
-- +goose StatementBegin
CREATE TABLE listings (
    id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    address TEXT NOT NULL,
    price INT NOT NULL,
    beds INT NOT NULL,
    baths INT NOT NULL,
    sq_ft INT NOT NULL,
    description TEXT,
    agent_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS listings;
-- +goose StatementEnd
