--UP
CREATE TABLE role_permissions (
                                  id UUID PRIMARY KEY,
                                  slug VARCHAR(64) UNIQUE,
                                  tenant_id UUID,
                                  name VARCHAR(32) UNIQUE,
                                  role_id UUID NOT NULL,
                                  permission_id UUID NOT NULL
);

--UP
ALTER TABLE role_permissions
    ADD COLUMN is_active BOOLEAN,
    ADD COLUMN created_by_id UUID references users(id),
    ADD COLUMN updated_by_id UUID references users(id),
    ADD COLUMN deleted_by_id UUID references users(id),
    ADD COLUMN created_at TIMESTAMP,
    ADD COLUMN updated_at TIMESTAMP,
    ADD COLUMN deleted_at TIMESTAMP,
    ADD UNIQUE (tenant_id, role_id, permission_id);

--DOWN
DROP TABLE role_permissions;
