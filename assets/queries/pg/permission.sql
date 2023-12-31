-- Model: Permission
-- Table: permissions

-- GetAll
SELECT * FROM permissions WHERE tenant_id = :tenant_id AND deleted_at IS NULL;

-- GetOne
SELECT * FROM permissions WHERE tenant_id = :tenant_id AND slug = :slug AND deleted_at IS NULL LIMIT 1;

-- Create
INSERT INTO permissions (tenant_id, slug, name, description, is_active, created_by_id, updated_by_id, created_at, updated_at)
VALUES (:tenant_id, :slug, :name, :description, :is_active, :created_by_id, :updated_by_id, :created_at, :updated_at);

-- Update
UPDATE permissions SET name = :name, description = :description, is_active = :is_active, updated_by_id = :updated_by_id, updated_at = :updated_at
WHERE tenant_id = :tenant_id AND slug = :slug;

-- SoftDelete
UPDATE permissions SET deleted_by = :deleted_by, deleted_at = :deleted_at WHERE tenant_id = :tenant_id AND slug = :slug;

-- Delete
DELETE FROM permissions WHERE tenant_id = :tenant_id AND slug = :slug;

-- Purge
DELETE FROM permissions WHERE tenant_id = :tenant_id AND deleted_at IS NOT NULL;
