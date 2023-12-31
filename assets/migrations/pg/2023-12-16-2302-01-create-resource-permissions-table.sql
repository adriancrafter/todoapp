--UP
CREATE TABLE resource_permissions (
                                      id UUID PRIMARY KEY,
                                      slug VARCHAR(64) UNIQUE,
                                      tenant_id VARCHAR(128),
                                      name VARCHAR(32) UNIQUE,
                                      resource_id UUID NOT NULL,
                                      permission_id UUID NOT NULL
);

--UP
ALTER TABLE resource_permissions
    ADD COLUMN is_active BOOLEAN,
    ADD COLUMN is_deleted BOOLEAN,
    ADD COLUMN created_by_id UUID,
    ADD COLUMN updated_by_id UUID,
    ADD COLUMN created_at TIMESTAMP,
    ADD COLUMN updated_at TIMESTAMP,
    ADD UNIQUE (tenant_id, resource_id, permission_id);

--DOWN
DROP TABLE resource_permissions;
