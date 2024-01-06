--UP
CREATE TABLE roles (
                       id UUID PRIMARY KEY,
                       slug VARCHAR(64) UNIQUE,
                       tenant_id UUID,
                       name VARCHAR(32) UNIQUE,
                       description TEXT NULL
);

--UP
ALTER TABLE roles
    ADD COLUMN is_active BOOLEAN,
    ADD COLUMN created_by_id UUID references users(id),
    ADD COLUMN updated_by_id UUID references users(id),
    ADD COLUMN deleted_by_id UUID references users(id),
    ADD COLUMN created_at TIMESTAMP,
    ADD COLUMN updated_at TIMESTAMP,
    ADD COLUMN deleted_at TIMESTAMP;

--DOWN
DROP TABLE roles;
