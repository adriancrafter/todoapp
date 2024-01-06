--UP
CREATE TABLE account_roles (
                               id UUID PRIMARY KEY,
                               slug VARCHAR(64) UNIQUE,
                               tenant_id UUID,
                               name VARCHAR(32) UNIQUE,
                               account_id UUID NOT NULL,
                               role_id UUID NOT NULL
);

--UP
ALTER TABLE account_roles
    ADD COLUMN is_active BOOLEAN,
    ADD COLUMN created_by_id UUID references users(id),
    ADD COLUMN updated_by_id UUID references users(id),
    ADD COLUMN deleted_by_id UUID references users(id),
    ADD COLUMN created_at TIMESTAMP,
    ADD COLUMN updated_at TIMESTAMP,
    ADD COLUMN deleted_at TIMESTAMP,
    ADD UNIQUE (tenant_id, account_id, role_id);

--DOWN
DROP TABLE account_roles;
