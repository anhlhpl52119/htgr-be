-- +goose Up
-- +goose StatementBegin
-- Uuidv4
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- auto update updated_at
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ language 'plpgsql';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP FUNCTION IF EXISTS update_updated_at_column CASCADE;
DROP EXTENSION IF EXISTS "uuid-ossp";
-- +goose StatementEnd