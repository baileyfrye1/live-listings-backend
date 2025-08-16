-- +goose Up
-- +goose StatementBegin
CREATE TABLE listings (
    id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    address TEXT NOT NULL,
    price INT NOT NULL CHECK (price > 0),
    beds INT NOT NULL CHECK (beds >= 0),
    baths INT NOT NULL CHECK (baths >= 0),
    sq_ft INT NOT NULL CHECK (sq_ft > 0),
    description TEXT,
    agent_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- indexes
CREATE INDEX idx_listings_agent_id ON listings(agent_id);
CREATE INDEX idx_listings_price ON listings(price);
CREATE INDEX idx_listings_address ON listings(address);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS listings;
-- +goose StatementEnd
