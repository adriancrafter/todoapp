--UP
CREATE TABLE account_roles (
                               id UUID PRIMARY KEY,
                               slug VARCHAR(64) UNIQUE,
                               tenant_id VARCHAR(128),
                               name VARCHAR(32) UNIQUE,
                               account_id UUID NOT NULL,
                               role_id UUID NOT NULL
);

--UP
ALTER TABLE account_roles
    ADD COLUMN is_active BOOLEAN,
    ADD COLUMN is_deleted BOOLEAN,
    ADD COLUMN created_by_id UUID,
    ADD COLUMN updated_by_id UUID,
    ADD COLUMN created_at TIMESTAMP,
    ADD COLUMN updated_at TIMESTAMP,
    ADD UNIQUE (tenant_id, account_id, role_id);

--DOWN
DROP TABLE account_roles;
