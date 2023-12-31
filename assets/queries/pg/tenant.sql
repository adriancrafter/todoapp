-- Model: Tenant
-- Table: tenants

-- GetAll
SELECT * FROM tenants WHERE deleted_at IS NULL;

-- GetOne
SELECT * FROM tenants WHERE id = :id AND deleted_at IS NULL LIMIT 1;

-- Create
INSERT INTO tenants (id, tenant_id, slug, name, fantasy_name, profile_id, is_active, owner_id, created_by_id, updated_by_id, created_at, updated_at)
VALUES (:id, :tenant_id, :slug, :name, :fantasy_name, :profile_id, :is_active, :owner_id, :created_by_id, :updated_by_id, :created_at, :updated_at);

-- Update
UPDATE tenants SET tenant_id = :tenant_id, slug = :slug, name = :name, fantasy_name = :fantasy_name, profile_id = :profile_id, is_active = :is_active, owner_id = :owner_id, updated_by_id = :updated_by_id, updated_at = :updated_at
WHERE id = :id;

-- SoftDelete
UPDATE tenants SET deleted_by = :deleted_by, deleted_at = :deleted_at WHERE id = :id;

-- Delete
DELETE FROM tenants WHERE id = :id;

-- Purge
DELETE FROM tenants WHERE deleted_at IS NOT NULL;
